package tabview

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Next key.Binding
	Prev key.Binding
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Next, k.Prev}}
}

func DefaultKeys() KeyMap {
	return KeyMap{
		Next: key.NewBinding(
			key.WithKeys("]"),
			key.WithHelp("]", "next tab"),
		),
		Prev: key.NewBinding(
			key.WithKeys("["),
			key.WithHelp("[", "previous tab"),
		),
	}
}
