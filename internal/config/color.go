package config

import (
	"fmt"
	"strconv"
	"strings"
)

type ColorConfig struct {
	Text   *string `kdl:"text"`
	Accent *string `kdl:"accent"`
	Muted  *string `kdl:"muted"`
	Error  *string `kdl:"error"`
	Notif  *string `kdl:"notif"`
}

func DefaultColorConfig() *ColorConfig {
	text := CNone
	accent := CBlue
	muted := CBrightBlack
	error := CRed
	notif := CYellow

	return &ColorConfig{
		Text:   &text,
		Accent: &accent,
		Muted:  &muted,
		Error:  &error,
		Notif:  &notif,
	}
}

func (c *ColorConfig) merge(src *ColorConfig) []error {
	var errs []error

	errs = append(errs, mergeColor(c.Text, src.Text, "text"))
	errs = append(errs, mergeColor(c.Accent, src.Accent, "accent"))
	errs = append(errs, mergeColor(c.Error, src.Error, "error"))
	errs = append(errs, mergeColor(c.Muted, src.Muted, "muted"))
	errs = append(errs, mergeColor(c.Notif, src.Notif, "notif"))

	return errs
}

func mergeColor(dst *string, src *string, tag string) error {
	if src == nil {
		return nil
	}

	color := *src

	if color == defaultKeyword {
		return nil
	}

	err := validateColor(color)
	if err != nil {
		return fmt.Errorf("%s color: %w", tag, err)
	}
	*dst = *src
	return nil
}

func validateColor(color string) error {
	c := strings.ToLower(color)
	if c == defaultKeyword {
		return nil
	}

	if ValidHex(c) {
		return nil
	}
	if ValidCfgColor(c) {
		return nil
	}
	return fmt.Errorf("unknown color: %q", color)
}

const (
	CBlack         = "black"
	CRed           = "red"
	CGreen         = "green"
	CYellow        = "yellow"
	CBlue          = "blue"
	CMagenta       = "magenta"
	CCyan          = "cyan"
	CWhite         = "white"
	CBrightBlack   = "bright_black"
	CBrightRed     = "bright_red"
	CBrightGreen   = "bright_green"
	CBrightYellow  = "bright_yellow"
	CBrightBlue    = "bright_blue"
	CBrightMagenta = "bright_magenta"
	CBrightCyan    = "bright_cyan"
	CBrightWhite   = "bright_white"
	CNone          = "none"
)

func ValidCfgColor(color string) bool {
	switch color {
	case CBlack, CRed, CGreen, CYellow, CBlue, CMagenta, CCyan, CWhite,
		CBrightBlack, CBrightRed, CBrightGreen, CBrightYellow,
		CBrightBlue, CBrightMagenta, CBrightCyan, CBrightWhite,
		CNone:
		return true
	default:
		return false
	}
}

func ValidHex(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	_, err := strconv.ParseUint(color[1:], 16, 64)
	return err == nil
}
