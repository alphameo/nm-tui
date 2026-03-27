// Package toggle provides toggling buttons
package toggle

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Symbols struct {
	Activated   string
	Deactivated string
}

var DefaultSymbols = &Symbols{
	Activated:   "[⏺]",
	Deactivated: "[ ]",
}

type Model struct {
	value   bool
	focus   bool
	Symbols *Symbols

	Keys *KeyMap
}

func New(initial bool) *Model {
	return &Model{
		value: initial,
		Keys:  defaultKeys,
	}
}

func (t *Model) SetValue(value bool) {
	t.value = value
}

func (t *Model) Value() bool {
	return t.value
}

func (t *Model) Init() tea.Cmd {
	return nil
}

func (t *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, t.Keys.Toggle) {
			t.value = !t.value
		}
	}
	return t, nil
}

func (t *Model) View() string {
	symbols := t.Symbols
	if symbols == nil {
		symbols = DefaultSymbols
	}
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
