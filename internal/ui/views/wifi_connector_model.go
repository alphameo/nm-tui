package views

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WifiConnectorModel struct {
	ssid     string
	password textinput.Model
	nm       infra.NetworkManager
}

func NewWifiConnector(networkManager infra.NetworkManager) *WifiConnectorModel {
	p := textinput.New()
	p.Focus()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"
	return &WifiConnectorModel{password: p, nm: networkManager}
}

func (m *WifiConnectorModel) setNew(ssid string) {
	m.ssid = ssid
	pw, err := m.nm.GetWifiPassword(ssid)
	if err == nil {
		m.password.SetValue(pw)
	}
}

func (m WifiConnectorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WifiConnectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			pw := m.password.Value()
			return m, tea.Sequence(
				SetPopupActivity(false),
				m.Connect(m.ssid, pw),
			)
		case tea.KeyCtrlR:
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
		}
	}

	var cmd tea.Cmd
	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m WifiConnectorModel) View() string {
	inputField := lipgloss.
		NewStyle().
		BorderStyle(styles.BorderStyle).
		Render(m.password.View())
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "SSID: %s\n%v", m.ssid, inputField)
	return sb.String()
}

type WifiConnectionMsg struct {
	err  error
	ssid string
}

func (m *WifiConnectorModel) Connect(ssid, password string) tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorState(Connecting),
		func() tea.Msg {
			err := m.nm.ConnectWifi(ssid, password)
			return WifiConnectionMsg{err: err, ssid: ssid}
		},
		SetWifiIndicatorState(None))
}
