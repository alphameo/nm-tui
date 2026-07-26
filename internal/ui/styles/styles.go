// Package styles provides styling for ui
package styles

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
)

var (
	TextColor   color.Color
	AccentColor color.Color
	MutedColor  color.Color
	ErrorColor  color.Color
	NotifColor  color.Color

	DefaultStyle lipgloss.Style
	AccentStyle  lipgloss.Style
	BoldStyle    lipgloss.Style

	Border               lipgloss.Border = lipgloss.RoundedBorder()
	BorderedStyle        lipgloss.Style
	BorderedFocusedStyle lipgloss.Style

	TableStyle     table.Styles
	DataTableStyle table.Styles

	TabViewStyles tabview.Styles

	OverlayStyle       lipgloss.Style
	NotifBorderedStyle lipgloss.Style

	ToggleStyle        lipgloss.Style
	ToggleFocusedStyle lipgloss.Style

	Spinner spinner.Spinner = spinner.Line
)

func Init(colors config.ColorConfig) error {
	err := convertColorConfig(colors)
	if err != nil {
		return err
	}

	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor)
	AccentStyle = lipgloss.NewStyle().Foreground(AccentColor)
	BoldStyle = DefaultStyle.Bold(true)

	BorderedStyle = DefaultStyle.Border(Border)
	BorderedFocusedStyle = lipgloss.NewStyle().Inherit(BorderedStyle).BorderForeground(AccentColor)
	BorderOffset = lipgloss.Width(Border.Left) * 2
	TabBarHeight = BorderOffset + 1

	TableStyle = tableStyle()
	DataTableStyle = dataTableStyle()

	TabViewStyles = *tabview.GenerateStyles(&BorderedStyle)

	OverlayStyle = DefaultStyle.
		Border(Border).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(2, 4).
		BorderForeground(AccentColor)

	NotifBorderedStyle = OverlayStyle.BorderForeground(NotifColor)

	ToggleStyle = DefaultStyle.Margin(0, 1)
	ToggleFocusedStyle = ToggleStyle.Foreground(AccentColor)

	SymbolColoredError = DefaultStyle.Foreground(ErrorColor).Render(SymbolError)

	return nil
}

func tableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(MutedColor).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(TextColor).
		Background(AccentColor).
		Bold(false)
	return style
}

func dataTableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(MutedColor).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(TextColor).
		Bold(false)
	return style
}

func convertColorConfig(colors config.ColorConfig) error {
	color, err := resolveCfgColor(*colors.Text)
	if err != nil {
		return err
	}
	TextColor = lipgloss.Color(color)

	color, err = resolveCfgColor(*colors.Accent)
	if err != nil {
		return err
	}
	AccentColor = lipgloss.Color(color)

	color, err = resolveCfgColor(*colors.Muted)
	if err != nil {
		return err
	}
	MutedColor = lipgloss.Color(color)

	color, err = resolveCfgColor(*colors.Error)
	if err != nil {
		return err
	}
	ErrorColor = lipgloss.Color(color)

	color, err = resolveCfgColor(*colors.Notif)
	if err != nil {
		return err
	}
	NotifColor = lipgloss.Color(color)

	return nil
}

func resolveCfgColor(cfgColor string) (string, error) {
	c := strings.ToLower(cfgColor)
	switch c {
	case config.CBlack:
		return "0", nil
	case config.CRed:
		return "1", nil
	case config.CGreen:
		return "2", nil
	case config.CYellow:
		return "3", nil
	case config.CBlue:
		return "4", nil
	case config.CMagenta:
		return "5", nil
	case config.CCyan:
		return "6", nil
	case config.CWhite:
		return "7", nil
	case config.CBrightBlack:
		return "8", nil
	case config.CBrightRed:
		return "9", nil
	case config.CBrightGreen:
		return "10", nil
	case config.CBrightYellow:
		return "11", nil
	case config.CBrightBlue:
		return "12", nil
	case config.CBrightMagenta:
		return "13", nil
	case config.CBrightCyan:
		return "14", nil
	case config.CBrightWhite:
		return "15", nil
	case config.CNone:
		return "", nil
	}

	if config.ValidHex(c) {
		return c, nil
	}

	return "", fmt.Errorf("color not resolved: %s", cfgColor)
}
