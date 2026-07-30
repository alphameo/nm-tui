// Package styles provides styling for ui
package styles

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

var (
	TextColor   color.Color
	BgColor     color.Color
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

	InputCursor      tea.CursorShape
	InputCursorBlink bool
	InputStyle       textinput.Styles

	TabViewStyles tabview.Styles

	OverlayStyle       lipgloss.Style
	NotifBorderedStyle lipgloss.Style

	ToggleStyle        lipgloss.Style
	ToggleFocusedStyle lipgloss.Style

	Spinner spinner.Spinner = spinner.Line
)

func Init(cfg config.Config) error {
	err := initIcons(*cfg.Icons)
	if err != nil {
		return err
	}
	err = initColors(*cfg.Colors)
	if err != nil {
		return err
	}

	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor).Background(BgColor)
	AccentStyle = DefaultStyle.Foreground(AccentColor).Bold(true)
	BoldStyle = DefaultStyle.Bold(true)

	BorderedStyle = DefaultStyle.Border(Border).BorderForeground(TextColor).BorderBackground(BgColor)
	BorderedFocusedStyle = BorderedStyle.BorderForeground(AccentColor)
	BorderOffset = lipgloss.Width(Border.Left) * 2
	TabBarHeight = BorderOffset + 1

	TableStyle = tableStyle()
	DataTableStyle = dataTableStyle()

	InputCursor = tea.CursorBar
	InputCursorBlink = true
	InputStyle = inputStyle()

	TabViewStyles = *tabview.GenerateStyles(&BorderedStyle)

	OverlayStyle = DefaultStyle.
		Border(Border).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(2, 4).
		BorderForeground(AccentColor).
		BorderBackground(BgColor)

	NotifBorderedStyle = OverlayStyle.BorderForeground(NotifColor)

	ToggleStyle = DefaultStyle.Margin(0, 1)
	ToggleFocusedStyle = ToggleStyle.Foreground(AccentColor)

	SymbolColoredError = DefaultStyle.Foreground(ErrorColor).Render(SymbolError)

	return nil
}

func tableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(Border).
		BorderForeground(MutedColor).
		BorderBackground(BgColor).
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
		BorderStyle(Border).
		BorderForeground(MutedColor).
		BorderBackground(BgColor).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(TextColor).
		Background(BgColor).
		Bold(false)
	return style
}

func inputStyle() textinput.Styles {
	return textinput.Styles{
		Focused: textinput.StyleState{
			Placeholder: DefaultStyle,
			Suggestion:  DefaultStyle,
			Prompt:      DefaultStyle,
			Text:        DefaultStyle,
		},
		Blurred: textinput.StyleState{
			Placeholder: DefaultStyle,
			Suggestion:  DefaultStyle,
			Prompt:      DefaultStyle,
			Text:        DefaultStyle,
		},
		Cursor: textinput.CursorStyle{
			Color: TextColor,
			Shape: InputCursor,
			Blink: InputCursorBlink,
		},
	}
}

func initIcons(icons config.IconConfig) error {
	border, err := resolveBorderStyle(*icons.BorderStyle)
	if err != nil {
		return err
	}
	Border = border

	spinner, err := resolveSpinnerStyle(*icons.SpinnerStyle)
	if err != nil {
		return err
	}
	Spinner = spinner

	SymbolsToggle = toggle.Symbols{Activated: *icons.ToggleOn, Deactivated: *icons.ToggleOff}
	SymbolPwHiddenChar = []rune(*icons.PwHiddenChar)[0]
	SymbolError = *icons.Error
	SymbolCheck = *icons.Check
	SymbolConnection = *icons.Connection
	SymbolSignal = *icons.Signal
	SymbolAccessPoint = *icons.AccessPoint
	SymbolInfra = *icons.Infra
	SymbolMesh = *icons.Mesh
	SymbolAdHoc = *icons.AdHoc

	return nil
}

func initColors(colors config.ColorConfig) error {
	color, err := resolveCfgColor(*colors.Text)
	if err != nil {
		return err
	}
	TextColor = lipgloss.Color(color)

	color, err = resolveCfgColor(*colors.Background)
	if err != nil {
		return err
	}
	BgColor = lipgloss.Color(color)

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
	case config.ColorBlack:
		return "0", nil
	case config.ColorRed:
		return "1", nil
	case config.ColorGreen:
		return "2", nil
	case config.ColorYellow:
		return "3", nil
	case config.ColorBlue:
		return "4", nil
	case config.ColorMagenta:
		return "5", nil
	case config.ColorCyan:
		return "6", nil
	case config.ColorWhite:
		return "7", nil
	case config.ColorBrightBlack:
		return "8", nil
	case config.ColorBrightRed:
		return "9", nil
	case config.ColorBrightGreen:
		return "10", nil
	case config.ColorBrightYellow:
		return "11", nil
	case config.ColorBrightBlue:
		return "12", nil
	case config.ColorBrightMagenta:
		return "13", nil
	case config.ColorBrightCyan:
		return "14", nil
	case config.ColorBrightWhite:
		return "15", nil
	case config.ColorNone:
		return "", nil
	}

	if config.ValidHex(c) {
		return c, nil
	}

	return "", fmt.Errorf("color not resolved: %q", cfgColor)
}

func resolveBorderStyle(borderStyle string) (lipgloss.Border, error) {
	switch borderStyle {
	case config.BorderASCII:
		return lipgloss.ASCIIBorder(), nil
	case config.BorderMarkdown:
		return lipgloss.MarkdownBorder(), nil

	case config.BorderRounded:
		return lipgloss.RoundedBorder(), nil
	case config.BorderSquare:
		return lipgloss.NormalBorder(), nil
	case config.BorderThickSquare:
		return lipgloss.ThickBorder(), nil
	case config.BorderDoubleSquare:
		return lipgloss.DoubleBorder(), nil
	case config.BorderBlock:
		return lipgloss.BlockBorder(), nil
	case config.BorderOuterHalfBlock:
		return lipgloss.OuterHalfBlockBorder(), nil
	case config.BorderInnerHalfBlock:
		return lipgloss.InnerHalfBlockBorder(), nil
	default:
		return lipgloss.Border{}, fmt.Errorf("border style not resolved: %q", borderStyle)
	}
}

func resolveSpinnerStyle(spinnerStyle string) (spinner.Spinner, error) {
	switch spinnerStyle {
	case config.SpinnerLine:
		return spinner.Line, nil
	case config.SpinnerEllipsis:
		return spinner.Ellipsis, nil

	case config.SpinnerDot:
		return spinner.Dot, nil
	case config.SpinnerMiniDot:
		return spinner.MiniDot, nil
	case config.SpinnerJump:
		return spinner.Jump, nil
	case config.SpinnerPulse:
		return spinner.Pulse, nil
	case config.SpinnerPoints:
		return spinner.Points, nil
	case config.SpinnerMeter:
		return spinner.Meter, nil
	case config.SpinnerHamburger:
		return spinner.Hamburger, nil
	default:
		return spinner.Spinner{}, fmt.Errorf("spinner style not resolved: %q", spinnerStyle)
	}
}
