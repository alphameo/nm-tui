package infra

import "errors"

var (
	ErrScanAvailableWifi error = errors.New("failed scanning wifi networks")
	ErrConnectWifi       error = errors.New("failed connecting to wifi network")

	ErrScanStoredWifi    error = errors.New("failed retrieving stored wifi networks")
	ErrConnectStoredWifi error = errors.New("failed connecting to stored wifi network")

	ErrDisconnectWifi error = errors.New("failed disconnecting from wifi network")

	ErrScanConnectedWifi          error = errors.New("failed retrieving connected wifi networks")
	ErrGetWifiPassword            error = errors.New("failed retrieving wifi network password")
	ErrGetWifiSSID                error = errors.New("failed retrieving wifi network ssid")
	ErrGetWifiAutoconnect         error = errors.New("failed retrieving wifi network autoconnect state")
	ErrGetWifiAutoconnectPriority error = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGetWifiActivity            error = errors.New("failed retrieving wifi network activity state")
	ErrGetWifiInfo                error = errors.New("failed retrieving wifi network information")

	ErrUpdateWifiInfo      error = errors.New("failed modifying wifi network information")
	ErrUpdateWifiInfoField error = errors.New("failed modifying wifi network information field")

	ErrDeleteWifi error = errors.New("failed deleting wifi connection")

	ErrConnectVPN error = errors.New("failed connecting to VPN")

	ErrGetConnectivityStatus error = errors.New("failed retrieving networking status")
	ErrEnableNetworking      error = errors.New("failed enabling networking")
	ErrDisableNetworking     error = errors.New("failed disabling networking")

	ErrGetRadioStaus error = errors.New("failed retrieving radio status")
	ErrGetWifiStaus  error = errors.New("failed retrieving wifi radio status")
	ErrEnableWifi    error = errors.New("failed enabling wifi radio")
	ErrDisableWifi   error = errors.New("failed disabling wifi radio")
	ErrGetWWANStaus  error = errors.New("failed retrieving wwan radio status")
	ErrEnableWWAN    error = errors.New("failed enabling wwan radio")
	ErrDisableWWAN   error = errors.New("failed disabling wwan radio")

	ErrCreateHotspot    error = errors.New("failed creating hotspot")
	ErrConnectHotspot   error = errors.New("failed creating hotspot")
	ErrDisconectHotspot error = errors.New("failed creating hotspot")
)
