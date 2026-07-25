package config

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/calico32/kdl-go"
)

type KeyBinding []string

func keyBinding(keys ...string) *KeyBinding {
	k := KeyBinding(keys)
	return &k
}

func (k *KeyBinding) UnmarshalKDL(node *kdl.Node) error {
	args := node.Arguments()
	*k = make(KeyBinding, len(args))
	for i, arg := range args {
		(*k)[i] = arg.String()
	}
	return nil
}

type KeyConfig struct {
	Toggle        *KeyBinding `kdl:"toggle"`
	Rescan        *KeyBinding `kdl:"rescan"`
	RescanFocused *KeyBinding `kdl:"rescan_focused"`
	FocusFirst    *KeyBinding `kdl:"focus_first"`
	FocusSecond   *KeyBinding `kdl:"focus_second"`

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
		Toggle:        keyBinding("space"),
		Rescan:        keyBinding("r"),
		RescanFocused: keyBinding("ctrl+r"),
		FocusFirst:    keyBinding("1"),
		FocusSecond:   keyBinding("2"),
		Main: &MainKeys{
			NextTab:   keyBinding("]"),
			PrevTab:   keyBinding("["),
			FocusNext: keyBinding("tab"),
			FocusPrev: keyBinding("shift+tab"),
			Quit:      keyBinding("esc", "ctrl+c", "q", "ctrl+q"),
		},
		Dialog: &DialogKeys{
			FocusDown:          keyBinding("ctrl+j"),
			FocusUp:            keyBinding("tab"),
			TogglePWVisibility: keyBinding("ctrl+p"),
			Accept:             keyBinding("ctrl+enter"),
			Close:              keyBinding("ctrl+q"),
		},
		Wifi: &WifiKeys{
			CreateProfile:     keyBinding("a", "c"),
			OpenCaptivePortal: keyBinding("l"),
			EnableHotspot:     keyBinding("ctrl+h"),
			CreateHotspot:     keyBinding("h"),
		},
		WifiAvailable: &WifiAvailableKeys{
			Connect: keyBinding("enter"),
		},
		WifiSaved: &WifiSavedKeys{
			Edit:       keyBinding("enter"),
			Connect:    keyBinding("space"),
			Disconnect: keyBinding("ctrl+space"),
			Delete:     keyBinding("d", "delete"),
		},
	}
}

func (k *KeyConfig) merge(src *KeyConfig) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(k.Toggle, src.Toggle, "toggle")...)
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
	errs = append(errs, mergeKeyList(m.NextTab, src.NextTab, "main.next_tab")...)
	errs = append(errs, mergeKeyList(m.PrevTab, src.PrevTab, "main.prev_tab")...)
	errs = append(errs, mergeKeyList(m.FocusNext, src.FocusNext, "main.focus_next")...)
	errs = append(errs, mergeKeyList(m.FocusPrev, src.FocusPrev, "main.focus_prev")...)
	errs = append(errs, mergeKeyList(m.Quit, src.Quit, "main.quit")...)
	return errs
}

func (d *DialogKeys) merge(src *DialogKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(d.FocusDown, src.FocusDown, "dialog.focus_down")...)
	errs = append(errs, mergeKeyList(d.FocusUp, src.FocusUp, "dialog.focus_up")...)
	errs = append(errs, mergeKeyList(d.TogglePWVisibility, src.TogglePWVisibility, "dialog.toggle_pw_visibility")...)
	errs = append(errs, mergeKeyList(d.Accept, src.Accept, "dialog.accept")...)
	errs = append(errs, mergeKeyList(d.Close, src.Close, "dialog.close")...)
	return errs
}

func (w *WifiKeys) merge(src *WifiKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(w.CreateProfile, src.CreateProfile, "wifi.create_profile")...)
	errs = append(errs, mergeKeyList(w.OpenCaptivePortal, src.OpenCaptivePortal, "wifi.open_network_login")...)
	errs = append(errs, mergeKeyList(w.EnableHotspot, src.EnableHotspot, "wifi.enable_hotspot")...)
	errs = append(errs, mergeKeyList(w.CreateHotspot, src.CreateHotspot, "wifi.create_hotspot")...)
	return errs
}

func (a *WifiAvailableKeys) merge(src *WifiAvailableKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(a.Connect, src.Connect, "wifi_available.connect")...)
	return errs
}

func (s *WifiSavedKeys) merge(src *WifiSavedKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(s.Edit, src.Edit, "wifi_saved.edit")...)
	errs = append(errs, mergeKeyList(s.Connect, src.Connect, "wifi_saved.connect")...)
	errs = append(errs, mergeKeyList(s.Disconnect, src.Disconnect, "wifi_saved.disconnect")...)
	errs = append(errs, mergeKeyList(s.Delete, src.Delete, "wifi_saved.delete")...)
	return errs
}

func mergeKeyList(dst *KeyBinding, src *KeyBinding, tag string) []error {
	if src == nil {
		return nil
	}

	var errs []error
	for _, v := range *src {
		if !validKeyName(v) {
			errs = append(errs, fmt.Errorf("invalid key %s: %q", tag, v))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	*dst = *src
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

func initKeyConfig() {
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
