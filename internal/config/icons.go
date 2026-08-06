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

type IconConfig struct {
	NerdPreset   *bool   `kdl:"nerd_preset"`
	BorderStyle  *string `kdl:"border_style"`
	SpinnerStyle *string `kdl:"spinner_style"`
	ToggleOff    *string `kdl:"toggle_off"`
	ToggleOn     *string `kdl:"toggle_on"`
	PwHiddenChar *string `kdl:"password_hidden_character"`
	Error        *string `kdl:"error"`
	Check        *string `kdl:"check"`
	Connection   *string `kdl:"connection"`
	Signal       *string `kdl:"signal"`
	AccessPoint  *string `kdl:"access_point"`
	Infra        *string `kdl:"infra"`
	Mesh         *string `kdl:"mesh"`
	AdHoc        *string `kdl:"ad_hoc"`
}

func DefaultNerdIconConfig() *IconConfig {
	nerd := true
	border := BorderRounded
	spinner := SpinnerMeter
	toggleOff := " "
	toggleOn := " "
	pwHiddenChar := "•"
	err := "✗"
	check := ""
	signal := ""
	connection := "󱘖"
	accessPoint := "󰀃"
	infra := "🖳"
	mesh := ""
	adHoc := ""
	return &IconConfig{
		NerdPreset:   &nerd,
		BorderStyle:  &border,
		SpinnerStyle: &spinner,
		ToggleOff:    &toggleOff,
		ToggleOn:     &toggleOn,
		PwHiddenChar: &pwHiddenChar,
		Error:        &err,
		Check:        &check,
		Connection:   &connection,
		Signal:       &signal,
		AccessPoint:  &accessPoint,
		Infra:        &infra,
		Mesh:         &mesh,
		AdHoc:        &adHoc,
	}
}

func DefaultNonNerdIconConfig() *IconConfig {
	nerd := false
	border := BorderASCII
	spinner := SpinnerLine
	toggleOff := "[ ]"
	toggleOn := "[x]"
	pwHiddenChar := "*"
	err := "!"
	check := "v"
	signal := "sig"
	connection := "con"
	accessPoint := "ap"
	infra := "infr"
	mesh := "#"
	adHoc := "ah"
	return &IconConfig{
		NerdPreset:   &nerd,
		BorderStyle:  &border,
		SpinnerStyle: &spinner,
		ToggleOff:    &toggleOff,
		ToggleOn:     &toggleOn,
		PwHiddenChar: &pwHiddenChar,
		Error:        &err,
		Check:        &check,
		Connection:   &connection,
		Signal:       &signal,
		AccessPoint:  &accessPoint,
		Infra:        &infra,
		Mesh:         &mesh,
		AdHoc:        &adHoc,
	}
}

func DefaultIconConfig() *IconConfig {
	return DefaultNonNerdIconConfig()
}

func (c *IconConfig) merge(src *IconConfig) []error {
	var errs []error
	collect := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	collect(mergeBorderStyle(c.BorderStyle, src.BorderStyle))
	collect(mergeSpinnerStyle(c.SpinnerStyle, src.SpinnerStyle))
	collect(mergeSymbol(c.PwHiddenChar, src.PwHiddenChar, "password_symbol"))
	collect(mergeIcon(c.ToggleOff, src.ToggleOff, "toggle_off"))
	collect(mergeIcon(c.ToggleOn, src.ToggleOn, "toggle_on"))
	collect(mergeIcon(c.Error, src.Error, "error"))
	collect(mergeIcon(c.Check, src.Check, "check"))
	collect(mergeIcon(c.Connection, src.Connection, "connection"))
	collect(mergeIcon(c.Signal, src.Signal, "signal"))
	collect(mergeIcon(c.AccessPoint, src.AccessPoint, "access_point"))
	collect(mergeIcon(c.Infra, src.Infra, "infra"))
	collect(mergeIcon(c.Mesh, src.Mesh, "mesh"))
	collect(mergeIcon(c.AdHoc, src.AdHoc, "ad_hoc"))

	return errs
}

func mergeBorderStyle(dst *string, src *string) error {
	if src == nil {
		return nil
	}

	border := *src

	if border == defaultKeyword {
		return nil
	}

	err := validateBorderStyle(border)
	if err != nil {
		return fmt.Errorf("border style: %w", err)
	}

	*dst = *src
	return nil
}

func mergeSpinnerStyle(dst *string, src *string) error {
	if src == nil {
		return nil
	}

	spinner := *src

	if spinner == defaultKeyword {
		return nil
	}

	err := validateSpinnerStyle(spinner)
	if err != nil {
		return fmt.Errorf("spinner style: %w", err)
	}

	*dst = *src
	return nil
}

func mergeIcon(dst *string, src *string, tag string) error {
	if src == nil {
		return nil
	}

	icon := *src

	if icon == defaultKeyword {
		return nil
	}

	if len(icon) == 0 {
		return fmt.Errorf("empty %s icon", tag)
	}
	*dst = *src
	return nil
}

func mergeSymbol(dst *string, src *string, tag string) error {
	if src == nil {
		return nil
	}

	symbol := *src

	if symbol == defaultKeyword {
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
