package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type wifiAvailableConfig struct {
	connColIdx              int
	ssidColIdx              int
	securityColIdx          int
	signalColIdx            int
	securityWidthProportion float32
	minSignalColWidth       int
}

var wifiAvailableCfg = wifiAvailableConfig{
	connColIdx:     0,
	ssidColIdx:     1,
	securityColIdx: 2,
	signalColIdx:   3,

	securityWidthProportion: 0.3,
	minSignalColWidth:       3,
}

type wifiAvailableState int

const (
	AvailableNil wifiAvailableState = iota
	AvailableScanning
	AvailableConnecting
	AvailableCreating
	AvailableDone
)

func (s *wifiAvailableState) String() string {
	switch *s {
	case AvailableScanning:
		return "Scanning"
	case AvailableConnecting:
		return "Connecting"
	case AvailableCreating:
		return "Creating Connection"
	case AvailableDone:
		return styles.SymbolCheck
	default:
		return "Undefined"
	}
}

type wifiAvailableKeyMap struct {
	rescan  key.Binding
	connect key.Binding
}

type WifiAvailableModel struct {
	dataTable table.Model

	indicatorSpinner     spinner.Model
	indicatorState       wifiAvailableState
	indicatorStateHeight int

	focus bool

	keys wifiAvailableKeyMap

	wm infra.WifiManager

	width  int
	height int
}

func NewWifiAvailableModel(keys wifiAvailableKeyMap, wifiManager infra.WifiManager) *WifiAvailableModel {
	cols := make([]table.Column, 4)
	cols[wifiAvailableCfg.connColIdx] = table.Column{
		Title: styles.SymbolConnection,
		Width: len(styles.SymbolConnection),
	}
	cols[wifiAvailableCfg.ssidColIdx] = table.Column{Title: "SSID"}
	cols[wifiAvailableCfg.securityColIdx] = table.Column{Title: "Security"}
	cols[wifiAvailableCfg.signalColIdx] = table.Column{
		Title: styles.SymbolSignal,
		Width: max(wifiAvailableCfg.minSignalColWidth, len(styles.SymbolSignal)),
	}

	initTableStyle := styles.DataTableStyle
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(initTableStyle),
	)

	s := spinner.New()
	s.Spinner = styles.Spinner

	model := &WifiAvailableModel{
		dataTable: t,

		indicatorSpinner: s,
		indicatorState:   AvailableDone,

		keys: keys,
		wm:   wifiManager,
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

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	secColWidth := int(float32(width) * wifiAvailableCfg.securityWidthProportion)
	signalColWidth := m.dataTable.Columns()[wifiAvailableCfg.signalColIdx].Width
	conColWidth := m.dataTable.Columns()[wifiAvailableCfg.connColIdx].Width
	ssidWidth := width - signalColWidth - tableUtilityOffset - conColWidth - secColWidth

	m.dataTable.Columns()[wifiAvailableCfg.securityColIdx].Width = secColWidth
	m.dataTable.Columns()[wifiAvailableCfg.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiAvailableModel) Width() int {
	return m.width
}

func (m *WifiAvailableModel) Height() int {
	return m.height
}

func (m *WifiAvailableModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.SetStyles(styles.TableStyle)
	m.dataTable.Focus()
	return nil
}

func (m *WifiAvailableModel) Blur() {
	m.focus = false
	m.dataTable.Blur()
	m.dataTable.SetStyles(styles.DataTableStyle)
}

func (m *WifiAvailableModel) Focused() bool {
	return m.focus
}

func (m *WifiAvailableModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiAvailableModel) Update(msg tea.Msg) (*WifiAvailableModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.rescan):
			if m.indicatorState != AvailableDone {
				return m, nil
			}
			return m, m.RescanCmd()
		case key.Matches(msg, m.keys.connect):
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, OpenConnectorCmd(row[wifiAvailableCfg.ssidColIdx])
			}
			return m, nil
		}
	case WifiAvialableStateMsg:
		return m, m.setStateCmd(wifiAvailableState(msg))
	case RescanWifiAvailableMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != AvailableDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiAvailableModel) View() string {
	view := m.dataTable.View()
	statusline := m.indicatorView()
	view = lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)

	var style lipgloss.Style
	if m.focus {
		style = styles.BorderedFocusedStyle
	} else {
		style = styles.BorderedStyle
	}
	view = renderer.RenderWithTitleAndKeybind(
		view,
		"Available networks",
		"1",
		style,
		styles.AccentColor,
	)
	return view
}

func (m *WifiAvailableModel) indicatorView() string {
	var view string
	if m.indicatorState != AvailableDone {
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
		m.setStateCmd(AvailableScanning),
		func() tea.Msg {
			list, err := m.wm.ScanWifis(context.Background())
			if err != nil {
				return tea.Batch(
					m.setStateCmd(AvailableDone),
					NotifyCmd("Cannot scan available wifi networks"),
				)
			}
			rows := []table.Row{}
			for _, wifiNet := range list {
				var connectionFlag string
				if wifiNet.Active {
					connectionFlag = styles.SymbolCheck
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
			return m.setStateCmd(AvailableDone)
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
	if state == AvailableDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func SetWifiAvailableStateCmd(state wifiAvailableState) tea.Cmd {
	return func() tea.Msg {
		return WifiAvialableStateMsg(state)
	}
}
