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
	stateColIdx             int
	ssidColIdx              int
	securityColIdx          int
	signalColIdx            int
	ssidColTitle            string
	securityColTitle        string
	stateColTitle           string
	securityWidthProportion float32
	minSignalColWidth       int
}

var availableNetworksCfg = availableNetworksConfig{
	stateColIdx:    0,
	ssidColIdx:     1,
	securityColIdx: 2,
	signalColIdx:   3,

	ssidColTitle:     "SSID",
	securityColTitle: "Security",
	stateColTitle:    "State",

	securityWidthProportion: 0.3,
	minSignalColWidth:       3,
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
	cols := make([]table.Column, 4)
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
	conColWidth := m.dataTable.Columns()[availableNetworksCfg.stateColIdx].Width
	ssidWidth := width - signalColWidth - tableUtilityOffset - conColWidth - secColWidth

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
	for _, wifiNet := range list {
		var connectionFlag string
		if wifiNet.Active {
			connectionFlag = styles.SymbolCheck
		} else if wifiNet.ProfileExists {
			connectionFlag = styles.SymbolSaved
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

	cmds := []tea.Cmd{SetNetworksStateCmd(NetsDone)}
	if err != nil {
		cmds = append(cmds, NotifyCmd("Cannot scan available wifi networks"))
	}
	return tea.Batch(cmds...)
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
