package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

type connectorKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	connect            key.Binding
}

func (k *connectorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.connect}
}

func (k *connectorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.connect}}
}

func connectorKeys() *connectorKeyMap {
	return &connectorKeyMap{
		connect: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "connect"),
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
}

type ConnectorModel struct {
	ssid string

	name     textinput.Model
	password textinput.Model

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *connectorKeyMap

	nm infra.WifiManager
}

func NewConnectorModel(keys *connectorKeyMap, networkManager infra.WifiManager) *ConnectorModel {
	name := textinput.New()
	name.SetWidth(20)
	name.Prompt = ""
	name.Placeholder = "Name"

	pw := textinput.New()
	pw.SetWidth(20)
	pw.Prompt = ""
	pw.EchoMode = textinput.EchoPassword
	pw.EchoCharacter = styles.PWCharacter
	pw.Placeholder = "Password"
	pw.Validate = passwordValidator
	pw.Err = passwordValidator(pw.Value())

	model := &ConnectorModel{
		ssid: "",

		name: name,

		password: pw,

		keys: keys,

		nm: networkManager,
	}

	inp := []Focusable{
		&model.name,
		&model.password,
	}
	model.focuses = inp

	return model
}

func (m *ConnectorModel) setNew(ssid string) tea.Cmd {
	m.ssid = ssid

	m.name.SetValue(ssid)
	m.focusIdx = 0

	m.password.Reset()
	pw, err := m.nm.GetWifiPassword(context.Background(), ssid)
	if err == nil {
		m.password.SetValue(pw)
	}
	m.password.Blur()

	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectorModel) Update(msg tea.Msg) (*ConnectorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	switch {
	case m.name.Focused():
		upd, cmd := m.name.Update(msg)
		m.name = upd
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(msg)
		m.password = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *ConnectorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ConnectorModel) handleKey(keyMsg tea.KeyPressMsg) (*ConnectorModel, tea.Cmd) {
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
		if m.password.Err == nil {
			return m, nil
		}
		return m, tea.Sequence(
			ClosePopupCmd(),
			m.connectToWifiCmd(),
		)
	}

	switch {
	case m.name.Focused():
		upd, cmd := m.name.Update(keyMsg)
		m.name = upd
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(keyMsg)
		m.password = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *ConnectorModel) View() string {
	ssid := m.ssid
	ssid = lipgloss.JoinHorizontal(lipgloss.Center, "SSID      ", ssid)

	name := styles.ViewBorderedFocusable(&m.name)
	name = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", name)

	password := styles.ViewInputWithValidation(&m.password)
	password = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", password)

	fields := []string{
		ssid,
		name,
		password,
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		fields...,
	)

	style := styles.OverlayStyle
	view = style.Render(view)
	view = compositor.Compose(
		styles.NetworkConnectorTitle,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
	return view
}

func (m *ConnectorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *ConnectorModel) connectToWifiCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(AvailableConnecting),
		func() tea.Msg {
			err := m.nm.ConnectWifi(
				context.Background(),
				m.name.Value(),
				m.ssid,
				m.password.Value(),
			)
			if err != nil {
				return tea.Batch(
					SetWifiAvailableStateCmd(AvailableDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot connect to %s via given password:\n%v",
						m.ssid, err,
					)),
					RescanWifiCmd(0),
				)
			}
			return tea.Batch(
				SetWifiAvailableStateCmd(AvailableDone),
				RescanWifiCmd(0),
			)
		},
	)
}
