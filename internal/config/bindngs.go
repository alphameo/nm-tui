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

type MainKeys struct {
	Quit KeyBinding `kdl:"quit"`
}

type PopupKeys struct {
	Close KeyBinding `kdl:"close"`
}

type TabsKeys struct {
	Next KeyBinding `kdl:"next"`
	Prev KeyBinding `kdl:"prev"`
}

type ToggleKeys struct {
	Toggle KeyBinding `kdl:"toggle"`
}

type WifiKeys struct {
	NextWindow        KeyBinding `kdl:"next_window"`
	FirstWindow       KeyBinding `kdl:"first_window"`
	SecondWindow      KeyBinding `kdl:"second_window"`
	Rescan            KeyBinding `kdl:"rescan"`
	Create            KeyBinding `kdl:"create"`
	OpenCaptivePortal KeyBinding `kdl:"open_captive_portal"`
	EnableHotspot     KeyBinding `kdl:"enable_hotspot"`
	CreateHotspot     KeyBinding `kdl:"create_hotspot"`
}

type WifiAvailableKeys struct {
	Rescan  KeyBinding `kdl:"rescan"`
	Connect KeyBinding `kdl:"connect"`
}

type WifiSavedKeys struct {
	Edit       KeyBinding `kdl:"edit"`
	Connect    KeyBinding `kdl:"connect"`
	Disconnect KeyBinding `kdl:"disconnect"`
	Rescan     KeyBinding `kdl:"rescan"`
	Delete     KeyBinding `kdl:"delete"`
}

type NetworkingKeys struct {
	Up     KeyBinding `kdl:"up"`
	Down   KeyBinding `kdl:"down"`
	Toggle KeyBinding `kdl:"toggle"`
	Rescan KeyBinding `kdl:"rescan"`
}

type ConnectorKeys struct {
	TogglePWVisibility KeyBinding `kdl:"toggle_pw_visibility"`
	Up                 KeyBinding `kdl:"up"`
	Down               KeyBinding `kdl:"down"`
	Connect            KeyBinding `kdl:"connect"`
}

type ProfileCreatorKeys struct {
	TogglePWVisibility KeyBinding `kdl:"toggle_pw_visibility"`
	Up                 KeyBinding `kdl:"up"`
	Down               KeyBinding `kdl:"down"`
	Create             KeyBinding `kdl:"create"`
}

type ProfileEditorKeys struct {
	TogglePWVisibility KeyBinding `kdl:"toggle_pw_visibility"`
	Up                 KeyBinding `kdl:"up"`
	Down               KeyBinding `kdl:"down"`
	Save               KeyBinding `kdl:"save"`
}

type HotspotCreatorKeys struct {
	TogglePWVisibility KeyBinding `kdl:"toggle_pw_visibility"`
	Up                 KeyBinding `kdl:"up"`
	Down               KeyBinding `kdl:"down"`
	Create             KeyBinding `kdl:"create"`
}
