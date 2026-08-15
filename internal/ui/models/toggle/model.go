// Package toggle provides toggling buttons
package toggle

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Styles struct {
	Focused lipgloss.Style
	Blured  lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{Focused: lipgloss.NewStyle(), Blured: lipgloss.NewStyle()}
}

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

	Styles Styles
}

func New() Model {
	return Model{
		value:   false,
		Symbols: DefaultSymbols(),
		Keys:    DefaultKeys(),
		Styles:  DefaultStyles(),
	}
}

func (t *Model) SetValue(value bool) {
	t.value = value
}

func (t *Model) Value() bool {
	return t.value
}

func (t Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !t.focus {
		return t, nil
	}

	if msg, ok := msg.(tea.KeyPressMsg); ok {
		if key.Matches(msg, t.Keys.Toggle) {
			t.value = !t.value
		}
	}
	return t, nil
}

func (t Model) View() string {
	var view string
	if t.value {
		view = t.Symbols.Activated
	} else {
		view = t.Symbols.Deactivated
	}
	style := *t.activeStyle()
	return style.Render(view)
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

func (t *Model) activeStyle() *lipgloss.Style {
	if t.focus {
		return &t.Styles.Focused
	}
	return &t.Styles.Blured
}
