package connections

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/controls"
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
	signalColWidth      int     = 3
	conFlagColWidth             = 1
	securityWidthPart   float32 = 0.3
	minSecurityColWidth         = 8
	minSsidWidth                = 4
	indicatorHeight             = 1
)

func (s *wifiState) String() string {
	switch *s {
	case Scanning:
		return "Scanning"
	case Connecting:
		return "Connecting"
	case None:
		return ""
	default:
		return "Undefined!!!"
	}
}

type WifiAvailableModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   wifiState
	connector        WifiConnectorModel
	pSsidCol         *table.Column
	pSecurityCol     *table.Column
}

func NewWifiAvailable() *WifiAvailableModel {
	cols := []table.Column{
		{Title: "󱘖", Width: conFlagColWidth},
		{Title: "SSID"},
		{Title: "Security"},
		{Title: "", Width: signalColWidth},
	}
	ssidCol := &cols[1]
	securityCol := &cols[2]
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)
	s := spinner.New()
	con := *NewWifiConnector()
	m := &WifiAvailableModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   Scanning,
		connector:        con,
		pSsidCol:         ssidCol,
		pSecurityCol:     securityCol,
	}
	return m
}

func (m *WifiAvailableModel) Resize(width, height int) {
	width -= styles.BorderOffset
	height -= styles.BorderOffset
	height -= indicatorHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	security := max(int(float32(width)*securityWidthPart), minSecurityColWidth)
	ssidWidth := width - signalColWidth - tableUtilityOffset - conFlagColWidth - security
	m.pSecurityCol.Width = security
	m.pSsidCol.Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiAvailableModel) Init() tea.Cmd {
	return m.UpdateRows()
}

func (m *WifiAvailableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if m.indicatorState != None {
				return m, nil
			}
			return m, m.UpdateRows()
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				m.connector.setNew(row[1])
				return m, tea.Sequence(
					controls.SetPopupActivity(true),
					controls.SetPopupContent(m.connector),
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
			cmd = m.UpdateRows()
		} else {
			cmd = tea.Sequence(
				controls.Notify(msg.err.Error()),
			)
		}
		return m, tea.Sequence(cmd, SetWifiIndicatorState(None))
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

	var symbol string
	if m.indicatorState != None {
		symbol = fmt.Sprintf("%s %s", m.indicatorState.String(), m.indicatorSpinner.View())
	} else {
		symbol = "󰄬"
	}
	statusline := lipgloss.Place(m.dataTable.Width(), 1, lipgloss.Center, lipgloss.Center, symbol)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n%s", view, statusline)
	return sb.String()
}

type scannedRowsMsg []table.Row

func (m WifiAvailableModel) UpdateRows() tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorState(Scanning),
		func() tea.Msg {
			list, err := nmcli.WifiScan()
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
		SetWifiIndicatorState(None))
}

type WifiIndicatorStateMsg wifiState

func SetWifiIndicatorState(state wifiState) tea.Cmd {
	return func() tea.Msg {
		return WifiIndicatorStateMsg(state)
	}
}

type WifiConnectionMsg struct {
	err  error
	ssid string
}

func WifiConnect(ssid, password string) tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorState(Connecting),
		func() tea.Msg {
			err := nmcli.WifiConnect(ssid, password)
			return WifiConnectionMsg{err: err, ssid: ssid}
		},
		SetWifiIndicatorState(None))
}
