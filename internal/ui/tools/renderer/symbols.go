package renderer

import (
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

var (
	errSymbol = "✗"
	ErrSymbol = styles.DefaultStyle.Foreground(styles.ErrorColor).Render(errSymbol)
)
