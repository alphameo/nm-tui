package views

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	wifiState                int
	wifiAvailableColumnIndex int
)

const (
	Scanning wifiState = iota
	Connecting
	None
	signalColWidth      int                      = 3
	conFlagColWidth     int                      = 1
	securityWidthPart   float32                  = 0.3
	minSecurityColWidth int                      = 8
	minSSIDWidth        int                      = 4
	indicatorHeight     int                      = 1
	ssidAvailCol        wifiAvailableColumnIndex = 1
	securityCol         wifiAvailableColumnIndex = 2
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
	connector        WifiConnectorModel
	nm               infra.NetworkManager
	width            int
	height           int
}

func NewWifiAvailable(networkManager infra.NetworkManager) *WifiAvailableModel {
	cols := []table.Column{
		{Title: "󱘖", Width: conFlagColWidth},
		{Title: "SSID"},
		{Title: "Security"},
		{Title: "", Width: signalColWidth},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)
	s := spinner.New()
	con := *NewWifiConnector(networkManager)
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

	height -= indicatorHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	security := max(int(float32(width)*securityWidthPart), minSecurityColWidth)
	ssidWidth := width - signalColWidth - tableUtilityOffset - conFlagColWidth - security
	m.dataTable.Columns()[securityCol].Width = security
	m.dataTable.Columns()[ssidAvailCol].Width = ssidWidth
	m.dataTable.UpdateViewport()
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
				m.connector.setNew(row[1])
				return m, tea.Sequence(
					SetPopupActivityCmd(true),
					SetPopupContentCmd(m.connector, "Wi-Fi Connector"),
				)
			}
			return m, nil
		}
	case scannedRowsMsg:
		m.dataTable.SetRows(msg)
		return m, nil
	case WifiIndicatorStateMsg:
		m.indicatorState = wifiState(msg)
		if m.indicatorState == None {
			return m, nil
		}
		return m, m.indicatorSpinner.Tick
	case WifiConnectionMsg:
		var cmd tea.Cmd
		if msg.err == nil {
			cmd = m.UpdateRowsCmd()
		} else {
			cmd = NotifyCmd(msg.err.Error())
		}
		return m, tea.Sequence(cmd, SetWifiIndicatorStateCmd(None))
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
	statusline = lipgloss.Place(m.dataTable.Width(), 1, lipgloss.Center, lipgloss.Center, statusline)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n%s", view, statusline)
	return sb.String()
}

type scannedRowsMsg []table.Row

func (m WifiAvailableModel) UpdateRowsCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorStateCmd(Scanning),
		func() tea.Msg {
			list, err := m.nm.GetAvailableWifi()
			if err != nil {
				logger.Errln(fmt.Errorf("error: %s", err.Error()))
			}
			rows := []table.Row{}
			for _, wifiNet := range list {
				var connectionFlag string
				if wifiNet.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{connectionFlag, wifiNet.SSID, wifiNet.Security, fmt.Sprint(wifiNet.Signal)})
			}
			return scannedRowsMsg(rows)
		},
		SetWifiIndicatorStateCmd(None))
}

type WifiIndicatorStateMsg wifiState

func SetWifiIndicatorStateCmd(state wifiState) tea.Cmd {
	return func() tea.Msg {
		return WifiIndicatorStateMsg(state)
	}
}
