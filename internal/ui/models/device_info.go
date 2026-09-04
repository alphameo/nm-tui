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
	viewport viewport.Model

	devMngr infra.DeviceManager
	Style   lipgloss.Style
}

func NewDeviceInfoModel(deviceManager infra.DeviceManager) *DeviceInfoModel {
	info := newDefaultViewport()
	model := &DeviceInfoModel{
		viewport: info,
		devMngr:  deviceManager,
		Style:    lipgloss.NewStyle(),
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

	m.viewport.SetContent(info)
	_ = m.viewport.GotoTop()
	return nil
}

func (m *DeviceInfoModel) Resize(width, height int) {
	m.Style = m.Style.Width(width).Height(height)

	border := m.Style.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
}

func (m *DeviceInfoModel) Init() tea.Cmd {
	m.viewport.GotoTop()
	return nil
}

func (m *DeviceInfoModel) Update(msg tea.Msg) (*DeviceInfoModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *DeviceInfoModel) View() string {
	view := m.viewport.View()

	view = m.Style.Render(view)
	title := styles.DefaultStyle.Render(renderer.RenderTitle(deviceInfoCfg.title))
	return compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
}
