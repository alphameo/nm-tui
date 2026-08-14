// Package nm provides NetworkManager api
package nm

import (
	"context"
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

func (n *CLI) ListDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	args := []string{"-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device", "status"}
	out, err := n.run(ctx, infra.ErrListDevices, args...)
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

func (n *CLI) ScanNetworks(ctx context.Context) ([]infra.AvailableNetwork, error) {
	args := []string{
		"-t", "-f", "SSID,IN-USE,SECURITY,SIGNAL",
		"device", "wifi", "list", "--rescan", "yes",
	}
	out, err := n.run(ctx, infra.ErrScanNetworks, args...)
	if err != nil {
		return nil, err
	}

	var res []infra.AvailableNetwork
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		ssid := parts[0]
		if ssid == "" {
			continue
		}

		signal, err := strconv.Atoi(parts[3])
		if err != nil {
			signal = 0
		}
		res = append(res, infra.AvailableNetwork{
			SSID:     ssid,
			Active:   parts[1] == "*",
			Security: parts[2],
			Signal:   signal,
		})
	}
	return res, nil
}

func (n *CLI) ListProfiles(ctx context.Context) ([]infra.NetworkProfileShort, error) {
	args := []string{"-t", "-f", "NAME,STATE", "connection", "show"}
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
		if len(parts) < 2 {
			continue
		}
		if parts[0] == "lo" {
			continue
		}

		name := parts[0]
		ssid, err := n.getWifiSSID(ctx, name)
		if err != nil {
			ssid = ""
		}
		wg.Add(1)
		wifi := infra.NetworkProfileShort{
			Name:   name,
			SSID:   ssid,
			Active: parts[1] == "activated",
			Mode:   infra.NetworkNil,
		}
		res = append(res, wifi)
		go func(idx int) {
			defer wg.Done()
			mode, err := n.getNetMode(ctx, name)
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
	out, err := n.run(ctx, infra.ErrGetWifiPassword, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getWifiSSID(ctx context.Context, id string) (string, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "802-11-wireless.ssid",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetWifiSSID, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (n *CLI) getWifiAutoconnect(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "connection.autoconnect",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetWifiAutoconnect, args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "yes", nil
}

func (n *CLI) getWifiAutoconnectPriority(ctx context.Context, id string) (int, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "connection.autoconnect-priority",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetWifiAutoconnectPriority, args...)
	if err != nil {
		return 0, err
	}
	autoconnectResp := strings.TrimSpace(string(out))
	autoconnectPriority, err := strconv.Atoi(autoconnectResp)
	if err != nil {
		return 0, fmt.Errorf("%w %s: %w", infra.ErrGetWifiAutoconnectPriority, id, err)
	}
	return autoconnectPriority, nil
}

func (n *CLI) getWifiActive(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "GENERAL.STATE",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetWifiActivity, args...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "activated", nil
}

func (n *CLI) getNetMode(ctx context.Context, id string) (infra.NetworkMode, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "802-11-wireless.mode",
		"connection", "show", id,
	}
	out, err := n.run(ctx, infra.ErrGetNetMode, args...)
	if err != nil {
		return infra.NetworkNil, err
	}
	res := strings.TrimSpace(string(out))
	var mode infra.NetworkMode
	switch res {
	case "infrastructure":
		mode = infra.NetworkInfra
	case "ap":
		mode = infra.NetworkAccessPoint
	case "adhoc":
		mode = infra.NetworkAdHoc
	case "mesh":
		mode = infra.NetworkMesh
	}
	if mode == infra.NetworkNil {
		return infra.NetworkNil, fmt.Errorf("%w for %s: got %s", infra.ErrParseNetMode, id, res)
	}
	return mode, nil
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

func (n *CLI) GetProfile(ctx context.Context, id string) (infra.NetworkProfile, error) {
	var errs []error
	info := infra.NetworkProfile{
		Name: id,
	}
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(6)

	go func() {
		defer wg.Done()
		ssid, err := n.getWifiSSID(ctx, id)
		setFetchResult(&mu, &errs, &info.SSID, ssid, err)
	}()

	go func() {
		defer wg.Done()
		password, err := n.GetProfilePassword(ctx, id)
		setFetchResult(&mu, &errs, &info.Password, password, err)
	}()

	go func() {
		defer wg.Done()
		autoconnect, err := n.getWifiAutoconnect(ctx, id)
		setFetchResult(&mu, &errs, &info.Autoconnect, autoconnect, err)
	}()

	go func() {
		defer wg.Done()
		autoconnectPriority, err := n.getWifiAutoconnectPriority(ctx, id)
		setFetchResult(&mu, &errs, &info.AutoconnectPriority, autoconnectPriority, err)
	}()

	go func() {
		defer wg.Done()
		activated, err := n.getWifiActive(ctx, id)
		setFetchResult(&mu, &errs, &info.Active, activated, err)
	}()

	go func() {
		defer wg.Done()
		mode, err := n.getNetMode(ctx, id)
		setFetchResult(&mu, &errs, &info.Mode, mode, err)
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

func (n *CLI) UpdateProfile(ctx context.Context, id string, info infra.UpdateProfile) error {
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
		"connection", "modify",
		id, "connection.id", info.Name,
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
