package infra

type WifiScanned struct {
	SSID     string
	Active   bool
	Security string
	Signal   int
}

type WifiStored struct {
	SSID   string
	Active bool
}

type NetworkManager interface {
	// ScanWifi shows list of wifi-networks able to be connected
	ScanWifi() ([]*WifiScanned, error)

	// GetStoredWifi shows list of stored connections and highlights the active one
	GetStoredWifi() ([]*WifiStored, error)

	// ConnectWifi connects to wifi-network with given ssid using given password.
	ConnectWifi(ssid, password string) error

	// ConnectSavedWifi connects to wifi-network with given ssid if its password is saved.
	ConnectSavedWifi(ssid string) error

	// DisconnectFromWifi() disconnects from wifi-network with given ssid.
	DisconnectFromWifi(ssid string) error

	// GetConnectedWifi gives table of saved connections.
	GetConnectedWifi() ([]string, error)

	// GetWifiPassword gives password of saved wifi-network with given ssid.
	GetWifiPassword(ssid string) (string, error)

	// DeleteWifiConnection removes wifi-network with given ssid from saved connections.
	DeleteWifiConnection(ssid string) error

	// ConnectVPN connects to vpn with given vpnName
	ConnectVPN(vpnName string) error
}
