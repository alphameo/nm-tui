package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiConnectorModel struct {
	wifiName string
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
		switch msg.Type {
		case tea.KeyEnter:
			pw := m.password.Value()
			return m, tea.Sequence(
				SetPopupActivityCmd(false),
				m.ConnectToWifiCmd(m.wifiName, pw),
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
	inputField := styles.BorderedStyle.Render(m.password.View())
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "SSID: %s\n%v", m.wifiName, inputField)
	return sb.String()
}

func (m *WifiConnectorModel) ConnectToWifiCmd(ssid, password string) tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorStateCmd(Connecting),
		func() tea.Msg {
			err := m.nm.ConnectWifi(ssid, password)
			if err != nil {
				return NotifyCmd(err.Error())
			}
			return nil
		},
		UpdateWifiAvailableCmd(),
		SetWifiIndicatorStateCmd(None))
}
