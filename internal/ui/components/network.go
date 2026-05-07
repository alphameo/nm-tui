package components

import (
	"context"
	"fmt"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type networkKeyMap struct{}

type NetworkModel struct {
	dataTable *table.Model

	deviceColIdx int
	typeColIdx   int
	connColIdx   int
	stateColIdx  int

	deviceWidthProportion float32
	minDeviceColWidth     int
	minConnColWidth       int

	keys *networkKeyMap

	nm infra.NetworkManager

	height int
	width  int
}

func NewNetworkModel(networkManager infra.NetworkManager, keys *networkKeyMap) *NetworkModel {
	cols := []table.Column{
		{Title: "Device"},
		{Title: "Type", Width: 11},
		{Title: "Connection"},
		{Title: "State", Width: 22},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithStyles(styles.TableStyle),
	)
	return &NetworkModel{
		dataTable:    &t,
		deviceColIdx: 0,
		typeColIdx:   1,
		connColIdx:   2,
		stateColIdx:  3,

		deviceWidthProportion: float32(0.4),
		minDeviceColWidth:     6,
		minConnColWidth:       10,

		nm:   networkManager,
		keys: keys,
	}
}

func (m *NetworkModel) Resize(width, height int) {
	m.height = height
	m.width = width

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	deviceColWidth := max(int(float32(width)*m.deviceWidthProportion), m.minDeviceColWidth)
	typeColWidth := m.dataTable.Columns()[m.typeColIdx].Width
	stateWidth := m.dataTable.Columns()[m.stateColIdx].Width

	connWidth := width - typeColWidth - deviceColWidth - tableUtilityOffset - stateWidth
	m.dataTable.Columns()[m.deviceColIdx].Width = deviceColWidth
	m.dataTable.Columns()[m.connColIdx].Width = connWidth
	m.dataTable.UpdateViewport()
}

func (m *NetworkModel) Width() int {
	return m.width
}

func (m *NetworkModel) Height() int {
	return m.height
}

func (m NetworkModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *NetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	// upd, cmd := m.dataTable.Update(msg)
	// m.dataTable = &upd
	// return m, cmd
	return m, nil
}

func (m *NetworkModel) handleKey(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	upd, cmd := m.dataTable.Update(keyMsg)
	m.dataTable = &upd
	return m, cmd
}

func (m *NetworkModel) View() string {
	tableStyle := styles.BorderedStyle
	table := tableStyle.Render(m.dataTable.View())
	radioStatus, _ := m.nm.GetRadioStatus(context.Background())
	wwan := fmt.Sprintf("WWAN:  %s", renderer.RenderEnabledStatus(radioStatus.EnabledWWAN))
	wifi := fmt.Sprintf("Wi-Fi: %s", renderer.RenderEnabledStatus(radioStatus.EnabledWifi))

	conStatus, _ := m.nm.GetConnectivityStatus(context.Background())
	connectivity := fmt.Sprintf("Connectivity status: %s", conStatus)

	return lipgloss.JoinVertical(lipgloss.Center, table, wwan, wifi, connectivity)
}

func (m *NetworkModel) RescanCmd() tea.Cmd {
	return func() tea.Msg {
		list, err := m.nm.GetNetworkDevices(context.Background())
		if err != nil {
			return NotifyCmd("Cannot get network devices")
		}

		rows := []table.Row{}
		for _, device := range list {
			rows = append(rows, table.Row{
				device.Device,
				device.Type,
				device.Connection,
				device.State,
			})
		}
		m.dataTable.SetRows(rows)
		m.dataTable.GotoTop()
		m.dataTable.UpdateViewport()
		return NilCmd
	}
}
