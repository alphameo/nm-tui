package components

import (
	"fmt"
	"log/slog"

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

	keys *wifiStoredKeyMap

	nm infra.NetworkManager

	width  int
	height int
}

type WifiStoredColumnIndex int

const (
	storedSSIDColumn WifiStoredColumnIndex = 1
	storedNameColumn WifiStoredColumnIndex = 2
)

func NewWifiStoredModel(networkManager infra.NetworkManager, keys *keyMapManager) *WifiStoredModel {
	cols := []table.Column{
		{Title: "󱘖", Width: connectionFlagColumnWidth},
		{Title: "SSID"},
		{Title: "Name"},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)
	s := spinner.New()
	info := NewStoredInfoModel(networkManager, keys)

	return &WifiStoredModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   DoneInStored,
		storedInfo:       info,
		keys:             keys.wifiStored,
		nm:               networkManager,
	}
}

func (m *WifiStoredModel) Resize(width, height int) {
	m.width = width
	m.height = height

	height -= indicatorStateHeight

	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	computedWidth := width - tableUtilityOffset - connectionFlagColumnWidth
	possibleNameWidth := computedWidth / 2
	ssidWidth := max(computedWidth-possibleNameWidth, minSSIDWidth)
	nameWidth := computedWidth - ssidWidth
	m.dataTable.Columns()[storedNameColumn].Width = nameWidth
	m.dataTable.Columns()[storedSSIDColumn].Width = ssidWidth
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
			info, err := m.nm.GetWifiInfo(row[storedNameColumn])
			if err != nil {
				return m, NotifyCmd(err.Error())
			}
			m.storedInfo.setNew(info)
			return m, tea.Sequence(
				SetPopupContentCmd(m.storedInfo, "Stored Wi-Fi info"),
				SetPopupActivityCmd(true),
			)

		case " ":
			return m, m.connectToSelectedCmd()

		case "shift+ ":
			return m,
				m.disconnectFromSelectedCmd()
		case "r":
			return m, RescanWifiStoredCmd()
		case "d":
			return m, m.deleteSelectedCmd()
		}
	case RescanWifiStoredMsg:
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

	var statusline string
	if m.indicatorState != DoneInStored {
		statusline = fmt.Sprintf(
			"%s %s",
			m.indicatorState.String(),
			m.indicatorSpinner.View(),
		)
	} else {
		statusline = m.indicatorState.String()
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		statusline,
	)
}

type RescanWifiStoredMsg struct{}

func RescanWifiStoredCmd() tea.Cmd {
	return func() tea.Msg {
		return RescanWifiStoredMsg{}
	}
}

func (m *WifiStoredModel) RescanCmd() tea.Cmd {
	return tea.Sequence(
		m.setStateCmd(ScanningStored),
		func() tea.Msg {
			list, err := m.nm.GetStoredWifi()
			if err != nil {
				slog.Error(err.Error())
				return NotifyCmd(err.Error())
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
			err := m.nm.ConnectStoredWifi(m.dataTable.SelectedRow()[storedNameColumn])
			if err != nil {
				return NotifyCmd(err.Error())
			}
			return tea.BatchMsg{
				m.setStateCmd(DoneInStored),
				m.gotoTop(),
				RescanWifiCmd(),
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
		err := m.nm.DisconnectFromWifi(m.dataTable.SelectedRow()[storedNameColumn])
		if err != nil {
			return NotifyCmd(err.Error())
		}
		return NilMsg{}
	}
}

func (m *WifiStoredModel) deleteSelectedCmd() tea.Cmd {
	row := m.dataTable.SelectedRow()
	return func() tea.Msg {
		err := m.nm.DeleteWifiConnection(row[storedNameColumn])
		if err != nil {
			return NotifyCmd(err.Error())
		}
		cursor := m.dataTable.Cursor()
		if cursor == len(m.dataTable.Rows())-1 {
			m.dataTable.SetCursor(cursor - 1)
		}
		return NilMsg{}
	}
}
