package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

const ErrorSymbol = "✗"

var (
	BorderOffset int = lipgloss.Width(Border.Left) * 2
	TabBarHeight int = BorderOffset + 1

	ErrorSymbolColored    = DefaultStyle.Foreground(ErrorColor).Render(ErrorSymbol)
	ProfileCreatorTitle   = renderer.RenderTitle("Create Network profile")
	HotspotCreatorTitle   = renderer.RenderTitle("Create Hotspot")
	NetworkConnectorTitle = renderer.RenderTitle("Connect to Network")
	SavedNetworkInfoTitle = renderer.RenderTitle("Saved network info")
)
