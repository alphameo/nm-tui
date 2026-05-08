package components

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

var wifiAvailableKeys = &wifiAvailableKeyMap{
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan"),
	),
	openConnector: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open connector"),
	),
}

type WifiAvailableModel struct {
	dataTable table.Model

	connColIdx     int
	ssidColIdx     int
	securityColIdx int
	signalColIdx   int

	securityWidthProportion float32
	minSecurityColumnWidth  int
	minSSIDWidth            int

	indicatorStateHeight int

	indicatorSpinner spinner.Model
	indicatorState   wifiAvailableState

	connector *WifiConnectorModel

	keys *wifiAvailableKeyMap

	wm infra.WifiManager

	width  int
	height int
}

func NewWifiAvailableModel(wifiConnector *WifiConnectorModel, keys *wifiAvailableKeyMap, wifiManager infra.WifiManager) *WifiAvailableModel {
	cols := []table.Column{
		{Title: "󱘖", Width: 1},
		{Title: "SSID"},
		{Title: "Security"},
		{Title: "", Width: 3},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(styles.TableStyle),
	)

	s := spinner.New()

	model := &WifiAvailableModel{
		dataTable: t,

		connColIdx:     0,
		ssidColIdx:     1,
		securityColIdx: 2,
		signalColIdx:   3,

		securityWidthProportion: 0.3,
		minSecurityColumnWidth:  8,
		minSSIDWidth:            4,

		indicatorSpinner: s,
		indicatorState:   DoneInAvailable,

		connector: wifiConnector,
		wm:        wifiManager,
		keys:      keys,
	}

	model.bakeSizes()

	return model
}

func (m *WifiAvailableModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *WifiAvailableModel) Resize(width, height int) {
	m.width = width
	m.height = height

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	secColWidth := max(int(float32(width)*m.securityWidthProportion), m.minSecurityColumnWidth)
	signalColWidth := m.dataTable.Columns()[m.signalColIdx].Width
	conColWidth := m.dataTable.Columns()[m.connColIdx].Width

	ssidWidth := width - signalColWidth - tableUtilityOffset - conColWidth - secColWidth
	m.dataTable.Columns()[m.securityColIdx].Width = secColWidth
	m.dataTable.Columns()[m.ssidColIdx].Width = ssidWidth
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
		return m.handleKey(msg)
	case WifiAvialableStateMsg:
		return m, m.setStateCmd(wifiAvailableState(msg))
	case RescanWifiAvailableMsg:
		time.Sleep(msg.delay)
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

func (m *WifiAvailableModel) handleKey(keyMsg tea.KeyMsg) (*WifiAvailableModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.rescan):
		if m.indicatorState != DoneInAvailable {
			return m, nil
		}
		return m, m.RescanCmd()
	case key.Matches(keyMsg, m.keys.openConnector):
		row := m.dataTable.SelectedRow()
		if row != nil {
			return m, m.callConnector(row[m.ssidColIdx])
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(keyMsg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiAvailableModel) View() string {
	view := m.dataTable.View()

	statusline := m.indicatorView()
	return lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)
}

func (m *WifiAvailableModel) indicatorView() string {
	var view string
	if m.indicatorState != DoneInAvailable {
		view = fmt.Sprintf(
			"%s %s",
			m.indicatorState.String(),
			m.indicatorSpinner.View(),
		)
	} else {
		view = m.indicatorState.String()
	}
	return view
}

func (m *WifiAvailableModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ScanningAvailable),
		func() tea.Msg {
			list, err := m.wm.ScanWifis(context.Background())
			if err != nil {
				return tea.BatchMsg{
					m.setStateCmd(DoneInAvailable),
					NotifyCmd("Cannot scan available wifi networks"),
				}
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
			return m.setStateCmd(DoneInAvailable)
		},
	)
}

type RescanWifiAvailableMsg struct {
	delay time.Duration
}

func RescanWifiAvailableCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiAvailableMsg{delay: delay}
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
	return tea.Batch(
		m.connector.setNew(wifiName),
		OpenPopup(m.connector, "Wi-Fi network Connector"),
	)
}
