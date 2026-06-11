package components

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/tabview"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type wifiKeyMap struct {
	nextWindow        key.Binding
	firstWindow       key.Binding
	secondWindow      key.Binding
	rescan            key.Binding
	create            key.Binding
	openCaptivePortal key.Binding
	enableHotspot     key.Binding
	createHotspot     key.Binding
}

func (k *wifiKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.nextWindow,
		k.firstWindow,
		k.secondWindow,
		k.rescan,
		k.create,
		k.openCaptivePortal,
		k.enableHotspot,
		k.createHotspot,
	}
}

func (k *wifiKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.nextWindow,
		k.firstWindow,
		k.secondWindow,
		k.rescan,
		k.create,
		k.openCaptivePortal,
		k.enableHotspot,
		k.createHotspot,
	}}
}

var wifiKeys = &wifiKeyMap{
	nextWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	firstWindow: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "first window"),
	),
	secondWindow: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "second window"),
	),
	rescan: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^r", "rescan"),
	),
	create: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "create new"),
	),
	openCaptivePortal: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("^p", "open captive portal"),
	),
	enableHotspot: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("h", "enable quick hotspot"),
	),
	createHotspot: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("H", "create hotspot"),
	),
}

type WifiModel struct {
	wifiAvailable  *WifiAvailableModel
	availableStyle *lipgloss.Style

	wifiSaved  *WifiSavedModel
	savedStyle *lipgloss.Style

	focuses        []Focusable // used for batch operations for wifi models
	focusWindowIdx int

	profileCreator *ProfileCreatorModel
	hotspotCreator *HotspotCreatorModel

	wm infra.WifiManager

	keys *wifiKeyMap

	width  int
	height int
}

func NewWifiModel(
	wifiAvailable *WifiAvailableModel,
	wifiSaved *WifiSavedModel,
	profileCreator *ProfileCreatorModel,
	hotspotCreator *HotspotCreatorModel,
	keys *wifiKeyMap,
	wifiManager infra.WifiManager,
) *WifiModel {
	availableStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)
	savedStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)
	w := &WifiModel{
		wifiAvailable:  wifiAvailable,
		availableStyle: &availableStyle,

		wifiSaved:  wifiSaved,
		savedStyle: &savedStyle,

		profileCreator: profileCreator,
		hotspotCreator: hotspotCreator,

		wm: wifiManager,

		keys: keys,
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

	availableBluredStyle := m.availableStyle.
		Width(m.wifiAvailable.Width()).
		Height(m.wifiAvailable.Height())
	m.availableStyle = &availableBluredStyle

	savedBluredStyle := m.savedStyle.
		Width(m.wifiSaved.Width()).
		Height(m.wifiSaved.Height())
	m.savedStyle = &savedBluredStyle
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

func (m *WifiModel) Init() tea.Cmd {
	return tea.Batch(
		m.wifiAvailable.Init(),
		m.wifiSaved.Init(),
	)
}

func (m *WifiModel) Update(msg tea.Msg) (*WifiModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
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

func (m *WifiModel) handleKey(keyMsg tea.KeyPressMsg) (*WifiModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.nextWindow):
		m.focuses[m.focusWindowIdx].Blur()
		m.focusWindowIdx = (m.focusWindowIdx + 1) % len(m.focuses)
		m.focuses[m.focusWindowIdx].Focus()
	case key.Matches(keyMsg, m.keys.firstWindow):
		m.focusWindowIdx = 0
		m.wifiSaved.Blur()
		m.wifiAvailable.Focus()
	case key.Matches(keyMsg, m.keys.secondWindow):
		m.focusWindowIdx = 1
		m.wifiSaved.Focus()
		m.wifiAvailable.Blur()
	case key.Matches(keyMsg, m.keys.rescan):
		return m, tea.Batch(
			RescanWifiSavedCmd(0),
			RescanWifiAvailableCmd(0),
		)
	case key.Matches(keyMsg, m.keys.create):
		return m, m.profileCreator.open()
	case key.Matches(keyMsg, m.keys.createHotspot):
		return m, m.hotspotCreator.open()
	case key.Matches(keyMsg, m.keys.openCaptivePortal):
		return m, func() tea.Msg {
			err := infra.OpenCaptivePortal(context.Background())
			if err != nil {
				return NotifyCmd("Failed open captive portal")
			}
			return NotifyCmd("Opening captive portal")
		}
	case key.Matches(keyMsg, m.keys.enableHotspot):
		return m, m.enableQuickHotspot()
	}

	var cmd tea.Cmd
	switch m.focusWindowIdx {
	case 0:
		m.wifiAvailable, cmd = m.wifiAvailable.Update(keyMsg)
	case 1:
		m.wifiSaved, cmd = m.wifiSaved.Update(keyMsg)
	}
	return m, cmd
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
		err := m.wm.EnableQuickWifiHotspot(context.Background())
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Failed enabling quick wifi hotspot:\n%v", err))
		}
		return RescanWifiCmd(0)
	}
}
