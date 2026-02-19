package components

import (
	"github.com/charmbracelet/bubbles/key"
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

var defaultToggleKeys = &toggleKeyMap{
	toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "toggle"),
	),
}

type toggleKeyMap struct {
	toggle key.Binding
}

func (k *toggleKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggle}
}

func (k *toggleKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.toggle}}
}

type ToggleModel struct {
	value   bool
	Focused bool
	Symbols *ToggleModelSymbols

	Keys *toggleKeyMap
}

func NewToggleModel(initial bool) *ToggleModel {
	return &ToggleModel{
		value: initial,
		Keys:  defaultToggleKeys,
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
		if key.Matches(msg, t.Keys.toggle) {
			t.value = !t.value
		}
	}
	return t, nil
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
