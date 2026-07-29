package styles

import (
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

const (
	SymbolError       = "✗"
	SymbolCheck       = ""
	SymbolConnection  = "󱘖"
	SymbolSignal      = ""
	SymbolAccessPoint = "󰀃"
	SymbolInfra       = "🖳"
	SymbolMesh        = ""
	SymbolAdHoc       = ""
	SymbolPwHiddenChar = '•'
)

var (
	BorderOffset int
	TabBarHeight int

	SymbolColoredError string
	SymbolsToggle      = toggle.Symbols{Activated: " ", Deactivated: " "}
)
