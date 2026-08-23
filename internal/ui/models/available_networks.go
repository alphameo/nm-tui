package models

import (
	"context"
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type availableNetworksConfig struct {
	stateColIdx    int
	ssidColIdx     int
	securityColIdx int
	signalColIdx   int
	bandColIdx     int
	rateColIdx     int
	modeColIdx     int
	deviceColIdx   int

	ssidColTitle     string
	securityColTitle string
	stateColTitle    string
	modeColTitle     string
	rateColTitle     string
	bandColTitle     string
	deviceColTitle   string

	securityWidthProportion float32
	minSignalColWidth       int
	minRateColWidth         int
	minBandColWidth         int
}

var availableNetworksCfg = availableNetworksConfig{
	stateColIdx:    0,
	modeColIdx:     1,
	ssidColIdx:     2,
	securityColIdx: 3,
	bandColIdx:     4,
	rateColIdx:     5,
	deviceColIdx:   6,
	signalColIdx:   7,

	ssidColTitle:     "SSID",
	securityColTitle: "Security",
	stateColTitle:    "State",
	modeColTitle:     "Mode",
	rateColTitle:     "Mb/s",
	bandColTitle:     "GHz",
	deviceColTitle:   "Device",

	securityWidthProportion: 0.3,
	minSignalColWidth:       3,
	minRateColWidth:         4,
	minBandColWidth:         3,
}

type availableNetworksKeyMap struct {
	connect    key.Binding
	activate   key.Binding
	deactivate key.Binding
}

type AvailableNetworksModel struct {
	dataTable          table.Model
	focusedTableStyles table.Styles
	bluredTableStyles  table.Styles

	focus bool

	keys availableNetworksKeyMap

	netMngr infra.NetworksManager

	focusedStyle lipgloss.Style
	bluredStyle  lipgloss.Style
}

func NewAvailableNetworksModel(
	keys availableNetworksKeyMap,
	networksManager infra.NetworksManager,
) *AvailableNetworksModel {
	cols := make([]table.Column, 8)
	cols[availableNetworksCfg.stateColIdx] = table.Column{
		Title: availableNetworksCfg.stateColTitle,
		Width: len(availableNetworksCfg.stateColTitle),
	}
	cols[availableNetworksCfg.ssidColIdx] = table.Column{
		Title: availableNetworksCfg.ssidColTitle,
		Width: len(availableNetworksCfg.ssidColTitle),
	}
	cols[availableNetworksCfg.securityColIdx] = table.Column{
		Title: availableNetworksCfg.securityColTitle,
		Width: len(availableNetworksCfg.securityColTitle),
	}
	cols[availableNetworksCfg.signalColIdx] = table.Column{
		Title: styles.SymbolSignal,
		Width: max(availableNetworksCfg.minSignalColWidth, len(styles.SymbolSignal)),
	}
	cols[availableNetworksCfg.modeColIdx] = table.Column{
		Title: availableNetworksCfg.modeColTitle,
		Width: len(availableNetworksCfg.modeColTitle),
	}
	cols[availableNetworksCfg.bandColIdx] = table.Column{
		Title: availableNetworksCfg.bandColTitle,
		Width: max(len(availableNetworksCfg.bandColTitle), availableNetworksCfg.minBandColWidth),
	}
	cols[availableNetworksCfg.rateColIdx] = table.Column{
		Title: availableNetworksCfg.rateColTitle,
		Width: max(len(availableNetworksCfg.rateColTitle), availableNetworksCfg.minRateColWidth),
	}
	cols[availableNetworksCfg.deviceColIdx] = table.Column{
		Title: availableNetworksCfg.deviceColTitle,
		Width: len(availableNetworksCfg.deviceColTitle),
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)

	model := &AvailableNetworksModel{
		dataTable:          t,
		focusedTableStyles: table.DefaultStyles(),
		bluredTableStyles:  table.DefaultStyles(),

		keys:         keys,
		netMngr:      networksManager,
		focusedStyle: lipgloss.NewStyle(),
		bluredStyle:  lipgloss.NewStyle(),
	}

	return model
}

func (m *AvailableNetworksModel) Resize(width, height int) {
	m.focusedStyle = m.focusedStyle.Width(width).Height(height)
	m.bluredStyle = m.bluredStyle.Width(width).Height(height)

	border := m.focusedStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	secColWidth := int(float32(width) * availableNetworksCfg.securityWidthProportion)
	signalColWidth := m.dataTable.Columns()[availableNetworksCfg.signalColIdx].Width
	stateColWidth := m.dataTable.Columns()[availableNetworksCfg.stateColIdx].Width
	bandColWidth := m.dataTable.Columns()[availableNetworksCfg.bandColIdx].Width
	rateColWidth := m.dataTable.Columns()[availableNetworksCfg.rateColIdx].Width
	modeColWidth := m.dataTable.Columns()[availableNetworksCfg.modeColIdx].Width
	deviceColWidth := m.dataTable.Columns()[availableNetworksCfg.deviceColIdx].Width
	ssidWidth := width - signalColWidth - tableUtilityOffset - stateColWidth - secColWidth - modeColWidth - rateColWidth - bandColWidth - deviceColWidth

	m.dataTable.Columns()[availableNetworksCfg.securityColIdx].Width = secColWidth
	m.dataTable.Columns()[availableNetworksCfg.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *AvailableNetworksModel) Width() int { return m.focusedStyle.GetWidth() }

func (m *AvailableNetworksModel) Height() int { return m.focusedStyle.GetHeight() }

func (m *AvailableNetworksModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.SetStyles(m.focusedTableStyles)
	m.dataTable.Focus()
	return nil
}

func (m *AvailableNetworksModel) Blur() {
	m.focus = false
	m.dataTable.Blur()
	m.dataTable.SetStyles(m.bluredTableStyles)
}

func (m *AvailableNetworksModel) Focused() bool {
	return m.focus
}

func (m *AvailableNetworksModel) activeStyle() *lipgloss.Style {
	if m.focus {
		return &m.focusedStyle
	}
	return &m.bluredStyle
}

func (m *NetworkProfilesModel) UpdateTable() {
	if m.focus {
		m.dataTable.SetStyles(m.focusedTableStyles)
	} else {
		m.dataTable.SetStyles(m.bluredTableStyles)
	}
}

func (m *NetworkProfilesModel) SetTableStyles(focused, blured table.Styles) {
	m.bluredTableStyles = blured
	m.focusedTableStyles = focused
	m.UpdateTable()
}

func (m *AvailableNetworksModel) Update(msg tea.Msg) (*AvailableNetworksModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.connect):
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, OpenConnectorCmd(row[availableNetworksCfg.ssidColIdx])
			}
			return m, nil
		case key.Matches(msg, m.keys.activate):
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, m.activateConnToSelectedCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.deactivate):
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, m.deactivateConnToSelectedCmd()
			}
			return m, nil
		}
	case AvailableNetworksStateMsg:
		return m, SetNetworksStateCmd(networksState(msg))
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *AvailableNetworksModel) View() string {
	view := m.dataTable.View()

	style := m.activeStyle()
	view = renderer.RenderWithTitleAndKeybind(
		view,
		"Available networks",
		"1",
		*style,
		styles.AccentColor,
	)
	return view
}

func (m *AvailableNetworksModel) setAvailable(list []AvailableNetwork, err error) tea.Cmd {
	rows := []table.Row{}
	for _, net := range list {
		var state string
		if net.Active {
			state = styles.SymbolCheck
		} else if net.ProfileExists {
			state = styles.SymbolSaved
		}
		row := make([]string, 8)
		row[availableNetworksCfg.stateColIdx] = state
		row[availableNetworksCfg.ssidColIdx] = net.SSID
		row[availableNetworksCfg.securityColIdx] = net.SecurityMode
		row[availableNetworksCfg.signalColIdx] = strconv.Itoa(net.Signal)
		row[availableNetworksCfg.rateColIdx] = formatRate(net.Rate)
		row[availableNetworksCfg.bandColIdx] = formatBand(net.Band)
		row[availableNetworksCfg.deviceColIdx] = net.LookingDevice
		row[availableNetworksCfg.modeColIdx] = ConvertNetworkMode(net.NetworkMode)

		rows = append(rows, row)
	}

	m.dataTable.SetRows(rows)
	m.dataTable.GotoTop()
	m.dataTable.UpdateViewport()

	cmds := []tea.Cmd{SetNetworksStateCmd(NetsDone)}
	if err != nil {
		cmds = append(cmds, NotifyCmd("Cannot scan available wifi networks"))
	}
	return tea.Batch(cmds...)
}

func formatBand(band float64) string {
	if band <= 0 {
		return ""
	}
	return fmt.Sprintf("%g", band)
}

func formatRate(rate float64) string {
	if rate <= 0 {
		return ""
	}
	return fmt.Sprintf("%g", rate)
}

type AvailableNetworksStateMsg networksState

func SetAvailableNetworksStateCmd(state networksState) tea.Cmd {
	return func() tea.Msg {
		return AvailableNetworksStateMsg(state)
	}
}

func (m *AvailableNetworksModel) activateConnToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		SetNetworksStateCmd(NetsActivating),
		func() tea.Msg {
			ssid := m.dataTable.SelectedRow()[availableNetworksCfg.ssidColIdx]
			err := m.netMngr.TryActivateNetwork(context.Background(), ssid)
			if err != nil {
				return tea.Batch(
					SetNetworksStateCmd(NetsDone),
					NotifyCmd(fmt.Sprintf("Cannot activate connection to network with SSID=%q\n"+
						"Try connect via profile", ssid)),
				)
			}
			return tea.Batch(
				SetNetworksStateCmd(NetsDone),
				RescanNetworksCmd(),
			)
		},
	)
}

func (m *AvailableNetworksModel) deactivateConnToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		SetNetworksStateCmd(NetsDeactivating),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[availableNetworksCfg.ssidColIdx]
			err := m.netMngr.DeactivateProfile(context.Background(), name)
			if err != nil {
				return tea.Batch(
					SetNetworksStateCmd(NetsDone),
					NotifyCmd(
						fmt.Sprintf("Error while deactivating connection to network with SSID=%q\n"+
							"try disconnect via profile (The profile name and SSID may differ)", name),
					),
				)
			}
			return tea.Batch(
				SetNetworksStateCmd(NetsDone),
				RescanNetworksCmd(),
			)
		})
}
