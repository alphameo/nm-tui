package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WifiModel struct {
	windows            []SizedModel // used for batch operations for wifi models
	wifiAvailable      *WifiAvailableModel
	wifiStored         *WifiStoredModel
	focusedWindowIndex int
	width              int
	height             int
}

func NewWifiModel(networkManager infra.NetworkManager) *WifiModel {
	a := NewWifiAvailableModel(networkManager)
	s := NewWifiStoredModel(networkManager)
	w := &WifiModel{wifiAvailable: a, wifiStored: s}

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
		switch msg.String() {
		case "tab":
			m.focusedWindowIndex = (m.focusedWindowIndex + 1) % len(m.windows)
		case "1":
			m.focusedWindowIndex = 0
		case "2":
			m.focusedWindowIndex = 1
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
	case updateWifiMsg:
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
		"1",
		&availableStyle,
		styles.AccentColor,
	)

	storedView := renderer.RenderWithTitleAndKeybind(
		m.wifiStored.View(),
		"Stored Wi-Fi",
		"2",
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

type updateWifiMsg struct{}

// UpdateWifiMsg is used to avoid extra instantiatons
var UpdateWifiMsg = updateWifiMsg{}

func UpdateWifiCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateWifiMsg
	}
}
