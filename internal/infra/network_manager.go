package infra

type AvailableWifi struct {
	SSID     string
	Active   bool
	Security string
	Signal   int
}

type SavedWifi struct {
	Name   string
	SSID   string
	Active bool
}

type WifiInfo struct {
	Name                string
	SSID                string
	Password            string
	Active              bool
	Autoconnect         bool
	AutoconnectPriority int
}

type UpdateWifiInfo struct {
	Name                string
	Password            string
	Autoconnect         bool
	AutoconnectPriority int
}

type RadioStatus struct {
	EnabledWifi bool
	EnabledWWAN bool
}

type ConnectivityStatus string

const (
	NetworkNone    = "none"
	NetworkPortal  = "portal"
	NetworkLimited = "limited"
	NetworkFull    = "full"
	NetworkUnknown = "unknown"
)

type NetworkDevice struct {
	Device     string
	Type       string
	State      string
	Connection string
}

type NetworkManager interface {
	// GetNetworkDevices returns network status of device
	GetNetworkDevices() ([]NetworkDevice, error)

	// ScanWifis shows list of wifi-networks able to be connected.
	ScanWifis() ([]AvailableWifi, error)

	// GetSavedWifis shows list of saved connections and highlights the active one.
	GetSavedWifis() ([]SavedWifi, error)

	// ConnectWifi creates connection with wifi-network.
	ConnectWifi(ssid, password string, hidden bool) error

	// ActivateWifi activates connection with wifi-network with given name.
	ActivateWifi(name string) error

	// DeactivateWifi deactivates connection with wifi-network with given name.
	DeactivateWifi(name string) error

	// GetSavedWifiSSIDs gives table of saved connections.
	GetSavedWifiSSIDs() ([]string, error)

	// GetWifiPassword gives password of saved wifi-network with given name.
	GetWifiPassword(name string) (string, error)

	// GetWifiInfo gives information about saved wifi-network with given name.
	GetWifiInfo(name string) (WifiInfo, error)

	// UpdateWifiInfo updates information about wifi-network with given name.
	UpdateWifiInfo(name string, info UpdateWifiInfo) error

	// DeleteWifiConnection removes wifi-network with given name from saved connections.
	DeleteWifiConnection(name string) error

	// CreateWifiConnection creates specified connection profile
	CreateWifiConnection(id, ssid, password, device string, hidden bool) error

	// CreateHotspot creates new hotspot
	CreateHotspot(device string, id string, password string, hidden bool) error

	// GetWifiStatus returns status of wifi on device
	GetWifiStatus() (bool, error)

	// GetWWANStatus returns status of Wireless Wide Area Network on device
	GetWWANStatus() (bool, error)

	// GetRadioStatus returns status of wifi and Wireless Wide Area Network on device
	GetRadioStatus() (RadioStatus, error)

	// EnableWifi enables wifi on device
	EnableWifi() error

	// DisableWifi disables wifi on device
	DisableWifi() error

	// EnableWWAN enables Wireless Wide Area Network on device
	EnableWWAN() error

	// DisableWWAN disables Wireless Wide Area Network on device
	DisableWWAN() error

	// GetConnectivityStatus returns status of networking on device
	GetConnectivityStatus() (ConnectivityStatus, error)

	// EnableNetworking enables all networking on device
	EnableNetworking() error

	// DisableNetworking disables all networking on device
	DisableNetworking() error
}
