package models

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

func ConvertNetworkMode(mode infra.NetworkMode) string {
	switch mode {
	case infra.NetworkAccessPoint:
		return styles.SymbolAccessPoint
	case infra.NetworkInfra:
		return styles.SymbolInfra
	case infra.NetworkMesh:
		return styles.SymbolMesh
	case infra.NetworkAdHoc:
		return styles.SymbolAdHoc
	default:
		return "?"
	}
}

func convertAvailableNetwork(record infra.AvailableNetwork) AvailableNetwork {
	return AvailableNetwork{
		SSID:          record.SSID,
		Active:        record.Active,
		SecurityMode:  record.SecurityMode,
		Signal:        record.Signal,
		Band:          record.Band,
		Rate:          record.Rate,
		LookingDevice: record.LookingDevice,
		NetworkMode:   record.NetworkMode,
	}
}

func convertAvailableNetworks(records []infra.AvailableNetwork) []AvailableNetwork {
	out := make([]AvailableNetwork, len(records))
	for i, record := range records {
		out[i] = convertAvailableNetwork(record)
	}
	return out
}

func convertNetworkProfileShort(record infra.NetworkProfileShort) NetworkProfileShort {
	return NetworkProfileShort{
		Name:   record.Name,
		SSID:   record.SSID,
		Active: record.Active,
		Mode:   ConvertNetworkMode(record.Mode),
	}
}

func convertNetworkProfileShorts(records []infra.NetworkProfileShort) []NetworkProfileShort {
	out := make([]NetworkProfileShort, len(records))
	for i, record := range records {
		out[i] = convertNetworkProfileShort(record)
	}
	return out
}
