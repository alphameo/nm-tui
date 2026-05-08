package components

import (
	"context"
	"fmt"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type networkKeyMap struct{}

type NetworkModel struct {
	devicesTable *table.Model

	deviceColIdx int
	typeColIdx   int
	connColIdx   int
	stateColIdx  int

	deviceWidthProportion float32
	minDeviceColWidth     int
	minConnColWidth       int

	wwan         *toggle.Model
	wifi         *toggle.Model
	connectivity string

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
		devicesTable: &t,
		deviceColIdx: 0,
		typeColIdx:   1,
		connColIdx:   2,
		stateColIdx:  3,

		deviceWidthProportion: float32(0.4),
		minDeviceColWidth:     6,
		minConnColWidth:       10,

		wwan: toggle.New(false),
		wifi: toggle.New(false),

		nm:   networkManager,
		keys: keys,
	}
}

func (m *NetworkModel) Resize(width, height int) {
	m.height = height
	m.width = width

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	m.devicesTable.SetWidth(width)
	m.devicesTable.SetHeight(height - 4)

	tableUtilityOffset := len(m.devicesTable.Columns()) * 2

	deviceColWidth := max(int(float32(width)*m.deviceWidthProportion), m.minDeviceColWidth)
	typeColWidth := m.devicesTable.Columns()[m.typeColIdx].Width
	stateWidth := m.devicesTable.Columns()[m.stateColIdx].Width

	connWidth := width - typeColWidth - deviceColWidth - tableUtilityOffset - stateWidth
	m.devicesTable.Columns()[m.deviceColIdx].Width = deviceColWidth
	m.devicesTable.Columns()[m.connColIdx].Width = connWidth
	m.devicesTable.UpdateViewport()
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
	upd, cmd := m.devicesTable.Update(msg)
	m.devicesTable = &upd
	return m, cmd
}

func (m *NetworkModel) handleKey(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	upd, cmd := m.devicesTable.Update(keyMsg)
	m.devicesTable = &upd
	return m, cmd
}

func (m *NetworkModel) View() string {
	tableStyle := styles.BorderedStyle
	table := tableStyle.Render(m.devicesTable.View())
	wwan := "WWAN:  "
	wwan = lipgloss.JoinHorizontal(lipgloss.Center, wwan, m.wwan.View())

	wifi := "Wi-Fi: "
	wifi = lipgloss.JoinHorizontal(lipgloss.Center, wifi, m.wifi.View())

	connectivity := fmt.Sprintf("Connectivity status: %s", m.connectivity)

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
		m.devicesTable.SetRows(rows)
		m.devicesTable.GotoTop()
		m.devicesTable.UpdateViewport()

		radioStatus, err := m.nm.GetRadioStatus(context.Background())
		if err != nil {
			return NotifyCmd("Cannot get radio status")
		}
		m.wwan.SetValue(radioStatus.EnabledWWAN)
		m.wifi.SetValue(radioStatus.EnabledWifi)

		conStatus, err := m.nm.GetConnectivityStatus(context.Background())
		if err != nil {
			return NotifyCmd("Cannot get connection status")
		}
		m.connectivity = string(conStatus)

		return NilCmd
	}
}
