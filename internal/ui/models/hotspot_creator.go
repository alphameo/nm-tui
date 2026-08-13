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
	ssid textinput.Model

	name textinput.Model

	password textinput.Model

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys hotspotCreatorKeyMap

	netMngr infra.NetworksManager
}

func NewHotspotCreatorModel(keys hotspotCreatorKeyMap, networksManager infra.NetworksManager) *HotspotCreatorModel {
	ssid := newDefaultInput()
	ssid.Placeholder = "SSID"

	name := newDefaultInput()
	name.Placeholder = "Name"

	pw := newDefaultPassword()
	pw.Placeholder = "Password"
	pw.Validate = passwordValidator
	pw.Err = passwordValidator(pw.Value())

	model := &HotspotCreatorModel{
		ssid:     ssid,
		name:     name,
		password: pw,
		keys:     keys,
		netMngr:  networksManager,
	}

	inp := []Focusable{
		&model.ssid,
		&model.name,
		&model.password,
	}
	model.focuses = inp

	return model
}

func (m *HotspotCreatorModel) Reset() tea.Cmd {
	m.ssid.Reset()
	m.focusIdx = 0

	m.name.Reset()
	m.name.Blur()

	m.password.Reset()
	m.password.Blur()

	return m.focuses[m.focusIdx].Focus()
}

func (m *HotspotCreatorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *HotspotCreatorModel) Update(msg tea.Msg) (*HotspotCreatorModel, tea.Cmd) {
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

	view = styles.OverlayStyle.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(hotspotCreatorCfg.title))
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

func (m *HotspotCreatorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *HotspotCreatorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
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
