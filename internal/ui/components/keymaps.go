package components

import (
	"charm.land/bubbles/v2/key"
	"github.com/alphameo/nm-tui/internal/ui/components/tabview"
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
)

func NewKeyMap(keys []string, keyHelp, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyHelp, desc),
	)
}

type keyMapManager struct {
	main          *mainKeyMap
	popup         *popupKeyMap
	tabs          *tabview.KeyMap
	toggle        *toggle.KeyMap
	network       *networkKeyMap
	wifi          *wifiKeyMap
	wifiSaved     *wifiSavedKeyMap
	wifiSavedInfo *wifiSavedInfoKeyMap
	wifiAvailable *wifiAvailableKeyMap
	wifiConnector *wifiConnectorKeyMap
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
	main:          mainKeys,
	popup:         popupKeys,
	tabs:          tabview.DefaultKeys,
	toggle:        toggle.DefaultKeys,
	network:       networkKeys,
	wifi:          wifiKeys,
	wifiSaved:     wifiSavedKeys,
	wifiSavedInfo: wifiSavedInfoKeys,
	wifiAvailable: wifiAvailableKeys,
	wifiConnector: wifiConnectorKeys,
}
