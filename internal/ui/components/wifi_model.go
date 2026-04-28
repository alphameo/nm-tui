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

type WifiModel struct {
	windows        []SizedModel // used for batch operations for wifi models
	wifiAvailable  *WifiAvailableModel
	wifiSaved      *WifiSavedModel
	focusWindowIdx int

	keys *wifiKeyMap

	width  int
	height int
}

func NewWifiModel(wifiAvailable *WifiAvailableModel, wifiSaved *WifiSavedModel, keys *wifiKeyMap, networkManager infra.WifiManager) *WifiModel {
	w := &WifiModel{
		wifiAvailable: wifiAvailable,
		wifiSaved:     wifiSaved,
		keys:          keys,
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
		switch {
		case key.Matches(msg, m.keys.nextWindow):
			m.focusWindowIdx = (m.focusWindowIdx + 1) % len(m.windows)
		case key.Matches(msg, m.keys.firstWindow):
			m.focusWindowIdx = 0
		case key.Matches(msg, m.keys.secondWindow):
			m.focusWindowIdx = 1
		case key.Matches(msg, m.keys.rescan):
			return m, tea.Batch(
				RescanWifiSavedCmd(0),
				RescanWifiAvailableCmd(0),
			)
		case key.Matches(msg, m.keys.openCaptivePortal):
			return m, func() tea.Msg {
				err := infra.OpenCaptivePortal(context.Background())
				if err != nil {
					return NotifyCmd("Failed open captive portal")
				}
				return NotifyCmd("Opening captive portal")
			}
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
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

func (m *WifiModel) View() string {
	availableStyle := styles.BorderedStyle.
		Width(m.wifiAvailable.Width()).
		Height(m.wifiAvailable.Height())

	savedStyle := styles.BorderedStyle.
		Width(m.wifiSaved.Width()).
		Height(m.wifiSaved.Height())

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

func (m *WifiModel) handleKeyMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	switch m.focusWindowIdx {
	case 0:
		upd, cmd = m.wifiAvailable.Update(msg)
		m.wifiAvailable = upd.(*WifiAvailableModel)
	case 1:
		upd, cmd = m.wifiSaved.Update(msg)
		m.wifiSaved = upd.(*WifiSavedModel)
	}
	return cmd
}

type RescanWifiMsg struct {
	delay time.Duration
}

func RescanWifiCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiMsg{delay: delay}
	}
}
