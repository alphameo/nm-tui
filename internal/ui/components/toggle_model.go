package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ToggleModelSymbols struct {
	Activated   string
	Deactivated string
}

var DefaultToggleModelSymbols = &ToggleModelSymbols{
	Activated:   "[⏺]",
	Deactivated: "[ ]",
}

type ToggleModel struct {
	value   bool
	Focused bool
	Symbols *ToggleModelSymbols
}

func NewToggleModel(initial bool) *ToggleModel {
	return &ToggleModel{value: initial}
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
		if msg.String() == " " {
			t.value = !t.value
		}
	}
	return t, nil
}

func RenderCheckbox(value bool) string {
	if value {
		return "[⏺]"
	} else {
		return "[ ]"
	}
}

func (t *ToggleModel) View() string {
	symbols := t.Symbols
	if symbols == nil {
		symbols = DefaultToggleModelSymbols
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
