// Package styles provides styling for ui
package styles

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/components/tabview"
)

var (
	TextColor   = lipgloss.Color("#cbcbcb")
	AccentColor = lipgloss.Color("#865fff")
	MutedColor  = lipgloss.Color("#595959")
	ErrorColor  = lipgloss.Color("#ff0000")
	NotifColor  = lipgloss.Color("#e4bf7a")

	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor)

	Border        = lipgloss.RoundedBorder()
	BorderedStyle = DefaultStyle.Border(Border)

	TableStyle     = tableStyle()
	DataTableStyle = dataTableStyle()

	TabViewStyles = tabview.GenerateStyles(&BorderedStyle)

	OverlayStyle = DefaultStyle.
			Border(Border).
			Align(lipgloss.Center, lipgloss.Center).
			Width(100).
			Height(10)
	NotifBorderedStyle = OverlayStyle.BorderForeground(NotifColor)
)

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
