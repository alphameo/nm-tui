package components

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type wifiSavedInfoKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	submit             key.Binding
}

func (k wifiSavedInfoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.submit}
}

func (k wifiSavedInfoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.submit}}
}

var wifiSavedInfoKeys = &wifiSavedInfoKeyMap{
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
	submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

type WifiSavedInfoModel struct {
	ssid string

	active bool

	name      textinput.Model
	nameBak   string
	nameStyle *lipgloss.Style

	password textinput.Model
	pwStyle  *lipgloss.Style

	autoconnect   *toggle.Model
	autoconnStyle *lipgloss.Style

	autoconnPriority   textinput.Model
	autoconnPriorStyle *lipgloss.Style

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *wifiSavedInfoKeyMap

	nm infra.WifiManager
}

func NewSavedInfoModel(keys *wifiSavedInfoKeyMap, networkManager infra.WifiManager) *WifiSavedInfoModel {
	name := textinput.New()
	name.SetWidth(20)
	name.Prompt = ""
	name.Placeholder = "name"
	nameStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	pw := textinput.New()
	pw.SetWidth(20)
	pw.Prompt = ""
	pw.EchoMode = textinput.EchoPassword
	pw.EchoCharacter = '•'
	pw.Placeholder = "password"
	pwStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	autoconn := toggle.New(false)
	autoconnStyle := lipgloss.NewStyle().Inherit(styles.DefaultStyle)

	autoconnPrior := textinput.New()
	autoconnPrior.SetWidth(4)
	autoconnPrior.Prompt = ""
	autoconnPrior.Validate = autoconnectPriorityValidator
	autoconnPriorStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)

	model := &WifiSavedInfoModel{
		name:      name,
		nameStyle: &nameStyle,

		password: pw,
		pwStyle:  &pwStyle,

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

func (m *WifiSavedInfoModel) setNew(info infra.WifiInfo) tea.Cmd {
	m.ssid = info.SSID

	m.active = info.Active

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

func (m *WifiSavedInfoModel) Init() tea.Cmd {
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiSavedInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *WifiSavedInfoModel) handleKey(keyMsg tea.KeyPressMsg) (*WifiSavedInfoModel, tea.Cmd) {
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
	case key.Matches(keyMsg, m.keys.submit):
		return m, tea.Sequence(
			SetPopupActivityCmd(false),
			m.saveWifiInfoCmd(),
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

func (m *WifiSavedInfoModel) View() tea.View {
	sb := strings.Builder{}
	fmt.Fprintf(
		&sb,
		"SSID      %s%s",
		m.ssid,
		m.connectionView(),
	)

	ssid := sb.String()

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

	pw := m.password.View()
	pwStyle := *m.pwStyle
	if m.password.Focused() {
		pwStyle = pwStyle.BorderForeground(styles.AccentColor)
	}
	pw = pwStyle.Render(pw)
	pw = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		pw,
	)

	autoconn := m.autoconnect.View().Content
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
		autoconPriorityErrView := renderer.ErrorSymbolColored
		autoconnPrior = lipgloss.JoinHorizontal(
			lipgloss.Center,
			autoconnPrior,
			autoconPriorityErrView,
		)
	}

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		ssid,
		name,
		pw,
		autoconn,
		autoconnPrior,
	))
}

func (m *WifiSavedInfoModel) connectionView() string {
	if m.active {
		return styles.DefaultStyle.
			Foreground(styles.AccentColor).
			Render(" (connected)")
	} else {
		return ""
	}
}

func (m *WifiSavedInfoModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiSavedInfoModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiSavedInfoModel) saveWifiInfoCmd() tea.Cmd {
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
