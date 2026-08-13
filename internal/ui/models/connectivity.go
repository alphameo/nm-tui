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
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type connectivityConfig struct {
	togglersStyle lipgloss.Style
	deviceColIdx  int
	typeColIdx    int
	connColIdx    int
	stateColIdx   int

	deviceWidthProportion float32
	typeWidthProportion   float32
	stateWidthProportion  float32
}

var connectivityCfg = connectivityConfig{
	togglersStyle:         lipgloss.NewStyle().Margin(1, 0),
	deviceColIdx:          0,
	typeColIdx:            1,
	connColIdx:            2,
	stateColIdx:           3,
	deviceWidthProportion: 0.2,
	typeWidthProportion:   0.15,
	stateWidthProportion:  0.3,
}

type connectivityState int

const (
	ConnectivityNil connectivityState = iota
	ConnectivityScanning
	ConnectivityTogglingWifi
	ConnectivityTogglingWWAN
	ConnectivityTogglingNetworking
	ConnectivityDone
)

func (s *connectivityState) String() string {
	switch *s {
	case ConnectivityScanning:
		return "Scanning"
	case ConnectivityTogglingWWAN:
		return "Toggling WWAN"
	case ConnectivityTogglingWifi:
		return "Toggling Wi-Fi"
	case ConnectivityTogglingNetworking:
		return "Toggling Wi-Fi"
	case ConnectivityDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type connectivityKeyMap struct {
	prev   key.Binding
	next   key.Binding
	rescan key.Binding
	toggle key.Binding
}

type ConnectivityModel struct {
	devicesTable table.Model

	wwan       toggle.Model
	wifi       toggle.Model
	networking toggle.Model

	connectivity string

	indicatorSpinner spinner.Model
	indicatorState   connectivityState

	focus bool

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys connectivityKeyMap

	connMngr infra.ConnectivityManager

	height int
	width  int
}

func NewConnectivityModel(keys connectivityKeyMap, connectivityManager infra.ConnectivityManager) *ConnectivityModel {
	cols := make([]table.Column, 4)
	cols[connectivityCfg.deviceColIdx] = table.Column{Title: "Device"}
	cols[connectivityCfg.typeColIdx] = table.Column{Title: "Type"}
	cols[connectivityCfg.connColIdx] = table.Column{Title: "Connection"}
	cols[connectivityCfg.stateColIdx] = table.Column{Title: "State"}

	t := table.New(
		table.WithColumns(cols),
		table.WithStyles(styles.DataTableStyle),
	)

	wwan := toggle.New()
	wwan.Symbols = styles.SymbolsToggle

	wifi := toggle.New()
	wifi.Symbols = styles.SymbolsToggle

	networking := toggle.New()
	networking.Symbols = styles.SymbolsToggle

	s := spinner.New()
	s.Spinner = styles.Spinner

	model := &ConnectivityModel{
		devicesTable:     t,
		indicatorSpinner: s,
		indicatorState:   ConnectivityDone,

		wwan:       wwan,
		wifi:       wifi,
		networking: networking,

		connectivity: "",

		connMngr: connectivityManager,
		keys:     keys,
	}

	focuses := []Focusable{
		&model.wwan,
		&model.wifi,
		&model.networking,
	}
	model.focuses = focuses

	return model
}

func (m *ConnectivityModel) Resize(width, height int) {
	m.height = height
	m.width = width

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	m.devicesTable.SetWidth(width)
	m.devicesTable.SetHeight(height - 9)

	tableUtilityOffset := len(m.devicesTable.Columns()) * 2

	deviceColWidth := int(float32(width) * connectivityCfg.deviceWidthProportion)
	typeColWidth := int(float32(width) * connectivityCfg.typeWidthProportion)
	stateWidth := int(float32(width) * connectivityCfg.stateWidthProportion)
	connWidth := width - typeColWidth - deviceColWidth - tableUtilityOffset - stateWidth

	m.devicesTable.Columns()[connectivityCfg.deviceColIdx].Width = deviceColWidth
	m.devicesTable.Columns()[connectivityCfg.typeColIdx].Width = typeColWidth
	m.devicesTable.Columns()[connectivityCfg.stateColIdx].Width = stateWidth
	m.devicesTable.Columns()[connectivityCfg.connColIdx].Width = connWidth
	m.devicesTable.UpdateViewport()
}

func (m *ConnectivityModel) Width() int {
	return m.width
}

func (m *ConnectivityModel) Height() int {
	return m.height
}

func (m *ConnectivityModel) Title() string {
	return "Connectivity"
}

func (m *ConnectivityModel) Focus() {
	m.focus = true
}

func (m *ConnectivityModel) Blur() {
	m.focus = false
}

func (m *ConnectivityModel) Focused() bool {
	return m.focus
}

func (m *ConnectivityModel) Init() tea.Cmd {
	return tea.Batch(
		m.RescanCmd(),
		m.focuses[m.focusIdx].Focus(),
	)
}

func (m *ConnectivityModel) Update(msg tea.Msg) (*ConnectivityModel, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.next):
			return m, m.focusNextCmd()
		case key.Matches(msg, m.keys.prev):
			return m, m.focusPrevCmd()
		case key.Matches(msg, m.keys.rescan):
			return m, m.RescanCmd()
		case key.Matches(msg, m.keys.toggle):
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
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != ConnectivityDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.wwan, cmd = m.wwan.Update(msg)
	cmds = append(cmds, cmd)

	m.wifi, cmd = m.wifi.Update(msg)
	cmds = append(cmds, cmd)

	m.networking, cmd = m.networking.Update(msg)
	cmds = append(cmds, cmd)

	m.devicesTable, cmd = m.devicesTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ConnectivityModel) UpdateAsTab(msg tea.Msg) (tabview.TabModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ConnectivityModel) View() string {
	table := styles.BorderedStyle.Render(m.devicesTable.View())

	wwan := styles.ViewToggle(m.wwan)
	wwan = lipgloss.JoinHorizontal(lipgloss.Center, "WWAN       ", wwan)

	wifi := styles.ViewToggle(m.wifi)
	wifi = lipgloss.JoinHorizontal(lipgloss.Center, "Wi-Fi      ", wifi)

	networking := styles.ViewToggle(m.networking)
	networking = lipgloss.JoinHorizontal(lipgloss.Center, "Networking ", networking)

	connectivity := styles.BoldStyle.Render(m.connectivity)
	connectivity = fmt.Sprintf("Connectivity %s", connectivity)

	statusline := m.indicatorView()

	togglers := lipgloss.JoinVertical(
		lipgloss.Left,
		wwan,
		wifi,
		networking,
		"",
		connectivity,
	)
	togglers = connectivityCfg.togglersStyle.Render(togglers)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		table,
		togglers,
		statusline,
	)
}

func (m *ConnectivityModel) indicatorView() string {
	var view string
	if m.indicatorState != ConnectivityDone {
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

func (m *ConnectivityModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ConnectivityScanning),
		func() tea.Msg {
			list, err := m.connMngr.ListDevices(context.Background())
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

			radioStatus, err := m.connMngr.GetRadioStatus(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get radio status")
			}
			m.wwan.SetValue(radioStatus.EnabledWWAN)
			m.wifi.SetValue(radioStatus.EnabledWifi)

			networkingStatus, err := m.connMngr.IsNetworkingEnabled(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get networking status")
			}
			m.networking.SetValue(networkingStatus)

			conStatus, err := m.connMngr.GetConnectivityStatus(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get connection status")
			}
			m.connectivity = conStatus.String()

			return m.setStateCmd(ConnectivityDone)
		},
	)
}

func (m *ConnectivityModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectivityModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectivityModel) setStateCmd(state connectivityState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == ConnectivityDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func (m *ConnectivityModel) toggleWWAN() tea.Cmd {
	if m.indicatorState != ConnectivityDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(ConnectivityTogglingWWAN),
		func() tea.Msg {
			var err error
			if m.wwan.Value() {
				err = m.connMngr.DisableWWAN(context.Background())
			} else {
				err = m.connMngr.EnableWWAN(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling WWAN")
			}

			return m.RescanCmd()
		},
	)
}

func (m *ConnectivityModel) toggleWIFI() tea.Cmd {
	if m.indicatorState != ConnectivityDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(ConnectivityTogglingWifi),
		func() tea.Msg {
			var err error
			if m.wifi.Value() {
				err = m.connMngr.DisableWifi(context.Background())
			} else {
				err = m.connMngr.EnableWifi(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling Wi-Fi")
			}

			return m.RescanCmd()
		},
	)
}

func (m *ConnectivityModel) toggleNetworking() tea.Cmd {
	if m.indicatorState != ConnectivityDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(ConnectivityTogglingNetworking),
		func() tea.Msg {
			var err error
			if m.networking.Value() {
				err = m.connMngr.DisableNetworking(context.Background())
			} else {
				err = m.connMngr.EnableNetworking(context.Background())
			}
			if err != nil {
				return NotifyCmd("Failed toggling networking")
			}

			return m.RescanCmd()
		},
	)
}
