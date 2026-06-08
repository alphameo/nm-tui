package infra

import "context"

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
