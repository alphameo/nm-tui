package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiStoredInfoInputIndex int

const (
	nameFocus wifiStoredInfoInputIndex = iota
	passwordFocus
	autoconnectFocus
	autoconnectPriorityFocus
)

type Focusable interface {
	Focus() tea.Cmd
	Blur()
}

type WifiStoredInfoModel struct {
	ssid                string
	active              bool
	name                string
	nameInput           textinput.Model
	password            textinput.Model
	autoconnect         *ToggleModel
	autoconnectPriority textinput.Model
	inputs              []Focusable // used for batch operations on input focusable elements
	focusedInputIndex   wifiStoredInfoInputIndex
	nm                  infra.NetworkManager
}

func NewStoredInfoModel(networkManager infra.NetworkManager) *WifiStoredInfoModel {
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

	t := NewToggleModel(false)

	ap := textinput.New()
	ap.Width = 4
	ap.Prompt = ""
	ap.Validate = autoconnectPriorityValidator

	model := &WifiStoredInfoModel{
		nameInput:           n,
		password:            p,
		autoconnect:         t,
		autoconnectPriority: ap,
		nm:                  networkManager,
	}
	inp := []Focusable{&model.nameInput, &model.password, model.autoconnect, &model.autoconnectPriority}
	model.inputs = inp

	return model
}

func (m *WifiStoredInfoModel) setNew(info *infra.WifiInfo) {
	m.ssid = info.SSID
	m.name = info.Name
	m.nameInput.SetValue(info.Name)
	m.password.SetValue(info.Password)
	m.active = info.Active
	m.autoconnect.SetValue(info.Autoconnect)
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
			return m, m.focusNextCmd()
		case "ctrl+k":
			return m, m.focusPrevCmd()
		case "enter":
			return m, tea.Batch(
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

	nameView := m.nameInput.View()
	if m.focusedInputIndex == nameFocus {
		nameView = inputStyle.BorderForeground(styles.AccentColor).Render(nameView)
	} else {
		nameView = inputStyle.Render(nameView)
	}
	nameView = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", nameView)

	passwordView := m.password.View()
	if m.focusedInputIndex == passwordFocus {
		passwordView = inputStyle.BorderForeground(styles.AccentColor).Render(passwordView)
	} else {
		passwordView = inputStyle.Render(passwordView)
	}
	passwordView = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", passwordView)

	autoconnectCheckboxView := m.autoconnect.View()
	if m.focusedInputIndex == autoconnectFocus {
		autoconnectCheckboxView = styles.DefaultStyle.Foreground(styles.AccentColor).Render(autoconnectCheckboxView)
	} else {
		autoconnectCheckboxView = styles.DefaultStyle.Render(autoconnectCheckboxView)
	}
	autoconnectCheckboxView = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect          ", autoconnectCheckboxView)

	autoconPriorityView := m.autoconnectPriority.View()
	if m.focusedInputIndex == autoconnectPriorityFocus {
		autoconPriorityView = inputStyle.BorderForeground(styles.AccentColor).Render(autoconPriorityView)
	} else {
		autoconPriorityView = inputStyle.Render(autoconPriorityView)
	}
	autoconPriorityView = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect priority ", autoconPriorityView)
	if m.autoconnectPriority.Err != nil {
		autoconPriorityErrView := renderer.ErrorSymbolColored
		autoconPriorityView = lipgloss.JoinHorizontal(lipgloss.Center, autoconPriorityView, autoconPriorityErrView)
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
	switch m.focusedInputIndex {
	case nameFocus:
		upd, cmd := m.nameInput.Update(key)
		m.nameInput = upd
		return m, cmd
	case passwordFocus:
		upd, cmd := m.password.Update(key)
		m.password = upd
		return m, cmd
	case autoconnectFocus:
		upd, cmd := m.autoconnect.Update(key)
		m.autoconnect = upd
		return m, cmd
	case autoconnectPriorityFocus:
		upd, cmd := m.autoconnectPriority.Update(key)
		m.autoconnectPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) handleMsg(msg tea.Msg) (*WifiStoredInfoModel, tea.Cmd) {
	switch m.focusedInputIndex {
	case nameFocus:
		upd, cmd := m.nameInput.Update(msg)
		m.nameInput = upd
		return m, cmd
	case passwordFocus:
		upd, cmd := m.password.Update(msg)
		m.password = upd
		return m, cmd
	case autoconnectFocus:
		upd, cmd := m.autoconnect.Update(msg)
		m.autoconnect = upd
		return m, cmd
	case autoconnectPriorityFocus:
		upd, cmd := m.autoconnectPriority.Update(msg)
		m.autoconnectPriority = upd
		return m, cmd
	default:
		return m, nil
	}
}

func (m *WifiStoredInfoModel) connectionView() string {
	if m.active {
		return styles.DefaultStyle.Foreground(styles.AccentColor).Render(" (connected)")
	} else {
		return ""
	}
}

func (m *WifiStoredInfoModel) focusNextCmd() tea.Cmd {
	if int(m.focusedInputIndex) >= len(m.inputs)-1 {
		return nil
	}
	m.inputs[m.focusedInputIndex].Blur()
	m.focusedInputIndex++
	return m.inputs[m.focusedInputIndex].Focus()
}

func (m *WifiStoredInfoModel) focusPrevCmd() tea.Cmd {
	if m.focusedInputIndex <= 0 {
		return nil
	}
	m.inputs[m.focusedInputIndex].Blur()
	m.focusedInputIndex--
	return m.inputs[m.focusedInputIndex].Focus()
}

func (m *WifiStoredInfoModel) saveWifiInfoCmd() tea.Cmd {
	return CmdChain(
		func() tea.Msg {
			ap, err := strconv.Atoi(m.autoconnectPriority.Value())
			if err != nil {
				return NotifyCmd(err.Error())
			}
			info := &infra.UpdateWifiInfo{
				Name:                m.nameInput.Value(),
				Password:            m.password.Value(),
				Autoconnect:         m.autoconnect.Value(),
				AutoconnectPriority: ap,
			}
			err = m.nm.UpdateWifiInfo(m.name, info)
			if err != nil {
				return NotifyCmd(err.Error())
			}
			return UpdateMsg
		},
		UpdateWifiCmd(),
	)
}

func autoconnectPriorityValidator(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("priority parsing error: %w", err)
	}
	return nil
}
