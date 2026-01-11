// Package infra provides interaction infrastructure layer
package infra

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
)

type Nmcli struct{}

func NewNMCLI() *Nmcli {
	return &Nmcli{}
}

const nmcliCmdName = "nmcli"

func (Nmcli) ScanWifi() ([]*WifiScanned, error) {
	// CMD: nmcli -t -f SSID,IN-USE,SECURITY,SIGNAL dev wifi
	args := []string{"-t", "-f", "SSID,IN-USE,SECURITY,SIGNAL", "dev", "wifi"}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err != nil {
		logger.Errf("Error scanning available wifi-networks (%s %s): %s\n", nmcliCmdName, args, err.Error())
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
	logger.Informf("Got list of available wifi-networks (%s %s)\n", nmcliCmdName, args)
	return res, nil
}

func (Nmcli) GetStoredWifi() ([]*WifiStored, error) {
	// CMD: nmcli -t -f NAME,STATE connection show
	args := []string{"-t", "-f", "NAME,STATE", "connection", "show"}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err != nil {
		logger.Errf("Error retreiving stored wifi-networks (%s %s): %s\n", nmcliCmdName, args, err.Error())
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
			SSID:   parts[0],
			Active: parts[1] == "activated",
		})
	}
	logger.Informf("Got list of stored wifi-networks (%s %s)\n", nmcliCmdName, args)
	return res, nil
}

func (n Nmcli) ConnectWifi(ssid, password string) error {
	// CMD: nmcli device wifi connect "<SSID>" password "<PASSWORD>"
	n.DeleteWifiConnection(ssid) // FIX: after nmcli 1.48.10 connection via password not able with saved networks
	args := []string{"device", "wifi", "connect", ssid, "password", password}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err == nil {
		logger.Informf("Connected to wifi %s (%s %s): %s", ssid, nmcliCmdName, args, string(out))
	} else {
		logger.Errf("Error connecting to wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
	}
	return err
}

func (Nmcli) ConnectSavedWifi(ssid string) error {
	// CMD: nmcli connection up "<SSID>"
	args := []string{"connection", "up", ssid}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err == nil {
		logger.Informf("Connected to saved wifi %s (%s %s): %s", ssid, nmcliCmdName, args, string(out))
	} else {
		logger.Errf("Error connecting to saved wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
	}
	return err
}

func (Nmcli) DisconnectFromWifi(ssid string) error {
	// CMD: nmcli connection down "<SSID>"
	args := []string{"connection", "down", ssid}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err == nil {
		logger.Informf("Disconnected from wifi %s (%s %s): %s", ssid, nmcliCmdName, args, string(out))
	} else {
		logger.Errf("Error disconnecting from wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
	}
	return err
}

func (Nmcli) GetConnectedWifi() ([]string, error) {
	// CMD: nmcli -t -f NAME connection show
	args := []string{"-t", "-f", "NAME", "connection", "show"}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err != nil {
		logger.Errf("Error retreiving list of connected wifi-networks (%s %s): %s\n", nmcliCmdName, args, err.Error())
		return nil, err
	}
	res := strings.Split(string(out), "\n")
	logger.Informf("Got list of connetcted wifi-networks (%s %s)\n", nmcliCmdName, args)
	return res, nil
}

func (Nmcli) GetWifiPassword(ssid string) (string, error) {
	// CMD: nmcli -s -g 802-11-wireless-security.psk connection show "<SSID>"
	args := []string{"-s", "-g", "802-11-wireless-security.psk", "connection", "show", ssid}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err != nil {
		logger.Errf("Error retrieving password to wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
		return "", err
	}
	pw := strings.Trim(string(out), " \n")
	logger.Informf("Got password to wifi %s (%s %s)\n", ssid, nmcliCmdName, args)
	return pw, nil
}

func (Nmcli) GetWifiInfo(ssid string) (*WifiInfo, error) {
	// CMD: nmcli -s -m tabular -t -f connection.id,802-11-wireless.ssid,802-11-wireless-security.psk,connection.autoconnect,connection.autoconnect-priority,GENERAL.STATE connection show
	args := []string{
		"-s",
		"-m",
		"tabular",
		"-t",
		"-f",
		"connection.id,802-11-wireless.ssid,802-11-wireless-security.psk,connection.autoconnect,connection.autoconnect-priority,GENERAL.STATE",
		"connection",
		"show",
		ssid,
	}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err != nil {
		logger.Errf("Error retrieving information about wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	autoconnectPriority, err := strconv.Atoi(lines[4])
	if err != nil {
		logger.Errf("Error retrieving information about wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
		return nil, err
	}

	active := len(lines) == 5

	logger.Informf("Got information about wifi %s (%s %s)\n", ssid, nmcliCmdName, args)
	return &WifiInfo{
		ID:                  lines[0],
		SSID:                lines[1],
		Password:            lines[2],
		Autoconnect:         lines[3] == "yes",
		AutoconnectPriority: autoconnectPriority,
		Active:              active,
	}, nil
}

func (Nmcli) DeleteWifiConnection(ssid string) error {
	// CMD: nmcli connection delete "<SSID>"
	args := []string{"connection", "delete", ssid}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err == nil {
		logger.Informf("Connection to wifi %s was deleted (%s %s): %s", ssid, nmcliCmdName, args, string(out))
	} else {
		logger.Errf("Error deleting connection to wifi %s (%s %s): %s\n", ssid, nmcliCmdName, args, err.Error())
	}
	return err
}

func (Nmcli) ConnectVPN(vpnName string) error {
	// CMD: nmcli connection up id "<VPN_NAME>"
	args := []string{"connection", "up", "id", vpnName}
	out, err := exec.Command(nmcliCmdName, args...).Output()
	if err == nil {
		logger.Informf("Connected to VPN %s (%s %s): %s", vpnName, nmcliCmdName, args, string(out))
	} else {
		logger.Errf("Error connecting to VPN %s (%s %s): %s\n", vpnName, nmcliCmdName, args, err.Error())
	}
	return err
}
