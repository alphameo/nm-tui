package infra

import "errors"

var (
	ErrGetDeviceStatuses error = errors.New("failed retrieving network device status")

	ErrCreateWifiConnection error = errors.New("failed creation of wifi connection")

	ErrScanWifis   error = errors.New("failed scanning wifi networks")
	ErrConnectWifi error = errors.New("failed connecting to wifi network")

	ErrGetSavedWifis    error = errors.New("failed retrieving saved wifi networks")
	ErrConnectSavedWifi error = errors.New("failed connecting to saved wifi network")

	ErrDisconnectWifi error = errors.New("failed disconnecting from wifi network")

	ErrGetSavedWifiSSIDs          error = errors.New("failed retrieving saved wifi SSIDs")
	ErrGetWifiPassword            error = errors.New("failed retrieving wifi network password")
	ErrGetWifiSSID                error = errors.New("failed retrieving wifi network ssid")
	ErrGetWifiAutoconnect         error = errors.New("failed retrieving wifi network autoconnect state")
	ErrGetWifiAutoconnectPriority error = errors.New("failed retrieving wifi network autoconnect priority")
	ErrGetWifiActivity            error = errors.New("failed retrieving wifi network activity state")
	ErrGetWifiInfo                error = errors.New("failed retrieving wifi network information")

	ErrUpdateWifiInfo      error = errors.New("failed modifying wifi network information")
	ErrUpdateWifiInfoField error = errors.New("failed modifying wifi network information field")

	ErrDeleteWifi error = errors.New("failed deleting wifi connection")

	ErrGetConnectivityStatus error = errors.New("failed retrieving networking status")
	ErrEnableNetworking      error = errors.New("failed enabling networking")
	ErrDisableNetworking     error = errors.New("failed disabling networking")

	ErrGetRadioStatus error = errors.New("failed retrieving radio status")
	ErrGetWifiStatus  error = errors.New("failed retrieving wifi radio status")
	ErrEnableWifi     error = errors.New("failed enabling wifi radio")
	ErrDisableWifi    error = errors.New("failed disabling wifi radio")
	ErrGetWWANStatus  error = errors.New("failed retrieving wwan radio status")
	ErrEnableWWAN     error = errors.New("failed enabling wwan radio")
	ErrDisableWWAN    error = errors.New("failed disabling wwan radio")

	ErrCreateHotspot error = errors.New("failed creating hotspot")
)
