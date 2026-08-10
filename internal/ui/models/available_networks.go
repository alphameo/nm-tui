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

type availableNetworksConfig struct {
	connColIdx              int
	ssidColIdx              int
	securityColIdx          int
	signalColIdx            int
	securityWidthProportion float32
	minSignalColWidth       int
}

var availableNetworksCfg = availableNetworksConfig{
	connColIdx:     0,
	ssidColIdx:     1,
	securityColIdx: 2,
	signalColIdx:   3,

	securityWidthProportion: 0.3,
	minSignalColWidth:       3,
}

type availableNetworksState int

const (
	AvailableNetsNil availableNetworksState = iota
	AvailableNetsScanning
	AvailableNetsConnecting
	AvailableNetsCreating
	AvailableNetsDone
)

func (s *availableNetworksState) String() string {
	switch *s {
	case AvailableNetsScanning:
		return "Scanning"
	case AvailableNetsConnecting:
		return "Connecting"
	case AvailableNetsCreating:
		return "Creating Connection Profile"
	case AvailableNetsDone:
		return styles.SymbolCheck
	default:
		return "Undefined"
	}
}

type availableNetworksKeyMap struct {
	rescan  key.Binding
	connect key.Binding
}

type AvailableNetworksModel struct {
	dataTable table.Model

	indicatorSpinner     spinner.Model
	indicatorState       availableNetworksState
	indicatorStateHeight int

	focus bool

	keys availableNetworksKeyMap

	netMngr infra.NetworksManager

	width  int
	height int
}

func NewAvailableNetworksModel(keys availableNetworksKeyMap, networksManager infra.NetworksManager) *AvailableNetworksModel {
	cols := make([]table.Column, 4)
	cols[availableNetworksCfg.connColIdx] = table.Column{
		Title: styles.SymbolConnection,
		Width: len(styles.SymbolConnection),
	}
	cols[availableNetworksCfg.ssidColIdx] = table.Column{Title: "SSID"}
	cols[availableNetworksCfg.securityColIdx] = table.Column{Title: "Security"}
	cols[availableNetworksCfg.signalColIdx] = table.Column{
		Title: styles.SymbolSignal,
		Width: max(availableNetworksCfg.minSignalColWidth, len(styles.SymbolSignal)),
	}

	initTableStyle := styles.DataTableStyle
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(initTableStyle),
	)

	s := spinner.New()
	s.Spinner = styles.Spinner

	model := &AvailableNetworksModel{
		dataTable: t,

		indicatorSpinner: s,
		indicatorState:   AvailableNetsDone,

		keys:    keys,
		netMngr: networksManager,
	}

	model.bakeSizes()

	return model
}

func (m *AvailableNetworksModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *AvailableNetworksModel) Resize(width, height int) {
	m.width = width
	m.height = height

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	secColWidth := int(float32(width) * availableNetworksCfg.securityWidthProportion)
	signalColWidth := m.dataTable.Columns()[availableNetworksCfg.signalColIdx].Width
	conColWidth := m.dataTable.Columns()[availableNetworksCfg.connColIdx].Width
	ssidWidth := width - signalColWidth - tableUtilityOffset - conColWidth - secColWidth

	m.dataTable.Columns()[availableNetworksCfg.securityColIdx].Width = secColWidth
	m.dataTable.Columns()[availableNetworksCfg.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *AvailableNetworksModel) Width() int {
	return m.width
}

func (m *AvailableNetworksModel) Height() int {
	return m.height
}

func (m *AvailableNetworksModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.SetStyles(styles.TableStyle)
	m.dataTable.Focus()
	return nil
}

func (m *AvailableNetworksModel) Blur() {
	m.focus = false
	m.dataTable.Blur()
	m.dataTable.SetStyles(styles.DataTableStyle)
}

func (m *AvailableNetworksModel) Focused() bool {
	return m.focus
}

func (m *AvailableNetworksModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *AvailableNetworksModel) Update(msg tea.Msg) (*AvailableNetworksModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.rescan):
			if m.indicatorState != AvailableNetsDone {
				return m, nil
			}
			return m, m.RescanCmd()
		case key.Matches(msg, m.keys.connect):
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, OpenConnectorCmd(row[availableNetworksCfg.ssidColIdx])
			}
			return m, nil
		}
	case AvailableNetworksStateMsg:
		return m, m.setStateCmd(availableNetworksState(msg))
	case RescanAvailableNetworksMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != AvailableNetsDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *AvailableNetworksModel) View() string {
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

func (m *AvailableNetworksModel) indicatorView() string {
	var view string
	if m.indicatorState != AvailableNetsDone {
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

func (m *AvailableNetworksModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(AvailableNetsScanning),
		func() tea.Msg {
			list, err := m.netMngr.ScanNetworks(context.Background())
			if err != nil {
				return tea.Batch(
					m.setStateCmd(AvailableNetsDone),
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
			return m.setStateCmd(AvailableNetsDone)
		},
	)
}

type RescanAvailableNetworksMsg struct {
	delay time.Duration
}

func RescanAvailableNetworksCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanAvailableNetworksMsg{delay: delay}
	}
}

type AvailableNetworksStateMsg availableNetworksState

func (m *AvailableNetworksModel) setStateCmd(state availableNetworksState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return nil
	}
	if state == AvailableNetsDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func SetAvailableNetworksStateCmd(state availableNetworksState) tea.Cmd {
	return func() tea.Msg {
		return AvailableNetworksStateMsg(state)
	}
}
