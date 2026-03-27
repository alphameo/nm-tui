package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiConnectorKeyMap struct {
	togglePasswordVisibility key.Binding
	up                       key.Binding
	down                     key.Binding
	connect                  key.Binding
}

func (k *wifiConnectorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePasswordVisibility, k.up, k.down, k.connect}
}

func (k *wifiConnectorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePasswordVisibility, k.up, k.down, k.connect}}
}

type WifiConnectorModel struct {
	wifiName string
	password textinput.Model
	hidden   *toggle.Model

	inputs            []Focusable // used for batch operations on input focusable elements
	focusedInputIndex int

	keys *wifiConnectorKeyMap

	nm infra.NetworkManager
}

func NewWifiConnector(keys *wifiConnectorKeyMap, networkManager infra.NetworkManager) *WifiConnectorModel {
	p := textinput.New()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"

	t := toggle.New(false)

	model := &WifiConnectorModel{
		password: p,
		hidden:   t,
		keys:     keys,
		nm:       networkManager,
	}

	inp := []Focusable{
		&model.password,
		model.hidden,
	}
	model.inputs = inp

	return model
}

func (m *WifiConnectorModel) setNew(wifiName string) tea.Cmd {
	m.wifiName = wifiName

	m.password.Reset()
	pw, err := m.nm.GetWifiPassword(wifiName)
	if err == nil {
		m.password.SetValue(pw)
	}
	m.password.Blur()

	m.hidden.SetValue(false)
	m.hidden.Blur()

	return m.inputs[0].Focus()
}

func (m *WifiConnectorModel) Init() tea.Cmd {
	return m.inputs[0].Focus()
}

func (m *WifiConnectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.down):
			return m, m.focusNextCmd()
		case key.Matches(msg, m.keys.up):
			return m, m.focusPrevCmd()
		case key.Matches(msg, m.keys.togglePasswordVisibility):
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case key.Matches(msg, m.keys.connect):
			return m, tea.Sequence(
				SetPopupActivityCmd(false),
				m.connectToWifiCmd(),
			)
		default:
			return m.handleKey(msg)
		}
	default:
		return m.handleMsg(msg)
	}
}

func (m *WifiConnectorModel) View() string {
	sb := strings.Builder{}
	inputStyle := styles.BorderedStyle

	fmt.Fprintf(&sb, "SSID      %s", m.wifiName)
	wifiName := sb.String()

	password := m.password.View()
	if m.password.Focused() {
		password = inputStyle.
			BorderForeground(styles.AccentColor).
			Render(password)
	} else {
		password = inputStyle.Render(password)
	}
	password = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		password,
	)

	hiddenCheckboxView := m.hidden.View()
	if m.hidden.Focused() {
		hiddenCheckboxView = styles.DefaultStyle.
			Foreground(styles.AccentColor).
			Render(hiddenCheckboxView)
	} else {
		hiddenCheckboxView = styles.DefaultStyle.
			Render(hiddenCheckboxView)
	}
	hiddenView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Hidden ",
		hiddenCheckboxView,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wifiName,
		password,
		hiddenView,
	)
}

func (m *WifiConnectorModel) handleKey(key tea.KeyMsg) (*WifiConnectorModel, tea.Cmd) {
	switch {
	case m.password.Focused():
		upd, cmd := m.password.Update(key)
		m.password = upd
		return m, cmd
	case m.hidden.Focused():
		upd, cmd := m.hidden.Update(key)
		m.hidden = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiConnectorModel) handleMsg(msg tea.Msg) (*WifiConnectorModel, tea.Cmd) {
	switch {
	case m.password.Focused():
		upd, cmd := m.password.Update(msg)
		m.password = upd
		return m, cmd
	case m.hidden.Focused():
		upd, cmd := m.hidden.Update(msg)
		m.hidden = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiConnectorModel) focusNextCmd() tea.Cmd {
	if int(m.focusedInputIndex) >= len(m.inputs)-1 {
		return nil
	}
	m.inputs[m.focusedInputIndex].Blur()
	m.focusedInputIndex++
	return m.inputs[m.focusedInputIndex].Focus()
}

func (m *WifiConnectorModel) focusPrevCmd() tea.Cmd {
	if m.focusedInputIndex <= 0 {
		return nil
	}
	m.inputs[m.focusedInputIndex].Blur()
	m.focusedInputIndex--
	return m.inputs[m.focusedInputIndex].Focus()
}

func (m *WifiConnectorModel) connectToWifiCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(ConnectingAvailable),
		func() tea.Msg {
			ssid := m.wifiName
			password := m.password.Value()
			err := m.nm.ConnectWifi(ssid, password, m.hidden.Value())
			if err != nil {
				return tea.BatchMsg{
					SetWifiAvailableStateCmd(DoneInAvailable),
					NotifyCmd(fmt.Sprintf(
						"Cannot connect to %s via given password",
						ssid,
					)),
					RescanWifiCmd(0),
				}
			}
			return tea.BatchMsg{
				SetWifiAvailableStateCmd(DoneInAvailable),
				RescanWifiCmd(0),
			}
		},
	)
}
