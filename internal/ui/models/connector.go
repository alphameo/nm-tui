package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/models/focus"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type connectorConfig struct {
	title   string
	nameIdx int
	pwIdx   int
}

var connectorCfg = connectorConfig{
	title:   "Connect to Network",
	nameIdx: 0,
	pwIdx:   1,
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

	focuses focus.Group

	keys connectorKeyMap

	netMngr infra.NetworksManager
	Style   lipgloss.Style
}

func NewConnectorModel(keys connectorKeyMap, networksManager infra.NetworksManager) *ConnectorModel {
	model := &ConnectorModel{
		ssid:     "",
		name:     newDefaultNameInput(),
		password: newDefaultPasswordInput(),
		keys:     keys,
		netMngr:  networksManager,
		Style:    lipgloss.NewStyle(),
	}

	inp := make([]focus.Focusable, 2)
	inp[connectorCfg.nameIdx] = &model.name
	inp[connectorCfg.pwIdx] = &model.password

	model.focuses = *focus.NewGroup(inp)

	return model
}

func (m *ConnectorModel) setNewNetworkCmd(ssid string) tea.Cmd {
	m.ssid = ssid

	m.name.SetValue(ssid)

	m.password.Reset()
	pw, err := m.netMngr.GetProfilePassword(context.Background(), ssid)
	if err == nil {
		m.password.SetValue(pw)
	}
	m.password.Blur()

	return m.focuses.SetFocusIdx(connectorCfg.pwIdx)
}

func (m *ConnectorModel) Init() tea.Cmd {
	return m.focuses.SetFocusIdx(connectorCfg.pwIdx)
}

func (m *ConnectorModel) Update(msg tea.Msg) (*ConnectorModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(msg, m.keys.next):
			return m, m.focuses.FocusCycleNextCmd()
		case key.Matches(msg, m.keys.prev):
			return m, m.focuses.FocusCyclePrevCmd()
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
				m.connectToNetworkCmd(),
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

	view = m.Style.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(connectorCfg.title))
	return compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
}

func (m *ConnectorModel) connectToNetworkCmd() tea.Cmd {
	return tea.Sequence(
		SetAvailableNetworksStateCmd(NetsConnecting),
		func() tea.Msg {
			err := m.netMngr.ConnectToNetwork(
				context.Background(),
				m.ssid,
				m.password.Value(),
			)
			if err != nil {
				return tea.Batch(
					SetAvailableNetworksStateCmd(NetsDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot connect to %s via given password:\n%v",
						m.ssid, err,
					)),
					QuickRescanNetworksCmd(),
				)
			}
			return tea.Batch(
				SetAvailableNetworksStateCmd(NetsDone),
				QuickRescanNetworksCmd(),
			)
		},
	)
}
