// Package styles provides styling for ui
package styles

import (
	"image/color"

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

	Border               lipgloss.Border
	BorderedStyle        lipgloss.Style
	BorderedFocusedStyle lipgloss.Style

	TableStyle     table.Styles
	DataTableStyle table.Styles

	TabViewStyles tabview.Styles

	OverlayStyle       lipgloss.Style
	NotifBorderedStyle lipgloss.Style

	ToggleStyle        lipgloss.Style
	ToggleFocusedStyle lipgloss.Style
)

func Init(colors config.ColorConfig) {
	TextColor = lipgloss.Color(colors.Text)
	AccentColor = lipgloss.Color(colors.Accent)
	MutedColor = lipgloss.Color(colors.Muted)
	ErrorColor = lipgloss.Color(colors.Error)
	NotifColor = lipgloss.Color(colors.Notif)

	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor)
	AccentStyle = lipgloss.NewStyle().Foreground(AccentColor)
	BoldStyle = DefaultStyle.Bold(true)

	Border = lipgloss.RoundedBorder()
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

	ErrorSymbolColored = DefaultStyle.Foreground(ErrorColor).Render(ErrorSymbol)
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
