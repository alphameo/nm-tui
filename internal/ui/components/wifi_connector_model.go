package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiConnectorKeyMap struct {
	connect                  key.Binding
	togglePasswordVisibility key.Binding
}

func (k *wifiConnectorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.connect, k.togglePasswordVisibility}
}

func (k *wifiConnectorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.connect, k.togglePasswordVisibility}}
}

type WifiConnectorModel struct {
	wifiName string
	password textinput.Model

	keys *wifiConnectorKeyMap

	nm infra.NetworkManager
}

func NewWifiConnector(keys *wifiConnectorKeyMap, networkManager infra.NetworkManager) *WifiConnectorModel {
	p := textinput.New()
	p.Focus()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"
	return &WifiConnectorModel{
		password: p,
		keys:     keys,
		nm:       networkManager,
	}
}

func (m *WifiConnectorModel) setNew(wifiName string) {
	m.wifiName = wifiName
	pw, err := m.nm.GetWifiPassword(wifiName)
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
		switch {
		case key.Matches(msg, m.keys.connect):
			pw := m.password.Value()
			return m, tea.Batch(
				SetPopupActivityCmd(false),
				m.connectToWifiCmd(m.wifiName, pw),
			)
		case key.Matches(msg, m.keys.togglePasswordVisibility):
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
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "SSID      %s", m.wifiName)
	wifiName := sb.String()
	password := styles.BorderedStyle.Render(m.password.View())
	password = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		password,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wifiName,
		password,
	)
}

func (m *WifiConnectorModel) connectToWifiCmd(ssid, password string) tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(ConnectingAvailable),
		func() tea.Msg {
			err := m.nm.ConnectWifi(ssid, password)
			if err != nil {
				return NotifyCmd(err.Error())
			}
			return SetWifiAvailableStateCmd(DoneInAvailable)
		},
	)
}
