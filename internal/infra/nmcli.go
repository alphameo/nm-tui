// Package infra provides interaction infrastructure layer
package infra

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

type Nmcli struct{}

func NewNMCLI() *Nmcli {
	return &Nmcli{}
}

const NmcliCommandName = "nmcli"

func (Nmcli) GetAvailableWifi() ([]*WifiScanned, error) {
	args := []string{"-t", "-f", "SSID,IN-USE,SECURITY,SIGNAL", "dev", "wifi"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrScanningAvailableWifi.Error(),
			"err",
			err,
			"stderr",
			stderr,
		)
		return nil, fmt.Errorf("%w: %s", ErrScanningAvailableWifi, stderr)
	}

	var res []*WifiScanned
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

		signal, err := strconv.Atoi(parts[3])
		if err != nil {
			slog.Warn("parsing signal strength", "line", line, "error", err)
			signal = 0
		}
		res = append(res, &WifiScanned{
			SSID:     parts[0],
			Active:   parts[1] == "*",
			Security: parts[2],
			Signal:   signal,
		})
	}
	slog.Info("scanned available wifi networks")
	return res, nil
}

func (n Nmcli) GetStoredWifi() ([]*WifiStored, error) {
	args := []string{"-t", "-f", "NAME,STATE", "connection", "show"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		return nil, fmt.Errorf("%w: %s", ErrScanningStoredWifi, stderr)
	}

	var res []*WifiStored

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

		res = append(res, &WifiStored{
			Name:   parts[0],
			Active: parts[1] == "activated",
		})
	}

	for _, wifi := range res {
		ssid, err := n.GetWifiSSID(wifi.Name)
		if err != nil {
			slog.Warn(
				"failed to get ssid for stored wifi",
				"name",
				wifi.Name,
				"error",
				err,
			)
			continue
		}
		wifi.SSID = ssid
	}

	slog.Info("retrieved stored wifi networks")
	return res, nil
}

func (n Nmcli) ConnectWifi(ssid, password string) error {
	err := n.DeleteWifiConnection(ssid)
	if err != nil {
		return err
	}
	args := []string{"device", "wifi", "connect", ssid, "password", password}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrConnectingWifi.Error(),
			"ssid",
			ssid,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf("%w %s: %s", ErrConnectingWifi, ssid, stderr)
	}
	slog.Info("connected to wifi", "ssid", ssid, "output", string(out))
	return nil
}

func (Nmcli) ConnectStoredWifi(id string) error {
	args := []string{"connection", "up", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrConnectingStoredWifi.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf("%w %s: %s", ErrConnectingStoredWifi, id, stderr)
	}
	slog.Info("connected to saved wifi", "id", id, "output", string(out))
	return nil
}

func (Nmcli) DisconnectFromWifi(id string) error {
	args := []string{"connection", "down", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrDisconnectWifi.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf("%w %s: %s", ErrDisconnectWifi, id, stderr)
	}
	slog.Info("disconnected from wifi", "id", id, "output", string(out))
	return nil
}

func (Nmcli) GetConnectedWifi() ([]string, error) {
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		return nil, fmt.Errorf("%w: %s", ErrScanningConnectedWifi, stderr)
	}
	slog.Info("retrieved connected wifi networks")
	return strings.Split(string(out), "\n"), nil
}

func (Nmcli) GetWifiPassword(id string) (string, error) {
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"802-11-wireless-security.psk",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrGettingWifiPassword.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return "", fmt.Errorf("%w for %s: %s", ErrGettingWifiPassword, id, stderr)
	}
	slog.Info("retrieved wifi password", "id", id)
	return strings.Trim(string(out), " \n"), nil
}

func (Nmcli) GetWifiSSID(id string) (string, error) {
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"802-11-wireless.ssid",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrGettingWifiSSID.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return "", fmt.Errorf("%w %s: %s", ErrGettingWifiSSID, id, stderr)
	}
	slog.Info("retrieved wifi ssid", "id", id)
	return strings.Trim(string(out), " \n"), nil
}

func (Nmcli) GetWifiAutoconnect(id string) (bool, error) {
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"connection.autoconnect",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrGettingWifiAutoconnect.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return false, fmt.Errorf("%w for %s: %s", ErrGettingWifiAutoconnect, id, stderr)
	}
	return strings.Trim(string(out), " \n") == "yes", nil
}

func (Nmcli) GetWifiAutoconnectPriority(id string) (int, error) {
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"connection.autoconnect-priority",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrGettingWifiAutoconnectPriority.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return 0, fmt.Errorf(
			"%w for %s: %s",
			ErrGettingWifiAutoconnectPriority,
			id,
			stderr,
		)
	}
	autoconnectResp := strings.Trim(string(out), " \n")
	autoconnectPriority, err := strconv.Atoi(autoconnectResp)
	if err != nil {
		slog.Error(
			"parsing autoconnect priority",
			"id",
			id,
			"value",
			autoconnectResp,
			"error",
			err,
		)
		return 0, fmt.Errorf(
			"%w %s: %w",
			ErrGettingWifiAutoconnectPriority,
			id,
			err,
		)
	}
	return autoconnectPriority, nil
}

func (Nmcli) GetWifiActivity(id string) (bool, error) {
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"GENERAL.STATE",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrGettingWifiActivity.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return false, fmt.Errorf(
			"%w for %s: %s",
			ErrGettingWifiActivity,
			id,
			stderr,
		)
	}
	return strings.Trim(string(out), " \n") == "activated", nil
}

func (n *Nmcli) GetWifiInfo(id string) (*WifiInfo, error) {
	var errs []error
	ssid, err := n.GetWifiSSID(id)
	if err != nil {
		errs = append(errs, err)
	}

	password, err := n.GetWifiPassword(id)
	if err != nil {
		errs = append(errs, err)
	}

	autoconnect, err := n.GetWifiAutoconnect(id)
	if err != nil {
		errs = append(errs, err)
	}

	autoconnectPriority, err := n.GetWifiAutoconnectPriority(id)
	if err != nil {
		errs = append(errs, err)
	}

	activated, err := n.GetWifiActivity(id)
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
			ErrGettingWifiInfo.Error(),
			"id",
			id,
			"failed operations",
			bigErrStr,
		)
		return nil, fmt.Errorf(
			"%w for %s: %s",
			ErrGettingWifiInfo,
			id,
			bigErrStr,
		)
	}

	return &WifiInfo{
		Name:                id,
		SSID:                ssid,
		Password:            password,
		Autoconnect:         autoconnect,
		AutoconnectPriority: autoconnectPriority,
		Active:              activated,
	}, nil
}

// UpdateWifiInfo is not atomic
func (n Nmcli) UpdateWifiInfo(id string, info *UpdateWifiInfo) error {
	var errs []error

	err := n.updateWifiID(
		id,
		info.Name,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiPassword(
		info.Name,
		info.Password,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiAutoconnect(
		info.Name,
		info.Autoconnect,
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = n.updateWifiAutoconnectPriority(
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
			ErrGettingWifiInfo.Error(),
			"id",
			id,
			"failed operations",
			bigErrStr,
		)
		return fmt.Errorf(
			"%w for %s: %s",
			ErrUpdatingWifiInfo,
			id,
			bigErrStr,
		)
	}
	return nil
}

func (Nmcli) updateWifiInfoField(id, field, value string) error {
	args := []string{"connection", "modify", id, field, value}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrUpdatingWifiInfoField.Error(),
			"id",
			id,
			"field",
			field,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf(
			"%w %s for %s: %s",
			ErrUpdatingWifiInfoField,
			field,
			id,
			stderr,
		)
	}
	slog.Info(
		"modified wifi field",
		"id",
		id,
		"field",
		field,
		"output",
		string(out),
	)
	return nil
}

func (n Nmcli) updateWifiID(id, newID string) error {
	return n.updateWifiInfoField(id, "connection.id", newID)
}

func (n Nmcli) updateWifiPassword(id, password string) error {
	return n.updateWifiInfoField(id, "802-11-wireless-security.psk", password)
}

func (n Nmcli) updateWifiAutoconnect(id string, autoconnect bool) error {
	var autoconnectValue string
	if autoconnect {
		autoconnectValue = "yes"
	} else {
		autoconnectValue = "no"
	}

	return n.updateWifiInfoField(id, "connection.autoconnect", autoconnectValue)
}

func (n Nmcli) updateWifiAutoconnectPriority(id string, priority int) error {
	return n.updateWifiInfoField(
		id,
		"connection.autoconnect-priority",
		strconv.Itoa(priority),
	)
}

func (Nmcli) DeleteWifiConnection(id string) error {
	args := []string{"connection", "delete", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrDeletingWifi.Error(),
			"id",
			id,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf("%w %s: %s", ErrDeletingWifi, id, stderr)
	}
	slog.Info(
		"deleted wifi connection",
		"id",
		id,
		"output",
		string(out),
	)
	return nil
}

func (Nmcli) ConnectVPN(vpnName string) error {
	args := []string{"connection", "up", "id", vpnName}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(
			ErrConnectingVPN.Error(),
			"vpn",
			vpnName,
			"stderr",
			stderr,
			"error",
			err,
		)
		return fmt.Errorf("%w %s: %s", ErrConnectingVPN, vpnName, stderr)
	}
	slog.Info("connected to VPN", "vpn", vpnName, "output", string(out))
	return nil
}

func ExtractStderr(err error) string {
	var stderr string
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return stderr
}
