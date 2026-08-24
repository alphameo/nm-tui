package models

import (
	"strings"

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
	device            deviceKeyMap
	networkDevices    networkDevicesKeyMap
	networks          networksKeyMap
	networkProfiles   networkProfilesKeyMap
	profileEditor     profileEditorKeyMap
	availableNetworks availableNetworksKeyMap
	connector         connectorKeyMap
	profileCreator    profileCreatorKeyMap
	hotspotCreator    hotspotCreatorKeyMap
	help              helpKeyMap
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
		device: deviceKeyMap{
			prev:   NewKey(*keys.FocusPrev, "prev field"),
			next:   NewKey(*keys.FocusNext, "next field"),
			rescan: NewKey(*keys.Rescan, "rescan"),
		},
		networkDevices: networkDevicesKeyMap{
			showInfo: NewKey(*keys.NetworkDevices.ShowInfo, "show info"),
		},
		networks: networksKeyMap{
			winNext:           NewKey(*keys.FocusNext, "next window"),
			winPrev:           NewKey(*keys.FocusPrev, "prev window"),
			rescan:            NewKey(*keys.Rescan, "rescan"),
			createProfile:     NewKey(*keys.Networks.CreateProfile, "create profile"),
			createHotspot:     NewKey(*keys.Networks.CreateHotspot, "create hotspot"),
			quickHotspot:      NewKey(*keys.Networks.QuickHotspot, "quick hotspot"),
			openCaptivePortal: NewKey(*keys.Networks.OpenCaptivePortal, "login portal"),
			win1:              NewKey(*keys.Focus1, "1st window"),
			win2:              NewKey(*keys.Focus2, "2nd window"),
		},
		networkProfiles: networkProfilesKeyMap{
			edit:       NewKey(*keys.NetworkProfiles.Edit, "edit"),
			activate:   NewKey(*keys.NetworkProfiles.Activate, "activate"),
			deactivate: NewKey(*keys.NetworkProfiles.Deactivate, "deactivate"),
			delete:     NewKey(*keys.NetworkProfiles.Delete, "delete"),
		},
		availableNetworks: availableNetworksKeyMap{
			connect:    NewKey(*keys.AvailableNetworks.Connect, "connect"),
			activate:   NewKey(*keys.AvailableNetworks.Activate, "activate"),
			deactivate: NewKey(*keys.AvailableNetworks.Deactivate, "deactivate"),
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
		help: helpKeyMap{
			quit: NewKey(*keys.Main.Help, "quit help"),
		},
	}
}

func NewKey(keys []string, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(strings.Join(keys, "/"), desc),
	)
}
