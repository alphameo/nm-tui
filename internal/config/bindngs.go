package config

type Keys struct {
	Main           MainKeys           `kdl:"main"`
	Popup          PopupKeys          `kdl:"popup"`
	Tabs           TabsKeys           `kdl:"tabs"`
	Toggle         ToggleKeys         `kdl:"toggle"`
	Wifi           WifiKeys           `kdl:"wifi"`
	WifiAvailable  WifiAvailableKeys  `kdl:"wifi_available"`
	WifiSaved      WifiSavedKeys      `kdl:"wifi_saved"`
	Networking     NetworkingKeys     `kdl:"networking"`
	Connector      ConnectorKeys      `kdl:"connector"`
	ProfileCreator ProfileCreatorKeys `kdl:"profile_creator"`
	ProfileEditor  ProfileEditorKeys  `kdl:"profile_editor"`
	HotspotCreator HotspotCreatorKeys `kdl:"hotspot_creator"`
}

type KeyBinding struct {
	Keys []string `kdl:"keys,arguments"`
}
type (
	MainKeys struct {
		Quit KeyBinding `kdl:"quit"`
	}
	PopupKeys struct {
		Close KeyBinding `kdl:"close"`
	}
	TabsKeys           struct{ Next, Prev KeyBinding }
	ToggleKeys         struct{ Toggle KeyBinding }
	WifiKeys           struct{ NextWindow, FirstWindow, SecondWindow, Rescan, Create, OpenCaptivePortal, EnableHotspot, CreateHotspot KeyBinding }
	WifiAvailableKeys  struct{ Rescan, Connect KeyBinding }
	WifiSavedKeys      struct{ Edit, Connect, Disconnect, Rescan, Delete KeyBinding }
	NetworkingKeys     struct{ Up, Down, Toggle, Rescan KeyBinding }
	ConnectorKeys      struct{ TogglePWVisibility, Up, Down, Connect KeyBinding }
	ProfileCreatorKeys struct{ TogglePWVisibility, Up, Down, Create KeyBinding }
	ProfileEditorKeys  struct{ TogglePWVisibility, Up, Down, Save KeyBinding }
	HotspotCreatorKeys struct{ TogglePWVisibility, Up, Down, Create KeyBinding }
)
