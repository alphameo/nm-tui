package components

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

var hotspotCreatorKeys = &hotspotCreatorKeyMap{
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

type hotspotCreatorModel struct {
	title string

	ssid      textinput.Model
	ssidStyle *lipgloss.Style

	name      textinput.Model
	nameStyle *lipgloss.Style

	password textinput.Model
	pwStyle  *lipgloss.Style

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *hotspotCreatorKeyMap

	nm infra.WifiManager
}

func NewHotspotCreator(keys *hotspotCreatorKeyMap, networkManager infra.WifiManager) *hotspotCreatorModel {
	ssid := textinput.New()
	ssid.SetWidth(20)
	ssid.Prompt = ""
	ssid.Placeholder = "SSID"
	ssidStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	name := textinput.New()
	name.SetWidth(20)
	name.Prompt = ""
	name.Placeholder = "Name"
	nameStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	pw := textinput.New()
	pw.SetWidth(20)
	pw.Prompt = ""
	pw.EchoMode = textinput.EchoPassword
	pw.EchoCharacter = '•'
	pw.Placeholder = "Password"
	pwStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	model := &hotspotCreatorModel{
		title: renderer.RenderTitle("Create Wi-Fi hotspot"),

		ssid:      ssid,
		ssidStyle: &ssidStyle,

		name:      name,
		nameStyle: &nameStyle,

		password: pw,
		pwStyle:  &pwStyle,

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

func (m *hotspotCreatorModel) reset() tea.Cmd {
	m.ssid.Reset()
	m.focusIdx = 0

	m.name.Reset()
	m.name.Blur()

	m.password.Reset()
	m.password.Blur()

	return m.focuses[m.focusIdx].Focus()
}

func (m *hotspotCreatorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *hotspotCreatorModel) Update(msg tea.Msg) (*hotspotCreatorModel, tea.Cmd) {
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

func (m *hotspotCreatorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *hotspotCreatorModel) handleKey(keyMsg tea.KeyPressMsg) (*hotspotCreatorModel, tea.Cmd) {
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
		return m, tea.Sequence(
			SetPopupActivityCmd(false),
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

func (m *hotspotCreatorModel) View() string {
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
	nameStyle := *m.nameStyle
	if m.name.Focused() {
		nameStyle = nameStyle.BorderForeground(styles.AccentColor)
	}
	name = nameStyle.Render(name)
	name = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Name     ",
		name,
	)
	password := m.password.View()
	pwStyle := *m.pwStyle
	if m.password.Focused() {
		pwStyle = pwStyle.BorderForeground(styles.AccentColor)
	}
	password = pwStyle.Render(password)
	password = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		password,
	)

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

func (m *hotspotCreatorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *hotspotCreatorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *hotspotCreatorModel) createHotspotCmd() tea.Cmd {
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

func (m *hotspotCreatorModel) open() tea.Cmd {
	return tea.Batch(
		m.reset(),
		OpenPopup(m),
	)
}
