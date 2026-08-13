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
	FocusNext     *KeyBinding `kdl:"focus_next"`
	FocusPrev     *KeyBinding `kdl:"focus_prev"`
	Focus1        *KeyBinding `kdl:"focus_1"`
	Focus2        *KeyBinding `kdl:"focus_2"`
	Focus3        *KeyBinding `kdl:"focus_3"`
	Focus4        *KeyBinding `kdl:"focus_4"`
	Focus5        *KeyBinding `kdl:"focus_5"`
	Focus6        *KeyBinding `kdl:"focus_6"`
	Focus7        *KeyBinding `kdl:"focus_7"`
	Focus8        *KeyBinding `kdl:"focus_8"`
	Focus9        *KeyBinding `kdl:"focus_9"`
	Focus10       *KeyBinding `kdl:"focus_10"`

	Main   *MainKeys   `kdl:"main"`
	Dialog *DialogKeys `kdl:"dialog"`

	Networks          *NetworksKeys          `kdl:"networks"`
	AvailableNetworks *AvailableNetworksKeys `kdl:"available_networks"`
	NetworkProfiles   *NetworkProfilesKeys   `kdl:"network_profiles"`
}

type MainKeys struct {
	TabNext *KeyBinding `kdl:"next_tab"`
	TabPrev *KeyBinding `kdl:"prev_tab"`
	Quit    *KeyBinding `kdl:"quit"`
}

type DialogKeys struct {
	TogglePWVisibility *KeyBinding `kdl:"toggle_pw_visibility"`
	Accept             *KeyBinding `kdl:"accept"`
	Close              *KeyBinding `kdl:"close"`
}

type NetworksKeys struct {
	CreateProfile     *KeyBinding `kdl:"create_profile"`
	OpenCaptivePortal *KeyBinding `kdl:"open_network_login"`
	QuickHotspot      *KeyBinding `kdl:"quick_hotspot"`
	CreateHotspot     *KeyBinding `kdl:"create_hotspot"`
}

type AvailableNetworksKeys struct {
	Connect *KeyBinding `kdl:"connect"`
}

type NetworkProfilesKeys struct {
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
		FocusNext:     keyBinding("tab"),
		FocusPrev:     keyBinding("shift+tab"),
		Focus1:        keyBinding("1"),
		Focus2:        keyBinding("2"),
		Focus3:        keyBinding("3"),
		Focus4:        keyBinding("4"),
		Focus5:        keyBinding("5"),
		Focus6:        keyBinding("6"),
		Focus7:        keyBinding("7"),
		Focus8:        keyBinding("8"),
		Focus9:        keyBinding("9"),
		Focus10:       keyBinding("10"),
		Main: &MainKeys{
			TabNext: keyBinding("]"),
			TabPrev: keyBinding("["),
			Quit:    keyBinding("esc", "ctrl+c", "q", "ctrl+q"),
		},
		Dialog: &DialogKeys{
			TogglePWVisibility: keyBinding("ctrl+p"),
			Accept:             keyBinding("enter"),
			Close:              keyBinding("esc", "ctrl+q", "ctrl+c"),
		},
		Networks: &NetworksKeys{
			CreateProfile:     keyBinding("a", "c"),
			OpenCaptivePortal: keyBinding("l"),
			QuickHotspot:      keyBinding("ctrl+h"),
			CreateHotspot:     keyBinding("h"),
		},
		AvailableNetworks: &AvailableNetworksKeys{
			Connect: keyBinding("enter"),
		},
		NetworkProfiles: &NetworkProfilesKeys{
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
	errs = append(errs, mergeKeyList(&k.Toggle, src.Toggle, "toggle")...)
	errs = append(errs, mergeKeyList(&k.Rescan, src.Rescan, "rescan")...)
	errs = append(errs, mergeKeyList(&k.RescanFocused, src.RescanFocused, "rescan_focused")...)
	errs = append(errs, mergeKeyList(&k.FocusNext, src.FocusNext, "focus_next")...)
	errs = append(errs, mergeKeyList(&k.FocusPrev, src.FocusPrev, "focus_prev")...)
	errs = append(errs, mergeKeyList(&k.Focus1, src.Focus1, "focus_1")...)
	errs = append(errs, mergeKeyList(&k.Focus2, src.Focus2, "focus_2")...)
	errs = append(errs, mergeKeyList(&k.Focus3, src.Focus3, "focus_3")...)
	errs = append(errs, mergeKeyList(&k.Focus4, src.Focus4, "focus_4")...)
	errs = append(errs, mergeKeyList(&k.Focus5, src.Focus5, "focus_5")...)
	errs = append(errs, mergeKeyList(&k.Focus6, src.Focus6, "focus_6")...)
	errs = append(errs, mergeKeyList(&k.Focus7, src.Focus7, "focus_7")...)
	errs = append(errs, mergeKeyList(&k.Focus8, src.Focus8, "focus_8")...)
	errs = append(errs, mergeKeyList(&k.Focus9, src.Focus9, "focus_9")...)
	errs = append(errs, mergeKeyList(&k.Focus10, src.Focus10, "focus_10")...)

	errs = append(errs, k.Main.merge(src.Main)...)
	errs = append(errs, k.Dialog.merge(src.Dialog)...)
	errs = append(errs, k.Networks.merge(src.Networks)...)
	errs = append(errs, k.AvailableNetworks.merge(src.AvailableNetworks)...)
	errs = append(errs, k.NetworkProfiles.merge(src.NetworkProfiles)...)
	return errs
}

func (m *MainKeys) merge(src *MainKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(&m.TabNext, src.TabNext, "main.next_tab")...)
	errs = append(errs, mergeKeyList(&m.TabPrev, src.TabPrev, "main.prev_tab")...)
	errs = append(errs, mergeKeyList(&m.Quit, src.Quit, "main.quit")...)
	return errs
}

func (d *DialogKeys) merge(src *DialogKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(&d.TogglePWVisibility, src.TogglePWVisibility, "dialog.toggle_pw_visibility")...)
	errs = append(errs, mergeKeyList(&d.Accept, src.Accept, "dialog.accept")...)
	errs = append(errs, mergeKeyList(&d.Close, src.Close, "dialog.close")...)
	return errs
}

func (w *NetworksKeys) merge(src *NetworksKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(&w.CreateProfile, src.CreateProfile, "networks.create_profile")...)
	errs = append(errs, mergeKeyList(&w.OpenCaptivePortal, src.OpenCaptivePortal, "networks.open_network_login")...)
	errs = append(errs, mergeKeyList(&w.QuickHotspot, src.QuickHotspot, "networks.quick_hotspot")...)
	errs = append(errs, mergeKeyList(&w.CreateHotspot, src.CreateHotspot, "networks.create_hotspot")...)
	return errs
}

func (a *AvailableNetworksKeys) merge(src *AvailableNetworksKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(&a.Connect, src.Connect, "available_networks.connect")...)
	return errs
}

func (s *NetworkProfilesKeys) merge(src *NetworkProfilesKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, mergeKeyList(&s.Edit, src.Edit, "network_profiles.edit")...)
	errs = append(errs, mergeKeyList(&s.Connect, src.Connect, "network_profiles.connect")...)
	errs = append(errs, mergeKeyList(&s.Disconnect, src.Disconnect, "network_profiles.disconnect")...)
	errs = append(errs, mergeKeyList(&s.Delete, src.Delete, "network_profiles.delete")...)
	return errs
}

func mergeKeyList(dst **KeyBinding, src *KeyBinding, tag string) []error {
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

	*dst = src
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
