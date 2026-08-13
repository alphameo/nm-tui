package models

import (
	"charm.land/bubbles/v2/key"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

func NewKeyMap(keys []string, keyHelp, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyHelp, desc),
	)
}

type keyMaps struct {
	main              mainKeyMap
	tabs              tabview.KeyMap
	toggle            toggle.KeyMap
	connectivity      connectivityKeyMap
	networks          networksKeyMap
	networkProfiles   networkProfilesKeyMap
	profileEditor     profileEditorKeyMap
	availableNetworks availableNetworksKeyMap
	connector         connectorKeyMap
	profileCreator    profileCreatorKeyMap
	hotspotCreator    hotspotCreatorKeyMap
}

func (k *keyMaps) ShortHelp() []key.Binding {
	return []key.Binding{k.main.quit}
}

func (k *keyMaps) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.main.quit,
		},
		{
			k.tabs.Next,
			k.tabs.Prev,
		},
	}
}

func initKeys(keys config.KeyConfig) keyMaps {
	return keyMaps{
		main: mainKeyMap{
			quit:       NewKey(*keys.Main.Quit, "quit"),
			closePopup: NewKey(*keys.Dialog.Close, "close popup"),
			help:       NewKey(*keys.Main.Help, "help"),
		},
		tabs: tabview.KeyMap{
			Next: NewKey(*keys.Main.TabNext, "next tab"),
			Prev: NewKey(*keys.Main.TabPrev, "prev tab"),
		},
		toggle: toggle.KeyMap{
			Toggle: NewKey(*keys.Toggle, "toggle"),
		},
		connectivity: connectivityKeyMap{
			prev:   NewKey(*keys.FocusPrev, "prev field"),
			next:   NewKey(*keys.FocusNext, "next field"),
			rescan: NewKey(*keys.Rescan, "rescan"),
		},
		networks: networksKeyMap{
			winNext:           NewKey(*keys.FocusNext, "next window"),
			winPrev:           NewKey(*keys.FocusPrev, "prev window"),
			rescan:            NewKey(*keys.Rescan, "rescan all"),
			createProfile:     NewKey(*keys.Networks.CreateProfile, "create profile"),
			createHotspot:     NewKey(*keys.Networks.CreateHotspot, "create hotspot"),
			quickHotspot:      NewKey(*keys.Networks.QuickHotspot, "quick hotspot"),
			openCaptivePortal: NewKey(*keys.Networks.OpenCaptivePortal, "open login portal"),
			win1:              NewKey(*keys.Focus1, "1st window"),
			win2:              NewKey(*keys.Focus2, "2nd window"),
		},
		networkProfiles: networkProfilesKeyMap{
			rescan:     NewKey(*keys.RescanFocused, "rescan profiles"),
			edit:       NewKey(*keys.NetworkProfiles.Edit, "edit"),
			connect:    NewKey(*keys.NetworkProfiles.Connect, "connect"),
			disconnect: NewKey(*keys.NetworkProfiles.Disconnect, "disconnect"),
			delete:     NewKey(*keys.NetworkProfiles.Delete, "delete"),
		},
		availableNetworks: availableNetworksKeyMap{
			rescan:  NewKey(*keys.RescanFocused, "rescan available"),
			connect: NewKey(*keys.AvailableNetworks.Connect, "connect"),
		},
		profileEditor: profileEditorKeyMap{
			prev:               NewKey(*keys.FocusPrev, "prev field"),
			next:               NewKey(*keys.FocusNext, "next field"),
			save:               NewKey(*keys.Dialog.Accept, "save"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "pw visibility"),
		},
		connector: connectorKeyMap{
			prev:               NewKey(*keys.FocusPrev, "prev field"),
			next:               NewKey(*keys.FocusNext, "next field"),
			connect:            NewKey(*keys.Dialog.Accept, "connect"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "pw visibility"),
		},
		profileCreator: profileCreatorKeyMap{
			prev:               NewKey(*keys.FocusPrev, "prev field"),
			next:               NewKey(*keys.FocusNext, "next field"),
			create:             NewKey(*keys.Dialog.Accept, "create"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "pw visibility"),
		},
		hotspotCreator: hotspotCreatorKeyMap{
			prev:               NewKey(*keys.FocusPrev, "prev field"),
			next:               NewKey(*keys.FocusNext, "next field"),
			create:             NewKey(*keys.Dialog.Accept, "create"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "pw visibility"),
		},
	}
}

func NewKey(keys []string, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(HelpFromKeys(keys...), desc),
	)
}
