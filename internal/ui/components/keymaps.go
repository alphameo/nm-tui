package components

import (
	"github.com/alphameo/nm-tui/internal/ui/components/floating"
	"github.com/charmbracelet/bubbles/key"
)

func NewKeyMap(keys []string, keyHelp, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyHelp, desc),
	)
}

type keyMaps struct {
	main           *mainKeyMap
	floating       *floating.FloatingKeyMap
	tabs           *tabsKeyMap
	toggle         *toggleKeyMap
	wifi           *wifiKeyMap
	wifiStored     *wifiStoredKeyMap
	wifiStoredInfo *wifiStoredInfoKeyMap
	wifiAvailable  *wifiAvailableKeyMap
	wifiConnector  *wifiConnectorKeyMap
}

var defaultKeyMap = &keyMaps{
	main:           mainKeys,
	floating:       floatingKeys,
	tabs:           tabsKeys,
	toggle:         toggleKeys,
	wifi:           wifiKeys,
	wifiStored:     wifiStoredKeys,
	wifiStoredInfo: wifiStoredInfoKeys,
	wifiAvailable:  wifiAvailableKeys,
	wifiConnector:  wifiConnectorKeys,
}

var mainKeys = &mainKeyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/q/^Q/^C", "quit"),
	),
}

var floatingKeys = &floating.FloatingKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/^Q/^C", "quit"),
	),
}

var tabsKeys = &tabsKeyMap{
	tabNext: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next tab"),
	),
	tabPrev: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "previous tab"),
	),
}

var toggleKeys = &toggleKeyMap{
	toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "toggle"),
	),
}

var wifiStoredKeys = &wifiStoredKeyMap{
	edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit"),
	),
	connect: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "connect"),
	),
	disconnect: key.NewBinding(
		key.WithKeys("shift+"),
		key.WithHelp("shift+ ", "disconnect"),
	),
	update: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan stored"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

var wifiStoredInfoKeys = &wifiStoredInfoKeyMap{
	togglePasswordVisibility: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^R", "toggle password visibility"),
	),
	up: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^K", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("^J", "down"),
	),
	submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

var wifiKeys = &wifiKeyMap{
	nextWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	firstWindow: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "first window"),
	),
	secondWindow: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "second window"),
	),
	rescan: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^R", "rescan"),
	),
}

var wifiAvailableKeys = &wifiAvailableKeyMap{
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan"),
	),
	openConnector: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open connector"),
	),
}

var wifiConnectorKeys = &wifiConnectorKeyMap{
	connect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open connector"),
	),
	togglePasswordVisibility: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^r", "toggle password visibility"),
	),
}
