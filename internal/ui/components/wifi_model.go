package components

import (
	"context"
	"time"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiKeyMap struct {
	nextWindow        key.Binding
	firstWindow       key.Binding
	secondWindow      key.Binding
	rescan            key.Binding
	openCaptivePortal key.Binding
}

func (k *wifiKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.nextWindow,
		k.firstWindow,
		k.secondWindow,
		k.rescan,
		k.openCaptivePortal,
	}
}

func (k *wifiKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.nextWindow,
		k.firstWindow,
		k.secondWindow,
		k.rescan,
		k.openCaptivePortal,
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
	openCaptivePortal: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("^p", "open captive portal"),
	),
}

type WifiModel struct {
	wifiAvailable  *WifiAvailableModel
	availableStyle *lipgloss.Style

	wifiSaved  *WifiSavedModel
	savedStyle *lipgloss.Style

	windows        []SizedModel // used for batch operations for wifi models
	focusWindowIdx int

	keys *wifiKeyMap

	width  int
	height int
}

func NewWifiModel(wifiAvailable *WifiAvailableModel, wifiSaved *WifiSavedModel, keys *wifiKeyMap, networkManager infra.WifiManager) *WifiModel {
	availableStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)
	savedStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)
	w := &WifiModel{
		wifiAvailable:  wifiAvailable,
		availableStyle: &availableStyle,

		wifiSaved:  wifiSaved,
		savedStyle: &savedStyle,

		keys: keys,
	}

	wins := []SizedModel{w.wifiAvailable, w.wifiSaved}
	w.windows = wins
	return w
}

func (m *WifiModel) Resize(width, height int) {
	m.width = width
	m.height = height

	savedHeight := height / 2
	availableHeight := height - savedHeight

	width -= styles.BorderOffset
	savedHeight -= styles.BorderOffset
	availableHeight -= styles.BorderOffset

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
	var cmds []tea.Cmd
	for _, window := range m.windows {
		cmds = append(cmds, window.Init())
	}
	return tea.Batch(cmds...)
}

func (m *WifiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case RescanWifiMsg:
		return m, tea.Batch(
			RescanWifiSavedCmd(msg.delay),
			RescanWifiAvailableCmd(msg.delay),
		)
	}
	var cmds []tea.Cmd

	upd, cmd := m.wifiAvailable.Update(msg)
	m.wifiAvailable = upd.(*WifiAvailableModel)
	cmds = append(cmds, cmd)

	upd, cmd = m.wifiSaved.Update(msg)
	m.wifiSaved = upd.(*WifiSavedModel)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiModel) handleKey(keyMsg tea.KeyMsg) (*WifiModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.nextWindow):
		m.focusWindowIdx = (m.focusWindowIdx + 1) % len(m.windows)
	case key.Matches(keyMsg, m.keys.firstWindow):
		m.focusWindowIdx = 0
	case key.Matches(keyMsg, m.keys.secondWindow):
		m.focusWindowIdx = 1
	case key.Matches(keyMsg, m.keys.rescan):
		return m, tea.Batch(
			RescanWifiSavedCmd(0),
			RescanWifiAvailableCmd(0),
		)
	case key.Matches(keyMsg, m.keys.openCaptivePortal):
		return m, func() tea.Msg {
			err := infra.OpenCaptivePortal(context.Background())
			if err != nil {
				return NotifyCmd("Failed open captive portal")
			}
			return NotifyCmd("Opening captive portal")
		}
	}

	var cmd tea.Cmd
	var upd tea.Model
	switch m.focusWindowIdx {
	case 0:
		upd, cmd = m.wifiAvailable.Update(keyMsg)
		m.wifiAvailable = upd.(*WifiAvailableModel)
	case 1:
		upd, cmd = m.wifiSaved.Update(keyMsg)
		m.wifiSaved = upd.(*WifiSavedModel)
	}
	return m, cmd
}

func (m *WifiModel) View() string {
	availableStyle := *m.availableStyle
	savedStyle := *m.savedStyle
	if m.focusWindowIdx == 0 {
		availableStyle = availableStyle.BorderForeground(styles.AccentColor)
	} else {
		savedStyle = savedStyle.BorderForeground(styles.AccentColor)
	}

	availableView := renderer.RenderWithTitleAndKeybind(
		m.wifiAvailable.View(),
		"Available networks",
		m.keys.firstWindow.Help().Key,
		&availableStyle,
		styles.AccentColor,
	)

	savedView := renderer.RenderWithTitleAndKeybind(
		m.wifiSaved.View(),
		"Saved networks",
		m.keys.secondWindow.Help().Key,
		&savedStyle,
		styles.AccentColor,
	)

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
