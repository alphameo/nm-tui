package models_test

import (
	"reflect"
	"testing"

	"github.com/alphameo/nm-tui/internal/ui/models"
)

func TestCrossReferenceNetworks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		available []models.AvailableNetwork
		profiles  []models.NetworkProfileShort
		wantAvail []models.AvailableNetwork
		wantProf  []models.NetworkProfileShort
	}{
		{
			name: "matches available network to saved profile by ssid",
			available: []models.AvailableNetwork{
				{SSID: "home"},
				{SSID: "cafe"},
			},
			profiles: []models.NetworkProfileShort{
				{Name: "home-net", SSID: "home"},
				{Name: "office", SSID: "office"},
			},
			wantAvail: []models.AvailableNetwork{
				{SSID: "home", ProfileExists: true},
				{SSID: "cafe"},
			},
			wantProf: []models.NetworkProfileShort{
				{Name: "home-net", SSID: "home", Available: true},
				{Name: "office", SSID: "office"},
			},
		},
		{
			name: "duplicate access points of same ssid all marked",
			available: []models.AvailableNetwork{
				{SSID: "office", Signal: 40},
				{SSID: "office", Signal: 20},
			},
			profiles: []models.NetworkProfileShort{
				{Name: "office", SSID: "office"},
			},
			wantAvail: []models.AvailableNetwork{
				{SSID: "office", Signal: 40, ProfileExists: true},
				{SSID: "office", Signal: 20, ProfileExists: true},
			},
			wantProf: []models.NetworkProfileShort{
				{Name: "office", SSID: "office", Available: true},
			},
		},
		{
			name: "empty ssid never matches",
			available: []models.AvailableNetwork{
				{SSID: "home"},
				{SSID: ""},
			},
			profiles: []models.NetworkProfileShort{
				{Name: "hidden", SSID: ""},
				{Name: "home-net", SSID: "home"},
			},
			wantAvail: []models.AvailableNetwork{
				{SSID: "home", ProfileExists: true},
				{SSID: ""},
			},
			wantProf: []models.NetworkProfileShort{
				{Name: "hidden", SSID: ""},
				{Name: "home-net", SSID: "home", Available: true},
			},
		},
		{
			name:      "no overlap leaves flags unset",
			available: []models.AvailableNetwork{{SSID: "a"}},
			profiles:  []models.NetworkProfileShort{{Name: "b", SSID: "b"}},
			wantAvail: []models.AvailableNetwork{{SSID: "a"}},
			wantProf:  []models.NetworkProfileShort{{Name: "b", SSID: "b"}},
		},
		{
			name:      "empty inputs",
			available: nil,
			profiles:  nil,
			wantAvail: []models.AvailableNetwork{},
			wantProf:  []models.NetworkProfileShort{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotAvail, gotProf := models.CrossReferenceNetworks(tt.available, tt.profiles)
			assertEq(t, "available", gotAvail, tt.wantAvail)
			assertEq(t, "profiles", gotProf, tt.wantProf)
		})
	}
}

func assertEq(t *testing.T, name string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s mismatch:\n got: %#v\nwant: %#v", name, got, want)
	}
}
