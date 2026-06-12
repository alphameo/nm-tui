package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type hotspotCreatorKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	create             key.Binding
}

func (k *hotspotCreatorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.create}
}

func (k *hotspotCreatorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.create}}
}

func hotspotCreatorKeys() *hotspotCreatorKeyMap {
	return &hotspotCreatorKeyMap{
		create: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "create"),
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

type HotspotCreatorModel struct {
	title string

	ssid      textinput.Model
	ssidStyle *lipgloss.Style

	name components.Name

	password components.Password

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *hotspotCreatorKeyMap

	nm infra.WifiManager
}

func NewHotspotCreatorModel(keys *hotspotCreatorKeyMap, networkManager infra.WifiManager) *HotspotCreatorModel {
	ssid := textinput.New()
	ssid.SetWidth(20)
	ssid.Prompt = ""
	ssid.Placeholder = "SSID"
	ssidStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	model := &HotspotCreatorModel{
		title: renderer.RenderTitle("Create Wi-Fi hotspot"),

		ssid:      ssid,
		ssidStyle: &ssidStyle,

		name: components.DefaultName(),

		password: components.DefaultPassword(),

		keys: keys,

		nm: networkManager,
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
		return m.handleKey(msg)
	}

	switch {
	case m.ssid.Focused():
		upd, cmd := m.ssid.Update(msg)
		m.ssid = upd
		return m, cmd
	case m.name.Focused():
		upd, cmd := m.name.Update(msg)
		m.name = components.NewName(&upd)
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(msg)
		m.password = components.NewPassword(&upd)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *HotspotCreatorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *HotspotCreatorModel) handleKey(keyMsg tea.KeyPressMsg) (*HotspotCreatorModel, tea.Cmd) {
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
	case key.Matches(keyMsg, m.keys.create):
		if m.password.Err != nil {
			return m, nil
		}
		return m, tea.Sequence(
			ClosePopupCmd(),
			m.createHotspotCmd(),
		)
	}

	switch {
	case m.ssid.Focused():
		upd, cmd := m.ssid.Update(keyMsg)
		m.ssid = upd
		return m, cmd
	case m.name.Focused():
		upd, cmd := m.name.Update(keyMsg)
		m.name = components.NewName(&upd)
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(keyMsg)
		m.password = components.NewPassword(&upd)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *HotspotCreatorModel) View() string {
	ssid := m.ssid.View()
	ssidStyle := *m.ssidStyle
	if m.ssid.Focused() {
		ssidStyle = ssidStyle.BorderForeground(styles.AccentColor)
	}
	ssid = ssidStyle.Render(ssid)
	ssid = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"SSID     ",
		ssid,
	)

	name := m.name.View()
	name = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", name)

	password := m.password.View()
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
		m.title,
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

func (m *HotspotCreatorModel) createHotspotCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(AvailableCreating),
		func() tea.Msg {
			err := m.nm.CreateWifiHotspot(
				context.Background(),
				m.name.Value(),
				m.ssid.Value(),
				m.password.Value(),
			)
			if err != nil {
				return tea.Batch(
					SetWifiAvailableStateCmd(AvailableDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot create hotspot %s:\n%v",
						m.ssid.Value(), err,
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
