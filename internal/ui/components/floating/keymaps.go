package floating

import "github.com/charmbracelet/bubbles/key"

type FloatingKeyMap struct {
	Quit key.Binding
}

func (k FloatingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k FloatingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

var defaultFloatingKeys = &FloatingKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/^c", "quit"),
	),
}
