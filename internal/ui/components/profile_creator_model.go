package components

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

type profileCreatorKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	create             key.Binding
}

func (k *profileCreatorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.create}
}

func (k *profileCreatorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.create}}
}

var profileCreatorKeys = &profileCreatorKeyMap{
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

type ProfileCreatorModel struct {
	ssid      textinput.Model
	ssidStyle *lipgloss.Style

	name      textinput.Model
	nameStyle *lipgloss.Style

	password textinput.Model
	pwStyle  *lipgloss.Style

	hidden      *toggle.Model
	hiddenStyle *lipgloss.Style

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *profileCreatorKeyMap

	nm infra.WifiManager
}

func NewProfileCreator(keys *profileCreatorKeyMap, networkManager infra.WifiManager) *ProfileCreatorModel {
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

	hiddenStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)

	t := toggle.New(false)

	model := &ProfileCreatorModel{
		ssid:      ssid,
		ssidStyle: &ssidStyle,

		name:      name,
		nameStyle: &nameStyle,

		password: pw,
		pwStyle:  &pwStyle,

		hidden:      t,
		hiddenStyle: &hiddenStyle,

		keys: keys,

		nm: networkManager,
	}

	inp := []Focusable{
		&model.ssid,
		&model.name,
		&model.password,
		model.hidden,
	}
	model.focuses = inp
	model.focusIdx = 0

	return model
}

func (m *ProfileCreatorModel) reset() tea.Cmd {
	m.ssid.Reset()
	m.focusIdx = 0

	m.name.Reset()
	m.name.Blur()

	m.password.Reset()
	m.password.Blur()

	m.hidden.SetValue(false)
	m.hidden.Blur()

	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileCreatorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileCreatorModel) Update(msg tea.Msg) (*ProfileCreatorModel, tea.Cmd) {
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
	case m.hidden.Focused():
		upd, cmd := m.hidden.Update(msg)
		m.hidden = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *ProfileCreatorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ProfileCreatorModel) handleKey(keyMsg tea.KeyPressMsg) (*ProfileCreatorModel, tea.Cmd) {
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
			m.createWifiConnCmd(),
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
	case m.hidden.Focused():
		upd, cmd := m.hidden.Update(keyMsg)
		m.hidden = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *ProfileCreatorModel) View() string {
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

	hidden := m.hidden.View()
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

	style := styles.OverlayStyle
	view = style.Render(view)
	view = compositor.Compose(
		styles.ProfileCreatorTitle,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
	return view
}

func (m *ProfileCreatorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileCreatorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileCreatorModel) createWifiConnCmd() tea.Cmd {
	return tea.Sequence(
		SetWifiAvailableStateCmd(AvailableCreating),
		func() tea.Msg {
			err := m.nm.CreateWifiConnection(
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
					SetWifiAvailableStateCmd(AvailableDone),
					NotifyCmd(fmt.Sprintf(
						"Cannot create connection to %s%s:\n%v",
						hidden, m.ssid.Value(), err,
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

func (m *ProfileCreatorModel) open() tea.Cmd {
	return tea.Batch(
		m.reset(),
		OpenPopup(m),
	)
}
