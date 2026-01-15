package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiStoredInfoFocusIndex int

const (
	name wifiStoredInfoFocusIndex = iota
	password
	autoconnect
	autoconnectPriority
)

type switcher bool

func (*switcher) Focus() tea.Cmd {
	return nil
}

func (*switcher) Blur() {
}

type Focusable interface {
	Focus() tea.Cmd
	Blur()
}

type WifiStoredInfoModel struct {
	ssid                string
	active              bool
	name                textinput.Model
	password            textinput.Model
	autoconnect         switcher
	autoconnectPriority textinput.Model
	inputs              []Focusable
	focus               wifiStoredInfoFocusIndex
}

func NewStoredInfoModel() *WifiStoredInfoModel {
	n := textinput.New()
	n.Width = 20
	n.Prompt = ""
	n.Focus()
	n.Placeholder = "name"

	p := textinput.New()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "password"

	a := switcher(false)

	ap := textinput.New()
	ap.Width = 4
	ap.Prompt = ""

	model := &WifiStoredInfoModel{name: n, password: p, autoconnect: a, autoconnectPriority: ap}
	inp := []Focusable{&model.name, &model.password, &model.autoconnect, &model.autoconnectPriority}
	model.inputs = inp

	return model
}

func (m *WifiStoredInfoModel) setNew(info *infra.WifiInfo) {
	m.ssid = info.SSID
	m.name.SetValue(info.Name)
	m.password.SetValue(info.Password)
	m.active = info.Active
	m.autoconnect = switcher(info.Autoconnect)
	m.autoconnectPriority.SetValue(strconv.Itoa(info.AutoconnectPriority))
}

func (m *WifiStoredInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *WifiStoredInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case "ctrl+j":
			return m, m.focusNext()
		case "ctrl+k":
			return m, m.focusPrev()
		default:
			return m.handleKey(msg)
		}
	default:
		return m.handleDefaultMessage(msg)
	}
}

func (m *WifiStoredInfoModel) View() string {
	inputStyle := lipgloss.
		NewStyle().
		BorderStyle(styles.BorderStyle)

	nameView := m.name.View()
	if m.focus == name {
		nameView = inputStyle.BorderForeground(styles.AccentColor).Render(nameView)
	} else {
		nameView = inputStyle.Render(nameView)
	}
	nameView = lipgloss.JoinHorizontal(lipgloss.Center, "Name", nameView)

	passwordView := m.password.View()
	if m.focus == password {
		passwordView = inputStyle.BorderForeground(styles.AccentColor).Render(passwordView)
	} else {
		passwordView = inputStyle.Render(passwordView)
	}
	passwordView = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", passwordView)

	autoconnectCheckboxView := checkboxView(bool(m.autoconnect))
	if m.focus == autoconnect {
		autoconnectCheckboxView = lipgloss.NewStyle().Foreground(styles.AccentColor).Render(autoconnectCheckboxView)
	}

	autoconPriorityView := m.autoconnectPriority.View()
	if m.focus == autoconnectPriority {
		autoconPriorityView = inputStyle.BorderForeground(styles.AccentColor).Render(autoconPriorityView)
	} else {
		autoconPriorityView = inputStyle.Render(autoconPriorityView)
	}
	autoconPriorityView = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect priority ", autoconPriorityView)

	sb := strings.Builder{}
	fmt.Fprintf(
		&sb,
		"SSID: %s%s\n%s\n%s\nAutoconnect %s\n%s",
		m.ssid,
		m.connectionView(),
		nameView,
		passwordView,
		autoconnectCheckboxView,
		autoconPriorityView,
	)
	return sb.String()
}

func (m *WifiStoredInfoModel) handleKey(key tea.KeyMsg) (*WifiStoredInfoModel, tea.Cmd) {
	switch m.focus {
	case name:
		upd, cmd := m.name.Update(key)
		m.name = upd
		return m, cmd
	case password:
		upd, cmd := m.password.Update(key)
		m.password = upd
		return m, cmd
	case autoconnect:
		if key.String() == " " {
			m.autoconnect = !m.autoconnect
		}
		return m, nil
	case autoconnectPriority:
		upd, cmd := m.autoconnectPriority.Update(key)
		m.autoconnectPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) handleDefaultMessage(msg tea.Msg) (*WifiStoredInfoModel, tea.Cmd) {
	switch m.focus {
	case name:
		upd, cmd := m.name.Update(msg)
		m.name = upd
		return m, cmd
	case password:
		upd, cmd := m.password.Update(msg)
		m.password = upd
		return m, cmd
	case autoconnectPriority:
		upd, cmd := m.autoconnectPriority.Update(msg)
		m.autoconnectPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) connectionView() string {
	if m.active {
		return lipgloss.NewStyle().Foreground(styles.AccentColor).Render(" (connected)")
	} else {
		return ""
	}
}

func checkboxView(value bool) string {
	if value {
		return "[⏺]"
	} else {
		return "[ ]"
	}
}

func (m *WifiStoredInfoModel) focusNext() tea.Cmd {
	if int(m.focus) >= len(m.inputs)-1 {
		return nil
	}
	m.inputs[m.focus].Blur()
	m.focus++
	return m.inputs[m.focus].Focus()
}

func (m *WifiStoredInfoModel) focusPrev() tea.Cmd {
	if m.focus <= 0 {
		return nil
	}
	m.inputs[m.focus].Blur()
	m.focus--
	return m.inputs[m.focus].Focus()
}
