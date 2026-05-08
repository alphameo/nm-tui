package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type wifiConnectorKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	connect            key.Binding
}

func (k *wifiConnectorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.connect}
}

func (k *wifiConnectorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.connect}}
}

var wifiConnectorKeys = &wifiConnectorKeyMap{
	connect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open connector"),
	),
	up: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^k", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("^j", "down"),
	),
	togglePWVisibility: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^r", "toggle password visibility"),
	),
}

type WifiConnectorModel struct {
	name     string
	password textinput.Model
	pwStyle  *lipgloss.Style

	hidden      *toggle.Model
	hiddenStyle *lipgloss.Style

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *wifiConnectorKeyMap

	nm infra.WifiManager
}

func NewWifiConnector(keys *wifiConnectorKeyMap, networkManager infra.WifiManager) *WifiConnectorModel {
	p := textinput.New()
	p.SetWidth(20)
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"

	pwStyle := styles.BorderedStyle.Width(p.Width() + 1) // offset for blinking cursor

	hiddenStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)

	t := toggle.New(false)

	model := &WifiConnectorModel{
		password: p,
		pwStyle:  &pwStyle,

		hidden:      t,
		hiddenStyle: &hiddenStyle,

		keys: keys,

		nm: networkManager,
	}

	inp := []Focusable{
		&model.password,
		model.hidden,
	}
	model.focuses = inp

	return model
}

func (m *WifiConnectorModel) setNew(wifiName string) tea.Cmd {
	m.name = wifiName

	m.password.Reset()
	pw, err := m.nm.GetWifiPassword(context.Background(), wifiName)
	if err == nil {
		m.password.SetValue(pw)
	}
	m.password.Blur()

	m.hidden.SetValue(false)
	m.hidden.Blur()

	return m.focuses[0].Focus()
}

func (m *WifiConnectorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiConnectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case batchCmdMsg:
		// Handle batched commands
		return m, tea.Batch(msg.cmds...)
	}

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

func (m *WifiConnectorModel) handleKey(keyMsg tea.KeyPressMsg) (*WifiConnectorModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.down):
		return m, m.focusNextCmd()
	case key.Matches(keyMsg, m.keys.up):
		return m, m.focusPrevCmd()
	case key.Matches(keyMsg, m.keys.togglePWVisibility):
		if m.password.EchoMode == textinput.EchoPassword {
			m.password.EchoMode = textinput.EchoNormal
		} else {
			m.password.EchoMode = textinput.EchoPassword
		}
		return m, nil
	case key.Matches(keyMsg, m.keys.connect):
		return m, tea.Sequence(
			SetPopupActivityCmd(false),
			m.connectToWifiCmd(),
		)
	}

	switch {
	case m.password.Focused():
		upd, cmd := m.password.Update(keyMsg)
		m.password = upd
		return m, cmd
	case m.hidden.Focused():
		upd, cmd := m.hidden.Update(keyMsg)
		m.hidden = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiConnectorModel) View() tea.View {
	sb := strings.Builder{}
	pwStyle := *m.pwStyle

	fmt.Fprintf(&sb, "SSID      %s", m.name)
	wifiName := sb.String()

	password := m.password.View()
	if m.password.Focused() {
		pwStyle = pwStyle.BorderForeground(styles.AccentColor)
	}
	password = pwStyle.Render(password)
	password = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		password,
	)

	hidden := m.hidden.View().Content
	hiddenStyle := *m.hiddenStyle
	if m.hidden.Focused() {
		hiddenStyle = hiddenStyle.Foreground(styles.AccentColor)
	}
	hidden = hiddenStyle.Render(hidden)
	hidden = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Hidden ",
		hidden,
	)

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		wifiName,
		password,
		hidden,
	))
}

func (m *WifiConnectorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiConnectorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiConnectorModel) connectToWifiCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(AvailableConnecting),
		func() tea.Msg {
			ssid := m.name
			password := m.password.Value()
			err := m.nm.ConnectWifi(context.Background(), ssid, password, m.hidden.Value())
			if err != nil {
				return batchCmdMsg{
					cmds: []tea.Cmd{
						SetWifiAvailableStateCmd(AvailableDone),
						NotifyCmd(fmt.Sprintf(
							"Cannot connect to %s via given password",
							ssid,
						)),
						RescanWifiCmd(0),
					},
				}
			}
			return batchCmdMsg{
				cmds: []tea.Cmd{
					SetWifiAvailableStateCmd(AvailableDone),
					RescanWifiCmd(0),
				},
			}
		},
	)
}
