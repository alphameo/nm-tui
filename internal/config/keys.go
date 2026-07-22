package config

import (
	"strings"
)

type KeyBinding struct {
	Keys []string `kdl:"keys,arguments"`
}

type Keys struct {
	Toggle *KeyBinding `kdl:"toggle"`

	Main   *MainKeys   `kdl:"main"`
	Dialog *DialogKeys `kdl:"dialog"`

	Wifi          *WifiKeys          `kdl:"wifi"`
	WifiAvailable *WifiAvailableKeys `kdl:"wifi_available"`
	WifiSaved     *WifiSavedKeys     `kdl:"wifi_saved"`
}

func DefaultKeys() *Keys {
	return &Keys{
		Toggle: &KeyBinding{Keys: []string{"space"}},
		Main: &MainKeys{
			NextTab:   &KeyBinding{Keys: []string{"]"}},
			PrevTab:   &KeyBinding{Keys: []string{"["}},
			FocusNext: &KeyBinding{Keys: []string{"tab"}},
			FocusPrev: &KeyBinding{Keys: []string{"shift+tab"}},
			Quit:      &KeyBinding{Keys: []string{"esc", "ctrl+c", "q", "ctrl+q"}},
		},
		Dialog: &DialogKeys{
			FocusDown:          &KeyBinding{Keys: []string{"ctrl+j"}},
			FocusUp:            &KeyBinding{Keys: []string{"tab"}},
			TogglePWVisibility: &KeyBinding{Keys: []string{"ctrl+p"}},
			Accept:             &KeyBinding{Keys: []string{"ctrl+enter"}},
			Close:              &KeyBinding{Keys: []string{"ctrl+q"}},
		},
		Wifi: &WifiKeys{
			CreateProfile:     &KeyBinding{Keys: []string{"a", "c"}},
			OpenCaptivePortal: &KeyBinding{Keys: []string{"l"}},
			EnableHotspot:     &KeyBinding{Keys: []string{"ctrl+h"}},
			CreateHotspot:     &KeyBinding{Keys: []string{"h"}},
		},
		WifiAvailable: &WifiAvailableKeys{
			Connect: &KeyBinding{Keys: []string{"enter"}},
		},
		WifiSaved: &WifiSavedKeys{
			Edit:       &KeyBinding{Keys: []string{"enter"}},
			Connect:    &KeyBinding{Keys: []string{"space"}},
			Disconnect: &KeyBinding{Keys: []string{"ctrl+space"}},
			Delete:     &KeyBinding{Keys: []string{"d", "delete"}},
		},
	}
}

type MainKeys struct {
	NextTab   *KeyBinding `kdl:"next_tab"`
	PrevTab   *KeyBinding `kdl:"prev_tab"`
	FocusNext *KeyBinding `kdl:"focus_next"`
	FocusPrev *KeyBinding `kdl:"focus_prev"`
	Quit      *KeyBinding `kdl:"quit"`
}

type DialogKeys struct {
	FocusDown          *KeyBinding `kdl:"focus_down"`
	FocusUp            *KeyBinding `kdl:"focus_up"`
	TogglePWVisibility *KeyBinding `kdl:"toggle_pw_visibility"`
	Accept             *KeyBinding `kdl:"accept"`
	Close              *KeyBinding `kdl:"close"`
}

type WifiKeys struct {
	CreateProfile     *KeyBinding `kdl:"create_profile"`
	OpenCaptivePortal *KeyBinding `kdl:"open_network_login"`
	EnableHotspot     *KeyBinding `kdl:"enable_hotspot"`
	CreateHotspot     *KeyBinding `kdl:"create_hotspot"`
}

type WifiAvailableKeys struct {
	Connect *KeyBinding `kdl:"connect"`
}

type WifiSavedKeys struct {
	Edit       *KeyBinding `kdl:"edit"`
	Connect    *KeyBinding `kdl:"connect"`
	Disconnect *KeyBinding `kdl:"disconnect"`
	Delete     *KeyBinding `kdl:"delete"`
}

func HelpFromKeys(keys []string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
