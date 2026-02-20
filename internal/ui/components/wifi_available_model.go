package components

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiAvailableState int

const (
	ScanningAvailable wifiAvailableState = iota
	ConnectingAvailable
	DoneInAvailable
)

func (s *wifiAvailableState) String() string {
	switch *s {
	case ScanningAvailable:
		return "Scanning"
	case ConnectingAvailable:
		return "Connecting"
	case DoneInAvailable:
		return "󰄬"
	default:
		return "Undefined!!!"
	}
}

type wifiAvailableColumnIndex int

const (
	ssidAvailable  wifiAvailableColumnIndex = 1
	securityColumn wifiAvailableColumnIndex = 2
)

const (
	signalColumnWidth         int     = 3
	connectionFlagColumnWidth int     = 1
	securityWidthProportion   float32 = 0.3
	minSecurityColumnWidth    int     = 8
	minSSIDWidth              int     = 4
	indicatorStateHeight      int     = 1
)

type wifiAvailableKeyMap struct {
	rescan        key.Binding
	openConnector key.Binding
}

func (k *wifiAvailableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.rescan, k.openConnector}
}

func (k *wifiAvailableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.rescan, k.openConnector}}
}

type WifiAvailableModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   wifiAvailableState

	connector *WifiConnectorModel

	keys *wifiAvailableKeyMap

	nm infra.NetworkManager

	width  int
	height int
}

func NewWifiAvailableModel(wifiConnector *WifiConnectorModel, keys *wifiAvailableKeyMap, networkManager infra.NetworkManager) *WifiAvailableModel {
	cols := []table.Column{
		{Title: "󱘖", Width: connectionFlagColumnWidth},
		{Title: "SSID"},
		{Title: "Security"},
		{Title: "", Width: signalColumnWidth},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)
	s := spinner.New()
	return &WifiAvailableModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   DoneInAvailable,
		connector:        wifiConnector,
		keys:             keys,
		nm:               networkManager,
	}
}

func (m *WifiAvailableModel) Resize(width, height int) {
	m.width = width
	m.height = height

	height -= indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	security := max(int(float32(width)*securityWidthProportion), minSecurityColumnWidth)
	ssidWidth := width - signalColumnWidth - tableUtilityOffset - connectionFlagColumnWidth - security
	m.dataTable.Columns()[securityColumn].Width = security
	m.dataTable.Columns()[ssidAvailable].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiAvailableModel) Width() int {
	return m.width
}

func (m *WifiAvailableModel) Height() int {
	return m.height
}

func (m *WifiAvailableModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiAvailableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.rescan):
			if m.indicatorState != DoneInAvailable {
				return m, nil
			}
			return m, m.RescanCmd()
		case key.Matches(msg, m.keys.openConnector):
			row := m.dataTable.SelectedRow()
			if row != nil {
				slog.Debug("call")
				return m, m.callConnector(row[ssidAvailable])
			}
			return m, nil
		}
	case WifiAvialableStateMsg:
		return m, m.setStateCmd(wifiAvailableState(msg))
	case RescanWifiAvailableMsg:
		return m, m.RescanCmd()
	}

	var cmd tea.Cmd
	if m.indicatorState != DoneInAvailable {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}
	m.dataTable, cmd = m.dataTable.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiAvailableModel) View() string {
	view := m.dataTable.View()

	var statusline string
	if m.indicatorState != DoneInAvailable {
		statusline = fmt.Sprintf(
			"%s %s",
			m.indicatorState.String(),
			m.indicatorSpinner.View(),
		)
	} else {
		statusline = m.indicatorState.String()
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)
}

func (m *WifiAvailableModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ScanningAvailable),
		func() tea.Msg {
			list, err := m.nm.GetAvailableWifi()
			if err != nil {
				slog.Error(err.Error())
				return NotifyCmd(err.Error())
			}
			rows := []table.Row{}
			for _, wifiNet := range list {
				var connectionFlag string
				if wifiNet.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{
					connectionFlag,
					wifiNet.SSID,
					wifiNet.Security,
					strconv.Itoa(wifiNet.Signal),
				})
			}

			m.dataTable.SetRows(rows)
			m.dataTable.GotoTop()
			m.dataTable.UpdateViewport()
			return tea.BatchMsg{
				NilCmd,
				m.setStateCmd(DoneInAvailable),
			}
		},
	)
}

type RescanWifiAvailableMsg struct{}

func RescanWifiAvailableCmd() tea.Cmd {
	return func() tea.Msg {
		return RescanWifiAvailableMsg{}
	}
}

type WifiAvialableStateMsg wifiAvailableState

func (m *WifiAvailableModel) setStateCmd(state wifiAvailableState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return nil
	}
	if state == DoneInAvailable {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func SetWifiAvailableStateCmd(state wifiAvailableState) tea.Cmd {
	return func() tea.Msg {
		return WifiAvialableStateMsg(state)
	}
}

func (m *WifiAvailableModel) callConnector(wifiName string) tea.Cmd {
	m.connector.setNew(wifiName)
	slog.Debug("named")
	return OpenPopup(m.connector, "Wi-Fi Connector")
}
