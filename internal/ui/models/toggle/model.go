// Package toggle provides toggling buttons
package toggle

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type Symbols struct {
	Activated   string
	Deactivated string
}

func DefaultSymbols() Symbols {
	return Symbols{
		Activated:   "[x]",
		Deactivated: "[ ]",
	}
}

type Model struct {
	value   bool
	focus   bool
	Symbols Symbols

	Keys KeyMap
}

func New() *Model {
	return &Model{
		value:   false,
		Symbols: DefaultSymbols(),
		Keys:    DefaultKeys(),
	}
}

func (t *Model) SetValue(value bool) {
	t.value = value
}

func (t *Model) Value() bool {
	return t.value
}

func (t *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, t.Keys.Toggle) {
			t.value = !t.value
		}
	}
	return t, nil
}

func (t *Model) View() string {
	symbols := t.Symbols
	if t.value {
		return symbols.Activated
	} else {
		return symbols.Deactivated
	}
}

func (t *Model) Focus() tea.Cmd {
	t.focus = true
	return nil
}

func (t *Model) Blur() {
	t.focus = false
}

func (t *Model) Focused() bool {
	return t.focus
}
