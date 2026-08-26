package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/models/focus"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type deviceState int

const (
	DeviceNil deviceState = iota
	DeviceScanning
	DeviceTogglingWifi
	DeviceTogglingWWAN
	DeviceTogglingNetworking
	DeviceDone
)

func (s *deviceState) String() string {
	switch *s {
	case DeviceScanning:
		return "Scanning"
	case DeviceTogglingWWAN:
		return "Toggling WWAN"
	case DeviceTogglingWifi:
		return "Toggling Wi-Fi"
	case DeviceTogglingNetworking:
		return "Toggling Wi-Fi"
	case DeviceDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type deviceKeyMap struct {
	prev   key.Binding
	next   key.Binding
	rescan key.Binding
}

type DeviceModel struct {
	netDevices *NetworkDevicesModel

	wwan       toggle.Model
	wifi       toggle.Model
	networking toggle.Model

	connectivity string

	indicatorSpinner spinner.Model
	indicatorState   deviceState
	IndicatorStyle   lipgloss.Style

	ControlsStyle lipgloss.Style

	focus bool

	focuses focus.Group

	keys deviceKeyMap

	connMngr infra.DeviceManager

	Style lipgloss.Style
}

func NewDeviceModel(
	networkDevices *NetworkDevicesModel,
	keys deviceKeyMap,
	deviceManager infra.DeviceManager,
) *DeviceModel {
	wwan := newDefaultToggle()

	wifi := newDefaultToggle()

	networking := newDefaultToggle()

	s := newDefaultSpinner()

	model := &DeviceModel{
		netDevices:       networkDevices,
		indicatorSpinner: s,
		indicatorState:   DeviceDone,
		IndicatorStyle:   lipgloss.NewStyle(),

		wwan:       wwan,
		wifi:       wifi,
		networking: networking,

		connectivity: "",

		ControlsStyle: lipgloss.NewStyle().Margin(1, 0),
		connMngr:      deviceManager,
		keys:          keys,
		Style:         lipgloss.NewStyle(),
	}

	focuses := []focus.Focusable{
		networkDevices,
		&model.wwan,
		&model.wifi,
		&model.networking,
	}
	model.focuses = *focus.NewGroup(focuses)

	return model
}

func (m *DeviceModel) Resize(width, height int) {
	m.Style = m.Style.Width(width).Height(height)

	border := m.Style.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	controlsHeight := lipgloss.Height(m.controlsView())
	statuslineHeight := lipgloss.Height(m.indicatorView())
	netDevicesHeight := height - controlsHeight - statuslineHeight

	m.netDevices.Resize(width, netDevicesHeight)
}

func (m *DeviceModel) Width() int { return m.Style.GetWidth() }

func (m *DeviceModel) Height() int { return m.Style.GetHeight() }

func (m *DeviceModel) Title() string { return "Device" }

func (m *DeviceModel) Focus() { m.focus = true }

func (m *DeviceModel) Blur() { m.focus = false }

func (m *DeviceModel) Focused() bool { return m.focus }

func (m *DeviceModel) Init() tea.Cmd {
	return tea.Batch(
		m.netDevices.Init(),
		m.RescanCmd(),
		m.focuses.FocusCurrent(),
	)
}

func (m *DeviceModel) Update(msg tea.Msg) (*DeviceModel, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case RescanDeviceMsg:
		return m, m.RescanCmd()
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.next):
			return m, m.focuses.FocusCycleNextCmd()
		case key.Matches(msg, m.keys.prev):
			return m, m.focuses.FocusCyclePrevCmd()
		case key.Matches(msg, m.keys.rescan):
			return m, m.RescanCmd()
		}
		if key.Matches(msg, m.wwan.Keys.Toggle) && m.wwan.Focused() {
			return m, m.toggleWWAN()
		}
		if key.Matches(msg, m.wifi.Keys.Toggle) && m.wifi.Focused() {
			return m, m.toggleWIFI()
		}
		if key.Matches(msg, m.networking.Keys.Toggle) && m.networking.Focused() {
			return m, m.toggleNetworking()
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != DeviceDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.wwan, cmd = m.wwan.Update(msg)
	cmds = append(cmds, cmd)

	m.wifi, cmd = m.wifi.Update(msg)
	cmds = append(cmds, cmd)

	m.networking, cmd = m.networking.Update(msg)
	cmds = append(cmds, cmd)

	m.netDevices, cmd = m.netDevices.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *DeviceModel) UpdateAsTab(msg tea.Msg) (tabview.TabModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *DeviceModel) View() string {
	netDevices := m.netDevices.View()

	controls := m.controlsView()
	statusline := m.indicatorView()

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		netDevices,
		controls,
		statusline,
	)
	return m.Style.Render(view)
}

func (m *DeviceModel) indicatorView() string {
	var view string
	if m.indicatorState != DeviceDone {
		view = fmt.Sprintf(
			"%s %s",
			m.indicatorState.String(),
			m.indicatorSpinner.View(),
		)
	} else {
		view = m.indicatorState.String()
	}
	return m.IndicatorStyle.Render(view)
}

func (m *DeviceModel) controlsView() string {
	wwan := m.wwan.View()
	wwan = lipgloss.JoinHorizontal(lipgloss.Center, "WWAN       ", wwan)

	wifi := m.wifi.View()
	wifi = lipgloss.JoinHorizontal(lipgloss.Center, "Wi-Fi      ", wifi)

	networking := m.networking.View()
	networking = lipgloss.JoinHorizontal(lipgloss.Center, "Networking ", networking)

	connectivity := styles.BoldStyle.Render(m.connectivity)
	connectivity = fmt.Sprintf("Connectivity %s", connectivity)

	togglers := lipgloss.JoinVertical(
		lipgloss.Left,
		wwan,
		wifi,
		networking,
		"",
		connectivity,
	)

	return m.ControlsStyle.Render(togglers)
}

func (m *DeviceModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(DeviceScanning),
		func() tea.Msg {
			list, err := m.connMngr.ListNetworkDevices(context.Background())
			if err != nil {
				return NotifyCmd("Cannot get network devices")
			}

			m.netDevices.setProfiles(list)

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

			return m.setStateCmd(DeviceDone)
		},
	)
}

func (m *DeviceModel) setStateCmd(state deviceState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == DeviceDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func (m *DeviceModel) toggleWWAN() tea.Cmd {
	if m.indicatorState != DeviceDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(DeviceTogglingWWAN),
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

func (m *DeviceModel) toggleWIFI() tea.Cmd {
	if m.indicatorState != DeviceDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(DeviceTogglingWifi),
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

func (m *DeviceModel) toggleNetworking() tea.Cmd {
	if m.indicatorState != DeviceDone {
		return nil
	}
	return tea.Sequence(
		m.setStateCmd(DeviceTogglingNetworking),
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

type RescanDeviceMsg struct{}

func RescanDeviceCmd() tea.Cmd {
	return func() tea.Msg {
		return RescanDeviceMsg{}
	}
}
