package models

import (
	"context"
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

type profileEditorKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	save               key.Binding
}

func (k profileEditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.save}
}

func (k profileEditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.save}}
}

func profileEditorKeys() *profileEditorKeyMap {
	return &profileEditorKeyMap{
		togglePWVisibility: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("^r", "toggle password visibility"),
		),
		up: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("^k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("^j", "down"),
		),
		save: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type ProfileEditorModel struct {
	ssid   string
	active bool
	mode   string

	name    textinput.Model
	nameBak string

	password textinput.Model

	autoconnect *toggle.Model

	autoconnPriority textinput.Model

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *profileEditorKeyMap

	nm infra.WifiManager
}

func NewProfileEditorModel(keys *profileEditorKeyMap, networkManager infra.WifiManager) *ProfileEditorModel {
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

	autoconn := toggle.New()
	autoconn.Symbols = styles.ToggleSymbols

	autoconnPrior := textinput.New()
	autoconnPrior.SetWidth(4)
	autoconnPrior.Prompt = ""
	autoconnPrior.Validate = autoconnectPriorityValidator

	model := &ProfileEditorModel{
		ssid:             "",
		active:           false,
		mode:             "",
		name:             name,
		password:         pw,
		autoconnect:      autoconn,
		autoconnPriority: autoconnPrior,

		keys: keys,
		nm:   networkManager,
	}
	inp := []Focusable{
		&model.name,
		&model.password,
		model.autoconnect,
		&model.autoconnPriority,
	}
	model.focuses = inp

	return model
}

func (m *ProfileEditorModel) setNew(name string) tea.Cmd {
	info, err := m.nm.GetWifiInfo(context.Background(), name)
	if err != nil {
		return NotifyCmd(
			fmt.Sprintf("Cannot get information about %s", name),
		)
	}

	m.ssid = info.SSID

	m.active = info.Active

	m.mode = info.Mode.String()

	m.name.Reset()
	m.name.SetValue(info.Name)
	m.name.Blur()
	m.nameBak = info.Name

	m.password.Reset()
	m.password.EchoMode = textinput.EchoPassword
	m.password.SetValue(info.Password)
	m.password.Blur()

	m.autoconnect.SetValue(info.Autoconnect)
	m.autoconnect.Blur()

	m.autoconnPriority.Reset()
	m.autoconnPriority.SetValue(strconv.Itoa(info.AutoconnectPriority))
	m.autoconnPriority.Blur()

	m.focusIdx = 0

	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileEditorModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileEditorModel) Update(msg tea.Msg) (*ProfileEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	default:
		switch {
		case m.name.Focused():
			upd, cmd := m.name.Update(msg)
			m.name = upd
			return m, cmd
		case m.password.Focused():
			upd, cmd := m.password.Update(msg)
			m.password = upd
			return m, cmd
		case m.autoconnect.Focused():
			upd, cmd := m.autoconnect.Update(msg)
			m.autoconnect = upd
			return m, cmd
		case m.autoconnPriority.Focused():
			upd, cmd := m.autoconnPriority.Update(msg)
			m.autoconnPriority = upd
			return m, cmd
		default:
			return m, nil
		}
	}
}

func (m *ProfileEditorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ProfileEditorModel) handleKey(keyMsg tea.KeyPressMsg) (*ProfileEditorModel, tea.Cmd) {
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
	case key.Matches(keyMsg, m.keys.save):
		if m.password.Err != nil {
			return m, nil
		}
		return m, tea.Sequence(
			ClosePopupCmd(),
			m.saveProfileInfoCmd(),
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
	case m.autoconnect.Focused():
		upd, cmd := m.autoconnect.Update(keyMsg)
		m.autoconnect = upd
		return m, cmd
	case m.autoconnPriority.Focused():
		upd, cmd := m.autoconnPriority.Update(keyMsg)
		m.autoconnPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *ProfileEditorModel) View() string {
	ssid := m.ssid
	ssid = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"SSID     ",
		ssid,
		m.connectionView(),
	)

	name := styles.ViewInput(&m.name)
	name = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", name)

	password := styles.ViewInputWithValidation(&m.password)
	password = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", password)

	mode := styles.BoldStyle.Render(m.mode)
	mode = lipgloss.JoinHorizontal(lipgloss.Center, "Mode     ", mode)

	autoconn := styles.ViewToggle(m.autoconnect)
	autoconn = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect          ", autoconn)

	autoconnPrior := styles.ViewInputWithValidation(&m.autoconnPriority)
	autoconnPrior = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect priority ", autoconnPrior)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		ssid,
		mode,
		"",
		name,
		password,
		autoconn,
		autoconnPrior,
	)

	style := styles.OverlayStyle
	view = style.Render(view)
	view = compositor.Compose(
		styles.SavedNetworkInfoTitle,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)

	return view
}

func (m *ProfileEditorModel) connectionView() string {
	if m.active {
		return styles.AccentStyle.Render(" (connected)")
	} else {
		return ""
	}
}

func (m *ProfileEditorModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileEditorModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *ProfileEditorModel) saveProfileInfoCmd() tea.Cmd {
	return func() tea.Msg {
		ap, err := strconv.Atoi(m.autoconnPriority.Value())
		if err != nil {
			return NotifyCmd(
				fmt.Sprintf(
					"Error while updating info about %s: %s",
					m.nameBak,
					err.Error(),
				),
			)
		}
		info := infra.UpdateWifiInfo{
			Name:                m.name.Value(),
			Password:            m.password.Value(),
			Autoconnect:         m.autoconnect.Value(),
			AutoconnectPriority: ap,
		}
		err = m.nm.UpdateWifiInfo(context.Background(), m.nameBak, info)
		if err != nil {
			return NotifyCmd(fmt.Sprintf(
				"Cannot update information about %s",
				m.nameBak,
			))
		}
		return RescanWifiSavedCmd(0)
	}
}

func autoconnectPriorityValidator(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("priority parsing error: %w", err)
	}
	return nil
}
