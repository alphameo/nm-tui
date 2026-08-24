package logging

import (
	"context"
	"log/slog"

	"github.com/alphameo/nm-tui/internal/infra"
)

// NetworksMiddleware implements infra.NetworksManager by delegating to the
// wrapped implementation. Successes are logged at Debug level, failures at
// Error level along with the exit code of the failed command when the error
// is an [*exec.ExitError].
type NetworksMiddleware struct {
	middleware

	networks infra.NetworksManager
}

// NewNetworks returns a *NetworksMiddleware wrapping the given manager.
func NewNetworks(logger *slog.Logger, networks infra.NetworksManager) *NetworksMiddleware {
	return &NetworksMiddleware{
		middleware: middleware{logger: logger, prefix: "networks"},
		networks:   networks,
	}
}

func (m *NetworksMiddleware) ListNetworksWithRescan(ctx context.Context) ([]infra.AvailableNetwork, error) {
	return callResult(m.middleware, "scan_networks", func() ([]infra.AvailableNetwork, error) {
		return m.networks.ListNetworksWithRescan(ctx)
	})
}

func (m *NetworksMiddleware) ListNetworks(ctx context.Context) ([]infra.AvailableNetwork, error) {
	return callResult(m.middleware, "list_networks", func() ([]infra.AvailableNetwork, error) {
		return m.networks.ListNetworks(ctx)
	})
}

func (m *NetworksMiddleware) ListProfileNames(ctx context.Context) ([]string, error) {
	return callResult(m.middleware, "list_profile_names", func() ([]string, error) {
		return m.networks.ListProfileNames(ctx)
	})
}

func (m *NetworksMiddleware) ListProfiles(ctx context.Context) ([]infra.NetworkProfileShort, error) {
	return callResult(m.middleware, "list_profiles", func() ([]infra.NetworkProfileShort, error) {
		return m.networks.ListProfiles(ctx)
	})
}

func (m *NetworksMiddleware) ConnectToNetwork(ctx context.Context, ssid, password string) error {
	return m.call("connect_to_network", func() error {
		return m.networks.ConnectToNetwork(ctx, ssid, password)
	})
}

func (m *NetworksMiddleware) TryActivateNetwork(ctx context.Context, ssid string) error {
	return m.call("try_activate_network", func() error {
		return m.networks.TryActivateNetwork(ctx, ssid)
	})
}

func (m *NetworksMiddleware) CreateConnectionProfile(
	ctx context.Context, id, ssid, password string, hidden bool,
) error {
	return m.call("create_connection_profile", func() error {
		return m.networks.CreateConnectionProfile(ctx, id, ssid, password, hidden)
	})
}

func (m *NetworksMiddleware) CreateHotspotProfile(ctx context.Context, id string, ssid string, password string) error {
	return m.call("create_hotspot_profile", func() error {
		return m.networks.CreateHotspotProfile(ctx, id, ssid, password)
	})
}

func (m *NetworksMiddleware) QuickHotspot(ctx context.Context) error {
	return m.call("quick_hotspot", func() error {
		return m.networks.QuickHotspot(ctx)
	})
}

func (m *NetworksMiddleware) DeleteProfile(ctx context.Context, name string) error {
	return m.call("delete_profile", func() error {
		return m.networks.DeleteProfile(ctx, name)
	})
}

func (m *NetworksMiddleware) ActivateProfile(ctx context.Context, name string) error {
	return m.call("activate_profile", func() error {
		return m.networks.ActivateProfile(ctx, name)
	})
}

func (m *NetworksMiddleware) DeactivateProfile(ctx context.Context, name string) error {
	return m.call("deactivate_profile", func() error {
		return m.networks.DeactivateProfile(ctx, name)
	})
}

func (m *NetworksMiddleware) GetProfilePassword(ctx context.Context, name string) (string, error) {
	return callResult(m.middleware, "get_wifi_password", func() (string, error) {
		return m.networks.GetProfilePassword(ctx, name)
	})
}

func (m *NetworksMiddleware) GetProfile(ctx context.Context, name string) (infra.NetworkProfile, error) {
	return callResult(m.middleware, "get_profile", func() (infra.NetworkProfile, error) {
		return m.networks.GetProfile(ctx, name)
	})
}

func (m *NetworksMiddleware) UpdateProfile(ctx context.Context, name string, info infra.UpdateNetworkProfile) error {
	return m.call("update_profile", func() error {
		return m.networks.UpdateProfile(ctx, name, info)
	})
}
