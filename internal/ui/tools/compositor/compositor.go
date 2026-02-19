// Package compositor provides functions that place one TUI frame inside another
package compositor

import (
	"github.com/alphameo/nm-tui/internal/ui/components/floating"
)

func PlaceTitle(view, title string) string {
	return floating.Compose(title, view, floating.Begin, floating.Begin, 2, 0)
}
