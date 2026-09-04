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
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type networksKeyMap struct {
	winNext           key.Binding
	winPrev           key.Binding
	win1              key.Binding
	win2              key.Binding
	rescan            key.Binding
	createProfile     key.Binding
	openCaptivePortal key.Binding
	quickHotspot      key.Binding
	createHotspot     key.Binding
}

type networksState int

const (
	NetsNil networksState = iota
	NetsScanning
	NetsActivating
	NetsDeactivating
	NetsConnecting
	NetsCreating
	NetsDone
)

func (s *networksState) String() string {
	switch *s {
	case NetsScanning:
		return "Scanning"
	case NetsActivating:
		return "Activating"
	case NetsDeactivating:
		return "Deactivating"
	case NetsConnecting:
		return "Connecting"
	case NetsCreating:
		return "Creating Connection Profile"
	case NetsDone:
		return styles.SymbolCheck
	default:
		return "Undefined"
	}
}

type NetworksModel struct {
	available *AvailableNetworksModel
	profiles  *NetworkProfilesModel

	indicatorSpinner spinner.Model
	indicatorState   networksState
	IndicatorStyle   lipgloss.Style

	focus bool

	focuses focus.Group

	netMngr infra.NetworksManager
	portal  infra.CaptivePortalOpener

	keys networksKeyMap

	Style lipgloss.Style
}

func NewNetworksModel(
	wifiAvailable *AvailableNetworksModel,
	wifiSaved *NetworkProfilesModel,
	keys networksKeyMap,
	networksManager infra.NetworksManager,
	portalOpener infra.CaptivePortalOpener,
) *NetworksModel {
	s := newDefaultSpinner()
	w := &NetworksModel{
		available: wifiAvailable,
		profiles:  wifiSaved,

		indicatorSpinner: s,
		IndicatorStyle:   lipgloss.NewStyle(),

		netMngr: networksManager,
		portal:  portalOpener,
		keys:    keys,
		Style:   lipgloss.NewStyle(),
	}

	wins := []focus.Focusable{w.available, w.profiles}
	w.focuses = *focus.NewGroup(wins)
	return w
}

func (m *NetworksModel) Resize(width, height int) {
	m.Style = m.Style.Width(width).Height(height)

	border := m.Style.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	indicatorStateHeight := lipgloss.Height(m.indicatorView())
	height -= indicatorStateHeight

	savedHeight := height / 2
	availableHeight := height - savedHeight

	m.available.Resize(width, availableHeight)
	m.profiles.Resize(width, savedHeight)
}

func (m *NetworksModel) Width() int { return m.Style.GetWidth() }

func (m *NetworksModel) Height() int { return m.Style.GetHeight() }

func (m *NetworksModel) Title() string { return "Networks" }

func (m *NetworksModel) Focus() { m.focus = true }

func (m *NetworksModel) Blur() { m.focus = false }

func (m *NetworksModel) Focused() bool { return m.focus }

func (m *NetworksModel) Init() tea.Cmd {
	return tea.Batch(
		m.listNetsCmd(),
		m.focuses.SetFocusIdx(0),
	)
}

func (m *NetworksModel) Update(msg tea.Msg) (*NetworksModel, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.winNext):
			return m, m.focuses.FocusCycleNextCmd()
		case key.Matches(msg, m.keys.winPrev):
			return m, m.focuses.FocusCyclePrevCmd()
		case key.Matches(msg, m.keys.win1):
			return m, m.focuses.SetFocusIdx(0)
		case key.Matches(msg, m.keys.win2):
			return m, m.focuses.SetFocusIdx(1)
		case key.Matches(msg, m.keys.rescan):
			if m.indicatorState != NetsDone {
				return m, nil
			}
			return m, m.rescanCmd()
		case key.Matches(msg, m.keys.createProfile):
			return m, OpenProfileCreatorCmd()
		case key.Matches(msg, m.keys.createHotspot):
			return m, OpenHotspotCreatorCmd()
		case key.Matches(msg, m.keys.openCaptivePortal):
			return m, func() tea.Msg {
				err := m.portal.OpenCaptivePortal(context.Background())
				if err != nil {
					return NotifyCmd("Failed open captive portal")
				}
				return NotifyCmd("Opening captive portal")
			}
		case key.Matches(msg, m.keys.quickHotspot):
			return m, m.quickHotspot()
		}
	case RescanNetworksMsg:
		return m, m.rescanCmd()
	case QuickRescanNetworksMsg:
		return m, m.quickRescanCmd()
	case NetworksStateMsg:
		return m, m.setStateCmd(networksState(msg))
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.available, cmd = m.available.Update(msg)
	cmds = append(cmds, cmd)

	m.profiles, cmd = m.profiles.Update(msg)
	cmds = append(cmds, cmd)

	if m.indicatorState != NetsDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *NetworksModel) View() string {
	availableView := m.available.View()
	savedView := m.profiles.View()

	statusline := m.indicatorView()
	view := lipgloss.JoinVertical(
		lipgloss.Center,
		availableView,
		savedView,
		statusline,
	)
	return m.Style.Render(view)
}

func (m *NetworksModel) indicatorView() string {
	var view string
	if m.indicatorState != NetsDone {
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

type NetworksStateMsg networksState

func (m *NetworksModel) setStateCmd(state networksState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == NetsDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func SetNetworksStateCmd(state networksState) tea.Cmd {
	return func() tea.Msg {
		return NetworksStateMsg(state)
	}
}

type RescanNetworksMsg struct{}

func RescanNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		return RescanNetworksMsg{}
	}
}

type QuickRescanNetworksMsg struct{}

func QuickRescanNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		return QuickRescanNetworksMsg{}
	}
}

type NetworksRescannedMsg struct {
	Available   []AvailableNetwork
	Profiles    []NetworkProfileShort
	ScanErr     error
	ProfilesErr error
}

func (m *NetworksModel) listNetsCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetsScanning),
		func() tea.Msg {
			ctx := context.Background()
			availableRecords, scanErr := m.netMngr.ListNetworks(ctx)
			availables := convertAvailableNetworks(availableRecords)
			profileRecords, profilesErr := m.netMngr.ListProfiles(ctx)
			profiles := convertNetworkProfileShorts(profileRecords)
			availableExt, profilesExt := CrossReferenceNetworks(availables, profiles)
			return NetworksRescannedMsg{
				Available:   availableExt,
				Profiles:    profilesExt,
				ScanErr:     scanErr,
				ProfilesErr: profilesErr,
			}
		},
	)
}

func (m *NetworksModel) quickRescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetsScanning),
		func() tea.Msg {
			ctx := context.Background()
			availableRecords, scanErr := m.netMngr.ListNetworks(ctx)
			availables := convertAvailableNetworks(availableRecords)
			profileRecords, profilesErr := m.netMngr.ListProfiles(ctx)
			profiles := convertNetworkProfileShorts(profileRecords)
			availableExt, profilesExt := CrossReferenceNetworks(availables, profiles)
			return NetworksRescannedMsg{
				Available:   availableExt,
				Profiles:    profilesExt,
				ScanErr:     scanErr,
				ProfilesErr: profilesErr,
			}
		},
	)
}

func (m *NetworksModel) rescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetsScanning),
		func() tea.Msg {
			ctx := context.Background()
			availableRecords, scanErr := m.netMngr.ListNetworksWithRescan(ctx)
			availables := convertAvailableNetworks(availableRecords)
			profileRecords, profilesErr := m.netMngr.ListProfiles(ctx)
			profiles := convertNetworkProfileShorts(profileRecords)
			availableExt, profilesExt := CrossReferenceNetworks(availables, profiles)
			return NetworksRescannedMsg{
				Available:   availableExt,
				Profiles:    profilesExt,
				ScanErr:     scanErr,
				ProfilesErr: profilesErr,
			}
		},
	)
}

func (m *NetworksModel) quickHotspot() tea.Cmd {
	return func() tea.Msg {
		err := m.netMngr.QuickHotspot(context.Background())
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Failed enabling quick wifi hotspot:\n%v", err))
		}
		return QuickRescanNetworksCmd()
	}
}

type AvailableNetwork struct {
	SSID          string
	Active        bool
	SecurityMode  string
	Signal        int
	Band          float64
	Rate          float64
	LookingDevice string
	NetworkMode   infra.NetworkMode
	ProfileExists bool
}

type NetworkProfileShort struct {
	Name      string
	UUID      string
	SSID      string
	Active    bool
	Mode      string
	Available bool
}

// CrossReferenceNetworks matches available networks and saved profiles by SSID:
// marks each profile Available when its SSID is currently in range and each
// available network ProfileExists when a saved profile targets the same SSID.
func CrossReferenceNetworks(
	available []AvailableNetwork,
	profiles []NetworkProfileShort,
) ([]AvailableNetwork, []NetworkProfileShort) {
	availableSSIDs := make(map[string]struct{}, len(available))
	for _, net := range available {
		if net.SSID != "" {
			availableSSIDs[net.SSID] = struct{}{}
		}
	}

	profileSSIDs := make(map[string]struct{}, len(profiles))
	for _, profile := range profiles {
		if profile.SSID != "" {
			profileSSIDs[profile.SSID] = struct{}{}
		}
	}

	res := make([]AvailableNetwork, len(available))
	for i, net := range available {
		_, net.ProfileExists = profileSSIDs[net.SSID]
		res[i] = net
	}

	prof := make([]NetworkProfileShort, len(profiles))
	for i, profile := range profiles {
		_, profile.Available = availableSSIDs[profile.SSID]
		prof[i] = profile
	}

	return res, prof
}
