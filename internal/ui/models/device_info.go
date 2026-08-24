package models

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type deviceInfoConfig struct {
	title string
}

var deviceInfoCfg = deviceInfoConfig{
	title: "Saved network info",
}

type DeviceInfoModel struct {
	info viewport.Model

	devMngr infra.DeviceManager
	Style   lipgloss.Style
}

func NewDeviceInfoModel(deviceManager infra.DeviceManager) *DeviceInfoModel {
	info := viewport.New()
	model := &DeviceInfoModel{
		info:    info,
		devMngr: deviceManager,
		Style:   lipgloss.NewStyle(),
	}

	return model
}

func (m *DeviceInfoModel) setNewDevice(name string) tea.Cmd {
	info, err := m.devMngr.GetDeviceInfo(context.Background(), name)
	if err != nil {
		return NotifyCmd(
			fmt.Sprintf("Cannot get information about device %q", name),
		)
	}

	m.info.SetContent(info)
	return nil
}

func (m *DeviceInfoModel) Init() tea.Cmd {
	return nil
}

func (m *DeviceInfoModel) Update(msg tea.Msg) (*DeviceInfoModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.info, cmd = m.info.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *DeviceInfoModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *DeviceInfoModel) View() string {
	view := m.info.View()

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
