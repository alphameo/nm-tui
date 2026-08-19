package config

import (
	"fmt"
	"unicode/utf8"
)

const (
	BorderASCII    = "ascii"
	BorderMarkdown = "markdown"

	BorderRounded        = "rounded"
	BorderSquare         = "square"
	BorderThickSquare    = "thick_square"
	BorderDoubleSquare   = "double_square"
	BorderBlock          = "block"
	BorderOuterHalfBlock = "outer_half_block"
	BorderInnerHalfBlock = "inner_half_block"
)

const (
	SpinnerLine     = "line"
	SpinnerEllipsis = "ellipsis"

	SpinnerDot       = "dot"
	SpinnerMiniDot   = "mini_dot"
	SpinnerJump      = "jump"
	SpinnerPulse     = "pulse"
	SpinnerPoints    = "points"
	SpinnerMeter     = "meter"
	SpinnerHamburger = "hamburger"
)

const (
	CursorBar       = "bar"
	CursorUnderline = "underline"
	CursorBlock     = "block"
)

type IconConfig struct {
	NerdPreset       *bool   `kdl:"nerd_preset"`
	BorderStyle      *string `kdl:"border_style"`
	SpinnerStyle     *string `kdl:"spinner_style"`
	InputCursorShape *string `kdl:"input_cursor_shape"`
	ToggleOff        *string `kdl:"toggle_off"`
	ToggleOn         *string `kdl:"toggle_on"`
	PwHiddenChar     *string `kdl:"password_hidden_character"`
	Error            *string `kdl:"error"`
	Check            *string `kdl:"check"`
	Connection       *string `kdl:"connection"`
	Signal           *string `kdl:"signal"`
	Saved            *string `kdl:"saved"`
	Available        *string `kdl:"available"`
	AccessPoint      *string `kdl:"access_point"`
	Infra            *string `kdl:"infra"`
	Mesh             *string `kdl:"mesh"`
	AdHoc            *string `kdl:"ad_hoc"`
	Ellipsis         *string `kdl:"ellipsis"`
	Separator        *string `kdl:"separator"`
}

func DefaultNerdIconConfig() *IconConfig {
	return &IconConfig{
		NerdPreset:       new(true),
		BorderStyle:      new(BorderRounded),
		SpinnerStyle:     new(SpinnerMeter),
		InputCursorShape: new(CursorBar),
		ToggleOff:        new(" "),
		ToggleOn:         new(" "),
		PwHiddenChar:     new("•"),
		Error:            new("✗"),
		Check:            new(" "),
		Connection:       new("󱘖 "),
		Signal:           new(" "),
		Saved:            new(" "),
		Available:        new("⭗ "),
		AccessPoint:      new("󰀃 "),
		Infra:            new("🖳 "),
		Mesh:             new(" "),
		AdHoc:            new(""),
		Separator:        new("•"),
		Ellipsis:         new("…"),
	}
}

func DefaultNonNerdIconConfig() *IconConfig {
	return &IconConfig{
		NerdPreset:       new(false),
		BorderStyle:      new(BorderASCII),
		SpinnerStyle:     new(SpinnerLine),
		InputCursorShape: new(CursorBar),
		ToggleOff:        new("[ ]"),
		ToggleOn:         new("[x]"),
		PwHiddenChar:     new("*"),
		Error:            new("!"),
		Check:            new("v"),
		Connection:       new("con"),
		Signal:           new("sig"),
		Saved:            new("sav"),
		Available:        new("o"),
		AccessPoint:      new("ap"),
		Infra:            new("infr"),
		Mesh:             new("#"),
		AdHoc:            new("ah"),
		Separator:        new("|"),
		Ellipsis:         new("_"),
	}
}

func DefaultIconConfig() *IconConfig {
	return DefaultNonNerdIconConfig()
}

func (c *IconConfig) Merge(src *IconConfig) []error {
	var errs []error
	collect := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	collect(mergeBorderStyle(c.BorderStyle, src.BorderStyle))
	collect(mergeSpinnerStyle(c.SpinnerStyle, src.SpinnerStyle))
	collect(mergeCursorShape(c.InputCursorShape, src.InputCursorShape))
	collect(mergeSymbol(c.PwHiddenChar, src.PwHiddenChar, "password_symbol"))
	collect(mergeIcon(c.ToggleOff, src.ToggleOff, "toggle_off"))
	collect(mergeIcon(c.ToggleOn, src.ToggleOn, "toggle_on"))
	collect(mergeIcon(c.Error, src.Error, "error"))
	collect(mergeIcon(c.Check, src.Check, "check"))
	collect(mergeIcon(c.Connection, src.Connection, "connection"))
	collect(mergeIcon(c.Signal, src.Signal, "signal"))
	collect(mergeIcon(c.Saved, src.Saved, "saved"))
	collect(mergeIcon(c.Available, src.Available, "saved"))
	collect(mergeIcon(c.AccessPoint, src.AccessPoint, "access_point"))
	collect(mergeIcon(c.Infra, src.Infra, "infra"))
	collect(mergeIcon(c.Mesh, src.Mesh, "mesh"))
	collect(mergeIcon(c.AdHoc, src.AdHoc, "ad_hoc"))
	collect(mergeIcon(c.Ellipsis, src.Ellipsis, "ellipsis"))
	collect(mergeIcon(c.Separator, src.Separator, "separator"))

	return errs
}

func mergeBorderStyle(dst *string, src *string) error {
	if src == nil {
		return nil
	}

	border := *src

	if border == DefaultKeyword {
		return nil
	}

	if err := validateBorderStyle(border); err != nil {
		return fmt.Errorf("border_style: %w", err)
	}

	*dst = *src
	return nil
}

func mergeSpinnerStyle(dst *string, src *string) error {
	if src == nil {
		return nil
	}

	spinner := *src

	if spinner == DefaultKeyword {
		return nil
	}

	if err := validateSpinnerStyle(spinner); err != nil {
		return fmt.Errorf("spinner_style: %w", err)
	}

	*dst = *src
	return nil
}

func mergeCursorShape(dst *string, src *string) error {
	if src == nil {
		return nil
	}

	cursor := *src

	if cursor == DefaultKeyword {
		return nil
	}

	if err := validateCursorShape(cursor); err != nil {
		return fmt.Errorf("cursor_shape: %w", err)
	}

	*dst = *src
	return nil
}

func mergeIcon(dst *string, src *string, tag string) error {
	if src == nil {
		return nil
	}

	icon := *src

	if icon == DefaultKeyword {
		return nil
	}

	if len(icon) == 0 {
		return fmt.Errorf("empty %s icon", tag)
	}
	*dst = *src
	return nil
}

//nolint:unparam // NOTE: mergeSymbol is kept parameterized for future single-rune symbol merges
func mergeSymbol(dst *string, src *string, tag string) error {
	if src == nil {
		return nil
	}

	symbol := *src

	if symbol == DefaultKeyword {
		return nil
	}

	if len(symbol) == 0 {
		return fmt.Errorf("empty %s symbol", tag)
	}
	if utf8.RuneCountInString(symbol) > 1 {
		return fmt.Errorf("length of symbol %s > 1: %q", tag, symbol)
	}

	*dst = *src
	return nil
}

func validateBorderStyle(border string) error {
	switch border {
	case BorderASCII, BorderMarkdown,
		BorderRounded,
		BorderSquare, BorderThickSquare, BorderDoubleSquare,
		BorderBlock, BorderOuterHalfBlock, BorderInnerHalfBlock:
		return nil
	default:
		return fmt.Errorf("unknown border style: %q", border)
	}
}

func validateSpinnerStyle(spinner string) error {
	switch spinner {
	case SpinnerLine, SpinnerEllipsis,

		SpinnerDot, SpinnerMiniDot,
		SpinnerJump, SpinnerPulse, SpinnerPoints, SpinnerMeter, SpinnerHamburger:
		return nil
	default:
		return fmt.Errorf("unknown spinner style: %q", spinner)
	}
}

func validateCursorShape(cursor string) error {
	switch cursor {
	case CursorBar, CursorUnderline, CursorBlock:
		return nil
	default:
		return fmt.Errorf("unknown cursor shape style: %q", cursor)
	}
}
