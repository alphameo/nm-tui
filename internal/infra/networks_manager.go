package infra

import (
	"context"
	"errors"
)

type AvailableNetwork struct {
	SSID          string
	Active        bool
	SecurityMode  string
	Signal        int
	Band          float64
	Rate          float64
	LookingDevice string
	NetworkMode   NetworkMode
}

type NetworkProfileShort struct {
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

type NetworkProfile struct {
	Name                string
	SSID                string
	Password            string
	Active              bool
	Autoconnect         bool
	AutoconnectPriority int
	Mode                NetworkMode
}

type UpdateProfile struct {
	Name                string
	Password            string
	Autoconnect         bool
	AutoconnectPriority int
}

type RadioStatus struct {
	EnabledWifi bool
	EnabledWWAN bool
}

var (
	ErrCreateWifiConnection = errors.New("failed to create wifi connection")

	ErrScanNetworks       = errors.New("failed to list networks with rescan")
	ErrListNetworks       = errors.New("failed to list networks")
	ErrConnectToNetwork   = errors.New("failed to connect to network")
	ErrTryActivateNetwork = errors.New("failed to activate network")

	ErrListProfiles    = errors.New("failed to list network profiles")
	ErrActivateProfile = errors.New("failed connecting to saved wifi network")

	ErrDeactivateProfile = errors.New("failed disconnecting from wifi network")

	ErrListProfileNames           = errors.New("failed retrieving saved connection names")
	ErrGetWifiPassword            = errors.New("failed retrieving wifi network password")
	ErrGetWifiSSID                = errors.New("failed retrieving wifi network ssid")
	ErrGetWifiAutoconnect         = errors.New("failed retrieving wifi network autoconnect state")
	ErrGetWifiAutoconnectPriority = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGetWifiActivity            = errors.New("failed retrieving wifi network activity state")
	ErrGetProfile                 = errors.New("failed retrieving wifi network information")
	ErrGetNetMode                 = errors.New("failed retrieving network mode")

	ErrUpdateProfile = errors.New("failed modifying wifi network information")

	ErrDeleteProfile = errors.New("failed deleting wifi connection")

	ErrCreateHotspotProfile = errors.New("failed to create hotspot profile")
	ErrQuickHotspot         = errors.New("failed enabling quick hotspot")
)

type NetworksManager interface {
	// ListNetworksWithRescan returns list of networks able to be connected.
	ListNetworksWithRescan(ctx context.Context) ([]AvailableNetwork, error)

	// ListNetworks returns cached list of networks able to be connected.
	ListNetworks(ctx context.Context) ([]AvailableNetwork, error)

	// ListProfileNames returns names of saved connections.
	ListProfileNames(ctx context.Context) ([]string, error)

	// ListProfiles returns saved connections.
	ListProfiles(ctx context.Context) ([]NetworkProfileShort, error)

	// ConnectToNetwork creates network connection.
	ConnectToNetwork(ctx context.Context, ssid, password string) error

	// TryActivateNetwork connects to the network by SSID. Uses credentials from the corresponding profile if
	// corresponding profile exists, or cretes profile if does not.
	TryActivateNetwork(ctx context.Context, ssid string) error

	// CreateConnectionProfile creates specified connection profile.
	CreateConnectionProfile(ctx context.Context, name, ssid, password string, hidden bool) error

	// CreateHotspotProfile creates new hotspot profile.
	CreateHotspotProfile(ctx context.Context, name string, ssid string, password string) error

	// QuickHotspot activates hotspot, silently creating it's profile if not present.
	QuickHotspot(ctx context.Context) error

	// DeleteProfile removes network profile with given name from saved connections.
	DeleteProfile(ctx context.Context, name string) error

	// ActivateProfile activates network profile with given name: connects to network or enables hotspot.
	ActivateProfile(ctx context.Context, name string) error

	// DeactivateProfile deactivates network profile with given name: disconnects to network or disnables hotspot.
	DeactivateProfile(ctx context.Context, name string) error

	// GetProfilePassword gives password of saved network profile with given name.
	GetProfilePassword(ctx context.Context, name string) (string, error)

	// GetProfile gives information about saved network with given name.
	GetProfile(ctx context.Context, name string) (NetworkProfile, error)

	// UpdateProfile updates information about wifi-network with given name.
	UpdateProfile(ctx context.Context, name string, info UpdateProfile) error
}
