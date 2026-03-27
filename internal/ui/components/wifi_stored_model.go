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

type wifiStoredState int

const (
	ScanningStored wifiStoredState = iota
	ConnectingStored
	DisconnectingStored
	DoneInStored
)

func (s *wifiStoredState) String() string {
	switch *s {
	case ScanningStored:
		return "Scanning"
	case ConnectingStored:
		return "Connecting"
	case DisconnectingStored:
		return "Disconnecting"
	case DoneInStored:
		return "󰄬"
	default:
		return "Undefined!!!"
	}
}

type wifiStoredKeyMap struct {
	edit       key.Binding
	connect    key.Binding
	disconnect key.Binding
	update     key.Binding
	delete     key.Binding
}

func (k *wifiStoredKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.edit,
		k.connect,
		k.disconnect,
		k.update,
		k.delete,
	}
}

func (k *wifiStoredKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.edit,
		k.connect,
		k.disconnect,
		k.update,
		k.delete,
	}}
}

type WifiStoredModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   wifiStoredState

	storedInfo *WifiStoredInfoModel

	connColIdx  int
	ssidColIdx int
	nameColIdx int

	minSSIDWidth         int
	indicatorStateHeight int

	keys *wifiStoredKeyMap

	nm infra.NetworkManager

	width  int
	height int
}

type WifiStoredColumnIndex int

func NewWifiStoredModel(storedInfo *WifiStoredInfoModel, keys *wifiStoredKeyMap, networkManager infra.NetworkManager) *WifiStoredModel {
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

	model := &WifiStoredModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   DoneInStored,
		storedInfo:       storedInfo,
		keys:             keys,
		nm:               networkManager,

		connColIdx:  0,
		ssidColIdx: 1,
		nameColIdx: 2,

		minSSIDWidth: 4,
	}
	model.bakeSizes()

	return model
}

func (m *WifiStoredModel) bakeSizes() {
	state := m.indicatorView()
	m.indicatorStateHeight = lipgloss.Height(state)
}

func (m *WifiStoredModel) Resize(width, height int) {
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

func (m *WifiStoredModel) Width() int {
	return m.width
}

func (m *WifiStoredModel) Height() int {
	return m.height
}

func (m *WifiStoredModel) Init() tea.Cmd {
	return m.RescanCmd()
}

func (m *WifiStoredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.storedInfo.setNew(info),
				OpenPopup(m.storedInfo, "Stored Wi-Fi info"),
			)

		case " ":
			return m, m.connectToSelectedCmd()

		case "shift+ ":
			return m,
				m.disconnectFromSelectedCmd()
		case "r":
			return m, RescanWifiStoredCmd(0)
		case "d":
			return m, m.deleteSelectedCmd()
		}
	case RescanWifiStoredMsg:
		time.Sleep(msg.delay)
		return m, m.RescanCmd()
	case WifiStoredStateMsg:
		return m, m.setStateCmd(wifiStoredState(msg))
	}

	var cmd tea.Cmd
	if m.indicatorState != DoneInStored {
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

func (m *WifiStoredModel) View() string {
	view := m.dataTable.View()

	statusline := m.indicatorView()
	return lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)
}

func (m *WifiStoredModel) indicatorView() string {
	var view string
	if m.indicatorState != DoneInStored {
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

type RescanWifiStoredMsg struct {
	delay time.Duration
}

func RescanWifiStoredCmd(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		return RescanWifiStoredMsg{delay: delay}
	}
}

func (m *WifiStoredModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ScanningStored),
		func() tea.Msg {
			list, err := m.nm.GetStoredWifi()
			if err != nil {
				return tea.BatchMsg{
					NotifyCmd("Cannot get stored wifi networks"),
					m.setStateCmd(DoneInStored),
				}
			}
			rows := []table.Row{}
			for _, wifiStored := range list {
				var connectionFlag string
				if wifiStored.Active {
					connectionFlag = ""
				}
				rows = append(rows, table.Row{
					connectionFlag,
					wifiStored.SSID,
					wifiStored.Name,
				})
			}

			m.dataTable.SetRows(rows)

			return m.setStateCmd(DoneInStored)
		},
	)
}

type WifiStoredStateMsg wifiStoredState

func (m *WifiStoredModel) setStateCmd(state wifiStoredState) tea.Cmd {
	updCmd := func() tea.Msg {
		m.indicatorState = state
		return NilMsg{}
	}

	if state == DoneInStored {
		return updCmd
	} else {
		return tea.Sequence(updCmd, m.indicatorSpinner.Tick)
	}
}

func (m *WifiStoredModel) connectToSelectedCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ConnectingStored),
		func() tea.Msg {
			name := m.dataTable.SelectedRow()[m.nameColIdx]
			err := m.nm.ConnectStoredWifi(name)
			if err != nil {
				return tea.BatchMsg{
					m.setStateCmd(DoneInStored),
					NotifyCmd(fmt.Sprintf("Cannot connect to %s", name)),
				}
			}
			return tea.BatchMsg{
				m.setStateCmd(DoneInStored),
				m.gotoTop(),
				RescanWifiCmd(0),
			}
		},
	)
}

func (m *WifiStoredModel) gotoTop() tea.Cmd {
	return func() tea.Msg {
		m.dataTable.GotoTop()
		return NilCmd
	}
}

func (m *WifiStoredModel) disconnectFromSelectedCmd() tea.Cmd {
	return func() tea.Msg {
		name := m.dataTable.SelectedRow()[m.nameColIdx]
		err := m.nm.DisconnectFromWifi(name)
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

func (m *WifiStoredModel) deleteSelectedCmd() tea.Cmd {
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
