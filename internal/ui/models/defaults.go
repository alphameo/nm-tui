package models

import (
	"charm.land/bubbles/v2/textinput"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

func newDefaultInput() textinput.Model {
	input := textinput.New()
	input.SetStyles(styles.InputStyle)
	input.SetWidth(20)
	input.Prompt = ""
	return input
}

func newDefaultPassword() textinput.Model {
	pw := newDefaultInput()
	pw.EchoMode = textinput.EchoPassword
	pw.EchoCharacter = styles.SymbolPwHiddenChar
	return pw
}
