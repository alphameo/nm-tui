// Package logging wraps the infra managers and emits a structured log
// record for every call.
package logging

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"time"

	"github.com/alphameo/nm-tui/internal/infra"
)

// Middleware implements all three manager interfaces by delegating to the
// wrapped implementations. Successes are logged at Debug level, failures at
// Error level along with the exit code of the failed command when the error
// is an [*exec.ExitError].
type Middleware struct {
	logger  *slog.Logger
	wifi    infra.WifiManager
	network infra.NetworkManager
	portal  infra.CaptivePortalOpener
}

// New returns a *Middleware wrapping the given managers.
func New(logger *slog.Logger, wifi infra.WifiManager, network infra.NetworkManager, portal infra.CaptivePortalOpener) *Middleware {
	return &Middleware{logger: logger, wifi: wifi, network: network, portal: portal}
}

func (m *Middleware) call(operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	m.log(operation, start, err)
	return err
}

func callResult[T any](m *Middleware, operation string, fn func() (T, error)) (T, error) {
	start := time.Now()
	res, err := fn()
	m.log(operation, start, err)
	return res, err
}

func (m *Middleware) log(operation string, start time.Time, err error) {
	if err == nil {
		m.logger.Debug("manager call",
			"operation", operation,
			"duration", time.Since(start),
		)
		return
	}
	m.logger.Error("manager call",
		"operation", operation,
		"duration", time.Since(start),
		"error", err,
		"exit_code", exitCode(err),
	)
}

func exitCode(err error) int {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return exitErr.ExitCode()
	}
	return -1
}

func (m *Middleware) ScanWifis(ctx context.Context) ([]infra.AvailableWifi, error) {
	return callResult(m, "wifi.scan_wifis", func() ([]infra.AvailableWifi, error) {
		return m.wifi.ScanWifis(ctx)
	})
}

func (m *Middleware) GetSavedWifiSSIDs(ctx context.Context) ([]string, error) {
	return callResult(m, "wifi.get_saved_wifi_ssids", func() ([]string, error) {
		return m.wifi.GetSavedWifiSSIDs(ctx)
	})
}

func (m *Middleware) GetSavedWifis(ctx context.Context) ([]infra.SavedWifi, error) {
	return callResult(m, "wifi.get_saved_wifis", func() ([]infra.SavedWifi, error) {
		return m.wifi.GetSavedWifis(ctx)
	})
}

func (m *Middleware) ConnectWifi(ctx context.Context, ssid, password string) error {
	return m.call("wifi.connect", func() error {
		return m.wifi.ConnectWifi(ctx, ssid, password)
	})
}

func (m *Middleware) CreateWifiConnection(ctx context.Context, id, ssid, password string, hidden bool) error {
	return m.call("wifi.create_connection", func() error {
		return m.wifi.CreateWifiConnection(ctx, id, ssid, password, hidden)
	})
}

func (m *Middleware) CreateWifiHotspot(ctx context.Context, id string, ssid string, password string) error {
	return m.call("wifi.create_hotspot", func() error {
		return m.wifi.CreateWifiHotspot(ctx, id, ssid, password)
	})
}

func (m *Middleware) EnableQuickWifiHotspot(ctx context.Context) error {
	return m.call("wifi.enable_quick_hotspot", func() error {
		return m.wifi.EnableQuickWifiHotspot(ctx)
	})
}

func (m *Middleware) DeleteWifiConnection(ctx context.Context, name string) error {
	return m.call("wifi.delete_connection", func() error {
		return m.wifi.DeleteWifiConnection(ctx, name)
	})
}

func (m *Middleware) ActivateWifi(ctx context.Context, name string) error {
	return m.call("wifi.activate", func() error {
		return m.wifi.ActivateWifi(ctx, name)
	})
}

func (m *Middleware) DeactivateWifi(ctx context.Context, name string) error {
	return m.call("wifi.deactivate", func() error {
		return m.wifi.DeactivateWifi(ctx, name)
	})
}

func (m *Middleware) GetWifiPassword(ctx context.Context, name string) (string, error) {
	return callResult(m, "wifi.get_password", func() (string, error) {
		return m.wifi.GetWifiPassword(ctx, name)
	})
}

func (m *Middleware) GetWifiInfo(ctx context.Context, name string) (infra.NetworkInfo, error) {
	return callResult(m, "wifi.get_info", func() (infra.NetworkInfo, error) {
		return m.wifi.GetWifiInfo(ctx, name)
	})
}

func (m *Middleware) UpdateWifiInfo(ctx context.Context, name string, info infra.UpdateWifiInfo) error {
	return m.call("wifi.update_info", func() error {
		return m.wifi.UpdateWifiInfo(ctx, name, info)
	})
}

func (m *Middleware) OpenCaptivePortal(ctx context.Context) error {
	return m.call("portal.open", func() error {
		return m.portal.OpenCaptivePortal(ctx)
	})
}

func (m *Middleware) GetNetworkDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	return callResult(m, "network.get_devices", func() ([]infra.NetworkDevice, error) {
		return m.network.GetNetworkDevices(ctx)
	})
}

func (m *Middleware) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	return callResult(m, "network.get_connectivity_status", func() (infra.ConnectivityStatus, error) {
		return m.network.GetConnectivityStatus(ctx)
	})
}

func (m *Middleware) GetNetworking(ctx context.Context) (bool, error) {
	return callResult(m, "network.get_status", func() (bool, error) {
		return m.network.GetNetworking(ctx)
	})
}

func (m *Middleware) EnableNetworking(ctx context.Context) error {
	return m.call("network.enable", func() error {
		return m.network.EnableNetworking(ctx)
	})
}

func (m *Middleware) DisableNetworking(ctx context.Context) error {
	return m.call("network.disable", func() error {
		return m.network.DisableNetworking(ctx)
	})
}

func (m *Middleware) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
	return callResult(m, "network.get_radio_status", func() (infra.RadioStatus, error) {
		return m.network.GetRadioStatus(ctx)
	})
}

func (m *Middleware) EnableWWAN(ctx context.Context) error {
	return m.call("network.enable_wwan", func() error {
		return m.network.EnableWWAN(ctx)
	})
}

func (m *Middleware) DisableWWAN(ctx context.Context) error {
	return m.call("network.disable_wwan", func() error {
		return m.network.DisableWWAN(ctx)
	})
}

func (m *Middleware) EnableWifi(ctx context.Context) error {
	return m.call("network.enable_wifi", func() error {
		return m.network.EnableWifi(ctx)
	})
}

func (m *Middleware) DisableWifi(ctx context.Context) error {
	return m.call("network.disable_wifi", func() error {
		return m.network.DisableWifi(ctx)
	})
}
