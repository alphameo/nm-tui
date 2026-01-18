package renderer

import (
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

const ErrorSymbol = "✗"

var ErrorSymbolColored = styles.DefaultStyle.Foreground(styles.ErrorColor).Render(ErrorSymbol)
