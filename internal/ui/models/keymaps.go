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
			quit:       NewKey(*keys.Main.Quit),
			closePopup: NewKey(*keys.Dialog.Close),
		},
		tabs: tabview.KeyMap{
			Next: NewKey(*keys.Main.TabNext),
			Prev: NewKey(*keys.Main.TabPrev),
		},
		toggle: toggle.KeyMap{
			Toggle: NewKey(*keys.Toggle),
		},
		connectivity: connectivityKeyMap{
			prev:   NewKey(*keys.FocusPrev),
			next:   NewKey(*keys.FocusNext),
			rescan: NewKey(*keys.Rescan),
			toggle: NewKey(*keys.Toggle),
		},
		networks: networksKeyMap{
			winNext:           NewKey(*keys.FocusNext),
			winPrev:           NewKey(*keys.FocusPrev),
			rescan:            NewKey(*keys.Rescan),
			createProfile:     NewKey(*keys.Networks.CreateProfile),
			createHotspot:     NewKey(*keys.Networks.CreateHotspot),
			quickHotspot:      NewKey(*keys.Networks.QuickHotspot),
			openCaptivePortal: NewKey(*keys.Networks.OpenCaptivePortal),
			win1:              NewKey(*keys.Focus1),
			win2:              NewKey(*keys.Focus2),
		},
		networkProfiles: networkProfilesKeyMap{
			rescan:     NewKey(*keys.RescanFocused),
			edit:       NewKey(*keys.NetworkProfiles.Edit),
			connect:    NewKey(*keys.NetworkProfiles.Connect),
			disconnect: NewKey(*keys.NetworkProfiles.Disconnect),
			delete:     NewKey(*keys.NetworkProfiles.Delete),
		},
		availableNetworks: availableNetworksKeyMap{
			rescan:  NewKey(*keys.RescanFocused),
			connect: NewKey(*keys.AvailableNetworks.Connect),
		},
		profileEditor: profileEditorKeyMap{
			prev:               NewKey(*keys.FocusPrev),
			next:               NewKey(*keys.FocusNext),
			save:               NewKey(*keys.Dialog.Accept),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility),
		},
		connector: connectorKeyMap{
			prev:               NewKey(*keys.FocusPrev),
			next:               NewKey(*keys.FocusNext),
			connect:            NewKey(*keys.Dialog.Accept),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility),
		},
		profileCreator: profileCreatorKeyMap{
			prev:               NewKey(*keys.FocusPrev),
			next:               NewKey(*keys.FocusNext),
			create:             NewKey(*keys.Dialog.Accept),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility),
		},
		hotspotCreator: hotspotCreatorKeyMap{
			prev:               NewKey(*keys.FocusPrev),
			next:               NewKey(*keys.FocusNext),
			create:             NewKey(*keys.Dialog.Accept),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility),
		},
	}
}

func NewKey(keys []string) key.Binding {
	return key.NewBinding(key.WithKeys(keys...))
}
