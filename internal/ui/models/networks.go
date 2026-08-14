package models

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
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

type NetworksModel struct {
	available *AvailableNetworksModel
	profiles  *NetworkProfilesModel

	focus bool

	focuses        []Focusable // used for batch operations for wifi models
	focusWindowIdx int

	netMngr infra.NetworksManager
	portal  infra.CaptivePortalOpener

	keys networksKeyMap

	width  int
	height int
}

func NewNetworksModel(
	wifiAvailable *AvailableNetworksModel,
	wifiSaved *NetworkProfilesModel,
	keys networksKeyMap,
	networksManager infra.NetworksManager,
	portalOpener infra.CaptivePortalOpener,
) *NetworksModel {
	w := &NetworksModel{
		available: wifiAvailable,
		profiles:  wifiSaved,
		netMngr:   networksManager,
		portal:    portalOpener,
		keys:      keys,
	}

	wins := []Focusable{w.available, w.profiles}
	w.available.Focus()
	w.focuses = wins
	return w
}

func (m *NetworksModel) Resize(width, height int) {
	m.width = width
	m.height = height

	savedHeight := height / 2
	availableHeight := height - savedHeight

	m.available.Resize(width, availableHeight)
	m.profiles.Resize(width, savedHeight)
}

func (m *NetworksModel) Width() int {
	return m.width
}

func (m *NetworksModel) Height() int {
	return m.height
}

func (m *NetworksModel) Title() string {
	return "Networks"
}

func (m *NetworksModel) Focus() {
	m.focus = true
}

func (m *NetworksModel) Blur() {
	m.focus = false
}

func (m *NetworksModel) Focused() bool {
	return m.focus
}

func (m *NetworksModel) Init() tea.Cmd {
	return tea.Batch(
		m.available.Init(),
		m.profiles.Init(),
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
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = (m.focusWindowIdx + 1) % len(m.focuses)
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.winPrev):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = (len(m.focuses) + m.focusWindowIdx - 1) % len(m.focuses)
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.win1):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = 0
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.win2):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = 1
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.rescan):
			return m, tea.Batch(
				RescanNetworkProfilesCmd(0),
				RescanAvailableNetworksCmd(0),
			)
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
		return m, tea.Batch(
			RescanNetworkProfilesCmd(msg.delay),
			RescanAvailableNetworksCmd(msg.delay),
		)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.available, cmd = m.available.Update(msg)
	cmds = append(cmds, cmd)

	m.profiles, cmd = m.profiles.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NetworksModel) UpdateAsTab(msg tea.Msg) (tabview.TabModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *NetworksModel) View() string {
	availableView := m.available.View()
	savedView := m.profiles.View()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		availableView,
		savedView,
	)
}

type RescanNetworksMsg struct {
	delay time.Duration
}

func RescanNetworksCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanNetworksMsg{delay: delay}
	}
}

func (m *NetworksModel) quickHotspot() tea.Cmd {
	return func() tea.Msg {
		err := m.netMngr.QuickHotspot(context.Background())
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Failed enabling quick wifi hotspot:\n%v", err))
		}
		return RescanNetworksCmd(0)
	}
}
