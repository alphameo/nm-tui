// Package styles provides styling for ui
package styles

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/help"
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
	userTermFG string

	TextColor   color.Color
	BgColor     color.Color
	AccentColor color.Color
	MutedColor  color.Color
	ErrorColor  color.Color
	NotifColor  color.Color

	DefaultStyle lipgloss.Style
	AccentStyle  lipgloss.Style
	MutedStyle   lipgloss.Style
	BoldStyle    lipgloss.Style

	Border               lipgloss.Border = lipgloss.RoundedBorder()
	BorderedStyle        lipgloss.Style
	BorderedFocusedStyle lipgloss.Style

	TableStyles     table.Styles
	DataTableStyles table.Styles

	InputCursor tea.CursorShape
	InputStyles textinput.Styles

	ToggleStyles toggle.Styles

	HelpStyles help.Styles

	TabViewStyles tabview.Styles

	OverlayStyle       lipgloss.Style
	NotifBorderedStyle lipgloss.Style

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

	UpdateColorscheme()

	return nil
}

func SetTextColor(color color.Color) {
	TextColor = color
}

func UpdateColorscheme() {
	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor).Background(BgColor)
	AccentStyle = DefaultStyle.Foreground(AccentColor).Bold(true)
	MutedStyle = DefaultStyle.Foreground(MutedColor)
	BoldStyle = DefaultStyle.Bold(true)

	BorderedStyle = DefaultStyle.Border(Border).BorderForeground(TextColor).BorderBackground(BgColor)
	BorderedFocusedStyle = BorderedStyle.BorderForeground(AccentColor)

	TableStyles = tableStyles()
	DataTableStyles = dataTableStyles()

	InputStyles = inputStyles()

	ToggleStyles = toggleStyles()

	HelpStyles = helpStyles()

	TabViewStyles = tabview.GenerateStyles(BorderedStyle)

	OverlayStyle = DefaultStyle.
		Border(Border).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(2, 4).
		BorderForeground(AccentColor).
		BorderBackground(BgColor)

	NotifBorderedStyle = OverlayStyle.BorderForeground(NotifColor)

	SymbolColoredError = DefaultStyle.Foreground(ErrorColor).Render(SymbolError)
}

func tableStyles() table.Styles {
	return table.Styles{
		Selected: lipgloss.NewStyle().
			Foreground(TextColor).
			Background(AccentColor),
		Header: lipgloss.NewStyle().
			BorderStyle(Border).
			BorderForeground(MutedColor).
			BorderBackground(BgColor).
			BorderBottom(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}
}

func dataTableStyles() table.Styles {
	style := tableStyles()
	style.Selected = style.Selected.Background(BgColor)
	return style
}

func inputStyles() textinput.Styles {
	return textinput.Styles{
		Focused: textinput.StyleState{
			Placeholder: MutedStyle,
			Suggestion:  DefaultStyle,
			Prompt:      DefaultStyle,
			Text:        DefaultStyle,
		},
		Blurred: textinput.StyleState{
			Placeholder: MutedStyle,
			Suggestion:  DefaultStyle,
			Prompt:      DefaultStyle,
			Text:        DefaultStyle,
		},
		Cursor: textinput.CursorStyle{
			Color: TextColor,
			Shape: InputCursor,
			Blink: true,
		},
	}
}

func toggleStyles() toggle.Styles {
	return toggle.Styles{
		Focused: AccentStyle.Margin(0, 1),
		Blured:  DefaultStyle.Margin(0, 1),
	}
}

func helpStyles() help.Styles {
	style := help.DefaultDarkStyles()
	style.ShortKey = DefaultStyle
	style.ShortDesc = MutedStyle
	style.ShortSeparator = MutedStyle
	style.Ellipsis = MutedStyle
	style.FullKey = DefaultStyle
	style.FullDesc = MutedStyle
	style.FullSeparator = MutedStyle
	return style
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

	cursor, err := resolveCursorShape(*icons.InputCursorShape)
	if err != nil {
		return err
	}
	InputCursor = cursor

	SymbolsToggle = toggle.Symbols{Activated: *icons.ToggleOn, Deactivated: *icons.ToggleOff}
	SymbolPwHiddenChar = []rune(*icons.PwHiddenChar)[0]
	SymbolError = *icons.Error
	SymbolCheck = *icons.Check
	SymbolConnection = *icons.Connection
	SymbolSignal = *icons.Signal
	SymbolSaved = *icons.Saved
	SymbolAvailable = *icons.Available
	SymbolAccessPoint = *icons.AccessPoint
	SymbolInfra = *icons.Infra
	SymbolMesh = *icons.Mesh
	SymbolAdHoc = *icons.AdHoc
	SymbolEllipsis = *icons.Ellipsis
	SymbolSeparator = *icons.Separator

	return nil
}

func initColors(colors config.ColorConfig) error {
	var err error
	userTermFG, err = queryTerminalForegroundColor()
	if err != nil {
		return err
	}

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

	BgColor = lipgloss.Color("")

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
		return userTermFG, nil
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
		return spinnerMeter, nil
	case config.SpinnerHamburger:
		return spinner.Hamburger, nil
	default:
		return spinner.Spinner{}, fmt.Errorf("spinner style not resolved: %q", spinnerStyle)
	}
}

func resolveCursorShape(cursor string) (tea.CursorShape, error) {
	switch cursor {
	case config.CursorBar:
		return tea.CursorBar, nil
	case config.CursorUnderline:
		return tea.CursorUnderline, nil
	case config.CursorBlock:
		return tea.CursorBlock, nil
	default:
		return 0, fmt.Errorf("input cursor shape not resolved: %q", cursor)
	}
}
