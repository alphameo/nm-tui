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

type ToggleModel struct {
	value   bool
	Focused bool
	Symbols *Symbols

	Keys *KeyMap
}

func New(initial bool) *ToggleModel {
	return &ToggleModel{
		value: initial,
		Keys:  defaultKeys,
	}
}

func (t *ToggleModel) SetValue(value bool) {
	t.value = value
}

func (t *ToggleModel) Value() bool {
	return t.value
}

func (t *ToggleModel) Init() tea.Cmd {
	return nil
}

func (t *ToggleModel) Update(msg tea.Msg) (*ToggleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, t.Keys.Toggle) {
			t.value = !t.value
		}
	}
	return t, nil
}

func (t *ToggleModel) View() string {
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

func (t *ToggleModel) Focus() tea.Cmd {
	t.Focused = true
	return nil
}

func (t *ToggleModel) Blur() {
	t.Focused = false
}
