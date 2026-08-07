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
	Notif  *string `kdl:"notification"`
}

func DefaultColorConfig() *ColorConfig {
	return &ColorConfig{
		Text:   new(ColorNone),
		Accent: new(ColorBlue),
		Muted:  new(ColorBrightBlack),
		Error:  new(ColorRed),
		Notif:  new(ColorYellow),
	}
}

func (c *ColorConfig) merge(src *ColorConfig) []error {
	var errs []error
	collect := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	collect(mergeColor(c.Text, src.Text, "text"))
	collect(mergeColor(c.Accent, src.Accent, "accent"))
	collect(mergeColor(c.Error, src.Error, "error"))
	collect(mergeColor(c.Muted, src.Muted, "muted"))
	collect(mergeColor(c.Notif, src.Notif, "notification"))

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

	if err := validateColor(color); err != nil {
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
	ColorBlack         = "black"
	ColorRed           = "red"
	ColorGreen         = "green"
	ColorYellow        = "yellow"
	ColorBlue          = "blue"
	ColorMagenta       = "magenta"
	ColorCyan          = "cyan"
	ColorWhite         = "white"
	ColorBrightBlack   = "bright_black"
	ColorBrightRed     = "bright_red"
	ColorBrightGreen   = "bright_green"
	ColorBrightYellow  = "bright_yellow"
	ColorBrightBlue    = "bright_blue"
	ColorBrightMagenta = "bright_magenta"
	ColorBrightCyan    = "bright_cyan"
	ColorBrightWhite   = "bright_white"
	ColorNone          = "none"
)

func ValidCfgColor(color string) bool {
	switch color {
	case ColorBlack, ColorRed, ColorGreen, ColorYellow,
		ColorBlue, ColorMagenta, ColorCyan, ColorWhite,
		ColorBrightBlack, ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightMagenta, ColorBrightCyan, ColorBrightWhite,
		ColorNone:
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
