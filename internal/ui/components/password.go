package components

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type Password struct {
	*textinput.Model
}

func NewPassword(input *textinput.Model) Password {
	return Password{Model: input}
}

func DefaultPassword() Password {
	input := textinput.New()
	input.SetWidth(20)
	input.Prompt = ""
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'
	input.Placeholder = "Password"
	input.Validate = passwordValidator
	input.Err = passwordValidator(input.Value())
	return NewPassword(&input)
}

func (m *Password) View() string {
	view := m.Model.View()
	var style lipgloss.Style
	errIndicator := " "
	if m.Err != nil {
		errIndicator = styles.ErrorSymbolColored
	}
	if m.Focused() {
		style = styles.InputFieldFocusedStyle
	} else {
		style = styles.InputFieldStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Password ",
		view,
		errIndicator,
	)

	return view
}

var ErrPasswordFmt error = errors.New("wrong password format")

func passwordValidator(input string) error {
	if len(input) < 8 {
		return fmt.Errorf("%w: length < 8", ErrPasswordFmt)
	}
	return nil
}
