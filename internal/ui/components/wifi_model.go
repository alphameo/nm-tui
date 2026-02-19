package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiKeyMap struct {
	nextWindow   key.Binding
	firstWindow  key.Binding
	secondWindow key.Binding
	rescan       key.Binding
}

func (k *wifiKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.nextWindow, k.firstWindow, k.secondWindow, k.rescan}
}

func (k *wifiKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.nextWindow, k.firstWindow, k.secondWindow, k.rescan}}
}

type WifiModel struct {
	windows            []SizedModel // used for batch operations for wifi models
	wifiAvailable      *WifiAvailableModel
	wifiStored         *WifiStoredModel
	focusedWindowIndex int

	keys *wifiKeyMap

	width  int
	height int
}

func NewWifiModel(networkManager infra.NetworkManager, keys *wifiKeyMap) *WifiModel {
	a := NewWifiAvailableModel(networkManager, wifiAvailableKeys)
	s := NewWifiStoredModel(networkManager, wifiStoredKeys)
	w := &WifiModel{wifiAvailable: a, wifiStored: s, keys: keys}

	wins := []SizedModel{w.wifiAvailable, w.wifiStored}
	w.windows = wins
	return w
}

func (m *WifiModel) Resize(width, height int) {
	m.width = width
	m.height = height

	storedHeight := height / 2
	availableHeight := height - storedHeight

	width -= BorderOffset
	storedHeight -= BorderOffset
	availableHeight -= BorderOffset

	m.wifiAvailable.Resize(width, availableHeight)
	m.wifiStored.Resize(width, storedHeight)
}

func (m *WifiModel) Width() int {
	return m.width
}

func (m *WifiModel) Height() int {
	return m.height
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
			m.focusedWindowIndex = (m.focusedWindowIndex + 1) % len(m.windows)
		case key.Matches(msg, m.keys.firstWindow):
			m.focusedWindowIndex = 0
		case key.Matches(msg, m.keys.secondWindow):
			m.focusedWindowIndex = 1
		case key.Matches(msg, m.keys.rescan):
			return m, tea.Batch(
				UpdateWifiStoredCmd(),
				UpdateWifiAvailableCmd(),
			)
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
	case UpdateWifiMsg:
		return m, tea.Batch(
			UpdateWifiStoredCmd(),
			UpdateWifiAvailableCmd(),
		)
	}
	var cmds []tea.Cmd

	upd, cmd := m.wifiAvailable.Update(msg)
	m.wifiAvailable = upd.(*WifiAvailableModel)
	cmds = append(cmds, cmd)

	upd, cmd = m.wifiStored.Update(msg)
	m.wifiStored = upd.(*WifiStoredModel)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiModel) View() string {
	availableStyle := styles.BorderedStyle.
		Width(m.wifiAvailable.Width()).
		Height(m.wifiAvailable.Height())

	storedStyle := styles.BorderedStyle.
		Width(m.wifiStored.Width()).
		Height(m.wifiStored.Height())

	if m.focusedWindowIndex == 0 {
		availableStyle = availableStyle.BorderForeground(styles.AccentColor)
	} else {
		storedStyle = storedStyle.BorderForeground(styles.AccentColor)
	}

	availableView := renderer.RenderWithTitleAndKeybind(
		m.wifiAvailable.View(),
		"Available Wi-Fi",
		m.keys.firstWindow.Help().Key,
		&availableStyle,
		styles.AccentColor,
	)

	storedView := renderer.RenderWithTitleAndKeybind(
		m.wifiStored.View(),
		"Stored Wi-Fi",
		m.keys.secondWindow.Help().Key,
		&storedStyle,
		styles.AccentColor,
	)

	return lipgloss.JoinVertical(lipgloss.Center, availableView, storedView)
}

func (m *WifiModel) handleKeyMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	switch m.focusedWindowIndex {
	case 0:
		upd, cmd = m.wifiAvailable.Update(msg)
		m.wifiAvailable = upd.(*WifiAvailableModel)
	case 1:
		upd, cmd = m.wifiStored.Update(msg)
		m.wifiStored = upd.(*WifiStoredModel)
	}
	return cmd
}

// UpdateWifiMsg is a fictive struct, which used to send as tea.Msg instead of nil to trigger main window re-render
type UpdateWifiMsg struct{}

func UpdateWifiCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateWifiMsg{}
	}
}
