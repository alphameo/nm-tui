package components

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

type wifiSavedState int

const (
	SavedNil wifiSavedState = iota
	SavedScanning
	SavedConnecting
	SavedDisconnecting
	SavedDone
)

func (s *wifiSavedState) String() string {
	switch *s {
	case SavedScanning:
		return "Scanning"
	case SavedConnecting:
		return "Connecting"
	case SavedDisconnecting:
		return "Disconnecting"
	case SavedDone:
		return "󰄬"
	default:
		return "Undefined"
	}
}

type wifiSavedKeyMap struct {
	edit       key.Binding
	connect    key.Binding
	disconnect key.Binding
	rescan     key.Binding
	delete     key.Binding
}

func (k *wifiSavedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.edit,
		k.connect,
		k.disconnect,
		k.rescan,
		k.delete,
	}
}

func (k *wifiSavedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.edit,
		k.connect,
		k.disconnect,
		k.rescan,
		k.delete,
	}}
}

var wifiSavedKeys = &wifiSavedKeyMap{
	edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit"),
	),
	connect: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("󱁐", "connect"),
	),
	disconnect: key.NewBinding(
		key.WithKeys("ctrl+space"),
		key.WithHelp("^󱁐", "disconnect"),
	),
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan saved"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

type WifiSavedModel struct {
	dataTable  table.Model
	connColIdx int
	ssidColIdx int
	nameColIdx int
	modeColIdx int

	tableFocusedStyle *table.Styles
	tableBluredStyle  *table.Styles

	ssidWidthProportion  float32
	indicatorStateHeight int

	indicatorSpinner spinner.Model
	indicatorState   wifiSavedState

	focus        bool
	focusedStyle *lipgloss.Style
	bluredStyle  *lipgloss.Style

	savedInfo *WifiInfoModel

	keys *wifiSavedKeyMap

	nm infra.WifiManager

	width  int
	height int
}

func NewWifiSavedModel(savedInfo *WifiInfoModel, keys *wifiSavedKeyMap, networkManager infra.WifiManager) *WifiSavedModel {
	cols := []table.Column{
		{Title: "󱘖", Width: 1},
		{Title: "Mode", Width: 4},
		{Title: "SSID"},
		{Title: "Name"},
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

	model := &WifiSavedModel{
		dataTable: t,

		connColIdx: 0,
		modeColIdx: 1,
		ssidColIdx: 2,
		nameColIdx: 3,

		tableFocusedStyle: &styles.TableStyle,
		tableBluredStyle:  &initTableStyle,

		ssidWidthProportion: 0.5,

		indicatorSpinner: s,
		indicatorState:   SavedDone,

		focusedStyle: &focusedStyle,
		bluredStyle:  &bluredStyle,

		savedInfo: savedInfo,
		keys:      keys,
		nm:        networkManager,
	}
	model.bakeSizes()

	return model
}

func (m *WifiSavedModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *WifiSavedModel) Resize(width, height int) {
	m.width = width
	m.height = height

	width -= styles.BorderOffset
	height -= styles.BorderOffset

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2
	connWidth := m.dataTable.Columns()[m.connColIdx].Width
	modeWidth := m.dataTable.Columns()[m.modeColIdx].Width

	computedWidth := width - tableUtilityOffset - connWidth - modeWidth
	possibleNameWidth := int(float32(computedWidth) * m.ssidWidthProportion)
	ssidWidth := computedWidth - possibleNameWidth
	nameWidth := computedWidth - ssidWidth

	m.dataTable.Columns()[m.nameColIdx].Width = nameWidth
	m.dataTable.Columns()[m.ssidColIdx].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiSavedModel) Width() int {
	return m.width
}

func (m *WifiSavedModel) Height() int {
	return m.height
}

func (m *WifiSavedModel) Focus() tea.Cmd {
	m.focus = true
	m.dataTable.SetStyles(*m.tableFocusedStyle)
	return nil
}

func (m *WifiSavedModel) Blur() {
	m.focus = false
	m.dataTable.SetStyles(*m.tableBluredStyle)
}

func (m *WifiSavedModel) Focused() bool {
	return m.focus
}

func (m *WifiSavedModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiSavedModel) Update(msg tea.Msg) (*WifiSavedModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case RescanWifiSavedMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	case WifiSavedStateMsg:
		return m, m.setStateCmd(wifiSavedState(msg))
	}

	var cmd tea.Cmd
	if m.indicatorState != SavedDone {
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

func (m *WifiSavedModel) handleKey(keyMsg tea.KeyPressMsg) (*WifiSavedModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.edit):
		row := m.dataTable.SelectedRow()
		if row == nil {
			return m, nil
		}
		name := row[m.nameColIdx]
		info, err := m.nm.GetWifiInfo(context.Background(), name)
		if err != nil {
			return m, NotifyCmd(
				fmt.Sprintf("Cannot get information about %s", name),
			)
		}

		return m, m.savedInfo.open(info)

	case key.Matches(keyMsg, m.keys.connect):
		return m, m.connectToSelectedCmd()

	case key.Matches(keyMsg, m.keys.disconnect):
		return m, m.disconnectFromSelectedCmd()
	case key.Matches(keyMsg, m.keys.rescan):
		return m, RescanWifiSavedCmd(0)
	case key.Matches(keyMsg, m.keys.delete):
		return m, m.deleteSelectedCmd()
	}
	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(keyMsg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiSavedModel) View() string {
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
		"Saved networks",
		"2",
		style,
		style.GetBorderTopForeground(),
	)
	return view
}

func (m *WifiSavedModel) indicatorView() string {
	var view string
	if m.indicatorState != SavedDone {
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

type RescanWifiSavedMsg struct {
	delay time.Duration
}

func RescanWifiSavedCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiSavedMsg{delay: delay}
	}
}

func (m *WifiSavedModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(SavedScanning),
		func() tea.Msg {
			list, err := m.nm.GetSavedWifis(context.Background())
			if err != nil {
				return tea.Batch(
					NotifyCmd("Cannot get saved wifi networks"),
					m.setStateCmd(SavedDone),
				)
			}
			rows := []table.Row{}
			for _, wifiSaved := range list {
				var connectionFlag string
				if wifiSaved.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{
					connectionFlag,
					renderer.RenderNetworkModeIcon(wifiSaved.Mode),
					wifiSaved.SSID,
					wifiSaved.Name,
				})
			}

			m.dataTable.SetRows(rows)

			return m.setStateCmd(SavedDone)
		},
	)
}

type WifiSavedStateMsg wifiSavedState

func (m *WifiSavedModel) setStateCmd(state wifiSavedState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == SavedDone {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func (m *WifiSavedModel) connectToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(SavedConnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[m.nameColIdx]
			err := m.nm.ActivateWifi(context.Background(), name)
			if err != nil {
				return tea.Batch(
					m.setStateCmd(SavedDone),
					NotifyCmd(fmt.Sprintf("Cannot connect to %s", name)),
				)
			}
			return tea.Batch(
				m.setStateCmd(SavedDone),
				m.gotoTop(),
				RescanWifiCmd(0),
			)
		},
	)
}

func (m *WifiSavedModel) gotoTop() tea.Cmd {
	return func() tea.Msg {
		m.dataTable.GotoTop()
		return NilCmd
	}
}

func (m *WifiSavedModel) disconnectFromSelectedCmd() tea.Cmd {
	return tea.Sequence(m.setStateCmd(SavedDisconnecting),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[m.nameColIdx]
			err := m.nm.DeactivateWifi(context.Background(), name)
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

func (m *WifiSavedModel) deleteSelectedCmd() tea.Cmd {
	row := m.dataTable.SelectedRow()
	return func() tea.Msg {
		name := row[m.nameColIdx]
		err := m.nm.DeleteWifiConnection(context.Background(), name)
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
