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

type wifiKeyMap struct {
	nextWindow        key.Binding
	firstWindow       key.Binding
	secondWindow      key.Binding
	rescan            key.Binding
	createProfile     key.Binding
	openCaptivePortal key.Binding
	enableHotspot     key.Binding
	createHotspot     key.Binding
}

type WifiModel struct {
	wifiAvailable *AvailableNetworksModel
	wifiSaved     *NetworkProfilesModel

	focus bool

	focuses        []Focusable // used for batch operations for wifi models
	focusWindowIdx int

	wm     infra.WifiManager
	portal infra.CaptivePortalOpener

	keys wifiKeyMap

	width  int
	height int
}

func NewWifiModel(
	wifiAvailable *AvailableNetworksModel,
	wifiSaved *NetworkProfilesModel,
	keys wifiKeyMap,
	wifiManager infra.WifiManager,
	portalOpener infra.CaptivePortalOpener,
) *WifiModel {
	w := &WifiModel{
		wifiAvailable: wifiAvailable,
		wifiSaved:     wifiSaved,
		wm:            wifiManager,
		portal:        portalOpener,
		keys:          keys,
	}

	wins := []Focusable{w.wifiAvailable, w.wifiSaved}
	w.wifiAvailable.Focus()
	w.focuses = wins
	return w
}

func (m *WifiModel) Resize(width, height int) {
	m.width = width
	m.height = height

	savedHeight := height / 2
	availableHeight := height - savedHeight

	m.wifiAvailable.Resize(width, availableHeight)
	m.wifiSaved.Resize(width, savedHeight)
}

func (m *WifiModel) Width() int {
	return m.width
}

func (m *WifiModel) Height() int {
	return m.height
}

func (m *WifiModel) Title() string {
	return "Wi-Fi"
}

func (m *WifiModel) Focus() {
	m.focus = true
}

func (m *WifiModel) Blur() {
	m.focus = false
}

func (m *WifiModel) Focused() bool {
	return m.focus
}

func (m *WifiModel) Init() tea.Cmd {
	return tea.Batch(
		m.wifiAvailable.Init(),
		m.wifiSaved.Init(),
	)
}

func (m *WifiModel) Update(msg tea.Msg) (*WifiModel, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.nextWindow):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = (m.focusWindowIdx + 1) % len(m.focuses)
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.firstWindow):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = 0
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.secondWindow):
			m.focuses[m.focusWindowIdx].Blur()
			m.focusWindowIdx = 1
			m.focuses[m.focusWindowIdx].Focus()
		case key.Matches(msg, m.keys.rescan):
			return m, tea.Batch(
				RescanWifiSavedCmd(0),
				RescanWifiAvailableCmd(0),
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
		case key.Matches(msg, m.keys.enableHotspot):
			return m, m.enableQuickHotspot()
		}
	case RescanWifiMsg:
		return m, tea.Batch(
			RescanWifiSavedCmd(msg.delay),
			RescanWifiAvailableCmd(msg.delay),
		)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.wifiAvailable, cmd = m.wifiAvailable.Update(msg)
	cmds = append(cmds, cmd)

	m.wifiSaved, cmd = m.wifiSaved.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiModel) UpdateAsTab(msg tea.Msg) (tabview.TabModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *WifiModel) View() string {
	availableView := m.wifiAvailable.View()
	savedView := m.wifiSaved.View()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		availableView,
		savedView,
	)
}

type RescanWifiMsg struct {
	delay time.Duration
}

func RescanWifiCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiMsg{delay: delay}
	}
}

func (m *WifiModel) enableQuickHotspot() tea.Cmd {
	return func() tea.Msg {
		err := m.wm.QuickHotspot(context.Background())
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Failed enabling quick wifi hotspot:\n%v", err))
		}
		return RescanWifiCmd(0)
	}
}
