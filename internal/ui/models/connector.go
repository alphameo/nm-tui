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
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type connectorConfig struct {
	title string
}

var connectorCfg = connectorConfig{
	title: "Connect to Network",
}

type connectorKeyMap struct {
	togglePWVisibility key.Binding
	prev               key.Binding
	next               key.Binding
	connect            key.Binding
}

type ConnectorModel struct {
	ssid string

	name     textinput.Model
	password textinput.Model

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys connectorKeyMap

	netMngr infra.NetworksManager
}

func NewConnectorModel(keys connectorKeyMap, networksManager infra.NetworksManager) *ConnectorModel {
	name := newDefaultInput()
	name.Placeholder = "Name"

	pw := newDefaultPassword()
	pw.Placeholder = "Password"
	pw.Validate = passwordValidator
	pw.Err = passwordValidator(pw.Value())

	model := &ConnectorModel{
		ssid:     "",
		name:     name,
		password: pw,
		keys:     keys,
		netMngr:  networksManager,
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
	pw, err := m.netMngr.GetProfilePassword(context.Background(), ssid)
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
		switch {
		case key.Matches(msg, m.keys.next):
			return m, m.focusNextCmd()
		case key.Matches(msg, m.keys.prev):
			return m, m.focusPrevCmd()
		case key.Matches(msg, m.keys.togglePWVisibility):
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case key.Matches(msg, m.keys.connect):
			if m.password.Err != nil {
				return m, nil
			}
			return m, tea.Sequence(
				ClosePopupCmd(),
				m.connectToWifiCmd(),
			)
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)

	m.password, cmd = m.password.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ConnectorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
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

	view = styles.OverlayStyle.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(connectorCfg.title))
	view = compositor.Compose(
		title,
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
		SetAvailableNetworksStateCmd(AvailableNetsConnecting),
		func() tea.Msg {
			err := m.netMngr.ConnectToNetwork(
				context.Background(),
				m.ssid,
				m.password.Value(),
			)
			if err != nil {
				return tea.Batch(
					SetAvailableNetworksStateCmd(AvailableNetsDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot connect to %s via given password:\n%v",
						m.ssid, err,
					)),
					RescanNetworksCmd(0),
				)
			}
			return tea.Batch(
				SetAvailableNetworksStateCmd(AvailableNetsDone),
				RescanNetworksCmd(0),
			)
		},
	)
}
