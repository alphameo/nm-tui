package tabview

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	TabNext key.Binding
	TabPrev key.Binding
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.TabNext, k.TabPrev}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.TabNext, k.TabPrev}}
}

func DefaultKeys() KeyMap {
	return KeyMap{
		TabNext: key.NewBinding(
			key.WithKeys("]"),
			key.WithHelp("]", "next tab"),
		),
		TabPrev: key.NewBinding(
			key.WithKeys("["),
			key.WithHelp("[", "previous tab"),
		),
	}
}
