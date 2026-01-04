package infra

type WifiScanned struct {
	SSID     string
	Active   bool
	Security string
	Signal   int
}

type WifiStored struct {
	Active bool
	Name   string
}

type NetworkManager interface {
	// WifiScan shows list of wifi-networks able to be connected
	WifiScan() ([]WifiScanned, error)

	// WifiStoredConnections shows list of stored connections and highlights the active one
	WifiStoredConnections() ([]WifiStored, error)

	// WifiConnect connects to wifi-network with given ssid using given password.
	WifiConnect(ssid, password string) error

	// WifiConnectSaved connects to wifi-network with given ssid if its password is saved.
	WifiConnectSaved(ssid string) error

	// WifiGetConnected gives table of saved connections.
	WifiGetConnected() ([]string, error)

	// WifiGetPassword gives password of saved wifi-network with given ssid.
	WifiGetPassword(ssid string) (string, error)

	// WifiDeleteConnection removes wifi-network with given ssid from saved connections.
	WifiDeleteConnection(ssid string) error

	// VpnConnect connects to vpn with given vpnName
	VpnConnect(vpnName string) error
}
