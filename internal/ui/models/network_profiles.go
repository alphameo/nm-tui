package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type networkProfilesKeyMap struct {
	edit       key.Binding
	activate   key.Binding
	deactivate key.Binding
	delete     key.Binding
}

type NetworkProfilesModel struct {
	dataTable          table.Model
	focusedTableStyles table.Styles
	bluredTableStyles  table.Styles

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
		Title: "State",
		Width: len("State"),
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

	model := &NetworkProfilesModel{
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

func (m *NetworkProfilesModel) Resize(width, height int) {
	m.focusedStyle = m.focusedStyle.Width(width).Height(height)
	m.bluredStyle = m.bluredStyle.Width(width).Height(height)

	border := m.focusedStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tablePaddingOffset := len(m.dataTable.Columns()) * 2

	connWidth := m.dataTable.Columns()[networkProfilesCfg.connColIdx].Width
	modeWidth := m.dataTable.Columns()[networkProfilesCfg.modeColIdx].Width

	computedWidth := width - tablePaddingOffset - connWidth - modeWidth
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

func (m *NetworkProfilesModel) Update(msg tea.Msg) (*NetworkProfilesModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
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
			return m, m.activateConnToSelectedCmd()

		case key.Matches(msg, m.keys.deactivate):
			return m, m.deactivateConnToSelectedCmd()
		case key.Matches(msg, m.keys.delete):
			return m, m.deleteSelectedCmd()
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NetworkProfilesModel) View() string {
	view := m.dataTable.View()

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

func (m *NetworkProfilesModel) setProfiles(list []NetworkProfileShort, err error) tea.Cmd {
	rows := []table.Row{}
	for _, wifiSaved := range list {
		var connectionFlag string
		if wifiSaved.Active {
			connectionFlag = styles.SymbolCheck
		} else if wifiSaved.Available {
			connectionFlag = styles.SymbolAvailable
		}
		rows = append(rows, table.Row{
			connectionFlag,
			wifiSaved.Mode,
			wifiSaved.SSID,
			wifiSaved.Name,
		})
	}

	m.dataTable.SetRows(rows)

	cmds := []tea.Cmd{SetNetworksStateCmd(NetsDone)}
	if err != nil {
		cmds = append(cmds, NotifyCmd("Cannot get network profiles"))
	}
	return tea.Batch(cmds...)
}

func (m *NetworkProfilesModel) activateConnToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		SetNetworksStateCmd(NetsActivating),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.netMngr.ActivateProfile(context.Background(), name)
			if err != nil {
				return tea.Batch(
					SetNetworksStateCmd(NetsDone),
					NotifyCmd(fmt.Sprintf("Cannot connect to %q", name)),
				)
			}
			return tea.Batch(
				SetNetworksStateCmd(NetsDone),
				m.gotoTop(),
				RescanNetworksCmd(),
			)
		},
	)
}

func (m *NetworkProfilesModel) deactivateConnToSelectedCmd() tea.Cmd {
	return tea.Sequence(SetNetworksStateCmd(NetsDeactivating),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.netMngr.DeactivateProfile(context.Background(), name)
			if err != nil {
				return tea.Batch(
					SetNetworksStateCmd(NetsDone),
					NotifyCmd(
						fmt.Sprintf("Error while deactivating connection with %q", name),
					),
				)
			}
			return tea.Batch(
				SetNetworksStateCmd(NetsDone),
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
			return NotifyCmd(fmt.Sprintf("Error while deleting profile %q", name))
		}
		cursor := m.dataTable.Cursor()
		if cursor == len(m.dataTable.Rows())-1 {
			m.dataTable.SetCursor(cursor - 1)
		}
		return RescanNetworksCmd()
	}
}

func (m *NetworkProfilesModel) gotoTop() tea.Cmd {
	return func() tea.Msg {
		m.dataTable.GotoTop()
		return NilCmd()
	}
}
