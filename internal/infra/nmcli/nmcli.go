// Package nmcli provides nmcli api
package nmcli

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
)

type NMCLI struct{}

func New() *NMCLI {
	return &NMCLI{}
}

const CommandName = "nmcli"

func (NMCLI) GetNetworkDevices(ctx context.Context) ([]infra.NetworkDevice, error) {
	args := []string{"-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device", "status"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetNetworkDevices.Error(),
			"err", err,
			"stderr", stderr,
		)
		return nil, fmt.Errorf("%w: %s", infra.ErrGetNetworkDevices, stderr)
	}

	var res []infra.NetworkDevice
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			slog.Warn("malformed network device line", "line", line)
			continue
		}

		res = append(res, infra.NetworkDevice{
			Device:     parts[0],
			Type:       parts[1],
			State:      parts[2],
			Connection: parts[3],
		})
	}
	slog.Info("scanned network devices")
	return res, nil
}

func (NMCLI) ScanWifis(ctx context.Context) ([]infra.AvailableWifi, error) {
	args := []string{
		"-t", "-f", "SSID,IN-USE,SECURITY,SIGNAL",
		"device", "wifi", "list", "--rescan", "yes",
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrScanWifis.Error(),
			"err", err,
			"stderr", stderr,
		)
		return nil, fmt.Errorf("%w: %s", infra.ErrScanWifis, stderr)
	}

	var res []infra.AvailableWifi
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			slog.Warn("malformed wifi line", "line", line)
			continue
		}

		ssid := parts[0]
		if ssid == "" {
			continue
		}

		signal, err := strconv.Atoi(parts[3])
		if err != nil {
			slog.Warn("parsing signal strength", "line", line, "error", err)
			signal = 0
		}
		res = append(res, infra.AvailableWifi{
			SSID:     ssid,
			Active:   parts[1] == "*",
			Security: parts[2],
			Signal:   signal,
		})
	}
	slog.Info("scanned available wifi networks")
	return res, nil
}

func (n NMCLI) GetSavedWifis(ctx context.Context) ([]infra.SavedWifi, error) {
	args := []string{"-t", "-f", "NAME,STATE", "connection", "show"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		return nil, fmt.Errorf("%w: %s", infra.ErrGetSavedWifis, stderr)
	}

	var res []infra.SavedWifi

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
			slog.Warn(
				"failed to get ssid for saved wifi",
				"name", name,
				"ssid", ssid,
				"error", err,
			)
		}
		res = append(res, infra.SavedWifi{
			Name:   name,
			SSID:   ssid,
			Active: parts[1] == "activated",
		})
	}

	slog.Info("retrieved saved wifi networks")
	return res, nil
}

func (n NMCLI) CreateWifiConnection(ctx context.Context, id, ssid, password, device string, hidden bool) error {
	hiddenStr := "no"
	if hidden {
		hiddenStr = "yes"
	}
	args := []string{
		"connection", "add", "type", "wifi",
		"con-name", id,
		"ifname", device,
		"ssid", ssid,
		"wifi.hidden", hiddenStr,
		"wifi-sec.key-mgmt", "wpa-psk", // use "sae" on fail
		"wifi-sec.psk", password,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrCreateWifiConnection.Error(),
			"ssid", ssid,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrCreateWifiConnection, ssid, stderr)
	}
	slog.Info("created wifi connection",
		"id", id,
		"ssid", ssid,
		"device", device,
		"hidden", hidden,
		"output", string(out))
	return nil
}

func (n NMCLI) ConnectWifi(ctx context.Context, ssid, password string, hidden bool) error {
	hiddenStr := "no"
	if hidden {
		hiddenStr = "yes"
	}
	args := []string{
		"device", "wifi", "connect", ssid,
		"password", password,
		"hidden", hiddenStr,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrConnectWifi.Error(),
			"ssid", ssid,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrConnectWifi, ssid, stderr)
	}
	slog.Info("connected to wifi", "ssid", ssid, "output", string(out))
	return nil
}

func (NMCLI) ActivateWifi(ctx context.Context, id string) error {
	args := []string{"connection", "up", id}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrConnectSavedWifi.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrConnectSavedWifi, id, stderr)
	}
	slog.Info("connected to saved wifi", "id", id, "output", string(out))
	return nil
}

func (NMCLI) DeactivateWifi(ctx context.Context, id string) error {
	args := []string{"connection", "down", id}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrDisconnectWifi.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrDisconnectWifi, id, stderr)
	}
	slog.Info("disconnected from wifi", "id", id, "output", string(out))
	return nil
}

func (NMCLI) GetSavedWifiSSIDs(ctx context.Context) ([]string, error) {
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		return nil, fmt.Errorf("%w: %s", infra.ErrGetSavedWifiSSIDs, stderr)
	}
	slog.Info("retrieved saved wifi networks")
	return strings.Split(string(out), "\n"), nil
}

func (NMCLI) GetWifiPassword(ctx context.Context, id string) (string, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "802-11-wireless-security.psk",
		"connection", "show", id,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiPassword.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return "", fmt.Errorf("%w for %s: %s", infra.ErrGetWifiPassword, id, stderr)
	}
	slog.Info("retrieved wifi password", "id", id)
	return strings.Trim(string(out), " \n"), nil
}

func (NMCLI) getWifiSSID(ctx context.Context, id string) (string, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "802-11-wireless.ssid",
		"connection", "show", id,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiSSID.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return "", fmt.Errorf("%w %s: %s", infra.ErrGetWifiSSID, id, stderr)
	}
	slog.Info("retrieved wifi ssid", "id", id)
	return strings.Trim(string(out), " \n"), nil
}

func (NMCLI) getWifiAutoconnect(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "connection.autoconnect",
		"connection", "show", id,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiAutoconnect.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return false, fmt.Errorf("%w for %s: %s", infra.ErrGetWifiAutoconnect, id, stderr)
	}
	return strings.Trim(string(out), " \n") == "yes", nil
}

func (NMCLI) getWifiAutoconnectPriority(ctx context.Context, id string) (int, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "connection.autoconnect-priority",
		"connection", "show", id,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiAutoconnectPriority.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return 0, fmt.Errorf("%w for %s: %s", infra.ErrGetWifiAutoconnectPriority, id, stderr)
	}
	autoconnectResp := strings.Trim(string(out), " \n")
	autoconnectPriority, err := strconv.Atoi(autoconnectResp)
	if err != nil {
		slog.Error(
			"parsing autoconnect priority",
			"id", id,
			"value", autoconnectResp,
			"error", err,
		)
		return 0, fmt.Errorf("%w %s: %w", infra.ErrGetWifiAutoconnectPriority, id, err)
	}
	return autoconnectPriority, nil
}

func (NMCLI) getWifiActive(ctx context.Context, id string) (bool, error) {
	args := []string{
		"-s", "-m", "tabular",
		"-t", "-f", "GENERAL.STATE",
		"connection", "show", id,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiActivity.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return false, fmt.Errorf("%w for %s: %s", infra.ErrGetWifiActivity, id, stderr)
	}
	return strings.Trim(string(out), " \n") == "activated", nil
}

func (n *NMCLI) GetWifiInfo(ctx context.Context, id string) (infra.WifiInfo, error) {
	var errs []error
	ssid, err := n.getWifiSSID(ctx, id)
	if err != nil {
		errs = append(errs, err)
	}

	password, err := n.GetWifiPassword(ctx, id)
	if err != nil {
		errs = append(errs, err)
	}

	autoconnect, err := n.getWifiAutoconnect(ctx, id)
	if err != nil {
		errs = append(errs, err)
	}

	autoconnectPriority, err := n.getWifiAutoconnectPriority(ctx, id)
	if err != nil {
		errs = append(errs, err)
	}

	activated, err := n.getWifiActive(ctx, id)
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
		slog.Error(
			infra.ErrGetWifiInfo.Error(),
			"id", id,
			"failed operations", bigErrStr,
		)
		return infra.WifiInfo{}, fmt.Errorf("%w for %s: %s", infra.ErrGetWifiInfo, id, bigErrStr)
	}

	return infra.WifiInfo{
		Name:                id,
		SSID:                ssid,
		Password:            password,
		Autoconnect:         autoconnect,
		AutoconnectPriority: autoconnectPriority,
		Active:              activated,
	}, nil
}

// UpdateWifiInfo is not atomic
func (n NMCLI) UpdateWifiInfo(ctx context.Context, id string, info infra.UpdateWifiInfo) error {
	var errs []error

	err := n.updateWifiID(
		ctx,
		id,
		info.Name,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiPassword(
		ctx,
		info.Name,
		info.Password,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiAutoconnect(
		ctx,
		info.Name,
		info.Autoconnect,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiAutoconnectPriority(
		ctx,
		info.Name,
		info.AutoconnectPriority,
	)
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
		slog.Error(
			infra.ErrGetWifiInfo.Error(),
			"id", id,
			"failed operations", bigErrStr,
		)
		return fmt.Errorf("%w for %s: %s", infra.ErrUpdateWifiInfo, id, bigErrStr)
	}
	return nil
}

func (NMCLI) updateWifiInfoField(ctx context.Context, id, field, value string) error {
	args := []string{"connection", "modify", id, field, value}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrUpdateWifiInfoField.Error(),
			"id", id,
			"field", field,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s for %s: %s", infra.ErrUpdateWifiInfoField, field, id, stderr)
	}
	slog.Info(
		"modified wifi field",
		"id", id,
		"field", field,
		"output", string(out),
	)
	return nil
}

func (n NMCLI) updateWifiID(ctx context.Context, id, newID string) error {
	return n.updateWifiInfoField(ctx, id, "connection.id", newID)
}

func (n NMCLI) updateWifiPassword(ctx context.Context, id, password string) error {
	return n.updateWifiInfoField(ctx, id, "802-11-wireless-security.psk", password)
}

func (n NMCLI) updateWifiAutoconnect(ctx context.Context, id string, autoconnect bool) error {
	var autoconnectValue string
	if autoconnect {
		autoconnectValue = "yes"
	} else {
		autoconnectValue = "no"
	}

	return n.updateWifiInfoField(ctx, id, "connection.autoconnect", autoconnectValue)
}

func (n NMCLI) updateWifiAutoconnectPriority(ctx context.Context, id string, priority int) error {
	return n.updateWifiInfoField(
		ctx,
		id,
		"connection.autoconnect-priority",
		strconv.Itoa(priority),
	)
}

func (NMCLI) DeleteWifiConnection(ctx context.Context, id string) error {
	args := []string{"connection", "delete", id}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrDeleteWifi.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrDeleteWifi, id, stderr)
	}
	slog.Info(
		"deleted wifi connection",
		"id", id,
		"output", string(out),
	)
	return nil
}

func (NMCLI) GetWifiStatus(ctx context.Context) (bool, error) {
	args := []string{"radio", "wifi"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiStatus.Error(),
			"stderr", stderr,
			"error", err,
		)
		return false, fmt.Errorf("%w: %s", infra.ErrGetWifiStatus, stderr)
	}
	slog.Info("retrieved wifi status", "output", string(out))
	return strings.Trim(string(out), " \n") == "enabled", nil
}

func (NMCLI) GetWWANStatus(ctx context.Context) (bool, error) {
	args := []string{"radio", "wwan"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetWifiStatus.Error(),
			"stderr", stderr,
			"error", err,
		)
		return false, fmt.Errorf("%w: %s", infra.ErrGetWWANStatus, stderr)
	}
	slog.Info("retrieved wwan status", "output", string(out))
	return strings.Trim(string(out), " \n") == "enabled", nil
}

func (n NMCLI) GetRadioStatus(ctx context.Context) (infra.RadioStatus, error) {
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
		slog.Error(
			infra.ErrGetRadioStatus.Error(),
			"failed operations",
			bigErrStr,
		)
		return infra.RadioStatus{}, fmt.Errorf("%w: %s", infra.ErrGetWifiInfo, bigErrStr)
	}

	return infra.RadioStatus{
		EnabledWifi: wifi,
		EnabledWWAN: wwan,
	}, nil
}

func (NMCLI) EnableWifi(ctx context.Context) error {
	args := []string{"radio", "wifi", "on"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrEnableWifi.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrEnableWifi, stderr)
	}
	slog.Info("wifi radio enabled", "output", string(out))
	return nil
}

func (NMCLI) DisableWifi(ctx context.Context) error {
	args := []string{"radio", "wifi", "off"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrDisableWifi.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrDisableWifi, stderr)
	}
	slog.Info("wifi radio disabled", "output", string(out))
	return nil
}

func (NMCLI) EnableWWAN(ctx context.Context) error {
	args := []string{"radio", "wwan", "on"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrEnableWWAN.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrEnableWWAN, stderr)
	}
	slog.Info("wifi radio enabled", "output", string(out))
	return nil
}

func (NMCLI) DisableWWAN(ctx context.Context) error {
	args := []string{"radio", "wwan", "off"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrDisableWWAN.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrDisableWWAN, stderr)
	}
	slog.Info("wifi radio disabled", "output", string(out))
	return nil
}

func (NMCLI) EnableNetworking(ctx context.Context) error {
	args := []string{"networking", "on"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrEnableNetworking.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrEnableNetworking, stderr)
	}
	slog.Info("networking enabled", "output", string(out))
	return nil
}

func (NMCLI) DisableNetworking(ctx context.Context) error {
	args := []string{"networking", "off"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrDisableNetworking.Error(),
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w: %s", infra.ErrDisableNetworking, stderr)
	}
	slog.Info("networking disabled", "output", string(out))
	return nil
}

func (NMCLI) GetConnectivityStatus(ctx context.Context) (infra.ConnectivityStatus, error) {
	args := []string{"networking", "connectivity", "check"}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrGetConnectivityStatus.Error(),
			"stderr", stderr,
			"error", err,
		)
		return "", fmt.Errorf("%w: %s", infra.ErrGetConnectivityStatus, stderr)
	}
	slog.Info("retrieved wwan status", "output", string(out))
	return infra.ConnectivityStatus(strings.Trim(string(out), " \n")), nil
}

func (NMCLI) CreateHotspot(ctx context.Context, device string, id string, password string, hidden bool) error {
	hiddenStr := "no"
	if hidden {
		hiddenStr = "yes"
	}
	args := []string{
		"nmcli", "device", "wifi", "hotspot",
		"ifname", device,
		"ssid", id,
		"password", password,
		"hidden", hiddenStr,
	}
	out, err := exec.CommandContext(ctx, CommandName, args...).Output()
	if err != nil {
		stderr := infra.ExtractStderr(err)
		slog.Error(
			infra.ErrCreateHotspot.Error(),
			"id", id,
			"stderr", stderr,
			"error", err,
		)
		return fmt.Errorf("%w %s: %s", infra.ErrCreateHotspot, id, stderr)
	}
	slog.Info(
		"hotspot created",
		"device", device,
		"id", id,
		"output", string(out),
		"hidden", hiddenStr,
	)
	return nil
}
