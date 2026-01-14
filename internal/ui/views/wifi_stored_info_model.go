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

const wifiInfoInputsCount int = 4

type WifiStoredInfoModel struct {
	ssid                string
	name                textinput.Model
	password            textinput.Model
	active              bool
	autoconnect         bool
	autoconnectPriority textinput.Model
	focusIndex          int
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

	ap := textinput.New()
	ap.Width = 4
	ap.Prompt = ""

	return &WifiStoredInfoModel{name: n, password: p, autoconnectPriority: ap}
}

func (m *WifiStoredInfoModel) setNew(info *infra.WifiInfo) {
	m.ssid = info.SSID
	m.name.SetValue(info.Name)
	m.password.SetValue(info.Password)
	m.active = info.Active
	m.autoconnect = info.Autoconnect
	m.autoconnectPriority.SetValue(strconv.Itoa(info.AutoconnectPriority))
}

func (m *WifiStoredInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *WifiStoredInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m, nil
		case "ctrl+j":
			m.focusNext()
			return m, nil
		case "ctrl+k":
			m.focusPrev()
			return m, nil
		}
	}
	return m, nil
}

func (m *WifiStoredInfoModel) View() string {
	inputStyle := lipgloss.
		NewStyle().
		BorderStyle(styles.BorderStyle)

	nameView := m.name.View()
	if m.focusIndex == 0 {
		nameView = inputStyle.BorderForeground(styles.AccentColor).Render(nameView)
	} else {
		nameView = inputStyle.Render(nameView)
	}
	nameView = lipgloss.JoinHorizontal(lipgloss.Center, "Name", nameView)

	passwordView := m.password.View()
	if m.focusIndex == 1 {
		passwordView = inputStyle.BorderForeground(styles.AccentColor).Render(passwordView)
	} else {
		passwordView = inputStyle.Render(passwordView)
	}
	passwordView = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", passwordView)

	autoconnectCheckboxView := checkboxView(m.active)
	if m.focusIndex == 2 {
		autoconnectCheckboxView = lipgloss.NewStyle().Foreground(styles.AccentColor).Render(autoconnectCheckboxView)
	}

	autoconPriorityView := m.autoconnectPriority.View()
	if m.focusIndex == 3 {
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

func (m *WifiStoredInfoModel) focusNext() {
	if m.focusIndex < wifiInfoInputsCount {
		m.focusIndex++
	}
}

func (m *WifiStoredInfoModel) focusPrev() {
	if m.focusIndex > 0 {
		m.focusIndex--
	}
}
