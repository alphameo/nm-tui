package infra

import "errors"

var (
	ErrScanningAvailableWifi          error = errors.New("failed scanning wifi networks")
	ErrScanningStoredWifi             error = errors.New("failed retrieving stored wifi networks")
	ErrConnectingWifi                 error = errors.New("failed connecting to wifi network")
	ErrConnectingStoredWifi           error = errors.New("failed connecting to stored wifi network")
	ErrDisconnectWifi                 error = errors.New("failed disconnecting to wifi network")
	ErrScanningConnectedWifi          error = errors.New("failed retrieving connected wifi networks")
	ErrGettingWifiPassword            error = errors.New("failed retrieving wifi network password")
	ErrGettingWifiSSID                error = errors.New("failed retrieving wifi network ssid")
	ErrGettingWifiAutoconnect         error = errors.New("failed retrieving wifi network autoconnect state")
	ErrGettingWifiAutoconnectPriority error = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGettingWifiActivity            error = errors.New("failed retrieving wifi network activity state")
	ErrGettingWifiInfo                error = errors.New("failed retrieving wifi network information")
	ErrUpdatingWifiInfo               error = errors.New("failed retrieving wifi network information")
	ErrUpdatingWifiInfoField          error = errors.New("failed modifying wifi network information field")
	ErrDeletingWifi                   error = errors.New("failed deleting wifi connection")
	ErrConnectingVPN                  error = errors.New("failed connecting to VPN")
)
