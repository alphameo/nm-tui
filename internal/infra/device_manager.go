package infra

import (
	"context"
	"errors"
)

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

var (
	ErrListNetworkDevices = errors.New("failed to list network devices")

	ErrGetConnectivityStatus = errors.New("failed to get connectivity status")
	ErrParseConnectivity     = errors.New("failed to parse connectivity status")

	ErrIsNetworkingEnabled = errors.New("failed to recognize is networking enabled")
	ErrEnableNetworking    = errors.New("failed to enable networking")
	ErrDisableNetworking   = errors.New("failed to disable networking")

	ErrGetRadioStatus = errors.New("failed to get radio status")

	ErrGetWifiStatus = errors.New("failed to get wifi radio status")
	ErrEnableWifi    = errors.New("failed to enable wifi radio")
	ErrDisableWifi   = errors.New("failed to disable wifi radio")

	ErrGetWWANStatus = errors.New("failed to get wwan radio status")
	ErrEnableWWAN    = errors.New("failed to enable wwan radio")
	ErrDisableWWAN   = errors.New("failed to disable wwan radio")
)

type DeviceManager interface {
	// ListNetworkDevices returns info about network devices
	ListNetworkDevices(ctx context.Context) ([]NetworkDevice, error)

	// GetConnectivityStatus returns connectivity status of device
	GetConnectivityStatus(ctx context.Context) (ConnectivityStatus, error)

	// IsNetworkingEnabled returns networking status
	IsNetworkingEnabled(ctx context.Context) (bool, error)

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
