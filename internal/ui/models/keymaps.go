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
			k.tabs.TabNext,
			k.tabs.TabPrev,
		},
	}
}

func initKeys(keys config.KeyConfig) keyMaps {
	return keyMaps{
		main: mainKeyMap{
			quit:       NewKey(*keys.Main.Quit, "quit"),
			closePopup: NewKey(*keys.Dialog.Close, "close dialog"),
		},
		tabs: tabview.KeyMap{
			TabNext: NewKey(*keys.Main.NextTab, "next tab"),
			TabPrev: NewKey(*keys.Main.PrevTab, "prev tab"),
		},
		toggle: toggle.KeyMap{
			Toggle: NewKey(*keys.Toggle, "toggle"),
		},
		connectivity: connectivityKeyMap{
			up:     NewKey(*keys.Dialog.FocusUp, "focus up"),
			down:   NewKey(*keys.Dialog.FocusDown, "focus down"),
			rescan: NewKey(*keys.Rescan, "rescan"),
			toggle: NewKey(*keys.Toggle, "toggle"),
		},
		networks: networksKeyMap{
			nextWindow:        NewKey(*keys.Main.FocusNext, "focus next"),
			rescan:            NewKey(*keys.Rescan, "rescan networks"),
			createProfile:     NewKey(*keys.Networks.CreateProfile, "create profile"),
			createHotspot:     NewKey(*keys.Networks.CreateHotspot, "create hotspot"),
			quickHotspot:      NewKey(*keys.Networks.QuickHotspot, "quick hotspot"),
			openCaptivePortal: NewKey(*keys.Networks.OpenCaptivePortal, "open captive portal"),
			firstWindow:       NewKey(*keys.Focus1, "focus first window"),
			secondWindow:      NewKey(*keys.Focus2, "focus second window"),
		},
		networkProfiles: networkProfilesKeyMap{
			rescan:     NewKey(*keys.RescanFocused, "rescan network profiles"),
			edit:       NewKey(*keys.NetworkProfiles.Edit, "edit profile"),
			connect:    NewKey(*keys.NetworkProfiles.Connect, "connect"),
			disconnect: NewKey(*keys.NetworkProfiles.Disconnect, "disconnect"),
			delete:     NewKey(*keys.NetworkProfiles.Delete, "delete profile"),
		},
		availableNetworks: availableNetworksKeyMap{
			rescan:  NewKey(*keys.RescanFocused, "rescan available networks"),
			connect: NewKey(*keys.AvailableNetworks.Connect, "connect"),
		},
		profileEditor: profileEditorKeyMap{
			up:                 NewKey(*keys.Dialog.FocusUp, "focus up"),
			down:               NewKey(*keys.Dialog.FocusDown, "focus down"),
			save:               NewKey(*keys.Dialog.Accept, "save changes"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "toggle password visibility"),
		},
		connector: connectorKeyMap{
			up:                 NewKey(*keys.Dialog.FocusUp, "focus up"),
			down:               NewKey(*keys.Dialog.FocusDown, "focus down"),
			connect:            NewKey(*keys.Dialog.Accept, "connect"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "toggle password visibility"),
		},
		profileCreator: profileCreatorKeyMap{
			up:                 NewKey(*keys.Dialog.FocusUp, "focus up"),
			down:               NewKey(*keys.Dialog.FocusDown, "focus down"),
			create:             NewKey(*keys.Dialog.Accept, "create"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "toggle password visibility"),
		},
		hotspotCreator: hotspotCreatorKeyMap{
			up:                 NewKey(*keys.Dialog.FocusUp, "focus up"),
			down:               NewKey(*keys.Dialog.FocusDown, "focus down"),
			create:             NewKey(*keys.Dialog.Accept, "create"),
			togglePWVisibility: NewKey(*keys.Dialog.TogglePWVisibility, "toggle password visibility"),
		},
	}
}

func NewKey(keys []string, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(HelpFromKeys(keys...), desc),
	)
}

func HelpFromKeys(keys ...string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
