package toggle

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Toggle key.Binding
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Toggle}}
}

var DefaultKeys = &KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("󱁐", "toggle"),
	),
}
