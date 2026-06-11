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

type NetworkInfo struct {
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
	CreateWifiHotspot(ctx context.Context, id string, ssid string, password string) error

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
	GetWifiInfo(ctx context.Context, name string) (NetworkInfo, error)

	// UpdateWifiInfo updates information about wifi-network with given name.
	UpdateWifiInfo(ctx context.Context, name string, info UpdateWifiInfo) error
}
