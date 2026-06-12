package models

import (
	"charm.land/bubbles/v2/key"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

func NewKeyMap(keys []string, keyHelp, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyHelp, desc),
	)
}

type keyMapManager struct {
	main           *mainKeyMap
	popup          *popupKeyMap
	tabs           *tabview.KeyMap
	toggle         *toggle.KeyMap
	networking     *networkingKeyMap
	wifi           *wifiKeyMap
	wifiSaved      *wifiSavedKeyMap
	profileEditor  *profileEditorKeyMap
	wifiAvailable  *wifiAvailableKeyMap
	connector      *connectorKeyMap
	profileCreator *profileCreatorKeyMap
	hotspotCreator *hotspotCreatorKeyMap
}

func (k *keyMapManager) ShortHelp() []key.Binding {
	return []key.Binding{k.main.quit}
}

func (k *keyMapManager) FullHelp() [][]key.Binding {
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

var defaultKeyMap = &keyMapManager{
	main:           mainKeys,
	popup:          popupKeys,
	tabs:           tabview.DefaultKeys,
	toggle:         toggle.DefaultKeys,
	networking:     networkingKeys,
	wifi:           wifiKeys,
	wifiSaved:      wifiSavedKeys,
	wifiAvailable:  wifiAvailableKeys,
	profileEditor:  profileEditorKeys,
	connector:      connectorKeys,
	profileCreator: profileCreatorKeys,
	hotspotCreator: hotspotCreatorKeys,
}
