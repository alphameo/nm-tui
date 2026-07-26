package styles

import (
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

const (
	ErrorSymbol       = "✗"
	CheckSymbol       = ""
	PWCharacter       = '•'
	ConnectionSymbol  = "󱘖"
	SignalSymbol      = ""
	AccessPointSymbol = "󰀃"
	InfraSymbol       = "🖳"
	MeshSymbol        = ""
	AdHocSymbol       = ""
)

var (
	BorderOffset int
	TabBarHeight int

	ErrorSymbolColored string
	ToggleSymbols      = toggle.Symbols{Activated: " ", Deactivated: " "}
)
