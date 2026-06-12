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
	"github.com/alphameo/nm-tui/internal/ui/components"
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
	ssid string

	active      bool
	activeStyle *lipgloss.Style

	mode      string
	modeStyle *lipgloss.Style

	name    components.Name
	nameBak string

	password components.Password

	autoconnect   *toggle.Model
	autoconnStyle *lipgloss.Style

	autoconnPriority   textinput.Model
	autoconnPriorStyle *lipgloss.Style
	focuses            []Focusable // used for batch operations on input focusable elements
	focusIdx           int

	keys *profileEditorKeyMap

	nm infra.WifiManager
}

func NewProfileEditorModel(keys *profileEditorKeyMap, networkManager infra.WifiManager) *ProfileEditorModel {
	activeStyle := styles.DefaultStyle.Foreground(styles.AccentColor)

	modeStyle := styles.DefaultStyle.Bold(true)

	autoconn := components.DefaultToggle()
	autoconnStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)

	autoconnPrior := textinput.New()
	autoconnPrior.SetWidth(4)
	autoconnPrior.Prompt = ""
	autoconnPrior.Validate = autoconnectPriorityValidator
	autoconnPriorStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	model := &ProfileEditorModel{
		ssid: "",

		active:      false,
		activeStyle: &activeStyle,

		mode:      "",
		modeStyle: &modeStyle,

		name: components.DefaultName(),

		password: components.DefaultPassword(),

		autoconnect:   autoconn,
		autoconnStyle: &autoconnStyle,

		autoconnPriority:   autoconnPrior,
		autoconnPriorStyle: &autoconnPriorStyle,

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
			m.name = components.NewName(&upd)
			return m, cmd
		case m.password.Focused():
			upd, cmd := m.password.Update(msg)
			m.password = components.NewPassword(&upd)
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
		m.name = components.NewName(&upd)
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(keyMsg)
		m.password = components.NewPassword(&upd)
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

	mode := m.modeStyle.Render(m.mode)
	mode = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Mode     ",
		mode,
	)

	autoconn := m.autoconnect.View()
	autoconnStyle := *m.autoconnStyle
	if m.autoconnect.Focused() {
		autoconnStyle = autoconnStyle.
			Foreground(styles.AccentColor)
	}
	autoconn = autoconnStyle.Render(autoconn)
	autoconn = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Autoconnect          ",
		autoconn,
	)

	autoconnPrior := m.autoconnPriority.View()
	autoconnPriorStyle := *m.autoconnPriorStyle
	if m.autoconnPriority.Focused() {
		autoconnPriorStyle = autoconnPriorStyle.BorderForeground(styles.AccentColor)
	}
	autoconnPrior = autoconnPriorStyle.Render(autoconnPrior)
	autoconnPrior = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Autoconnect priority ",
		autoconnPrior,
	)
	if m.autoconnPriority.Err != nil {
		autoconPriorityErrView := styles.ErrorSymbolColored
		autoconnPrior = lipgloss.JoinHorizontal(
			lipgloss.Center,
			autoconnPrior,
			autoconPriorityErrView,
		)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		ssid,
		mode,
		m.name.View(),
		m.password.View(),
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
		return m.activeStyle.Render(" (connected)")
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
