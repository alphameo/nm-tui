package views

import (
	"github.com/alphameo/nm-tui/internal/infra"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxWindows int = 2

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

	width -= borderOffset
	storedHeight -= borderOffset
	availableHeight -= borderOffset

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
			m.winIndex = (m.winIndex + 1) % maxWindows
		default:
			cmd := m.handleKeyMsg(msg)
			return m, cmd
		}
	default:
		var cmds []tea.Cmd

		upd, cmd := m.available.Update(msg)
		m.available = upd.(*WifiAvailableModel)
		cmds = append(cmds, cmd)

		upd, cmd = m.stored.Update(msg)
		m.stored = upd.(*WifiStoredModel)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *WifiModel) View() string {
	defaultStyle := lipgloss.NewStyle().
		BorderStyle(BorderStyle).
		Width(m.available.width).
		Height(m.available.height)

	selectedStyle := defaultStyle.BorderForeground(lipgloss.Color("63"))
	var viewAvailable, viewStored string
	if m.winIndex == 0 {
		viewAvailable = selectedStyle.Render(m.available.View())
		viewStored = defaultStyle.Render(m.stored.View())
	} else {
		viewAvailable = defaultStyle.Render(m.available.View())
		viewStored = selectedStyle.Render(m.stored.View())
	}

	return viewAvailable + "\n" + viewStored
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
