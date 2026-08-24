// Package nm provides NetworkManager api
package nm

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/alphameo/nm-tui/internal/infra"
)

type CLI struct{}

func NewCLI() *CLI {
	return &CLI{}
}

const CommandName = "nmcli"

const (
	KeyMgmgtNone  string = "none"
	KeyMgmtWpaPsk string = "wpa-psk"
	// Add sae if errors.
)

// run executes nmcli with the given args and returns its stdout. On failure
// the returned error wraps opErr and the underlying [*exec.ExitError].
func (n *CLI) run(ctx context.Context, opErr error, args ...string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", opErr, err, infra.ExtractStderr(err))
	}
	return out, nil
}

func (n *CLI) ListNetworkDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	args := []string{"-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device", "status"}
	out, err := n.run(ctx, infra.ErrListNetworkDevices, args...)
	if err != nil {
		return nil, err
	}

	var res []infra.NetworkDevice
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		res = append(res, infra.NetworkDevice{
			Device:     parts[0],
			Type:       parts[1],
			State:      parts[2],
			Connection: parts[3],
		})
	}
	return res, nil
}

const listNetworksFlags = "SSID,IN-USE,SECURITY,SIGNAL,BAND,RATE,DEVICE,MODE"

func (n *CLI) ListNetworksWithRescan(ctx context.Context) ([]infra.AvailableNetwork, error) {
	args := []string{
		"-t", "-f", listNetworksFlags,
		"device", "wifi", "list", "--rescan", "yes",
	}
	out, err := n.run(ctx, infra.ErrScanNetworks, args...)
	if err != nil {
		return nil, err
	}

	return parseNetworks(string(out))
}

func (n *CLI) ListNetworks(ctx context.Context) ([]infra.AvailableNetwork, error) {
	args := []string{
		"-t", "-f", listNetworksFlags,
		"device", "wifi", "list",
	}
	out, err := n.run(ctx, infra.ErrListNetworks, args...)
	if err != nil {
		return nil, err
	}

	return parseNetworks(string(out))
}

func parseNetworks(networks string) ([]infra.AvailableNetwork, error) {
	var res []infra.AvailableNetwork
	var errs []error
	lines := strings.SplitSeq(networks, "\n")
	for line := range lines {
		var netErrs []error
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		ssid := parts[0]

		signal, err := parseSignal(parts[3])
		if err != nil {
			netErrs = append(netErrs, err)
		}

		band, err := parseBandGHz(parts[4])
		if err != nil {
			netErrs = append(netErrs, err)
		}

		rate, err := parseRateMbits(parts[5])
		if err != nil {
			netErrs = append(netErrs, err)
		}

		mode, err := parseNetworkMode(parts[7])
		if err != nil {
			netErrs = append(netErrs, err)
		}

		res = append(res, infra.AvailableNetwork{
			SSID:          ssid,
			Active:        parts[1] == "*",
			SecurityMode:  parts[2],
			Signal:        signal,
			Band:          band,
			Rate:          rate,
			LookingDevice: parts[6],
			NetworkMode:   mode,
		})
		if len(netErrs) > 0 {
			errs = append(errs, fmt.Errorf("error during parse %s: %w", ssid, errors.Join(netErrs...)))
		}
	}
	return res, errors.Join(errs...)
}

func parseSignal(token string) (int, error) {
	signal, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("can't convert to signal: got %q", token)
	}
	return signal, nil
}

func parseRateMbits(token string) (float64, error) {
	rate, ok := strings.CutSuffix(strings.TrimSpace(token), " Mbit/s")
	if !ok {
		return 0, fmt.Errorf("can't parse rate: got %q", token)
	}
	v, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0, fmt.Errorf("can't convert to rate: %w", err)
	}
	return v, nil
}

func parseBandGHz(token string) (float64, error) {
	rate, ok := strings.CutSuffix(strings.TrimSpace(token), " GHz")
	if !ok {
		return 0, fmt.Errorf("can't parse band: got %q", token)
	}
	v, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0, fmt.Errorf("can't convert to band: %w", err)
	}
	return v, nil
}

func (n *CLI) ListProfiles(ctx context.Context) ([]infra.NetworkProfileShort, error) {
	args := []string{"-t", "-f", "UUID,NAME,STATE", "connection", "show"}
	out, err := n.run(ctx, infra.ErrListProfiles, args...)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var res []infra.NetworkProfileShort
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		name := parts[1]
		if name == "lo" {
			continue
		}

		uuid := parts[0]
		ssid, err := n.getProfileSSID(ctx, uuid)
		if err != nil {
			ssid = ""
		}
		wg.Add(1)
		wifi := infra.NetworkProfileShort{
			UUID:   uuid,
			Name:   name,
			SSID:   ssid,
			Active: parts[2] == "activated",
			Mode:   infra.NetworkNil,
		}
		res = append(res, wifi)
		go func(idx int) {
			defer wg.Done()
			mode, err := n.getProfileMode(ctx, name)
			if err != nil {
				mode = infra.NetworkNil
			}
			mu.Lock()
			defer mu.Unlock()
			res[idx].Mode = mode
		}(len(res) - 1)
	}

	wg.Wait()

	return res, nil
}

func (n *CLI) CreateConnectionProfile(ctx context.Context, name, ssid, password string, hidden bool) error {
	hiddenStr := "no"
	if hidden {
		hiddenStr = "yes"
	}
	var keyMgmt string
	if len(password) == 0 {
		keyMgmt = KeyMgmgtNone
	} else {
		keyMgmt = KeyMgmtWpaPsk
	}
	args := []string{
		"connection", "add", "type", "wifi",
		"con-name", name,
		"ssid", ssid,
		"wifi.hidden", hiddenStr,
		"wifi-sec.key-mgmt", keyMgmt,
		"wifi-sec.psk", password,
	}
	_, err := n.run(ctx, infra.ErrCreateWifiConnection, args...)
	return err
}

func (n *CLI) ConnectToNetwork(ctx context.Context, ssid, password string) error {
	args := []string{
		"device", "wifi", "connect", ssid,
		"password", password,
	}
	_, err := n.run(ctx, infra.ErrConnectToNetwork, args...)
	return err
}

func (n *CLI) TryActivateNetwork(ctx context.Context, ssid string) error {
	args := []string{
		"device", "wifi", "connect", ssid,
	}
	_, err := n.run(ctx, infra.ErrTryActivateNetwork, args...)
	return err
}

func (n *CLI) ActivateProfile(ctx context.Context, id string) error {
	args := []string{"connection", "up", id}
	_, err := n.run(ctx, infra.ErrActivateProfile, args...)
	return err
}

func (n *CLI) DeactivateProfile(ctx context.Context, id string) error {
	args := []string{"connection", "down", id}
	_, err := n.run(ctx, infra.ErrDeactivateProfile, args...)
	return err
}

func (n *CLI) ListProfileNames(ctx context.Context) ([]string, error) {
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := n.run(ctx, infra.ErrListProfileNames, args...)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}

func (n *CLI) GetProfilePassword(ctx context.Context, id string) (string, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "802-11-wireless-security.psk",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetProfilePassword, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getProfileSSID(ctx context.Context, id string) (string, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "802-11-wireless.ssid",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: ssid", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getProfileID(ctx context.Context, uuid string) (string, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "connection.id",
		"connection", "show", uuid,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: ssid", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getProfileUUID(ctx context.Context, id string) (string, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "connection.uuid",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: ssid", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getProfileAutoconn(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "connection.autoconnect",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: autoconnection", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "yes", nil
}

func (n *CLI) getProfileHidden(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "802-11-wireless.hidden",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: autoconnection", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "yes", nil
}

func (n *CLI) getProfileAutoconnPriority(ctx context.Context, id string) (int, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "connection.autoconnect-priority",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: autoconnection priority", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return 0, err
	}
	autoconnectResp := strings.TrimSpace(string(out))
	autoconnectPriority, err := strconv.Atoi(autoconnectResp)
	if err != nil {
		return 0, fmt.Errorf("%w: autoconnection priority: %w", infra.ErrGetProfileProperty, err)
	}
	return autoconnectPriority, nil
}

func (n *CLI) getProfileActive(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "GENERAL.STATE",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: is active", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "activated", nil
}

func (n *CLI) getProfileMode(ctx context.Context, id string) (infra.NetworkMode, error) {
	args := []string{
		"-m", "tabular",
		"-t", "-f", "802-11-wireless.mode",
		"connection", "show", id,
	}
	out, err := n.run(ctx, fmt.Errorf("%w: mode", infra.ErrGetProfileProperty), args...)
	if err != nil {
		return infra.NetworkNil, err
	}
	res := strings.TrimSpace(string(out))
	mode, err := parseNetworkMode(res)
	if err != nil {
		return infra.NetworkNil, fmt.Errorf("%w: mode: %w", infra.ErrGetProfileProperty, err)
	}
	return mode, nil
}

func parseNetworkMode(mode string) (infra.NetworkMode, error) {
	var res infra.NetworkMode
	switch strings.ToLower(mode) {
	case "infrastructure", "infra":
		res = infra.NetworkInfra
	case "ap":
		res = infra.NetworkAccessPoint
	case "adhoc":
		res = infra.NetworkAdHoc
	case "mesh":
		res = infra.NetworkMesh
	}

	if res == infra.NetworkNil {
		return infra.NetworkNil, fmt.Errorf("can't parse network mode: got %q", mode)
	}
	return res, nil
}

// setFetchResult stores a fetched field into dst under mu, collecting any error.
func setFetchResult[T any](mu *sync.Mutex, errs *[]error, dst *T, value T, err error) {
	mu.Lock()
	defer mu.Unlock()
	if err != nil {
		*errs = append(*errs, err)
	}
	*dst = value
}

//nolint:funlen // close undevidable operations
func (n *CLI) GetProfile(ctx context.Context, id string) (infra.NetworkProfile, error) {
	var errs []error
	info := infra.NetworkProfile{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(9)

	go func() {
		defer wg.Done()
		ssid, err := n.getProfileSSID(ctx, id)
		setFetchResult(&mu, &errs, &info.SSID, ssid, err)
	}()

	go func() {
		defer wg.Done()
		uuid, err := n.getProfileUUID(ctx, id)
		setFetchResult(&mu, &errs, &info.UUID, uuid, err)
	}()

	go func() {
		defer wg.Done()
		uuid, err := n.getProfileID(ctx, id)
		setFetchResult(&mu, &errs, &info.Name, uuid, err)
	}()
	go func() {
		defer wg.Done()
		password, err := n.GetProfilePassword(ctx, id)
		setFetchResult(&mu, &errs, &info.Password, password, err)
	}()

	go func() {
		defer wg.Done()
		autoconnect, err := n.getProfileAutoconn(ctx, id)
		setFetchResult(&mu, &errs, &info.Autoconnect, autoconnect, err)
	}()

	go func() {
		defer wg.Done()
		autoconnectPriority, err := n.getProfileAutoconnPriority(ctx, id)
		setFetchResult(&mu, &errs, &info.AutoconnectPriority, autoconnectPriority, err)
	}()

	go func() {
		defer wg.Done()
		activated, err := n.getProfileActive(ctx, id)
		setFetchResult(&mu, &errs, &info.Active, activated, err)
	}()

	go func() {
		defer wg.Done()
		mode, err := n.getProfileMode(ctx, id)
		setFetchResult(&mu, &errs, &info.Mode, mode, err)
	}()

	go func() {
		defer wg.Done()
		hidden, err := n.getProfileHidden(ctx, id)
		setFetchResult(&mu, &errs, &info.Hidden, hidden, err)
	}()

	wg.Wait()

	if len(errs) != 0 {
		sb := strings.Builder{}
		for i, err := range errs {
			sb.WriteString(err.Error())
			if i != 0 {
				sb.WriteString("\n")
			}
		}
		bigErrStr := sb.String()
		return infra.NetworkProfile{}, fmt.Errorf("%w for %s: %s", infra.ErrGetProfile, id, bigErrStr)
	}

	return info, nil
}

func (n *CLI) UpdateProfile(ctx context.Context, id string, info infra.UpdateNetworkProfile) error {
	var keyMgmgt string
	if len(info.Password) == 0 {
		keyMgmgt = KeyMgmgtNone
	} else {
		keyMgmgt = KeyMgmtWpaPsk
	}
	var autoconnect string
	if info.Autoconnect {
		autoconnect = "yes"
	} else {
		autoconnect = "no"
	}
	args := []string{
		"connection", "modify", id,
		"connection.id", info.Name,
		"802-11-wireless-security.key-mgmt", keyMgmgt,
		"802-11-wireless-security.psk", info.Password,
		"connection.autoconnect", autoconnect,
		"connection.autoconnect-priority", strconv.Itoa(info.AutoconnectPriority),
	}
	_, err := n.run(ctx, infra.ErrUpdateProfile, args...)
	return err
}

func (n *CLI) DeleteProfile(ctx context.Context, id string) error {
	args := []string{"connection", "delete", id}
	_, err := n.run(ctx, infra.ErrDeleteProfile, args...)
	return err
}

func (n *CLI) GetWifiStatus(ctx context.Context) (bool, error) {
	args := []string{"radio", "wifi"}
	out, err := n.run(ctx, infra.ErrGetWifiStatus, args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "enabled", nil
}

func (n *CLI) GetWWANStatus(ctx context.Context) (bool, error) {
	args := []string{"radio", "wwan"}
	out, err := n.run(ctx, infra.ErrGetWWANStatus, args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "enabled", nil
}

func (n *CLI) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
	var errs []error
	wifi, err := n.GetWifiStatus(ctx)
	if err != nil {
		errs = append(errs, err)
	}
	wwan, err := n.GetWWANStatus(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		sb := strings.Builder{}
		for i, err := range errs {
			sb.WriteString(err.Error())
			if i != 0 {
				sb.WriteString("\n")
			}
		}
		bigErrStr := sb.String()
		return infra.RadioStatus{}, fmt.Errorf("%w: %s", infra.ErrGetRadioStatus, bigErrStr)
	}

	return infra.RadioStatus{
		EnabledWifi: wifi,
		EnabledWWAN: wwan,
	}, nil
}

func (n *CLI) EnableWifi(ctx context.Context) error {
	args := []string{"radio", "wifi", "on"}
	_, err := n.run(ctx, infra.ErrEnableWifi, args...)
	return err
}

func (n *CLI) DisableWifi(ctx context.Context) error {
	args := []string{"radio", "wifi", "off"}
	_, err := n.run(ctx, infra.ErrDisableWifi, args...)
	return err
}

func (n *CLI) EnableWWAN(ctx context.Context) error {
	args := []string{"radio", "wwan", "on"}
	_, err := n.run(ctx, infra.ErrEnableWWAN, args...)
	return err
}

func (n *CLI) DisableWWAN(ctx context.Context) error {
	args := []string{"radio", "wwan", "off"}
	_, err := n.run(ctx, infra.ErrDisableWWAN, args...)
	return err
}

func (n *CLI) IsNetworkingEnabled(ctx context.Context) (bool, error) {
	args := []string{"networking"}
	out, err := n.run(ctx, infra.ErrIsNetworkingEnabled, args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "enabled", nil
}

func (n *CLI) EnableNetworking(ctx context.Context) error {
	args := []string{"networking", "on"}
	_, err := n.run(ctx, infra.ErrEnableNetworking, args...)
	return err
}

func (n *CLI) DisableNetworking(ctx context.Context) error {
	args := []string{"networking", "off"}
	_, err := n.run(ctx, infra.ErrDisableNetworking, args...)
	return err
}

func (n *CLI) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	args := []string{"networking", "connectivity", "check"}
	out, err := n.run(ctx, infra.ErrGetConnectivityStatus, args...)
	if err != nil {
		return infra.ConnectvityNil, err
	}
	res := strings.TrimSpace(string(out))
	var mode infra.ConnectivityStatus
	switch res {
	case "none":
		mode = infra.ConnectivityNone
	case "portal":
		mode = infra.ConnectivityPortal
	case "limited":
		mode = infra.ConnectivityLimited
	case "full":
		mode = infra.ConnectivityFull
	case "unknown":
		mode = infra.ConnectivityUnknown
	}
	if mode == infra.ConnectvityNil {
		return infra.ConnectvityNil, fmt.Errorf("%w: got %s", infra.ErrParseConnectivity, res)
	}
	return mode, nil
}

func (n *CLI) CreateHotspotProfile(ctx context.Context, name string, ssid string, password string) error {
	args := []string{
		"device", "wifi", "hotspot",
		"con-name", name,
		"ssid", ssid,
		"password", password,
	}
	_, err := n.run(ctx, infra.ErrCreateHotspotProfile, args...)
	return err
}

func (n *CLI) QuickHotspot(ctx context.Context) error {
	args := []string{"device", "wifi", "hotspot"}
	_, err := n.run(ctx, infra.ErrQuickHotspot, args...)
	return err
}

func (n *CLI) GetDeviceInfo(ctx context.Context, name string) (string, error) {
	args := []string{"device", "show", name}
	out, err := n.run(ctx, infra.ErrGetDeviceInfo, args...)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
