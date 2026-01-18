package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	tea "github.com/charmbracelet/bubbletea"
)

const wifiWindowsCount int = 2

type WifiModel struct {
	available *WifiAvailableModel
	stored    *WifiStoredModel
	winIndex  int
	width     int
	height    int
}

func NewWifiModel(networkManager infra.NetworkManager) *WifiModel {
	wifiAvailable := NewWifiAvailable(networkManager)
	wifiStored := NewWifiStored(networkManager)

	return &WifiModel{available: wifiAvailable, stored: wifiStored}
}

func (m *WifiModel) Resize(width, height int) {
	m.width = width
	m.height = height

	storedHeight := height / 2
	availableHeight := height - storedHeight

	width -= BorderOffset
	storedHeight -= BorderOffset
	availableHeight -= BorderOffset

	m.available.Resize(width, availableHeight)
	m.stored.Resize(width, storedHeight)
}

func (m *WifiModel) Init() tea.Cmd {
	cmd := []tea.Cmd{m.available.Init(), m.stored.Init()}
	return tea.Batch(cmd...)
}

func (m *WifiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.winIndex = (m.winIndex + 1) % wifiWindowsCount
		case "1":
			m.winIndex = 0
		case "2":
			m.winIndex = 1
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
	}
	var cmds []tea.Cmd

	upd, cmd := m.available.Update(msg)
	m.available = upd.(*WifiAvailableModel)
	cmds = append(cmds, cmd)

	upd, cmd = m.stored.Update(msg)
	m.stored = upd.(*WifiStoredModel)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *WifiModel) View() string {
	availableStyle := styles.BorderedStyle.
		Width(m.available.width).
		Height(m.available.height)

	storedStyle := styles.BorderedStyle.
		Width(m.stored.width).
		Height(m.stored.height)

	if m.winIndex == 0 {
		availableStyle = availableStyle.BorderForeground(styles.AccentColor)
	} else {
		storedStyle = storedStyle.BorderForeground(styles.AccentColor)
	}

	availableView := renderer.RenderWithTitleAndKeybind(
		m.available.View(),
		"Available Wi-Fi",
		"1",
		&availableStyle,
		styles.AccentColor,
	)

	storedView := renderer.RenderWithTitleAndKeybind(
		m.stored.View(),
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
	switch m.winIndex {
	case 0:
		upd, cmd = m.available.Update(msg)
		m.available = upd.(*WifiAvailableModel)
	case 1:
		upd, cmd = m.stored.Update(msg)
		m.stored = upd.(*WifiStoredModel)
	}
	return cmd
}
