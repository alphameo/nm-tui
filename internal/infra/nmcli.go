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
	// CMD: nmcli -t -f SSID,IN-USE,SECURITY,SIGNAL dev wifi
	// TODO: nmcli -t -f SSID,IN-USE,SECURITY,SIGNAL,FREQ,RATE,BANDWIDTH dev wifi
	args := []string{"-t", "-f", "SSID,IN-USE,SECURITY,SIGNAL", "dev", "wifi"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error scanning available wifi-networks (%s %s): %s", NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}

	var res []*WifiScanned
	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		signal, _ := strconv.Atoi(parts[3])
		res = append(res, &WifiScanned{
			SSID:     parts[0],
			Active:   parts[1] == "*",
			Security: parts[2],
			Signal:   signal,
		})
	}
	msg := fmt.Sprintf("Got list of available wifi-networks (%s %s)", NmcliCommandName, args)
	slog.Info(msg)
	return res, nil
}

func (n Nmcli) GetStoredWifi() ([]*WifiStored, error) {
	// CMD: nmcli -t -f NAME,STATE connection show
	args := []string{"-t", "-f", "NAME,STATE", "connection", "show"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error retreiving stored wifi-networks (%s %s): %s", NmcliCommandName, args, err.Error())
		slog.Error(msg)
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
			wifi.SSID = ""
		}
		wifi.SSID = ssid
	}

	msg := fmt.Sprintf("Got list of stored wifi-networks (%s %s)", NmcliCommandName, args)
	slog.Info(msg)
	return res, nil
}

func (n Nmcli) ConnectWifi(ssid, password string) error {
	// CMD: nmcli device wifi connect "<SSID>" password "<PASSWORD>"
	err := n.DeleteWifiConnection(ssid) // FIX: after nmcli 1.48.10 connection via password not able with saved networks
	if err != nil {
		return err
	}
	args := []string{"device", "wifi", "connect", ssid, "password", password}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Connected to wifi %s (%s %s): %s", ssid, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error connecting to wifi %s (%s %s): %s", ssid, NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}
	return err
}

func (Nmcli) ConnectStoredWifi(id string) error {
	// CMD: nmcli connection up "<ID>"
	args := []string{"connection", "up", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Connected to saved wifi %s (%s %s): %s", id, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error connecting to saved wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Info(msg)
	}
	return err
}

func (Nmcli) DisconnectFromWifi(id string) error {
	// CMD: nmcli connection down "<ID>"
	args := []string{"connection", "down", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Disconnected from wifi %s (%s %s): %s", id, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error disconnecting from wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}
	return err
}

func (Nmcli) GetConnectedWifi() ([]string, error) {
	// CMD: nmcli -t -f NAME connection show
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error retreiving list of connected wifi-networks (%s %s): %s", NmcliCommandName, args, err.Error())
		slog.Error(msg)
		return nil, err
	}
	res := strings.Split(string(out), "\n")
	msg := fmt.Sprintf("Got list of connetcted wifi-networks (%s %s)", NmcliCommandName, args)
	slog.Info(msg)
	return res, nil
}

func (Nmcli) GetWifiPassword(id string) (string, error) {
	// CMD: nmcli -s -m tabular -t -f 802-11-wireless-security.psk connection show "<ID>"
	args := []string{"-s", "-m", "tabular", "-t", "-f", "802-11-wireless-security.psk", "connection", "show", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error retrieving password to wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Info(msg)
		return "", err
	}
	pw := strings.Trim(string(out), " \n")
	msg := fmt.Sprintf("Got password to wifi %s (%s %s)", id, NmcliCommandName, args)
	slog.Info(msg)
	return pw, nil
}

func (Nmcli) GetWifiSSID(id string) (string, error) {
	// CMD: nmcli -s -m tabular -t -f 802-11-wireless.ssid connection show "<ID>"
	args := []string{"-s", "-m", "tabular", "-t", "-f", "802-11-wireless.ssid", "connection", "show", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error retrieving ssid for wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
		return "", err
	}
	ssid := strings.Trim(string(out), " \n")
	msg := fmt.Sprintf("Got password to wifi %s (%s %s)", id, NmcliCommandName, args)
	slog.Info(msg)
	return ssid, nil
}

func (Nmcli) GetWifiInfo(id string) (*WifiInfo, error) {
	// CMD: nmcli -s -m tabular -t -f connection.id,802-11-wireless.ssid,802-11-wireless-security.psk,connection.autoconnect,connection.autoconnect-priority,GENERAL.STATE connection show "<ID>"
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"connection.id,802-11-wireless.ssid,802-11-wireless-security.psk,connection.autoconnect,connection.autoconnect-priority,GENERAL.STATE",
		"connection",
		"show",
		id,
	}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("Error retrieving information about wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	autoconnectPriority, err := strconv.Atoi(lines[4])
	if err != nil {
		msg := fmt.Sprintf("Error retrieving information about wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
		return nil, err
	}

	var active bool
	if len(lines) > 5 {
		active = lines[5] == "activated"
	}

	msg := fmt.Sprintf("Got information about wifi %s (%s %s)", id, NmcliCommandName, args)
	slog.Info(msg)
	return &WifiInfo{
		Name:                lines[0],
		SSID:                lines[1],
		Password:            lines[2],
		Autoconnect:         lines[3] == "yes",
		AutoconnectPriority: autoconnectPriority,
		Active:              active,
	}, nil
}

// UpdateWifiInfo is not atomic
func (n Nmcli) UpdateWifiInfo(id string, info *UpdateWifiInfo) error {
	err := n.updateWifiID(id, info.Name)
	if err != nil {
		return err
	}

	err = n.updateWifiPassword(info.Name, info.Password)
	if err != nil {
		return err
	}

	err = n.updateWifiAutoconnect(info.Name, info.Autoconnect)
	if err != nil {
		return err
	}

	err = n.updateWifiAutoconnectPriority(info.Name, info.AutoconnectPriority)
	if err != nil {
		return err
	}

	return err
}

func (Nmcli) updateWifiInfoField(id, field, value string) error {
	// CMD: nmcli connection modify "<ID>" "<field>" "<value>"
	args := []string{"connection", "modify", id, field, value}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Wifi %s was modified (%s %s): %s", id, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error modifying wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}

	return err
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
	return n.updateWifiInfoField(id, "connection.autoconnect-priority", strconv.Itoa(priority))
}

func (Nmcli) DeleteWifiConnection(id string) error {
	// CMD: nmcli connection delete "<ID>"
	args := []string{"connection", "delete", id}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Connection to wifi %s was deleted (%s %s): %s", id, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error deleting connection to wifi %s (%s %s): %s", id, NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}
	return err
}

func (Nmcli) ConnectVPN(vpnName string) error {
	// CMD: nmcli connection up id "<VPN_NAME>"
	args := []string{"connection", "up", "id", vpnName}
	out, err := exec.Command(NmcliCommandName, args...).Output()
	if err == nil {
		msg := fmt.Sprintf("Connected to VPN %s (%s %s): %s", vpnName, NmcliCommandName, args, string(out))
		slog.Info(msg)
	} else {
		msg := fmt.Sprintf("Error connecting to VPN %s (%s %s): %s", vpnName, NmcliCommandName, args, err.Error())
		slog.Error(msg)
	}
	return err
}
