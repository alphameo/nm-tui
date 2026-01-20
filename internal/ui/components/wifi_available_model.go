package components

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiState int

const (
	Scanning wifiState = iota
	Connecting
	None
)

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

func (s *wifiState) String() string {
	switch *s {
	case Scanning:
		return "Scanning"
	case Connecting:
		return "Connecting"
	case None:
		return "󰄬"
	default:
		return "Undefined!!!"
	}
}

type WifiAvailableModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   wifiState

	connector *WifiConnectorModel

	nm infra.NetworkManager

	width  int
	height int
}

func NewWifiAvailableModel(networkManager infra.NetworkManager) *WifiAvailableModel {
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
	con := NewWifiConnector(networkManager)
	m := &WifiAvailableModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   Scanning,
		connector:        con,
		nm:               networkManager,
	}
	return m
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
	return m.UpdateRowsCmd()
}

func (m *WifiAvailableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if m.indicatorState != None {
				return m, nil
			}
			return m, m.UpdateRowsCmd()
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, m.callConnector(row[ssidAvailable])
			}
			return m, nil
		}
	case WifiIndicatorStateMsg:
		return m, m.setWifiStateCmd(wifiAvailableState(msg))
	case updateWifiAvailableMsg:
		return m, m.UpdateRowsCmd()
	}

	var cmd tea.Cmd
	if m.indicatorState != None {
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
	if m.indicatorState != None {
		statusline = fmt.Sprintf("%s %s", m.indicatorState.String(), m.indicatorSpinner.View())
	} else {
		statusline = m.indicatorState.String()
	}
	return lipgloss.JoinVertical(lipgloss.Center, view, statusline)
}

func (m *WifiAvailableModel) UpdateRowsCmd() tea.Cmd {
	return tea.Sequence(
		m.setWifiStateCmd(Scanning),
		func() tea.Msg {
			list, err := m.nm.GetAvailableWifi()
			if err != nil {
				logger.Errln(fmt.Errorf("error: %s", err.Error()))
				return NotifyCmd(err.Error())
			}
			rows := []table.Row{}
			for _, wifiNet := range list {
				var connectionFlag string
				if wifiNet.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{connectionFlag, wifiNet.SSID, wifiNet.Security, fmt.Sprint(wifiNet.Signal)})
			}

			m.dataTable.SetRows(rows)
			m.dataTable.UpdateViewport()
			return UpdateMsg
		},
		m.setWifiStateCmd(None),
	)
}

type updateWifiAvailableMsg struct{}

// UpdateWifiAvailableMsg is used to avoid extra instantiatons
var UpdateWifiAvailableMsg = updateWifiAvailableMsg{}

func UpdateWifiAvailableCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateWifiAvailableMsg
	}
}

type WifiIndicatorStateMsg wifiState

func (m *WifiAvailableModel) setWifiStateCmd(state wifiAvailableState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return nil
	}
	if state == None {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func (m *WifiAvailableModel) callConnector(wifiName string) tea.Cmd {
	m.connector.setNew(wifiName)
	return tea.Sequence(
		SetPopupActivityCmd(true),
		SetPopupContentCmd(m.connector, "Wi-Fi Connector"),
	)
}

func SetWifiIndicatorStateCmd(state wifiState) tea.Cmd {
	return func() tea.Msg {
		return WifiIndicatorStateMsg(state)
	}
}
