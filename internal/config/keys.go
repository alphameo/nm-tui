package config

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"codeberg.org/shimeoki/kdly"
)

type KeyBinding []string

func (k *KeyBinding) UnmarshalKDLNode(n *kdly.Node) error {
	kb := make(KeyBinding, 0, len(n.Entries))
	for e := range n.Arguments() {
		v, ok := e.Value.(kdly.String)
		if !ok {
			return fmt.Errorf("keybinding is strings only")
		}

		str, err := v.Resolve()
		if err != nil {
			return fmt.Errorf("invalid keybinding string: %w", err)
		}

		kb = append(kb, str)
	}

	*k = kb
	return nil
}

type KeyConfig struct {
	Toggle    *KeyBinding `kdly:"toggle"`
	Rescan    *KeyBinding `kdly:"rescan"`
	FocusNext *KeyBinding `kdly:"focus_next"`
	FocusPrev *KeyBinding `kdly:"focus_prev"`
	Focus1    *KeyBinding `kdly:"focus_1"`
	Focus2    *KeyBinding `kdly:"focus_2"`
	Focus3    *KeyBinding `kdly:"focus_3"`
	Focus4    *KeyBinding `kdly:"focus_4"`
	Focus5    *KeyBinding `kdly:"focus_5"`
	Focus6    *KeyBinding `kdly:"focus_6"`
	Focus7    *KeyBinding `kdly:"focus_7"`
	Focus8    *KeyBinding `kdly:"focus_8"`
	Focus9    *KeyBinding `kdly:"focus_9"`
	Focus10   *KeyBinding `kdly:"focus_10"`

	Main   *MainKeys   `kdly:"main"`
	Dialog *DialogKeys `kdly:"dialog"`

	Networks          *NetworksKeys          `kdly:"networks"`
	NetworkDevices    *NetworkDevicesKeys    `kdly:"network_devices"`
	AvailableNetworks *AvailableNetworksKeys `kdly:"available_networks"`
	NetworkProfiles   *NetworkProfilesKeys   `kdly:"network_profiles"`
}

type MainKeys struct {
	Help    *KeyBinding `kdly:"help"`
	TabNext *KeyBinding `kdly:"next_tab"`
	TabPrev *KeyBinding `kdly:"prev_tab"`
	Quit    *KeyBinding `kdly:"quit"`
}

type DialogKeys struct {
	TogglePWVisibility *KeyBinding `kdly:"toggle_pw_visibility"`
	Accept             *KeyBinding `kdly:"accept"`
	Close              *KeyBinding `kdly:"close"`
}

type NetworksKeys struct {
	CreateProfile     *KeyBinding `kdly:"create_profile"`
	OpenCaptivePortal *KeyBinding `kdly:"open_network_login"`
	QuickHotspot      *KeyBinding `kdly:"quick_hotspot"`
	CreateHotspot     *KeyBinding `kdly:"create_hotspot"`
}

type NetworkDevicesKeys struct {
	ShowInfo *KeyBinding `kdly:"show_info"`
}

type AvailableNetworksKeys struct {
	Connect    *KeyBinding `kdly:"connect"`
	Activate   *KeyBinding `kdly:"activate"`
	Deactivate *KeyBinding `kdly:"deactivate"`
}

type NetworkProfilesKeys struct {
	Edit       *KeyBinding `kdly:"edit"`
	Activate   *KeyBinding `kdly:"activate"`
	Deactivate *KeyBinding `kdly:"deactivate"`
	Delete     *KeyBinding `kdly:"delete"`
}

func DefaultKeys() *KeyConfig {
	return &KeyConfig{
		Toggle:    &KeyBinding{"space"},
		Rescan:    &KeyBinding{"r"},
		FocusNext: &KeyBinding{"tab"},
		FocusPrev: &KeyBinding{"shift+tab"},
		Focus1:    &KeyBinding{"1"},
		Focus2:    &KeyBinding{"2"},
		Focus3:    &KeyBinding{"3"},
		Focus4:    &KeyBinding{"4"},
		Focus5:    &KeyBinding{"5"},
		Focus6:    &KeyBinding{"6"},
		Focus7:    &KeyBinding{"7"},
		Focus8:    &KeyBinding{"8"},
		Focus9:    &KeyBinding{"9"},
		Focus10:   &KeyBinding{"0"},
		Main: &MainKeys{
			Help:    &KeyBinding{"?"},
			TabNext: &KeyBinding{"]"},
			TabPrev: &KeyBinding{"["},
			Quit:    &KeyBinding{"esc", "ctrl+c", "q", "ctrl+q"},
		},
		Dialog: &DialogKeys{
			TogglePWVisibility: &KeyBinding{"ctrl+p"},
			Accept:             &KeyBinding{"enter"},
			Close:              &KeyBinding{"esc", "ctrl+q", "ctrl+c"},
		},
		Networks: &NetworksKeys{
			CreateProfile:     &KeyBinding{"a", "c"},
			OpenCaptivePortal: &KeyBinding{"l"},
			QuickHotspot:      &KeyBinding{"ctrl+h"},
			CreateHotspot:     &KeyBinding{"h"},
		},
		NetworkDevices: &NetworkDevicesKeys{
			ShowInfo: &KeyBinding{"enter"},
		},
		AvailableNetworks: &AvailableNetworksKeys{
			Connect:    &KeyBinding{"enter"},
			Activate:   &KeyBinding{"space"},
			Deactivate: &KeyBinding{"ctrl+space"},
		},
		NetworkProfiles: &NetworkProfilesKeys{
			Edit:       &KeyBinding{"enter"},
			Activate:   &KeyBinding{"space"},
			Deactivate: &KeyBinding{"ctrl+space"},
			Delete:     &KeyBinding{"d", "delete"},
		},
	}
}

func (k *KeyConfig) Merge(src *KeyConfig) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&k.Toggle, src.Toggle, "toggle")...)
	errs = append(errs, MergeKeyList(&k.Rescan, src.Rescan, "rescan")...)
	errs = append(errs, MergeKeyList(&k.FocusNext, src.FocusNext, "focus_next")...)
	errs = append(errs, MergeKeyList(&k.FocusPrev, src.FocusPrev, "focus_prev")...)
	errs = append(errs, MergeKeyList(&k.Focus1, src.Focus1, "focus_1")...)
	errs = append(errs, MergeKeyList(&k.Focus2, src.Focus2, "focus_2")...)
	errs = append(errs, MergeKeyList(&k.Focus3, src.Focus3, "focus_3")...)
	errs = append(errs, MergeKeyList(&k.Focus4, src.Focus4, "focus_4")...)
	errs = append(errs, MergeKeyList(&k.Focus5, src.Focus5, "focus_5")...)
	errs = append(errs, MergeKeyList(&k.Focus6, src.Focus6, "focus_6")...)
	errs = append(errs, MergeKeyList(&k.Focus7, src.Focus7, "focus_7")...)
	errs = append(errs, MergeKeyList(&k.Focus8, src.Focus8, "focus_8")...)
	errs = append(errs, MergeKeyList(&k.Focus9, src.Focus9, "focus_9")...)
	errs = append(errs, MergeKeyList(&k.Focus10, src.Focus10, "focus_10")...)

	errs = append(errs, k.Main.Merge(src.Main)...)
	errs = append(errs, k.Dialog.Merge(src.Dialog)...)
	errs = append(errs, k.Networks.Merge(src.Networks)...)
	errs = append(errs, k.NetworkDevices.Merge(src.NetworkDevices)...)
	errs = append(errs, k.AvailableNetworks.Merge(src.AvailableNetworks)...)
	errs = append(errs, k.NetworkProfiles.Merge(src.NetworkProfiles)...)
	return errs
}

func (m *MainKeys) Merge(src *MainKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&m.Help, src.Help, "rescan")...)
	errs = append(errs, MergeKeyList(&m.TabNext, src.TabNext, "main.next_tab")...)
	errs = append(errs, MergeKeyList(&m.TabPrev, src.TabPrev, "main.prev_tab")...)
	errs = append(errs, MergeKeyList(&m.Quit, src.Quit, "main.quit")...)
	return errs
}

func (d *DialogKeys) Merge(src *DialogKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&d.TogglePWVisibility, src.TogglePWVisibility, "dialog.toggle_pw_visibility")...)
	errs = append(errs, MergeKeyList(&d.Accept, src.Accept, "dialog.accept")...)
	errs = append(errs, MergeKeyList(&d.Close, src.Close, "dialog.close")...)
	return errs
}

func (w *NetworksKeys) Merge(src *NetworksKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&w.CreateProfile, src.CreateProfile, "networks.create_profile")...)
	errs = append(errs, MergeKeyList(&w.OpenCaptivePortal, src.OpenCaptivePortal, "networks.open_network_login")...)
	errs = append(errs, MergeKeyList(&w.QuickHotspot, src.QuickHotspot, "networks.quick_hotspot")...)
	errs = append(errs, MergeKeyList(&w.CreateHotspot, src.CreateHotspot, "networks.create_hotspot")...)
	return errs
}

func (m *NetworkDevicesKeys) Merge(src *NetworkDevicesKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&m.ShowInfo, src.ShowInfo, "show_info")...)
	return errs
}

func (a *AvailableNetworksKeys) Merge(src *AvailableNetworksKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&a.Connect, src.Connect, "available_networks.connect")...)
	errs = append(errs, MergeKeyList(&a.Activate, src.Activate, "available_networks.connect")...)
	errs = append(errs, MergeKeyList(&a.Deactivate, src.Deactivate, "available_networks.connect")...)
	return errs
}

func (s *NetworkProfilesKeys) Merge(src *NetworkProfilesKeys) []error {
	if src == nil {
		return nil
	}

	var errs []error
	errs = append(errs, MergeKeyList(&s.Edit, src.Edit, "network_profiles.edit")...)
	errs = append(errs, MergeKeyList(&s.Activate, src.Activate, "network_profiles.activate")...)
	errs = append(errs, MergeKeyList(&s.Deactivate, src.Deactivate, "network_profiles.deactivate")...)
	errs = append(errs, MergeKeyList(&s.Delete, src.Delete, "network_profiles.delete")...)
	return errs
}

func MergeKeyList(dst **KeyBinding, src *KeyBinding, tag string) []error {
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

func buildValidKey() map[string]bool {
	m := map[string]bool{
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

	for i := 1; i <= 63; i++ {
		m[fmt.Sprintf("f%d", i)] = true
	}

	return m
}

var validKey = buildValidKey()

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
