package models

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

func newDefaultInput() textinput.Model {
	input := textinput.New()
	input.SetStyles(styles.InputStyles)
	input.SetWidth(20)
	input.Prompt = ""
	return input
}

func newDefaultSSIDInput() textinput.Model {
	ssid := newDefaultInput()
	ssid.Placeholder = "SSID"
	return ssid
}

func newDefaultNameInput() textinput.Model {
	ssid := newDefaultInput()
	ssid.Placeholder = "Name"
	return ssid
}

func newDefaultPasswordInput() textinput.Model {
	pw := newDefaultInput()
	pw.EchoMode = textinput.EchoPassword
	pw.EchoCharacter = styles.SymbolPwHiddenChar
	pw.Placeholder = "Password"
	pw.Validate = passwordValidator
	pw.Err = passwordValidator(pw.Value())
	return pw
}

func newDefaultToggle() toggle.Model {
	t := toggle.New()
	t.Styles = styles.ToggleStyles
	t.Symbols = styles.SymbolsToggle
	return t
}

func newDefaultSpinner() spinner.Model {
	s := spinner.New()
	s.Style = styles.DefaultStyle
	s.Spinner = styles.Spinner
	return s
}

var ErrPasswordFmt = errors.New("wrong password format")

func passwordValidator(input string) error {
	l := len(input)
	if l > 0 && l < 8 {
		return fmt.Errorf("%w: length < 8", ErrPasswordFmt)
	}
	return nil
}

func newDefaultViewport() viewport.Model {
	viewport := viewport.New()
	viewport.Style = styles.DefaultStyle
	return viewport
}
