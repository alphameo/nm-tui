package infra

import (
	"errors"
	"os/exec"
)

var (
	ErrListDevices error = errors.New("failed to list network devices")

	ErrGetConnectivityStatus error = errors.New("failed to get connectivity status")
	ErrParseConnectivity     error = errors.New("failed to parse connectivity status")

	ErrIsNetworkingEnabled   error = errors.New("failed to recognize is networking enabled")
	ErrEnableNetworking      error = errors.New("failed to enable networking")
	ErrDisableNetworking     error = errors.New("failed to disable networking")

	ErrGetRadioStatus error = errors.New("failed to get radio status")

	ErrGetWifiStatus  error = errors.New("failed to get wifi radio status")
	ErrEnableWifi     error = errors.New("failed to enable wifi radio")
	ErrDisableWifi    error = errors.New("failed to disable wifi radio")

	ErrGetWWANStatus  error = errors.New("failed to get wwan radio status")
	ErrEnableWWAN     error = errors.New("failed to enable wwan radio")
	ErrDisableWWAN    error = errors.New("failed to disable wwan radio")

	ErrCreateWifiConnection error = errors.New("failed to create wifi connection")

	ErrScanNetworks     error = errors.New("failed to scan networks")
	ErrConnectToNetwork error = errors.New("failed to connect to network")

	ErrListProfiles    error = errors.New("failed to list network profiles")
	ErrActivateProfile error = errors.New("failed connecting to saved wifi network")

	ErrDeactivateProfile error = errors.New("failed disconnecting from wifi network")

	ErrListProfileNames           error = errors.New("failed retrieving saved connection names")
	ErrGetWifiPassword            error = errors.New("failed retrieving wifi network password")
	ErrGetWifiSSID                error = errors.New("failed retrieving wifi network ssid")
	ErrGetWifiAutoconnect         error = errors.New("failed retrieving wifi network autoconnect state")
	ErrGetWifiAutoconnectPriority error = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGetWifiActivity            error = errors.New("failed retrieving wifi network activity state")
	ErrGetProfile                 error = errors.New("failed retrieving wifi network information")
	ErrGetNetMode                 error = errors.New("failed retrieving network mode")
	ErrParseNetMode               error = errors.New("failed to parse network mode")

	ErrUpdateProfile error = errors.New("failed modifying wifi network information")

	ErrDeleteProfile error = errors.New("failed deleting wifi connection")

	ErrCreateHotspotProfile error = errors.New("failed to create hotspot profile")
	ErrQuickHotspot         error = errors.New("failed enabling quick hotspot")

	ErrOpenCaptivePortal   error = errors.New("failed to open captive portal")
	ErrGetGatewayIP        error = errors.New("failed to get gateway ip")
	ErrUnsupportedPlarform error = errors.New("unsupported platform")
)

func ExtractStderr(err error) string {
	var stderr string
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return stderr
}
