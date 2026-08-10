package models

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type networkProfilesState int

const (
	NetProfilesNil networkProfilesState = iota
	NetProfilesScanning
	NetProfilesConnecting
	NetProfilesDisconnecting
	NetProfilesDone
)

func (s *networkProfilesState) String() string {
	switch *s {
	case NetProfilesScanning:
		return "Scanning"
	case NetProfilesConnecting:
		return "Connecting"
	case NetProfilesDisconnecting:
		return "Disconnecting"
	case NetProfilesDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type networkProfilesKeyMap struct {
	edit       key.Binding
	connect    key.Binding
	disconnect key.Binding
	rescan     key.Binding
	delete     key.Binding
}

type NetworkProfilesModel struct {
	dataTable table.Model

	indicatorSpinner     spinner.Model
	indicatorState       networkProfilesState
	indicatorStateHeight int

	focus bool

	keys networkProfilesKeyMap

	nm infra.WifiManager

	width  int
	height int
}

type networkProfilesConfig struct {
	connColIdx int
	ssidColIdx int
	nameColIdx int
	modeColIdx int

	ssidWidthProportion float32
}

var networkProfilesCfg = networkProfilesConfig{
	connColIdx:          0,
	modeColIdx:          1,
	ssidColIdx:          2,
	nameColIdx:          3,
	ssidWidthProportion: 0.5,
}

func NewNetworkProfilesModel(keys networkProfilesKeyMap, networkManager infra.WifiManager) *NetworkProfilesModel {
	cols := make([]table.Column, 4)
	cols[networkProfilesCfg.connColIdx] = table.Column{Title: styles.SymbolConnection, Width: len(styles.SymbolConnection)}
	cols[networkProfilesCfg.modeColIdx] = table.Column{Title: "Mode", Width: 4}
	cols[networkProfilesCfg.ssidColIdx] = table.Column{Title: "SSID"}
	cols[networkProfilesCfg.nameColIdx] = table.Column{Title: "Name"}

	initTableStyle := styles.DataTableStyle
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(initTableStyle),
	)

	s := spinner.New()
	s.Spinner = styles.Spinner

	model := &NetworkProfilesModel{
		dataTable: t,

		indicatorSpinner: s,
		indicatorState:   NetProfilesDone,

		keys: keys,
		nm:   networkManager,
	}
	model.bakeSizes()

	return model
}

func (m *NetworkProfilesModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *NetworkProfilesModel) Resize(width, height int) {
	m.width = width
	m.height = height

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2
	connWidth := m.dataTable.Columns()[networkProfilesCfg.connColIdx].Width
	modeWidth := m.dataTable.Columns()[networkProfilesCfg.modeColIdx].Width

	computedWidth := width - tableUtilityOffset - connWidth - modeWidth
	possibleNameWidth := int(float32(computedWidth) * networkProfilesCfg.ssidWidthProportion)
	ssidWidth := computedWidth - possibleNameWidth
	nameWidth := computedWidth - ssidWidth

	m.dataTable.Columns()[networkProfilesCfg.nameColIdx].Width = nameWidth
	m.dataTable.Columns()[networkProfilesCfg.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *NetworkProfilesModel) Width() int {
	return m.width
}

func (m *NetworkProfilesModel) Height() int {
	return m.height
}

func (m *NetworkProfilesModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.Focus()
	m.dataTable.SetStyles(styles.TableStyle)
	return nil
}

func (m *NetworkProfilesModel) Blur() {
	m.focus = false
	m.dataTable.Blur()
	m.dataTable.SetStyles(styles.DataTableStyle)
}

func (m *NetworkProfilesModel) Focused() bool {
	return m.focus
}

func (m *NetworkProfilesModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *NetworkProfilesModel) Update(msg tea.Msg) (*NetworkProfilesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.edit):
			row := m.dataTable.SelectedRow()
			if row == nil {
				return m, nil
			}
			name := row[networkProfilesCfg.nameColIdx]

			return m, OpenProfileEditorCmd(name)

		case key.Matches(msg, m.keys.connect):
			return m, m.connectToSelectedCmd()

		case key.Matches(msg, m.keys.disconnect):
			return m, m.disconnectFromSelectedCmd()
		case key.Matches(msg, m.keys.rescan):
			return m, RescanWifiSavedCmd(0)
		case key.Matches(msg, m.keys.delete):
			return m, m.deleteSelectedCmd()
		}
	case RescanNetworkProfilesMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	case WifiSavedStateMsg:
		return m, m.setStateCmd(networkProfilesState(msg))
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.indicatorState != NetProfilesDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NetworkProfilesModel) View() string {
	view := m.dataTable.View()
	statusline := m.indicatorView()
	view = lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)

	var style lipgloss.Style
	if m.focus {
		style = styles.BorderedFocusedStyle
	} else {
		style = styles.BorderedStyle
	}
	view = renderer.RenderWithTitleAndKeybind(
		view,
		"Network profiles",
		"2",
		style,
		styles.AccentColor,
	)
	return view
}

func (m *NetworkProfilesModel) indicatorView() string {
	var view string
	if m.indicatorState != NetProfilesDone {
		view = fmt.Sprintf(
			"%s %s",
			m.indicatorState.String(),
			m.indicatorSpinner.View(),
		)
	} else {
		view = m.indicatorState.String()
	}
	return view
}

type RescanNetworkProfilesMsg struct {
	delay time.Duration
}

func RescanWifiSavedCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanNetworkProfilesMsg{delay: delay}
	}
}

func (m *NetworkProfilesModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetProfilesScanning),
		func() tea.Msg {
			list, err := m.nm.ListProfiles(context.Background())
			if err != nil {
				return tea.Batch(
					NotifyCmd("Cannot get saved wifi networks"),
					m.setStateCmd(NetProfilesDone),
				)
			}
			rows := []table.Row{}
			for _, wifiSaved := range list {
				var connectionFlag string
				if wifiSaved.Active {
					connectionFlag = styles.SymbolCheck
				}
				rows = append(rows, table.Row{
					connectionFlag,
					ViewNetworkMode(wifiSaved.Mode),
					wifiSaved.SSID,
					wifiSaved.Name,
				})
			}

			m.dataTable.SetRows(rows)

			return m.setStateCmd(NetProfilesDone)
		},
	)
}

func ViewNetworkMode(mode infra.NetworkMode) string {
	switch mode {
	case infra.NetworkAccessPoint:
		return styles.SymbolAccessPoint
	case infra.NetworkInfra:
		return styles.SymbolInfra
	case infra.NetworkMesh:
		return styles.SymbolMesh
	case infra.NetworkAdHoc:
		return styles.SymbolAdHoc
	default:
		return "?"
	}
}

type WifiSavedStateMsg networkProfilesState

func (m *NetworkProfilesModel) setStateCmd(state networkProfilesState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == NetProfilesDone {
		return updCmd
	}
	return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
}

func (m *NetworkProfilesModel) connectToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(NetProfilesConnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.nm.ActivateProfile(context.Background(), name)
			if err != nil {
				return tea.Batch(
					m.setStateCmd(NetProfilesDone),
					NotifyCmd(fmt.Sprintf("Cannot connect to %s", name)),
				)
			}
			return tea.Batch(
				m.setStateCmd(NetProfilesDone),
				m.gotoTop(),
				RescanWifiCmd(0),
			)
		},
	)
}

func (m *NetworkProfilesModel) gotoTop() tea.Cmd {
	return func() tea.Msg {
		m.dataTable.GotoTop()
		return NilCmd
	}
}

func (m *NetworkProfilesModel) disconnectFromSelectedCmd() tea.Cmd {
	return tea.Sequence(m.setStateCmd(NetProfilesDisconnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[networkProfilesCfg.nameColIdx]
			err := m.nm.DeactivateProfile(context.Background(), name)
			if err != nil {
				return NotifyCmd(
					fmt.Sprintf("Error while disconnecting from %s", name),
				)
			}
			return tea.Batch(
				m.gotoTop(),
				RescanWifiCmd(200*time.Millisecond),
			)
		})
}

func (m *NetworkProfilesModel) deleteSelectedCmd() tea.Cmd {
	row := m.dataTable.SelectedRow()
	return func() tea.Msg {
		name := row[networkProfilesCfg.nameColIdx]
		err := m.nm.DeleteProfile(context.Background(), name)
		if err != nil {
			return NotifyCmd(fmt.Sprintf("Error while deleting %s", name))
		}
		cursor := m.dataTable.Cursor()
		if cursor == len(m.dataTable.Rows())-1 {
			m.dataTable.SetCursor(cursor - 1)
		}
		return RescanWifiCmd(time.Millisecond * 200)
	}
}
