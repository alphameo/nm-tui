package config

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type KeyBinding struct {
	Keys []string `kdl:"keys,arguments"`
}

type KeyConfig struct {
	Toggle *KeyBinding `kdl:"toggle"`

	Main   *MainKeys   `kdl:"main"`
	Dialog *DialogKeys `kdl:"dialog"`

	Wifi          *WifiKeys          `kdl:"wifi"`
	WifiAvailable *WifiAvailableKeys `kdl:"wifi_available"`
	WifiSaved     *WifiSavedKeys     `kdl:"wifi_saved"`
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

func DefaultKeys() *KeyConfig {
	return &KeyConfig{
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

func (k *KeyConfig) merge(src *KeyConfig) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, k.Toggle.merge(src.Toggle)...)
	errs = append(errs, k.Main.merge(src.Main)...)
	errs = append(errs, k.Dialog.merge(src.Dialog)...)
	errs = append(errs, k.Wifi.merge(src.Wifi)...)
	errs = append(errs, k.WifiAvailable.merge(src.WifiAvailable)...)
	errs = append(errs, k.WifiSaved.merge(src.WifiSaved)...)
	return errs
}

func (m *MainKeys) merge(src *MainKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, m.NextTab.merge(src.NextTab)...)
	errs = append(errs, m.PrevTab.merge(src.PrevTab)...)
	errs = append(errs, m.FocusNext.merge(src.FocusNext)...)
	errs = append(errs, m.FocusPrev.merge(src.FocusPrev)...)
	errs = append(errs, m.Quit.merge(src.Quit)...)
	return errs
}

func (d *DialogKeys) merge(src *DialogKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, d.FocusDown.merge(src.FocusDown)...)
	errs = append(errs, d.FocusUp.merge(src.FocusUp)...)
	errs = append(errs, d.TogglePWVisibility.merge(src.TogglePWVisibility)...)
	errs = append(errs, d.Accept.merge(src.Accept)...)
	errs = append(errs, d.Close.merge(src.Close)...)
	return errs
}

func (w *WifiKeys) merge(src *WifiKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, w.CreateProfile.merge(src.CreateProfile)...)
	errs = append(errs, w.OpenCaptivePortal.merge(src.OpenCaptivePortal)...)
	errs = append(errs, w.EnableHotspot.merge(src.EnableHotspot)...)
	errs = append(errs, w.CreateHotspot.merge(src.CreateHotspot)...)
	return errs
}

func (a *WifiAvailableKeys) merge(src *WifiAvailableKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, a.Connect.merge(src.Connect)...)
	return errs
}

func (s *WifiSavedKeys) merge(src *WifiSavedKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, s.Edit.merge(src.Edit)...)
	errs = append(errs, s.Connect.merge(src.Connect)...)
	errs = append(errs, s.Disconnect.merge(src.Disconnect)...)
	errs = append(errs, s.Delete.merge(src.Delete)...)
	return errs
}

func (b *KeyBinding) merge(src *KeyBinding) []error {
	if src == nil {
		return nil
	}

	var errs []error
	for _, k := range src.Keys {
		if !validKeyName(k) {
			errs = append(errs, fmt.Errorf("invalid key: %q", k))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	b.Keys = src.Keys
	return nil
}

var validModifier = map[string]bool{
	"ctrl": true, "alt": true, "shift": true,
	"meta": true, "hyper": true, "super": true,
	"capslock": true, "scrolllock": true, "numlock": true,
}

var validKey = map[string]bool{
	"enter": true, "tab": true, "backspace": true, "esc": true, "space": true,
	"up": true, "down": true, "left": true, "right": true,
	"begin": true, "find": true, "insert": true, "delete": true, "select": true,
	"pgup": true, "pgdown": true, "home": true, "end": true,

	"equal": true, "mul": true, "plus": true, "comma": true,
	"minus": true, "period": true, "div": true, "sep": true,
	"0": true, "1": true, "2": true, "3": true, "4": true,
	"5": true, "6": true, "7": true, "8": true, "9": true,

	"capslock": true, "scrolllock": true, "numlock": true,
	"printscreen": true, "pause": true, "menu": true,

	"mediaplay": true, "mediapause": true, "mediaplaypause": true,
	"mediastop": true, "mediafastforward": true, "mediarewind": true,
	"medianext": true, "mediaprev": true, "mediarecord": true,

	"lowervol": true, "raisevol": true, "mute": true,

	"leftshift": true, "leftalt": true, "leftctrl": true,
	"leftsuper": true, "lefthyper": true, "leftmeta": true,
	"rightshift": true, "rightalt": true, "rightctrl": true,
	"rightsuper": true, "righthyper": true, "rightmeta": true,
	"isolevel3shift": true, "isolevel5shift": true,
}

func init() {
	for i := 1; i <= 63; i++ {
		validKey[fmt.Sprintf("f%d", i)] = true
	}
}

func validKeyName(s string) bool {
	if s == "" {
		return false
	}

	parts := strings.Split(s, "+")
	if len(parts) == 0 {
		return false
	}

	key := strings.ToLower(parts[len(parts)-1])

	if utf8.RuneCountInString(key) == 1 {
		r, _ := utf8.DecodeRuneInString(key)
		if unicode.IsPrint(r) {
			return true
		}
	}

	if !validKey[key] {
		return false
	}

	for _, m := range parts[:len(parts)-1] {
		if !validModifier[strings.ToLower(m)] {
			return false
		}
	}

	return true
}

func HelpFromKeys(keys []string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
