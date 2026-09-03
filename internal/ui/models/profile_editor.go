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
	"github.com/alphameo/nm-tui/internal/ui/models/focus"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type profileEditorConfig struct {
	title string
}

var profileEditorCfg = profileEditorConfig{
	title: "Saved network info",
}

type profileEditorKeyMap struct {
	togglePWVisibility key.Binding
	prev               key.Binding
	next               key.Binding
	save               key.Binding
}

type ProfileEditorModel struct {
	ssid   string
	uuid   string
	active bool
	hidden bool
	mode   string

	name    textinput.Model
	nameBak string

	password         textinput.Model
	autoconnect      toggle.Model
	autoconnPriority textinput.Model

	focuses focus.Group

	keys profileEditorKeyMap

	netMngr infra.NetworksManager
	Style   lipgloss.Style
}

func NewProfileEditorModel(keys profileEditorKeyMap, networksManager infra.NetworksManager) *ProfileEditorModel {
	autoconnPrior := newDefaultInput()
	autoconnPrior.SetWidth(4)
	autoconnPrior.Validate = autoconnectPriorityValidator

	model := &ProfileEditorModel{
		ssid:             "",
		uuid:             "",
		active:           false,
		hidden:           false,
		mode:             "",
		name:             newDefaultNameInput(),
		password:         newDefaultPasswordInput(),
		autoconnect:      newDefaultToggle(),
		autoconnPriority: autoconnPrior,

		keys:    keys,
		netMngr: networksManager,
		Style:   lipgloss.NewStyle(),
	}
	inp := []focus.Focusable{
		&model.name,
		&model.password,
		&model.autoconnect,
		&model.autoconnPriority,
	}
	model.focuses = *focus.NewGroup(inp)

	return model
}

func (m *ProfileEditorModel) setNewProfile(name string) tea.Cmd {
	info, err := m.netMngr.GetProfile(context.Background(), name)
	if err != nil {
		return NotifyCmd(
			fmt.Sprintf("Cannot get information about %s", name),
		)
	}

	m.ssid = info.SSID

	m.uuid = info.UUID

	m.active = info.Active

	m.hidden = info.Hidden

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

	return m.focuses.SetFocusIdx(0)
}

func (m *ProfileEditorModel) Init() tea.Cmd {
	return m.focuses.SetFocusIdx(0)
}

func (m *ProfileEditorModel) Update(msg tea.Msg) (*ProfileEditorModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(msg, m.keys.next):
			return m, m.focuses.FocusCycleNextCmd()
		case key.Matches(msg, m.keys.prev):
			return m, m.focuses.FocusCyclePrevCmd()
		case key.Matches(msg, m.keys.togglePWVisibility):
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case key.Matches(msg, m.keys.save):
			if m.password.Err != nil {
				return m, nil
			}
			return m, tea.Sequence(
				ClosePopupCmd(),
				m.saveProfileInfoCmd(),
			)
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)

	m.password, cmd = m.password.Update(msg)
	cmds = append(cmds, cmd)

	m.autoconnect, cmd = m.autoconnect.Update(msg)
	cmds = append(cmds, cmd)

	m.autoconnPriority, cmd = m.autoconnPriority.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProfileEditorModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *ProfileEditorModel) View() string {
	ssid := m.ssid
	ssid = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"SSID     ", ssid,
		m.connectionView(),
	)

	uuid := m.uuid
	uuid = lipgloss.JoinHorizontal(lipgloss.Center, "UUID     ", uuid)

	name := styles.ViewBorderedFocusable(&m.name)
	name = lipgloss.JoinHorizontal(lipgloss.Center, "Name     ", name)

	password := styles.ViewInputWithValidation(&m.password)
	password = lipgloss.JoinHorizontal(lipgloss.Center, "Password ", password)

	mode := styles.BoldStyle.Render(m.mode)
	mode = lipgloss.JoinHorizontal(lipgloss.Center, "Mode     ", mode)

	autoconn := m.autoconnect.View()
	autoconn = lipgloss.JoinHorizontal(lipgloss.Center, "Autoconnect  ", autoconn)

	autoconnPrior := styles.ViewInputWithValidation(&m.autoconnPriority)
	autoconnPrior = lipgloss.JoinHorizontal(lipgloss.Center, "  - priority ", autoconnPrior)

	hidden := m.hiddenView()
	hidden = lipgloss.JoinHorizontal(lipgloss.Center, "Hidden   ", hidden)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		ssid,
		uuid,
		hidden,
		mode,
		"",
		name,
		password,
		autoconn,
		autoconnPrior,
	)

	view = m.Style.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(profileEditorCfg.title))
	return compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
}

func (m *ProfileEditorModel) connectionView() string {
	if m.active {
		return styles.AccentStyle.Render(" (connected)")
	}
	return ""
}

func (m *ProfileEditorModel) hiddenView() string {
	if m.hidden {
		return styles.BoldStyle.Render("yes")
	}
	return styles.BoldStyle.Render("no")
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
		info := infra.UpdateNetworkProfile{
			Name:                m.name.Value(),
			Password:            m.password.Value(),
			Autoconnect:         m.autoconnect.Value(),
			AutoconnectPriority: ap,
		}
		err = m.netMngr.UpdateProfile(context.Background(), m.uuid, info)
		if err != nil {
			return NotifyCmd(fmt.Sprintf(
				"Cannot update information about %s",
				m.nameBak,
			))
		}
		return QuickRescanNetworksCmd()
	}
}

func autoconnectPriorityValidator(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("priority parsing error: %w", err)
	}
	return nil
}
