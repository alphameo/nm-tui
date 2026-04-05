package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focusable interface {
	Focused() bool
	Focus() tea.Cmd
	Blur()
}

type wifiStoredInfoKeyMap struct {
	togglePWVisibility key.Binding
	up                 key.Binding
	down               key.Binding
	submit             key.Binding
}

func (k wifiStoredInfoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.togglePWVisibility, k.up, k.down, k.submit}
}

func (k wifiStoredInfoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.togglePWVisibility, k.up, k.down, k.submit}}
}

type WifiStoredInfoModel struct {
	ssid    string
	active  bool
	nameBak string

	name            textinput.Model
	password        textinput.Model
	autoconnect     *toggle.Model
	autoconPriority textinput.Model

	focuses  []Focusable // used for batch operations on input focusable elements
	focusIdx int

	keys *wifiStoredInfoKeyMap

	nm infra.NetworkManager
}

func NewStoredInfoModel(keys *wifiStoredInfoKeyMap, networkManager infra.NetworkManager) *WifiStoredInfoModel {
	n := textinput.New()
	n.Width = 20
	n.Prompt = ""
	n.Placeholder = "name"

	p := textinput.New()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "password"

	t := toggle.New(false)

	ap := textinput.New()
	ap.Width = 4
	ap.Prompt = ""
	ap.Validate = autoconnectPriorityValidator

	model := &WifiStoredInfoModel{
		name:            n,
		password:        p,
		autoconnect:     t,
		autoconPriority: ap,
		keys:            keys,
		nm:              networkManager,
	}
	inp := []Focusable{
		&model.name,
		&model.password,
		model.autoconnect,
		&model.autoconPriority,
	}
	model.focuses = inp

	return model
}

func (m *WifiStoredInfoModel) setNew(info infra.WifiInfo) tea.Cmd {
	m.ssid = info.SSID
	m.nameBak = info.Name
	m.active = info.Active

	m.name.Reset()
	m.name.SetValue(info.Name)
	m.name.Blur()

	m.password.Reset()
	m.password.SetValue(info.Password)
	m.password.Blur()

	m.autoconnect.SetValue(info.Autoconnect)
	m.autoconnect.Blur()

	m.autoconPriority.Reset()
	m.autoconPriority.SetValue(strconv.Itoa(info.AutoconnectPriority))
	m.autoconPriority.Blur()

	m.focusIdx = 0

	return m.focuses[0].Focus()
}

func (m *WifiStoredInfoModel) Init() tea.Cmd {
	return m.focuses[0].Focus()
}

func (m *WifiStoredInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.down):
			return m, m.focusNextCmd()
		case key.Matches(msg, m.keys.up):
			return m, m.focusPrevCmd()
		case key.Matches(msg, m.keys.togglePWVisibility):
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case key.Matches(msg, m.keys.submit):
			return m, tea.Sequence(
				SetPopupActivityCmd(false),
				m.saveWifiInfoCmd(),
			)
		default:
			return m.handleKey(msg)
		}
	default:
		return m.handleMsg(msg)
	}
}

func (m *WifiStoredInfoModel) View() string {
	sb := strings.Builder{}
	fmt.Fprintf(
		&sb,
		"SSID      %s%s",
		m.ssid,
		m.connectionView(),
	)

	ssidView := sb.String()

	inputStyle := styles.BorderedStyle

	nameView := m.name.View()
	if m.name.Focused() {
		nameView = inputStyle.
			BorderForeground(styles.AccentColor).
			Width(m.name.Width + 1). // offset for blinking cursor
			Render(nameView)
	} else {
		nameView = inputStyle.
			Width(m.name.Width + 1). // offset for blinking cursor
			Render(nameView)
	}
	nameView = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Name     ",
		nameView,
	)

	passwordView := m.password.View()
	if m.password.Focused() {
		passwordView = inputStyle.
			Width(m.password.Width + 1). // offset for blinking cursor
			BorderForeground(styles.AccentColor).
			Render(passwordView)
	} else {
		passwordView = inputStyle.
			Width(m.password.Width + 1). // offset for blinking cursor
			Render(passwordView)
	}
	passwordView = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		passwordView,
	)

	autoconnectCheckboxView := m.autoconnect.View()
	if m.autoconnect.Focused() {
		autoconnectCheckboxView = styles.DefaultStyle.
			Foreground(styles.AccentColor).
			Render(autoconnectCheckboxView)
	} else {
		autoconnectCheckboxView = styles.DefaultStyle.
			Render(autoconnectCheckboxView)
	}
	autoconnectCheckboxView = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Autoconnect          ",
		autoconnectCheckboxView,
	)

	autoconPriorityView := m.autoconPriority.View()
	if m.autoconPriority.Focused() {
		autoconPriorityView = inputStyle.
			BorderForeground(styles.AccentColor).
			Render(autoconPriorityView)
	} else {
		autoconPriorityView = inputStyle.Render(autoconPriorityView)
	}
	autoconPriorityView = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Autoconnect priority ",
		autoconPriorityView,
	)
	if m.autoconPriority.Err != nil {
		autoconPriorityErrView := renderer.ErrorSymbolColored
		autoconPriorityView = lipgloss.JoinHorizontal(
			lipgloss.Center,
			autoconPriorityView,
			autoconPriorityErrView,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		ssidView,
		nameView,
		passwordView,
		autoconnectCheckboxView,
		autoconPriorityView,
	)
}

func (m *WifiStoredInfoModel) handleKey(key tea.KeyMsg) (*WifiStoredInfoModel, tea.Cmd) {
	switch {
	case m.name.Focused():
		upd, cmd := m.name.Update(key)
		m.name = upd
		return m, cmd
	case m.password.Focused():
		upd, cmd := m.password.Update(key)
		m.password = upd
		return m, cmd
	case m.autoconnect.Focused():
		upd, cmd := m.autoconnect.Update(key)
		m.autoconnect = upd
		return m, cmd
	case m.autoconPriority.Focused():
		upd, cmd := m.autoconPriority.Update(key)
		m.autoconPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) handleMsg(msg tea.Msg) (*WifiStoredInfoModel, tea.Cmd) {
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
	case m.autoconPriority.Focused():
		upd, cmd := m.autoconPriority.Update(msg)
		m.autoconPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) connectionView() string {
	if m.active {
		return styles.DefaultStyle.
			Foreground(styles.AccentColor).
			Render(" (connected)")
	} else {
		return ""
	}
}

func (m *WifiStoredInfoModel) focusNextCmd() tea.Cmd {
	if int(m.focusIdx) >= len(m.focuses)-1 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx++
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiStoredInfoModel) focusPrevCmd() tea.Cmd {
	if m.focusIdx <= 0 {
		return nil
	}
	m.focuses[m.focusIdx].Blur()
	m.focusIdx--
	return m.focuses[m.focusIdx].Focus()
}

func (m *WifiStoredInfoModel) saveWifiInfoCmd() tea.Cmd {
	return func() tea.Msg {
		ap, err := strconv.Atoi(m.autoconPriority.Value())
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
		err = m.nm.UpdateWifiInfo(m.nameBak, info)
		if err != nil {
			return NotifyCmd(fmt.Sprintf(
				"Cannot update information about %s",
				m.nameBak,
			))
		}
		return RescanWifiStoredCmd(0)
	}
}

func autoconnectPriorityValidator(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("priority parsing error: %w", err)
	}
	return nil
}
