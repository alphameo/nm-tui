package infra

import "errors"

var (
	ErrScanAvailableWifi          error = errors.New("failed scanning wifi networks")
	ErrScanStoredWifi             error = errors.New("failed retrieving stored wifi networks")
	ErrConnectWifi                error = errors.New("failed connecting to wifi network")
	ErrConnectStoredWifi          error = errors.New("failed connecting to stored wifi network")
	ErrDisconnectWifi             error = errors.New("failed disconnecting to wifi network")
	ErrScanConnectedWifi          error = errors.New("failed retrieving connected wifi networks")
	ErrGetWifiPassword            error = errors.New("failed retrieving wifi network password")
	ErrGetWifiSSID                error = errors.New("failed retrieving wifi network ssid")
	ErrGetWifiAutoconnect         error = errors.New("failed retrieving wifi network autoconnect state")
	ErrGetWifiAutoconnectPriority error = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGetWifiActivity            error = errors.New("failed retrieving wifi network activity state")
	ErrGetWifiInfo                error = errors.New("failed retrieving wifi network information")
	ErrUpdateWifiInfo             error = errors.New("failed retrieving wifi network information")
	ErrUpdateWifiInfoField        error = errors.New("failed modifying wifi network information field")
	ErrDeleteWifi                 error = errors.New("failed deleting wifi connection")
	ErrConnectVPN                 error = errors.New("failed connecting to VPN")
)
