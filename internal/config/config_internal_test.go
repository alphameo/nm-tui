package config

import "testing"

func TestValidateTime(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want bool
	}{
		{"positive", 1, true},
		{"large positive", 10000, true},
		{"zero", 0, false},
		{"negative", -5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTime(tt.in)
			if (err == nil) != tt.want {
				t.Errorf("validateTime(%d) error = %v, want ok=%v", tt.in, err, tt.want)
			}
		})
	}
}
