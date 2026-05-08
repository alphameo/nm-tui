package components

import (
	"context"
	"fmt"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type networkState int

const (
	NetworkScanning networkState = iota
	NetworkTogglingWifi
	NetworkTogglingWWAN
	NetworkTogglingNetworking
	NetworkDone
)

func (s *networkState) String() string {
	switch *s {
	case NetworkScanning:
		return "Scanning"
	case NetworkTogglingWWAN:
		return "Toggling WWAN"
	case NetworkTogglingWifi:
		return "Toggling Wi-Fi"
	case NetworkTogglingNetworking:
		return "Toggling Wi-Fi"
	case NetworkDone:
		return "󰄬"
	default:
		return "Undefined!!!"
	}
}

type networkKeyMap struct {
	up     key.Binding
	down   key.Binding
	rescan key.Binding
	toggle key.Binding
}

func (k networkKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down}
}

func (k networkKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.up, k.down}}
}

var networkKeys = &networkKeyMap{
	up: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^k", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("^j", "down"),
	),
	toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("󱁐", "toggle"),
	),
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan state"),
	),
}

type NetworkModel struct {
	devicesTable *table.Model

	deviceColIdx int
	typeColIdx   int
	connColIdx   int
	stateColIdx  int

	deviceWidthProportion float32
	minDeviceColWidth     int
	minConnColWidth       int

	wwan      *toggle.Model
	wwanStyle *lipgloss.Style

	wifi      *toggle.Model
	wifiStyle *lipgloss.Style

	networking      *toggle.Model
	networkingStyle *lipgloss.Style

	connectivity string

	indicatorSpinner spinner.Model
	indicatorState   networkState

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

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

	wwanStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)
	wifiStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)
	networkingStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)

	s := spinner.New()

	model := &NetworkModel{
		devicesTable: &t,
		deviceColIdx: 0,
		typeColIdx:   1,
		connColIdx:   2,
		stateColIdx:  3,

		deviceWidthProportion: float32(0.4),
		minDeviceColWidth:     6,
		minConnColWidth:       10,

		indicatorSpinner: s,
		indicatorState:   NetworkDone,

		wwan:      toggle.New(false),
		wwanStyle: &wwanStyle,

		wifi:      toggle.New(false),
		wifiStyle: &wifiStyle,

		networking:      toggle.New(false),
		networkingStyle: &networkingStyle,

		nm:   networkManager,
		keys: keys,
	}

	focuses := []Focusable{
		model.wwan,
		model.wifi,
		model.networking,
	}
	model.focuses = focuses

	return model
}

func (m *NetworkModel) Resize(width, height int) {
	m.height = height
	m.width = width

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	m.devicesTable.SetWidth(width)
	m.devicesTable.SetHeight(height - 5)

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

func (m *NetworkModel) Init() tea.Cmd {
	return tea.Batch(
		m.RescanCmd(),
		m.focuses[m.focusIdx].Focus(),
	)
}

func (m *NetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	var cmd tea.Cmd
	if m.indicatorState != NetworkDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}
	upd, cmd := m.devicesTable.Update(msg)
	m.devicesTable = &upd
	return m, cmd
}

func (m *NetworkModel) handleKey(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.down):
		return m, m.focusNextCmd()
	case key.Matches(keyMsg, m.keys.up):
		return m, m.focusPrevCmd()
	case key.Matches(keyMsg, m.keys.rescan):
		return m, m.RescanCmd()
	case key.Matches(keyMsg, m.keys.toggle):
		if m.wwan.Focused() {
			return m, m.toggleWWAN()
		}
		if m.wifi.Focused() {
			return m, m.toggleWIFI()
		}
		if m.networking.Focused() {
			return m, m.toggleNetworking()
		}
	}
	switch {
	case m.wwan.Focused():
		upd, cmd := m.wwan.Update(keyMsg)
		m.wwan = upd
		return m, cmd
	case m.wifi.Focused():
		upd, cmd := m.wifi.Update(keyMsg)
		m.wifi = upd
		return m, cmd
	case m.networking.Focused():
		upd, cmd := m.wifi.Update(keyMsg)
		m.wifi = upd
		return m, cmd
	case m.devicesTable.Focused():
		upd, cmd := m.devicesTable.Update(keyMsg)
		m.devicesTable = &upd
		return m, cmd
	}
	return m, nil
}

func (m *NetworkModel) View() string {
	table := styles.BorderedStyle.Render(m.devicesTable.View())

	wwan := m.wwan.View()
	wwanStyle := *m.wwanStyle
	if m.wwan.Focused() {
		wwanStyle = m.wwanStyle.Foreground(styles.AccentColor)
	}
	wwan = wwanStyle.Render(wwan)
	wwan = lipgloss.JoinHorizontal(lipgloss.Center, "WWAN:       ", wwan)

	wifi := m.wifi.View()
	wifiStyle := *m.wifiStyle
	if m.wifi.Focused() {
		wifiStyle = m.wifiStyle.Foreground(styles.AccentColor)
	}
	wifi = wifiStyle.Render(wifi)
	wifi = lipgloss.JoinHorizontal(lipgloss.Center, "Wi-Fi:      ", wifi)

	networking := m.networking.View()
	networkingStyle := *m.networkingStyle
	if m.networking.Focused() {
		networkingStyle = networkingStyle.Foreground(styles.AccentColor)
	}
	networking = networkingStyle.Render(networking)
	networking = lipgloss.JoinHorizontal(lipgloss.Center, "Networking: ", networking)

	connectivity := fmt.Sprintf("Connectivity status: %s", m.connectivity)

	statusline := m.indicatorView()

	togglers := lipgloss.JoinVertical(
		lipgloss.Left,
		wwan,
		wifi,
		networking,
		connectivity,
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		table,
		togglers,
		statusline,
	)
}

func (m *NetworkModel) indicatorView() string {
	var view string
	if m.indicatorState != NetworkDone {
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

func (m *NetworkModel) RescanCmd() tea.Cmd {
	return tea.Sequence(m.setStateCmd(NetworkScanning),
		func() tea.Msg {
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

			networkingStatus, err := m.nm.GetNetworking(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get networking status")
			}
			m.networking.SetValue(networkingStatus)

			conStatus, err := m.nm.GetConnectivityStatus(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get connection status")
			}
			m.connectivity = string(conStatus)

			return m.setStateCmd(NetworkDone)
		},
	)
}

func (m *NetworkModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *NetworkModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *NetworkModel) setStateCmd(state networkState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == NetworkDone {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func (m *NetworkModel) toggleWWAN() tea.Cmd {
	if m.indicatorState != NetworkDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(NetworkTogglingWWAN),
		func() tea.Msg {
			var err error
			if m.wwan.Value() {
				err = m.nm.DisableWWAN(context.Background())
			} else {
				err = m.nm.EnableWWAN(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling WWAN")
			}

			return m.RescanCmd()
		},
	)
}

func (m *NetworkModel) toggleWIFI() tea.Cmd {
	if m.indicatorState != NetworkDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(NetworkTogglingWifi),
		func() tea.Msg {
			var err error
			if m.wifi.Value() {
				err = m.nm.DisableWifi(context.Background())
			} else {
				err = m.nm.EnableWifi(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling Wi-Fi")
			}

			return m.RescanCmd()
		},
	)
}

func (m *NetworkModel) toggleNetworking() tea.Cmd {
	if m.indicatorState != NetworkDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(NetworkTogglingNetworking),
		func() tea.Msg {
			var err error
			if m.networking.Value() {
				err = m.nm.DisableNetworking(context.Background())
			} else {
				err = m.nm.EnableNetworking(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling networking")
			}

			return m.RescanCmd()
		},
	)
}
