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

type NetworkManager interface {
	// GetAvailableWifi shows list of wifi-networks able to be connected.
	GetAvailableWifi() ([]*WifiScanned, error)

	// GetStoredWifi shows list of stored connections and highlights the active one.
	GetStoredWifi() ([]*WifiStored, error)

	// ConnectWifi connects to wifi-network with given ssid using given password.
	ConnectWifi(ssid, password string) error

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

	// UpdateWifiInfo updates information about wifi-network with given name.
	UpdateWifiInfo(name string, info *UpdateWifiInfo) error

	// DeleteWifiConnection removes wifi-network with given name from saved connections.
	DeleteWifiConnection(name string) error

	// ConnectVPN connects to vpn with given vpnName.
	ConnectVPN(vpnName string) error
}
