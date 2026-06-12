package components

import "github.com/alphameo/nm-tui/internal/ui/models/toggle"

func DefaultToggle() *toggle.Model {
	t := toggle.New()
	t.Symbols.Activated = " "
	t.Symbols.Deactivated = " "
	return t
}
