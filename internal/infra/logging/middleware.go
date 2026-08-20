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
	logger   *slog.Logger
	networks infra.NetworksManager
	device   infra.DeviceManager
	portal   infra.CaptivePortalOpener
}

// New returns a *Middleware wrapping the given managers.
func New(
	logger *slog.Logger,
	networks infra.NetworksManager,
	device infra.DeviceManager,
	portal infra.CaptivePortalOpener,
) *Middleware {
	return &Middleware{logger: logger, networks: networks, device: device, portal: portal}
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

func (m *Middleware) ListNetworksWithRescan(ctx context.Context) ([]infra.AvailableNetwork, error) {
	return callResult(m, "networks.scan_networks", func() ([]infra.AvailableNetwork, error) {
		return m.networks.ListNetworksWithRescan(ctx)
	})
}

func (m *Middleware) ListNetworks(ctx context.Context) ([]infra.AvailableNetwork, error) {
	return callResult(m, "networks.list_networks", func() ([]infra.AvailableNetwork, error) {
		return m.networks.ListNetworks(ctx)
	})
}

func (m *Middleware) ListProfileNames(ctx context.Context) ([]string, error) {
	return callResult(m, "networks.list_profile_names", func() ([]string, error) {
		return m.networks.ListProfileNames(ctx)
	})
}

func (m *Middleware) ListProfiles(ctx context.Context) ([]infra.NetworkProfileShort, error) {
	return callResult(m, "networks.list_profiles", func() ([]infra.NetworkProfileShort, error) {
		return m.networks.ListProfiles(ctx)
	})
}

func (m *Middleware) ConnectToNetwork(ctx context.Context, ssid, password string) error {
	return m.call("networks.connect_to_network", func() error {
		return m.networks.ConnectToNetwork(ctx, ssid, password)
	})
}

func (m *Middleware) TryActivateNetwork(ctx context.Context, ssid string) error {
	return m.call("networks.try_activate_network", func() error {
		return m.networks.TryActivateNetwork(ctx, ssid)
	})
}

func (m *Middleware) CreateConnectionProfile(ctx context.Context, id, ssid, password string, hidden bool) error {
	return m.call("networks.create_connection_profile", func() error {
		return m.networks.CreateConnectionProfile(ctx, id, ssid, password, hidden)
	})
}

func (m *Middleware) CreateHotspotProfile(ctx context.Context, id string, ssid string, password string) error {
	return m.call("networks.create_hotspot_profile", func() error {
		return m.networks.CreateHotspotProfile(ctx, id, ssid, password)
	})
}

func (m *Middleware) QuickHotspot(ctx context.Context) error {
	return m.call("networks.quick_hotspot", func() error {
		return m.networks.QuickHotspot(ctx)
	})
}

func (m *Middleware) DeleteProfile(ctx context.Context, name string) error {
	return m.call("networks.delete_profile", func() error {
		return m.networks.DeleteProfile(ctx, name)
	})
}

func (m *Middleware) ActivateProfile(ctx context.Context, name string) error {
	return m.call("networks.activate_profile", func() error {
		return m.networks.ActivateProfile(ctx, name)
	})
}

func (m *Middleware) DeactivateProfile(ctx context.Context, name string) error {
	return m.call("networks.deactivate_profile", func() error {
		return m.networks.DeactivateProfile(ctx, name)
	})
}

func (m *Middleware) GetProfilePassword(ctx context.Context, name string) (string, error) {
	return callResult(m, "networks.get_wifi_password", func() (string, error) {
		return m.networks.GetProfilePassword(ctx, name)
	})
}

func (m *Middleware) GetProfile(ctx context.Context, name string) (infra.NetworkProfile, error) {
	return callResult(m, "networks.get_profile", func() (infra.NetworkProfile, error) {
		return m.networks.GetProfile(ctx, name)
	})
}

func (m *Middleware) UpdateProfile(ctx context.Context, name string, info infra.UpdateProfile) error {
	return m.call("networks.update_profile", func() error {
		return m.networks.UpdateProfile(ctx, name, info)
	})
}

func (m *Middleware) OpenCaptivePortal(ctx context.Context) error {
	return m.call("portal.open_captive_portal", func() error {
		return m.portal.OpenCaptivePortal(ctx)
	})
}

func (m *Middleware) ListNetworkDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	return callResult(m, "device.list_devices", func() ([]infra.NetworkDevice, error) {
		return m.device.ListNetworkDevices(ctx)
	})
}

func (m *Middleware) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	return callResult(m, "device.get_connectivity_status", func() (infra.ConnectivityStatus, error) {
		return m.device.GetConnectivityStatus(ctx)
	})
}

func (m *Middleware) IsNetworkingEnabled(ctx context.Context) (bool, error) {
	return callResult(m, "device.is_networking_enabled", func() (bool, error) {
		return m.device.IsNetworkingEnabled(ctx)
	})
}

func (m *Middleware) EnableNetworking(ctx context.Context) error {
	return m.call("device.enable_networking", func() error {
		return m.device.EnableNetworking(ctx)
	})
}

func (m *Middleware) DisableNetworking(ctx context.Context) error {
	return m.call("device.disable_networking", func() error {
		return m.device.DisableNetworking(ctx)
	})
}

func (m *Middleware) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
	return callResult(m, "device.get_radio_status", func() (infra.RadioStatus, error) {
		return m.device.GetRadioStatus(ctx)
	})
}

func (m *Middleware) EnableWWAN(ctx context.Context) error {
	return m.call("device.enable_wwan", func() error {
		return m.device.EnableWWAN(ctx)
	})
}

func (m *Middleware) DisableWWAN(ctx context.Context) error {
	return m.call("device.disable_wwan", func() error {
		return m.device.DisableWWAN(ctx)
	})
}

func (m *Middleware) EnableWifi(ctx context.Context) error {
	return m.call("device.enable_wifi", func() error {
		return m.device.EnableWifi(ctx)
	})
}

func (m *Middleware) DisableWifi(ctx context.Context) error {
	return m.call("device.disable_wifi", func() error {
		return m.device.DisableWifi(ctx)
	})
}
