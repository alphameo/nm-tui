package infra

type WifiScanned struct {
	SSID     string
	Active   bool
	Security string
	Signal   int
}

type WifiStored struct {
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

type DeviceStatus struct {
	Device     string
	Type       string
	State      string
	Connection string
}

type NetworkManager interface {
	// GetDeviceStatus returns network status of device
	GetDeviceStatuses() ([]DeviceStatus, error)

	// GetAvailableWifi shows list of wifi-networks able to be connected.
	GetAvailableWifi() ([]*WifiScanned, error)

	// GetStoredWifi shows list of stored connections and highlights the active one.
	GetStoredWifi() ([]*WifiStored, error)

	// ConnectWifi connects to wifi-network with given ssid using given password.
	ConnectWifi(ssid, password string, hidden bool) error

	// ConnectStoredWifi connects to wifi-network with given name if its password is saved.
	ConnectStoredWifi(name string) error

	// DisconnectFromWifi disconnects from wifi-network with given name.
	DisconnectFromWifi(name string) error

	// GetConnectedWifi gives table of saved connections.
	GetConnectedWifi() ([]string, error)

	// GetWifiPassword gives password of saved wifi-network with given name.
	GetWifiPassword(name string) (string, error)

	// GetWifiInfo gives information about saved wifi-network with given name.
	GetWifiInfo(name string) (*WifiInfo, error)

	// GetWWANStatus returns status of wifi on device
	GetWifiStatus() (bool, error)

	// GetWWANStatus returns status of Wireless Wide Area Network on device
	GetWWANStatus() (bool, error)

	// GetRadioStatus returns status of wifi and Wireless Wide Area Network on device
	GetRadioStatus() (RadioStatus, error)

	// EnableWifi enables wifi on device
	EnableWifi() error

	// EnableWWAN enables Wireless Wide Area Network on device
	EnableWWAN() error

	// DisableWifi disables wifi on device
	DisableWifi() error

	// DisableWWAN disables Wireless Wide Area Network on device
	DisableWWAN() error

	// UpdateWifiInfo updates information about wifi-network with given name.
	UpdateWifiInfo(name string, info *UpdateWifiInfo) error

	// DeleteWifiConnection removes wifi-network with given name from saved connections.
	DeleteWifiConnection(name string) error

	// GetNetworking returns status of networking on device
	GetConnectivityStatus() (ConnectivityStatus, error)

	// EnableNetworking enables all networking on device
	EnableNetworking() error

	// DisableNetworking disables all networking on device
	DisableNetworking() error

	// CreateWifiConnection creates specified connection profile
	CreateWifiConnection(id, ssid, password, device string, hidden bool) error

	// CreateHotspot creates new hotspot
	CreateHotspot(device string, id string, password string, hidden bool) error
}
