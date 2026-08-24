package models

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

type networkDevicesConfig struct {
	deviceColIdx int
	typeColIdx   int
	connColIdx   int
	stateColIdx  int

	deviceColTitle string
	typeColTitle   string
	connColTitle   string
	stateColTitle  string

	deviceWidthProportion float32
	typeWidthProportion   float32
	stateWidthProportion  float32
}

var networkDevicesCfg = networkDevicesConfig{
	deviceColIdx: 0,
	typeColIdx:   1,
	connColIdx:   2,
	stateColIdx:  3,

	deviceColTitle: "Device",
	typeColTitle:   "Type",
	connColTitle:   "Connection",
	stateColTitle:  "State",

	deviceWidthProportion: 0.2,
	typeWidthProportion:   0.15,
	stateWidthProportion:  0.3,
}

type networkDevicesKeyMap struct {
	showInfo key.Binding
}

type NetworkDevicesModel struct {
	table              table.Model
	focusedTableStyles table.Styles
	bluredTableStyles  table.Styles

	focus bool

	keys networkDevicesKeyMap

	connMngr infra.DeviceManager

	focusedStyle lipgloss.Style
	bluredStyle  lipgloss.Style
}

func NewNetworkDevicesModel(keys networkDevicesKeyMap, deviceManager infra.DeviceManager) *NetworkDevicesModel {
	cols := make([]table.Column, 4)
	cols[networkDevicesCfg.deviceColIdx] = table.Column{
		Title: networkDevicesCfg.deviceColTitle,
		Width: len(networkDevicesCfg.deviceColTitle),
	}
	cols[networkDevicesCfg.typeColIdx] = table.Column{
		Title: networkDevicesCfg.typeColTitle,
		Width: len(networkDevicesCfg.typeColTitle),
	}
	cols[networkDevicesCfg.connColIdx] = table.Column{
		Title: networkDevicesCfg.connColTitle,
		Width: len(networkDevicesCfg.connColTitle),
	}
	cols[networkDevicesCfg.stateColIdx] = table.Column{
		Title: networkDevicesCfg.stateColTitle,
		Width: len(networkDevicesCfg.stateColTitle),
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithStyles(styles.DataTableStyles),
	)
	model := &NetworkDevicesModel{
		table:              t,
		focusedTableStyles: table.DefaultStyles(),
		bluredTableStyles:  table.DefaultStyles(),

		connMngr:     deviceManager,
		keys:         keys,
		focusedStyle: lipgloss.NewStyle(),
		bluredStyle:  lipgloss.NewStyle(),
	}
	return model
}

func (m *NetworkDevicesModel) Resize(width, height int) {
	m.focusedStyle = m.focusedStyle.Width(width).Height(height)
	m.bluredStyle = m.bluredStyle.Width(width).Height(height)

	border := m.focusedStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	m.table.SetWidth(width)
	m.table.SetHeight(height)

	// NOTE: from padding
	tablePaddingOffset := len(m.table.Columns()) * 2

	deviceColWidth := int(float32(width) * networkDevicesCfg.deviceWidthProportion)
	typeColWidth := int(float32(width) * networkDevicesCfg.typeWidthProportion)
	stateWidth := int(float32(width) * networkDevicesCfg.stateWidthProportion)
	connWidth := width - typeColWidth - deviceColWidth - tablePaddingOffset - stateWidth

	m.table.Columns()[networkDevicesCfg.deviceColIdx].Width = deviceColWidth
	m.table.Columns()[networkDevicesCfg.typeColIdx].Width = typeColWidth
	m.table.Columns()[networkDevicesCfg.stateColIdx].Width = stateWidth
	m.table.Columns()[networkDevicesCfg.connColIdx].Width = connWidth
	m.table.UpdateViewport()
}
func (m *NetworkDevicesModel) Width() int { return m.focusedStyle.GetWidth() }

func (m *NetworkDevicesModel) Height() int { return m.focusedStyle.GetHeight() }

func (m *NetworkDevicesModel) Focused() bool { return m.focus }

func (m *NetworkDevicesModel) Focus() tea.Cmd {
	m.focus = true
	m.table.SetStyles(m.focusedTableStyles)
	m.table.Focus()
	return nil
}

func (m *NetworkDevicesModel) Blur() {
	m.focus = false
	m.table.Blur()
	m.table.SetStyles(m.bluredTableStyles)
}

func (m *NetworkDevicesModel) activeStyle() *lipgloss.Style {
	if m.focus {
		return &m.focusedStyle
	}
	return &m.bluredStyle
}

func (m *NetworkDevicesModel) UpdateTable() {
	if m.focus {
		m.table.SetStyles(m.focusedTableStyles)
	} else {
		m.table.SetStyles(m.bluredTableStyles)
	}
}

func (m *NetworkDevicesModel) SetTableStyles(focused, blured table.Styles) {
	m.bluredTableStyles = blured
	m.focusedTableStyles = focused
	m.UpdateTable()
}

func (m *NetworkDevicesModel) Init() tea.Cmd {
	return nil
}

func (m *NetworkDevicesModel) Update(msg tea.Msg) (*NetworkDevicesModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		if !m.focus {
			return m, nil
		}
		if key.Matches(msg, m.keys.showInfo) {
			return m, OpenDeviceInfoCmd(m.table.SelectedRow()[networkDevicesCfg.deviceColIdx])
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NetworkDevicesModel) View() string {
	view := m.table.View()
	style := m.activeStyle()
	view = style.Render(view)
	title := styles.AccentStyle.Render("Network Devices")
	return compositor.Compose(
		title,
		view,
		compositor.Begin,
		compositor.Begin,
		2,
		0,
	)
}
