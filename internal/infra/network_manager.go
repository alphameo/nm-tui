package infra

import "context"

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
	Mode   NetworkMode
}

type NetworkMode int

const (
	NetworkNil NetworkMode = iota
	NetworkAccessPoint
	NetworkInfra
	NetworkMesh
	NetworkAdHoc
)

func (m NetworkMode) String() string {
	switch m {
	case NetworkAccessPoint:
		return "Access Point"
	case NetworkInfra:
		return "Infrastructure"
	case NetworkMesh:
		return "Mesh"
	case NetworkAdHoc:
		return "AdHoc"
	default:
		return "Undefined"
	}
}

func (m NetworkMode) Icon() string {
	switch m {
	case NetworkAccessPoint:
		return "󰀃"
	case NetworkInfra:
		return "🖳"
	case NetworkMesh:
		return ""
	case NetworkAdHoc:
		return ""
	default:
		return "?"
	}
}

type WifiInfo struct {
	Name                string
	SSID                string
	Password            string
	Active              bool
	Autoconnect         bool
	AutoconnectPriority int
	Mode                NetworkMode
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

type ConnectivityStatus int

const (
	ConnectvityNil ConnectivityStatus = iota
	ConnectivityNone
	ConnectivityPortal
	ConnectivityLimited
	ConnectivityFull
	ConnectivityUnknown
)

func (c ConnectivityStatus) String() string {
	switch c {
	case ConnectivityNone:
		return "None"
	case ConnectivityPortal:
		return "Portal"
	case ConnectivityLimited:
		return "Limited"
	case ConnectivityFull:
		return "Full"
	case ConnectivityUnknown:
		return "Unknown"
	default:
		return "Undefined"
	}
}

type NetworkDevice struct {
	Device     string
	Type       string
	State      string
	Connection string
}

type NetworkManager interface {
	// GetNetworkDevices returns info about network devices
	GetNetworkDevices(ctx context.Context) ([]NetworkDevice, error)

	// GetConnectivityStatus returns connectivity status of device
	GetConnectivityStatus(ctx context.Context) (ConnectivityStatus, error)

	// GetNetworking returns networking status
	GetNetworking(ctx context.Context) (bool, error)

	// EnableNetworking enables all networking on device
	EnableNetworking(ctx context.Context) error

	// DisableNetworking disables all networking on device
	DisableNetworking(ctx context.Context) error

	// GetRadioStatus returns status of wifi and Wireless Wide Area Network on device
	GetRadioStatus(ctx context.Context) (RadioStatus, error)
	// EnableWWAN enables Wireless Wide Area Network on device
	EnableWWAN(ctx context.Context) error

	// DisableWWAN disables Wireless Wide Area Network on device
	DisableWWAN(ctx context.Context) error

	// EnableWifi enables wifi on device
	EnableWifi(ctx context.Context) error

	// DisableWifi disables wifi on device
	DisableWifi(ctx context.Context) error
}

type WifiManager interface {
	// ScanWifis shows list of wifi-networks able to be connected.
	ScanWifis(ctx context.Context) ([]AvailableWifi, error)

	// GetSavedWifiSSIDs gives table of saved connections.
	GetSavedWifiSSIDs(ctx context.Context) ([]string, error)

	// GetSavedWifis shows list of saved connections and highlights the active one.
	GetSavedWifis(ctx context.Context) ([]SavedWifi, error)

	// ConnectWifi creates connection with wifi-network.
	ConnectWifi(ctx context.Context, id, ssid, password string) error

	// CreateWifiConnection creates specified connection profile
	CreateWifiConnection(ctx context.Context, id, ssid, password string, hidden bool) error

	// CreateWifiHotspot creates new hotspot
	CreateWifiHotspot(ctx context.Context, id string, password string, hidden bool) error

	// EanbleWifiHotspot creates new hotspot
	EnableQuickWifiHotspot(ctx context.Context) error

	// DeleteWifiConnection removes wifi-network with given name from saved connections.
	DeleteWifiConnection(ctx context.Context, name string) error

	// ActivateWifi activates connection with wifi-network with given name.
	ActivateWifi(ctx context.Context, name string) error

	// DeactivateWifi deactivates connection with wifi-network with given name.
	DeactivateWifi(ctx context.Context, name string) error

	// GetWifiPassword gives password of saved wifi-network with given name.
	GetWifiPassword(ctx context.Context, name string) (string, error)

	// GetWifiInfo gives information about saved wifi-network with given name.
	GetWifiInfo(ctx context.Context, name string) (WifiInfo, error)

	// UpdateWifiInfo updates information about wifi-network with given name.
	UpdateWifiInfo(ctx context.Context, name string, info UpdateWifiInfo) error
}
