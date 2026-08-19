package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type networkProfilesState int

const (
	NetProfilesNil networkProfilesState = iota
	NetProfilesScanning
	NetProfilesConnecting
	NetProfilesDisconnecting
	NetProfilesDone
)

func (s *networkProfilesState) String() string {
	switch *s {
	case NetProfilesScanning:
		return "Scanning"
	case NetProfilesConnecting:
		return "Connecting"
	case NetProfilesDisconnecting:
		return "Disconnecting"
	case NetProfilesDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type networkProfilesKeyMap struct {
	edit       key.Binding
	activate   key.Binding
	deactivate key.Binding
	rescan     key.Binding
	delete     key.Binding
}

type NetworkProfilesModel struct {
	dataTable          table.Model
	focusedTableStyles table.Styles
	bluredTableStyles  table.Styles

	indicatorSpinner spinner.Model
	indicatorState   networkProfilesState

	focus bool

	keys networkProfilesKeyMap

	netMngr infra.NetworksManager

	focusedStyle lipgloss.Style
	bluredStyle  lipgloss.Style
}

type networkProfilesConfig struct {
	connColIdx int
	modeColIdx int
	ssidColIdx int
	nameColIdx int

	modeColTitle string
	ssidColTitle string
	nameColTitle string

	ssidWidthProportion float32
}

var networkProfilesCfg = networkProfilesConfig{
	connColIdx: 0,
	modeColIdx: 1,
	ssidColIdx: 2,
	nameColIdx: 3,

	modeColTitle: "Mode",
	ssidColTitle: "SSID",
	nameColTitle: "Name",

	ssidWidthProportion: 0.5,
}

func NewNetworkProfilesModel(keys networkProfilesKeyMap, networksManager infra.NetworksManager) *NetworkProfilesModel {
	cols := make([]table.Column, 4)
	cols[networkProfilesCfg.connColIdx] = table.Column{
		Title: styles.SymbolConnection,
		Width: len(styles.SymbolConnection),
	}
	cols[networkProfilesCfg.modeColIdx] = table.Column{
		Title: networkProfilesCfg.modeColTitle,
		Width: len(networkProfilesCfg.modeColTitle),
	}
	cols[networkProfilesCfg.ssidColIdx] = table.Column{
		Title: networkProfilesCfg.ssidColTitle,
		Width: len(networkProfilesCfg.ssidColTitle),
	}
	cols[networkProfilesCfg.nameColIdx] = table.Column{
		Title: networkProfilesCfg.nameColTitle,
		Width: len(networkProfilesCfg.nameColTitle),
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)

	s := spinner.New()
	s.Spinner = styles.Spinner

	model := &NetworkProfilesModel{
		dataTable:          t,
		focusedTableStyles: table.DefaultStyles(),
		bluredTableStyles:  table.DefaultStyles(),

		indicatorSpinner: s,
		indicatorState:   NetProfilesDone,

		keys:         keys,
		netMngr:      networksManager,
		focusedStyle: lipgloss.NewStyle(),
		bluredStyle:  lipgloss.NewStyle(),
	}

	return model
}

func (m *NetworkProfilesModel) Resize(width, height int) {
	m.focusedStyle = m.focusedStyle.Width(width).Height(height)
	m.bluredStyle = m.bluredStyle.Width(width).Height(height)

	border := m.focusedStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	indicatorStateHeight := lipgloss.Height(m.indicatorView())
	height -= indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2
	connWidth := m.dataTable.Columns()[networkProfilesCfg.connColIdx].Width
	modeWidth := m.dataTable.Columns()[networkProfilesCfg.modeColIdx].Width

	computedWidth := width - tableUtilityOffset - connWidth - modeWidth
	possibleNameWidth := int(float32(computedWidth) * networkProfilesCfg.ssidWidthProportion)
	ssidWidth := computedWidth - possibleNameWidth
	nameWidth := computedWidth - ssidWidth

	m.dataTable.Columns()[networkProfilesCfg.nameColIdx].Width = nameWidth
	m.dataTable.Columns()[networkProfilesCfg.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *NetworkProfilesModel) Width() int { return m.focusedStyle.GetWidth() }

func (m *NetworkProfilesModel) Height() int { return m.focusedStyle.GetHeight() }

func (m *NetworkProfilesModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.Focus()
	m.dataTable.SetStyles(m.focusedTableStyles)
	return nil
}

func (m *NetworkProfilesModel) Blur() {
	m.focus = false
	m.dataTable.Blur()
	m.dataTable.SetStyles(m.bluredTableStyles)
}

func (m *NetworkProfilesModel) Focused() bool {
	return m.focus
}

func (m *NetworkProfilesModel) activeStyle() *lipgloss.Style {
	if m.focus {
		return &m.focusedStyle
	}
	return &m.bluredStyle
}

func (m *AvailableNetworksModel) UpdateTable() {
	if m.focus {
		m.dataTable.SetStyles(m.focusedTableStyles)
	} else {
		m.dataTable.SetStyles(m.bluredTableStyles)
	}
}

func (m *AvailableNetworksModel) SetTableStyles(focused, blured table.Styles) {
	m.bluredTableStyles = blured
	m.focusedTableStyles = focused
	m.UpdateTable()
}

func (m *NetworkProfilesModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *NetworkProfilesModel) Update(msg tea.Msg) (*NetworkProfilesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.edit):
			row := m.dataTable.SelectedRow()
			if row == nil {
				return m, nil
			}
			name := row[networkProfilesCfg.nameColIdx]

			return m, OpenProfileEditorCmd(name)

		case key.Matches(msg, m.keys.activate):
			return m, m.connectToSelectedCmd()

		case key.Matches(msg, m.keys.deactivate):
			return m, m.disconnectFromSelectedCmd()
		case key.Matches(msg, m.keys.rescan):
			return m, RescanNetworkProfilesCmd()
		case key.Matches(msg, m.keys.delete):
			return m, m.deleteSelectedCmd()
		}
	case RescanNetworkProfilesMsg:
		return m, m.RescanCmd()
	case WifiSavedStateMsg:
		return m, m.setStateCmd(networkProfilesState(msg))
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != NetProfilesDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NetworkProfilesModel) View() string {
	view := m.dataTable.View()
	statusline := m.indicatorView()
	view = lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)

	style := m.activeStyle()
	view = renderer.RenderWithTitleAndKeybind(
		view,
		"Network profiles",
		"2",
		*style,
		styles.AccentColor,
	)
	return view
}

func (m *NetworkProfilesModel) indicatorView() string {
	var view string
	if m.indicatorState != NetProfilesDone {
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

type RescanNetworkProfilesMsg struct{}

func RescanNetworkProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		return RescanNetworkProfilesMsg{}
	}
}

func (m *NetworkProfilesModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetProfilesScanning),
		func() tea.Msg {
			list, err := m.netMngr.ListProfiles(context.Background())
			if err != nil {
				return tea.Batch(
					NotifyCmd("Cannot get network profiles"),
					m.setStateCmd(NetProfilesDone),
				)
			}
			rows := []table.Row{}
			for _, wifiSaved := range list {
				var connectionFlag string
				if wifiSaved.Active {
					connectionFlag = styles.SymbolCheck
				}
				rows = append(rows, table.Row{
					connectionFlag,
					ViewNetworkMode(wifiSaved.Mode),
					wifiSaved.SSID,
					wifiSaved.Name,
				})
			}

			m.dataTable.SetRows(rows)

			return m.setStateCmd(NetProfilesDone)
		},
	)
}

func ViewNetworkMode(mode infra.NetworkMode) string {
	switch mode {
	case infra.NetworkAccessPoint:
		return styles.SymbolAccessPoint
	case infra.NetworkInfra:
		return styles.SymbolInfra
	case infra.NetworkMesh:
		return styles.SymbolMesh
	case infra.NetworkAdHoc:
		return styles.SymbolAdHoc
	default:
		return "?"
	}
}

type WifiSavedStateMsg networkProfilesState

func (m *NetworkProfilesModel) setStateCmd(state networkProfilesState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == NetProfilesDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func (m *NetworkProfilesModel) connectToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetProfilesConnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.netMngr.ActivateProfile(context.Background(), name)
			if err != nil {
				return tea.Batch(
					m.setStateCmd(NetProfilesDone),
					NotifyCmd(fmt.Sprintf("Cannot connect to %s", name)),
				)
			}
			return tea.Batch(
				m.setStateCmd(NetProfilesDone),
				m.gotoTop(),
				RescanNetworksCmd(),
			)
		},
	)
}

func (m *NetworkProfilesModel) gotoTop() tea.Cmd {
	return func() tea.Msg {
		m.dataTable.GotoTop()
		return NilCmd()
	}
}

func (m *NetworkProfilesModel) disconnectFromSelectedCmd() tea.Cmd {
	return tea.Sequence(m.setStateCmd(NetProfilesDisconnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.netMngr.DeactivateProfile(context.Background(), name)
			if err != nil {
				return NotifyCmd(
					fmt.Sprintf("Error while disconnecting from %s", name),
				)
			}
			return tea.Batch(
				m.gotoTop(),
				RescanNetworksCmd(),
			)
		})
}

func (m *NetworkProfilesModel) deleteSelectedCmd() tea.Cmd {
	row := m.dataTable.SelectedRow()
	return func() tea.Msg {
		name := row[networkProfilesCfg.nameColIdx]
		err := m.netMngr.DeleteProfile(context.Background(), name)
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Error while deleting %s", name))
		}
		cursor := m.dataTable.Cursor()
		if cursor == len(m.dataTable.Rows())-1 {
			m.dataTable.SetCursor(cursor - 1)
		}
		return RescanNetworksCmd()
	}
}
