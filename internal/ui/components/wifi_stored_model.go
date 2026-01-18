package components

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiStoredModel struct {
	dataTable  table.Model
	storedInfo *WifiStoredInfoModel
	nm         infra.NetworkManager
	width      int
	height     int
}

type WifiStoredColumnIndex int

const (
	ssidColumn WifiStoredColumnIndex = 1
	nameColumn WifiStoredColumnIndex = 2
)

func NewWifiStored(networkManager infra.NetworkManager) *WifiStoredModel {
	cols := []table.Column{
		{Title: "󱘖", Width: conectionFlagColumnWidth},
		{Title: "SSID"},
		{Title: "Name"},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle)
	s := NewStoredInfoModel(networkManager)

	return &WifiStoredModel{
		dataTable:  t,
		storedInfo: s,
		nm:         networkManager,
	}
}

func (m *WifiStoredModel) Resize(width, height int) {
	m.width = width
	m.height = height
	m.dataTable.SetWidth(width)
	m.dataTable.SetHeight(height)

	tableUtilityOffset := len(m.dataTable.Columns()) * 2

	computedWidth := width - tableUtilityOffset - conectionFlagColumnWidth
	possibleNameWidth := computedWidth / 2
	ssidWidth := max(computedWidth-possibleNameWidth, minSSIDWidth)
	nameWidth := computedWidth - ssidWidth
	m.dataTable.Columns()[nameColumn].Width = nameWidth
	m.dataTable.Columns()[ssidColumn].Width = ssidWidth
	m.dataTable.UpdateViewport()
}

func (m *WifiStoredModel) Init() tea.Cmd {
	return m.UpdateRowsCmd()
}

type updateWifiStoredMsg struct{}

// UpdateWifiStoredMsg is used to avoid extra instantiatons
var UpdateWifiStoredMsg = updateWifiStoredMsg{}

func (m *WifiStoredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				info, err := m.nm.GetWifiInfo(row[2])
				if err != nil {
					return m, NotifyCmd(err.Error())
				}
				m.storedInfo.setNew(info)
				return m, tea.Sequence(SetPopupActivityCmd(true), SetPopupContentCmd(m.storedInfo, "Stored Wi-Fi info"))
			}
			return m, nil
		case " ":
			return m, tea.Sequence(m.connectSelectedCmd(), m.UpdateRowsCmd())
		case "shift+ ":
			return m, tea.Sequence(m.disconnectFromSelectedCmd(), m.UpdateRowsCmd())
		case "r":
			return m, m.UpdateRowsCmd()
		case "d":
			row := m.dataTable.SelectedRow()
			cursor := m.dataTable.Cursor()
			if cursor == len(m.dataTable.Rows())-1 {
				m.dataTable.SetCursor(cursor - 1)
			}
			return m, tea.Sequence(
				func() tea.Msg {
					err := m.nm.DeleteWifiConnection(row[2])
					if err != nil {
						return NotifyCmd(err.Error())
					}
					return nil
				},
				m.UpdateRowsCmd())
		}
	case updateWifiStoredMsg:
		return m, m.UpdateRowsCmd()
	}

	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m *WifiStoredModel) View() string {
	view := m.dataTable.View()

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s", view)
	return sb.String()
}

func (m *WifiStoredModel) UpdateRowsCmd() tea.Cmd {
	return func() tea.Msg {
		list, err := m.nm.GetStoredWifi()
		if err != nil {
			logger.Errln(fmt.Errorf("error: %s", err.Error()))
			return NotifyCmd(err.Error())
		}
		rows := []table.Row{}
		for _, wifiStored := range list {
			var connectionFlag string
			if wifiStored.Active {
				connectionFlag = ""
			}
			rows = append(rows, table.Row{connectionFlag, wifiStored.SSID, wifiStored.Name})
		}

		m.dataTable.SetRows(rows)
		m.dataTable.UpdateViewport()
		return UpdateMsg
	}
}

func (m *WifiStoredModel) connectSelectedCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.nm.ConnectStoredWifi(m.dataTable.SelectedRow()[2])
		if err != nil {
			return NotifyCmd(err.Error())
		}
		m.dataTable.GotoTop()
		return nil
	}
}

func (m *WifiStoredModel) disconnectFromSelectedCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.nm.DisconnectFromWifi(m.dataTable.SelectedRow()[2])
		if err != nil {
			return NotifyCmd(err.Error())
		}
		return nil
	}
}

func UpdateWifiStoredRowsCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateWifiStoredMsg
	}
}
