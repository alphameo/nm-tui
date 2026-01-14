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

func (m WifiStoredInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WifiStoredInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m, nil
		}
	}
	return m, nil
}

func (m WifiStoredInfoModel) View() string {
	inputStyle := lipgloss.
		NewStyle().
		BorderStyle(styles.BorderStyle)

	nameBlock := lipgloss.JoinHorizontal(lipgloss.Center, "Name", inputStyle.Render(m.name.View()))

	passwordBlock := lipgloss.JoinHorizontal(lipgloss.Center, "Password ", inputStyle.Render(m.password.View()))
	autoconPriorityView := lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect priority ", inputStyle.Render(m.autoconnectPriority.View()))

	sb := strings.Builder{}
	fmt.Fprintf(
		&sb,
		"SSID: %s%s\n%s\n%s\nAutoconnect %s\n%s",
		m.ssid,
		m.connectionView(),
		nameBlock,
		passwordBlock,
		checkboxView(m.active),
		autoconPriorityView,
	)
	return sb.String()
}

func (m WifiStoredInfoModel) connectionView() string {
	if m.active {
		return " (connected)"
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

// func (m WifiStoredInfoModel) focusNext() {
// 	m.focusIndex = min(m.focusIndex + 1, m.)
// }
