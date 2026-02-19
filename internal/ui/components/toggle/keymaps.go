package toggle

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Toggle key.Binding
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Toggle}}
}

var defaultKeys = &KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "toggle"),
	),
}
