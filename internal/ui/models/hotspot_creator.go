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

type hotspotCreatorConfig struct {
	title string
}

var hotspotCreatorCfg = hotspotCreatorConfig{
	title: "Create Hotspot",
}

type hotspotCreatorKeyMap struct {
	togglePWVisibility key.Binding
	prev               key.Binding
	next               key.Binding
	create             key.Binding
}

type HotspotCreatorModel struct {
	ssid     textinput.Model
	name     textinput.Model
	password textinput.Model

	focuses focus.Group

	keys hotspotCreatorKeyMap

	netMngr infra.NetworksManager
	Style   lipgloss.Style
}

func NewHotspotCreatorModel(keys hotspotCreatorKeyMap, networksManager infra.NetworksManager) *HotspotCreatorModel {
	model := &HotspotCreatorModel{
		ssid:     newDefaultSSIDInput(),
		name:     newDefaultNameInput(),
		password: newDefaultPasswordInput(),
		keys:     keys,
		netMngr:  networksManager,
		Style:    lipgloss.NewStyle(),
	}

	inp := []focus.Focusable{
		&model.ssid,
		&model.name,
		&model.password,
	}
	model.focuses = *focus.NewGroup(inp)

	return model
}

func (m *HotspotCreatorModel) Reset() tea.Cmd {
	m.ssid.Reset()

	m.name.Reset()
	m.name.Blur()

	m.password.Reset()
	m.password.Blur()

	return m.focuses.SetFocusIdx(0)
}

func (m *HotspotCreatorModel) Init() tea.Cmd {
	return m.focuses.SetFocusIdx(0)
}

func (m *HotspotCreatorModel) Update(msg tea.Msg) (*HotspotCreatorModel, tea.Cmd) {
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
				m.createHotspotProfileCmd(),
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

	return m, tea.Batch(cmds...)
}

func (m *HotspotCreatorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *HotspotCreatorModel) View() string {
	ssid := styles.ViewBorderedFocusable(&m.ssid)
	ssid = lipgloss.JoinHorizontal(lipgloss.Center, "SSID     ", ssid)

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
	title := styles.DefaultStyle.Render(renderer.RenderTitle(hotspotCreatorCfg.title))
	return compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
}

func (m *HotspotCreatorModel) createHotspotProfileCmd() tea.Cmd {
	return tea.Sequence(
		SetAvailableNetworksStateCmd(AvailableNetsCreating),
		func() tea.Msg {
			err := m.netMngr.CreateHotspotProfile(
				context.Background(),
				m.name.Value(),
				m.ssid.Value(),
				m.password.Value(),
			)
			if err != nil {
				return tea.Batch(
					SetAvailableNetworksStateCmd(AvailableNetsDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot create hotspot %s:\n%v",
						m.ssid.Value(), err,
					)),
					RescanNetworksCmd(),
				)
			}
			return tea.Batch(
				SetAvailableNetworksStateCmd(AvailableNetsDone),
				RescanNetworksCmd(),
			)
		},
	)
}
