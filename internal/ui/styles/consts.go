package styles

import (
	"time"

	"charm.land/bubbles/v2/spinner"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

var (
	SymbolPwHiddenChar rune
	SymbolError        string
	SymbolCheck        string
	SymbolConnection   string
	SymbolSignal       string
	SymbolSaved        string
	SymbolAccessPoint  string
	SymbolInfra        string
	SymbolMesh         string
	SymbolAdHoc        string
)

var (
	SymbolColoredError string
	SymbolsToggle      toggle.Symbols
	SymbolEllipsis     string
	SymbolSeparator    string
)

var spinnerMeter = spinner.Spinner{
	Frames: []string{
		"▱ ▱ ▱ ",
		"▰ ▱ ▱ ",
		"▰ ▰ ▱ ",
		"▰ ▰ ▰ ",
		"▰ ▰ ▱ ",
		"▰ ▱ ▱ ",
		"▱ ▱ ▱ ",
	},
	FPS: time.Second / 7,
}
