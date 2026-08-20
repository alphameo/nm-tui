package logging

import (
	"context"
	"log/slog"

	"github.com/alphameo/nm-tui/internal/infra"
)

// DeviceMiddleware implements infra.DeviceManager by delegating to the wrapped
// implementation. Successes are logged at Debug level, failures at Error level
// along with the exit code of the failed command when the error is an
// [*exec.ExitError].
type DeviceMiddleware struct {
	middleware

	device infra.DeviceManager
}

// NewDevice returns a *DeviceMiddleware wrapping the given manager.
func NewDevice(logger *slog.Logger, device infra.DeviceManager) *DeviceMiddleware {
	return &DeviceMiddleware{
		middleware: middleware{logger: logger, prefix: "device"},
		device:     device,
	}
}

func (m *DeviceMiddleware) ListNetworkDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	return callResult(m.middleware, "list_devices", func() ([]infra.NetworkDevice, error) {
		return m.device.ListNetworkDevices(ctx)
	})
}

func (m *DeviceMiddleware) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	return callResult(m.middleware, "get_connectivity_status", func() (infra.ConnectivityStatus, error) {
		return m.device.GetConnectivityStatus(ctx)
	})
}

func (m *DeviceMiddleware) IsNetworkingEnabled(ctx context.Context) (bool, error) {
	return callResult(m.middleware, "is_networking_enabled", func() (bool, error) {
		return m.device.IsNetworkingEnabled(ctx)
	})
}

func (m *DeviceMiddleware) EnableNetworking(ctx context.Context) error {
	return m.call("enable_networking", func() error {
		return m.device.EnableNetworking(ctx)
	})
}

func (m *DeviceMiddleware) DisableNetworking(ctx context.Context) error {
	return m.call("disable_networking", func() error {
		return m.device.DisableNetworking(ctx)
	})
}

func (m *DeviceMiddleware) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
	return callResult(m.middleware, "get_radio_status", func() (infra.RadioStatus, error) {
		return m.device.GetRadioStatus(ctx)
	})
}

func (m *DeviceMiddleware) EnableWWAN(ctx context.Context) error {
	return m.call("enable_wwan", func() error {
		return m.device.EnableWWAN(ctx)
	})
}

func (m *DeviceMiddleware) DisableWWAN(ctx context.Context) error {
	return m.call("disable_wwan", func() error {
		return m.device.DisableWWAN(ctx)
	})
}

func (m *DeviceMiddleware) EnableWifi(ctx context.Context) error {
	return m.call("enable_wifi", func() error {
		return m.device.EnableWifi(ctx)
	})
}

func (m *DeviceMiddleware) DisableWifi(ctx context.Context) error {
	return m.call("disable_wifi", func() error {
		return m.device.DisableWifi(ctx)
	})
}
