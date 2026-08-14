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
	logger       *slog.Logger
	networks     infra.NetworksManager
	connectivity infra.ConnectivityManager
	portal       infra.CaptivePortalOpener
}

// New returns a *Middleware wrapping the given managers.
func New(
	logger *slog.Logger,
	networks infra.NetworksManager,
	connectivity infra.ConnectivityManager,
	portal infra.CaptivePortalOpener,
) *Middleware {
	return &Middleware{logger: logger, networks: networks, connectivity: connectivity, portal: portal}
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

func (m *Middleware) ScanNetworks(ctx context.Context) ([]infra.AvailableNetwork, error) {
	return callResult(m, "wifi.scan_networks", func() ([]infra.AvailableNetwork, error) {
		return m.networks.ScanNetworks(ctx)
	})
}

func (m *Middleware) ListProfileNames(ctx context.Context) ([]string, error) {
	return callResult(m, "wifi.list_profile_names", func() ([]string, error) {
		return m.networks.ListProfileNames(ctx)
	})
}

func (m *Middleware) ListProfiles(ctx context.Context) ([]infra.NetworkProfileShort, error) {
	return callResult(m, "wifi.list_profiles", func() ([]infra.NetworkProfileShort, error) {
		return m.networks.ListProfiles(ctx)
	})
}

func (m *Middleware) ConnectToNetwork(ctx context.Context, ssid, password string) error {
	return m.call("wifi.connect_to_network", func() error {
		return m.networks.ConnectToNetwork(ctx, ssid, password)
	})
}

func (m *Middleware) CreateConnectionProfile(ctx context.Context, id, ssid, password string, hidden bool) error {
	return m.call("wifi.create_connection_profile", func() error {
		return m.networks.CreateConnectionProfile(ctx, id, ssid, password, hidden)
	})
}

func (m *Middleware) CreateHotspotProfile(ctx context.Context, id string, ssid string, password string) error {
	return m.call("wifi.create_hotspot_profile", func() error {
		return m.networks.CreateHotspotProfile(ctx, id, ssid, password)
	})
}

func (m *Middleware) QuickHotspot(ctx context.Context) error {
	return m.call("wifi.quick_hotspot", func() error {
		return m.networks.QuickHotspot(ctx)
	})
}

func (m *Middleware) DeleteProfile(ctx context.Context, name string) error {
	return m.call("wifi.delete_profile", func() error {
		return m.networks.DeleteProfile(ctx, name)
	})
}

func (m *Middleware) ActivateProfile(ctx context.Context, name string) error {
	return m.call("wifi.activate_profile", func() error {
		return m.networks.ActivateProfile(ctx, name)
	})
}

func (m *Middleware) DeactivateProfile(ctx context.Context, name string) error {
	return m.call("wifi.deactivate_profile", func() error {
		return m.networks.DeactivateProfile(ctx, name)
	})
}

func (m *Middleware) GetProfilePassword(ctx context.Context, name string) (string, error) {
	return callResult(m, "wifi.get_wifi_password", func() (string, error) {
		return m.networks.GetProfilePassword(ctx, name)
	})
}

func (m *Middleware) GetProfile(ctx context.Context, name string) (infra.NetworkProfile, error) {
	return callResult(m, "wifi.get_profile", func() (infra.NetworkProfile, error) {
		return m.networks.GetProfile(ctx, name)
	})
}

func (m *Middleware) UpdateProfile(ctx context.Context, name string, info infra.UpdateProfile) error {
	return m.call("wifi.update_profile", func() error {
		return m.networks.UpdateProfile(ctx, name, info)
	})
}

func (m *Middleware) OpenCaptivePortal(ctx context.Context) error {
	return m.call("portal.open_captive_portal", func() error {
		return m.portal.OpenCaptivePortal(ctx)
	})
}

func (m *Middleware) ListDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	return callResult(m, "network.list_devices", func() ([]infra.NetworkDevice, error) {
		return m.connectivity.ListDevices(ctx)
	})
}

func (m *Middleware) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	return callResult(m, "network.get_connectivity_status", func() (infra.ConnectivityStatus, error) {
		return m.connectivity.GetConnectivityStatus(ctx)
	})
}

func (m *Middleware) IsNetworkingEnabled(ctx context.Context) (bool, error) {
	return callResult(m, "network.is_networking_enabled", func() (bool, error) {
		return m.connectivity.IsNetworkingEnabled(ctx)
	})
}

func (m *Middleware) EnableNetworking(ctx context.Context) error {
	return m.call("network.enable_networking", func() error {
		return m.connectivity.EnableNetworking(ctx)
	})
}

func (m *Middleware) DisableNetworking(ctx context.Context) error {
	return m.call("network.disable_networking", func() error {
		return m.connectivity.DisableNetworking(ctx)
	})
}

func (m *Middleware) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
	return callResult(m, "network.get_radio_status", func() (infra.RadioStatus, error) {
		return m.connectivity.GetRadioStatus(ctx)
	})
}

func (m *Middleware) EnableWWAN(ctx context.Context) error {
	return m.call("network.enable_wwan", func() error {
		return m.connectivity.EnableWWAN(ctx)
	})
}

func (m *Middleware) DisableWWAN(ctx context.Context) error {
	return m.call("network.disable_wwan", func() error {
		return m.connectivity.DisableWWAN(ctx)
	})
}

func (m *Middleware) EnableWifi(ctx context.Context) error {
	return m.call("network.enable_wifi", func() error {
		return m.connectivity.EnableWifi(ctx)
	})
}

func (m *Middleware) DisableWifi(ctx context.Context) error {
	return m.call("network.disable_wifi", func() error {
		return m.connectivity.DisableWifi(ctx)
	})
}
