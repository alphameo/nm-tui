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
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type profileCreatorConfig struct {
	title string
}

var profileCreatorCfg = profileCreatorConfig{
	title: "Create Network profile",
}

type profileCreatorKeyMap struct {
	togglePWVisibility key.Binding
	prev               key.Binding
	next               key.Binding
	create             key.Binding
}

type ProfileCreatorModel struct {
	ssid     textinput.Model
	name     textinput.Model
	password textinput.Model
	hidden   toggle.Model

	focuses focus.Group

	keys profileCreatorKeyMap

	netMngr infra.NetworksManager
	style   lipgloss.Style
}

func NewProfileCreatorModel(keys profileCreatorKeyMap, networksManager infra.NetworksManager) *ProfileCreatorModel {
	model := &ProfileCreatorModel{
		ssid:     newDefaultSSIDInput(),
		name:     newDefaultNameInput(),
		password: newDefaultPasswordInput(),
		hidden:   newDefaultToggle(),

		keys: keys,

		netMngr: networksManager,
		style:   lipgloss.NewStyle(),
	}

	inp := []focus.Focusable{
		&model.ssid,
		&model.name,
		&model.password,
		&model.hidden,
	}
	model.focuses = *focus.NewGroup(inp)

	return model
}

func (m *ProfileCreatorModel) Reset() tea.Cmd {
	m.ssid.Reset()

	m.name.Reset()
	m.name.Blur()

	m.password.Reset()
	m.password.Blur()

	m.hidden.SetValue(false)
	m.hidden.Blur()

	return m.focuses.SetFocusIdx(0)
}

func (m *ProfileCreatorModel) Init() tea.Cmd {
	return m.focuses.SetFocusIdx(0)
}

//nolint:dupl // intentionally similar to profile_editor for now; will diverge
func (m *ProfileCreatorModel) Update(msg tea.Msg) (*ProfileCreatorModel, tea.Cmd) {
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
		case key.Matches(msg, m.keys.create):
			if m.password.Err != nil {
				return m, nil
			}
			return m, tea.Sequence(
				ClosePopupCmd(),
				m.createWifiConnCmd(),
			)
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.ssid, cmd = m.ssid.Update(msg)
	cmds = append(cmds, cmd)

	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)

	m.password, cmd = m.password.Update(msg)
	cmds = append(cmds, cmd)

	m.hidden, cmd = m.hidden.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProfileCreatorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ProfileCreatorModel) View() string {
	ssid := styles.ViewBorderedFocusable(&m.ssid)
	ssid = lipgloss.JoinHorizontal(lipgloss.Center, "SSID     ", ssid)

	name := styles.ViewBorderedFocusable(&m.name)
	name = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", name)

	password := styles.ViewBorderedFocusable(&m.password)
	password = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", password)

	hidden := m.hidden.View()
	hidden = lipgloss.JoinHorizontal(lipgloss.Center, "Hidden ", hidden)

	fields := []string{
		ssid,
		name,
		password,
		hidden,
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		fields...,
	)

	view = m.style.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(profileCreatorCfg.title))
	return compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
}

func (m *ProfileCreatorModel) createWifiConnCmd() tea.Cmd {
	return tea.Sequence(
		SetAvailableNetworksStateCmd(AvailableNetsCreating),
		func() tea.Msg {
			err := m.netMngr.CreateConnectionProfile(
				context.Background(),
				m.name.Value(),
				m.ssid.Value(),
				m.password.Value(),
				m.hidden.Value(),
			)
			if err != nil {
				var hidden string
				if m.hidden.Value() {
					hidden = "hidden "
				}
				return tea.Batch(
					SetAvailableNetworksStateCmd(AvailableNetsDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot create connection to %s%s:\n%v",
						hidden, m.ssid.Value(), err,
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
