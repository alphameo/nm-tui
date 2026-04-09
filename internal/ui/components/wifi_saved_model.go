package components

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiSavedState int

const (
	ScanningSaved wifiSavedState = iota
	ConnectingSaved
	DisconnectingSaved
	DoneInSaved
)

func (s *wifiSavedState) String() string {
	switch *s {
	case ScanningSaved:
		return "Scanning"
	case ConnectingSaved:
		return "Connecting"
	case DisconnectingSaved:
		return "Disconnecting"
	case DoneInSaved:
		return "󰄬"
	default:
		return "Undefined!!!"
	}
}

type wifiSavedKeyMap struct {
	edit       key.Binding
	connect    key.Binding
	disconnect key.Binding
	update     key.Binding
	delete     key.Binding
}

func (k *wifiSavedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.edit,
		k.connect,
		k.disconnect,
		k.update,
		k.delete,
	}
}

func (k *wifiSavedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.edit,
		k.connect,
		k.disconnect,
		k.update,
		k.delete,
	}}
}

type WifiSavedModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   wifiSavedState

	savedInfo *WifiSavedInfoModel

	connColIdx int
	ssidColIdx int
	nameColIdx int

	minSSIDWidth         int
	indicatorStateHeight int

	keys *wifiSavedKeyMap

	nm infra.WifiManager

	width  int
	height int
}

func NewWifiSavedModel(savedInfo *WifiSavedInfoModel, keys *wifiSavedKeyMap, networkManager infra.WifiManager) *WifiSavedModel {
	cols := []table.Column{
		{Title: "󱘖", Width: 1},
		{Title: "SSID"},
		{Title: "Name"},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)

	s := spinner.New()

	model := &WifiSavedModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   DoneInSaved,
		savedInfo:        savedInfo,
		keys:             keys,
		nm:               networkManager,

		connColIdx: 0,
		ssidColIdx: 1,
		nameColIdx: 2,

		minSSIDWidth: 4,
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

	height -= m.indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2
	conColWidth := m.dataTable.Columns()[m.connColIdx].Width

	computedWidth := width - tableUtilityOffset - conColWidth
	possibleNameWidth := computedWidth / 2
	ssidWidth := max(computedWidth-possibleNameWidth, m.minSSIDWidth)
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

func (m *WifiSavedModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiSavedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			row := m.dataTable.SelectedRow()
			if row == nil {
				return m, nil
			}
			name := row[m.nameColIdx]
			info, err := m.nm.GetWifiInfo(name)
			if err != nil {
				return m, NotifyCmd(
					fmt.Sprintf("Cannot get information about %s", name),
				)
			}

			return m, tea.Batch(
				m.savedInfo.setNew(info),
				OpenPopup(m.savedInfo, "Saved Wi-Fi network info"),
			)

		case " ":
			return m, m.connectToSelectedCmd()

		case "shift+ ":
			return m,
				m.disconnectFromSelectedCmd()
		case "r":
			return m, RescanWifiSavedCmd(0)
		case "d":
			return m, m.deleteSelectedCmd()
		}
	case RescanWifiSavedMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	case WifiSavedStateMsg:
		return m, m.setStateCmd(wifiSavedState(msg))
	}

	var cmd tea.Cmd
	if m.indicatorState != DoneInSaved {
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

func (m *WifiSavedModel) View() string {
	view := m.dataTable.View()

	statusline := m.indicatorView()
	return lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)
}

func (m *WifiSavedModel) indicatorView() string {
	var view string
	if m.indicatorState != DoneInSaved {
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
		m.setStateCmd(ScanningSaved),
		func() tea.Msg {
			list, err := m.nm.GetSavedWifis()
			if err != nil {
				return tea.BatchMsg{
					NotifyCmd("Cannot get saved wifi networks"),
					m.setStateCmd(DoneInSaved),
				}
			}
			rows := []table.Row{}
			for _, wifiSaved := range list {
				var connectionFlag string
				if wifiSaved.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{
					connectionFlag,
					wifiSaved.SSID,
					wifiSaved.Name,
				})
			}

			m.dataTable.SetRows(rows)

			return m.setStateCmd(DoneInSaved)
		},
	)
}

type WifiSavedStateMsg wifiSavedState

func (m *WifiSavedModel) setStateCmd(state wifiSavedState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == DoneInSaved {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func (m *WifiSavedModel) connectToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ConnectingSaved),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[m.nameColIdx]
			err := m.nm.ActivateWifi(name)
			if err != nil {
				return tea.BatchMsg{
					m.setStateCmd(DoneInSaved),
					NotifyCmd(fmt.Sprintf("Cannot connect to %s", name)),
				}
			}
			return tea.BatchMsg{
				m.setStateCmd(DoneInSaved),
				m.gotoTop(),
				RescanWifiCmd(0),
			}
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
	return func() tea.Msg {
		name := m.dataTable.SelectedRow()[m.nameColIdx]
		err := m.nm.DeactivateWifi(name)
		if err != nil {
			return NotifyCmd(
				fmt.Sprintf("Error while disconnecting from %s", name),
			)
		}
		return tea.BatchMsg{
			m.gotoTop(),
			RescanWifiCmd(200 * time.Millisecond),
		}
	}
}

func (m *WifiSavedModel) deleteSelectedCmd() tea.Cmd {
	row := m.dataTable.SelectedRow()
	return func() tea.Msg {
		name := row[m.nameColIdx]
		err := m.nm.DeleteWifiConnection(name)
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
