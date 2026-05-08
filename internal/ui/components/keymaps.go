package components

import (
	"github.com/alphameo/nm-tui/internal/ui/components/toggle"
	"github.com/charmbracelet/bubbles/key"
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
	tabs          *tabsKeyMap
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
			k.tabs.tabNext,
			k.tabs.tabPrev,
		},
	}
}

var defaultKeyMap = &keyMapManager{
	main:          mainKeys,
	popup:         popupKeys,
	tabs:          tabsKeys,
	toggle:        toggle.DefaultKeys,
	network:       networkKeys,
	wifi:          wifiKeys,
	wifiSaved:     wifiSavedKeys,
	wifiSavedInfo: wifiSavedInfoKeys,
	wifiAvailable: wifiAvailableKeys,
	wifiConnector: wifiConnectorKeys,
}
