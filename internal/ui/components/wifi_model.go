package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiWindowIndex int

const (
	wifiAvailableWindowIndex WifiWindowIndex = iota
	wifiStoredWindowIndex
)

type WifiModel struct {
	windows            []SizedModel
	focusedWindowIndex WifiWindowIndex
	width              int
	height             int
}

func NewWifiModel(networkManager infra.NetworkManager) *WifiModel {
	wifiAvailable := NewWifiAvailableModel(networkManager)
	wifiStored := NewWifiStoredModel(networkManager)

	w := []SizedModel{wifiAvailable, wifiStored}

	return &WifiModel{windows: w}
}

func (m *WifiModel) Resize(width, height int) {
	m.width = width
	m.height = height

	storedHeight := height / 2
	availableHeight := height - storedHeight

	width -= BorderOffset
	storedHeight -= BorderOffset
	availableHeight -= BorderOffset

	m.windows[wifiAvailableWindowIndex].Resize(width, availableHeight)
	m.windows[wifiStoredWindowIndex].Resize(width, storedHeight)
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
			nextIndex := int(m.focusedWindowIndex + 1)
			m.focusedWindowIndex = WifiWindowIndex(nextIndex % len(m.windows))
		case "1":
			m.focusedWindowIndex = 0
		case "2":
			m.focusedWindowIndex = 1
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
	}
	var cmds []tea.Cmd

	upd, cmd := m.windows[wifiAvailableWindowIndex].Update(msg)
	m.windows[wifiAvailableWindowIndex] = upd.(*WifiAvailableModel)
	cmds = append(cmds, cmd)

	upd, cmd = m.windows[wifiStoredWindowIndex].Update(msg)
	m.windows[wifiStoredWindowIndex] = upd.(*WifiStoredModel)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiModel) View() string {
	availableStyle := styles.BorderedStyle.
		Width(m.windows[wifiAvailableWindowIndex].Width()).
		Height(m.windows[wifiAvailableWindowIndex].Height())

	storedStyle := styles.BorderedStyle.
		Width(m.windows[wifiStoredWindowIndex].Width()).
		Height(m.windows[wifiStoredWindowIndex].Height())

	if m.focusedWindowIndex == 0 {
		availableStyle = availableStyle.BorderForeground(styles.AccentColor)
	} else {
		storedStyle = storedStyle.BorderForeground(styles.AccentColor)
	}

	availableView := renderer.RenderWithTitleAndKeybind(
		m.windows[wifiAvailableWindowIndex].View(),
		"Available Wi-Fi",
		"1",
		&availableStyle,
		styles.AccentColor,
	)

	storedView := renderer.RenderWithTitleAndKeybind(
		m.windows[wifiStoredWindowIndex].View(),
		"Stored Wi-Fi",
		"2",
		&storedStyle,
		styles.AccentColor,
	)

	sb := strings.Builder{}
	fmt.Fprintf(
		&sb,
		"%s\n%s",
		availableView,
		storedView,
	)
	return sb.String()
}

func (m *WifiModel) handleKeyMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	switch m.focusedWindowIndex {
	case 0:
		upd, cmd = m.windows[wifiAvailableWindowIndex].Update(msg)
		m.windows[wifiAvailableWindowIndex] = upd.(*WifiAvailableModel)
	case 1:
		upd, cmd = m.windows[wifiStoredWindowIndex].Update(msg)
		m.windows[wifiStoredWindowIndex] = upd.(*WifiStoredModel)
	}
	return cmd
}
