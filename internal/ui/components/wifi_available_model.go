package components

import (
	"context"
	"fmt"
	"strconv"
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

type wifiAvailableState int

const (
	AvailableNil wifiAvailableState = iota
	AvailableScanning
	AvailableConnecting
	AvailableCreating
	AvailableDone
)

func (s *wifiAvailableState) String() string {
	switch *s {
	case AvailableScanning:
		return "Scanning"
	case AvailableConnecting:
		return "Connecting"
	case AvailableCreating:
		return "Creating Connection"
	case AvailableDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type wifiAvailableKeyMap struct {
	rescan  key.Binding
	connect key.Binding
}

func (k *wifiAvailableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.rescan, k.connect}
}

func (k *wifiAvailableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.rescan, k.connect}}
}

var wifiAvailableKeys = &wifiAvailableKeyMap{
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan"),
	),
	connect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "connect to selected"),
	),
}

type WifiAvailableModel struct {
	dataTable table.Model

	connColIdx     int
	ssidColIdx     int
	securityColIdx int
	signalColIdx   int

	tableFocusedStyle *table.Styles
	tableBluredStyle  *table.Styles

	securityWidthProportion float32

	indicatorStateHeight int

	indicatorSpinner spinner.Model
	indicatorState   wifiAvailableState

	focus        bool
	focusedStyle *lipgloss.Style
	bluredStyle  *lipgloss.Style

	connector *WifiConnectorModel

	keys *wifiAvailableKeyMap

	wm infra.WifiManager

	width  int
	height int
}

func NewWifiAvailableModel(wifiConnector *WifiConnectorModel, keys *wifiAvailableKeyMap, wifiManager infra.WifiManager) *WifiAvailableModel {
	cols := []table.Column{
		{Title: "󱘖", Width: 1},
		{Title: "SSID"},
		{Title: "Security"},
		{Title: "", Width: 3},
	}
	initTableStyle := styles.DataTableStyle
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(initTableStyle),
	)

	bluredStyle := lipgloss.NewStyle().Inherit(styles.BorderedStyle)
	focusedStyle := bluredStyle.BorderForeground(styles.AccentColor)

	s := spinner.New()

	model := &WifiAvailableModel{
		dataTable: t,

		connColIdx:     0,
		ssidColIdx:     1,
		securityColIdx: 2,
		signalColIdx:   3,

		tableFocusedStyle: &styles.TableStyle,
		tableBluredStyle:  &initTableStyle,

		securityWidthProportion: 0.3,

		indicatorSpinner: s,
		indicatorState:   AvailableDone,

		focusedStyle: &focusedStyle,
		bluredStyle:  &bluredStyle,

		connector: wifiConnector,
		keys:      keys,
		wm:        wifiManager,
	}

	model.bakeSizes()

	return model
}

func (m *WifiAvailableModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *WifiAvailableModel) Resize(width, height int) {
	m.width = width
	m.height = height

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	secColWidth := int(float32(width) * m.securityWidthProportion)
	signalColWidth := m.dataTable.Columns()[m.signalColIdx].Width
	conColWidth := m.dataTable.Columns()[m.connColIdx].Width
	ssidWidth := width - signalColWidth - tableUtilityOffset - conColWidth - secColWidth

	m.dataTable.Columns()[m.securityColIdx].Width = secColWidth
	m.dataTable.Columns()[m.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiAvailableModel) Width() int {
	return m.width
}

func (m *WifiAvailableModel) Height() int {
	return m.height
}

func (m *WifiAvailableModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.SetStyles(*m.tableFocusedStyle)
	return nil
}

func (m *WifiAvailableModel) Blur() {
	m.focus = false
	m.dataTable.SetStyles(*m.tableBluredStyle)
}

func (m *WifiAvailableModel) Focused() bool {
	return m.focus
}

func (m *WifiAvailableModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiAvailableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case WifiAvialableStateMsg:
		return m, m.setStateCmd(wifiAvailableState(msg))
	case RescanWifiAvailableMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	}

	var cmd tea.Cmd
	if m.indicatorState != AvailableDone {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}
	m.dataTable, cmd = m.dataTable.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiAvailableModel) handleKey(keyMsg tea.KeyPressMsg) (*WifiAvailableModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.rescan):
		if m.indicatorState != AvailableDone {
			return m, nil
		}
		return m, m.RescanCmd()
	case key.Matches(keyMsg, m.keys.connect):
		row := m.dataTable.SelectedRow()
		if row != nil {
			return m, m.connector.openConnector(row[m.ssidColIdx])
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(keyMsg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiAvailableModel) View() tea.View {
	view := m.dataTable.View()
	statusline := m.indicatorView()
	view = lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)

	var style *lipgloss.Style
	if m.focus {
		style = m.focusedStyle
	} else {
		style = m.bluredStyle
	}
	view = renderer.RenderWithTitleAndKeybind(
		view,
		"Available networks",
		"1",
		style,
		style.GetBorderTopForeground(),
	)
	return tea.NewView(view)
}

func (m *WifiAvailableModel) indicatorView() string {
	var view string
	if m.indicatorState != AvailableDone {
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

func (m *WifiAvailableModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(AvailableScanning),
		func() tea.Msg {
			list, err := m.wm.ScanWifis(context.Background())
			if err != nil {
				return tea.Batch(
					m.setStateCmd(AvailableDone),
					NotifyCmd("Cannot scan available wifi networks"),
				)
			}
			rows := []table.Row{}
			for _, wifiNet := range list {
				var connectionFlag string
				if wifiNet.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{
					connectionFlag,
					wifiNet.SSID,
					wifiNet.Security,
					strconv.Itoa(wifiNet.Signal),
				})
			}

			m.dataTable.SetRows(rows)
			m.dataTable.GotoTop()
			m.dataTable.UpdateViewport()
			return m.setStateCmd(AvailableDone)
		},
	)
}

type RescanWifiAvailableMsg struct {
	delay time.Duration
}

func RescanWifiAvailableCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiAvailableMsg{delay: delay}
	}
}

type WifiAvialableStateMsg wifiAvailableState

func (m *WifiAvailableModel) setStateCmd(state wifiAvailableState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return nil
	}
	if state == AvailableDone {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func SetWifiAvailableStateCmd(state wifiAvailableState) tea.Cmd {
	return func() tea.Msg {
		return WifiAvialableStateMsg(state)
	}
}
