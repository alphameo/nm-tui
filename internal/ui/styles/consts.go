package styles

import (
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

const (
	ErrorSymbol = "✗"
	CheckSymbol = ""
	PWCharacter = '•'
)

var (
	BorderOffset int
	TabBarHeight int

	ErrorSymbolColored string
	ToggleSymbols      = toggle.Symbols{Activated: " ", Deactivated: " "}

	ProfileCreatorTitle   = renderer.RenderTitle("Create Network profile")
	HotspotCreatorTitle   = renderer.RenderTitle("Create Hotspot")
	NetworkConnectorTitle = renderer.RenderTitle("Connect to Network")
	SavedNetworkInfoTitle = renderer.RenderTitle("Saved network info")
)
