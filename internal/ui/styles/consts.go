package styles

import (
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

const (
	SymbolError       = "✗"
	SymbolCheck       = ""
	CharacterPassword = '•'
	SymbolConnection  = "󱘖"
	SymbolSignal      = ""
	SymbolAccessPoint = "󰀃"
	SymbolInfra       = "🖳"
	SymbolMesh        = ""
	SymbolAdHoc       = ""
)

var (
	BorderOffset int
	TabBarHeight int

	SymbolColoredError string
	SymbolsToggle      = toggle.Symbols{Activated: " ", Deactivated: " "}
)
