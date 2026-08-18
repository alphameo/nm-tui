package config

import "testing"

func TestValidatePositiveTime(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			err := validatePositiveTime(tt.in)
			if (err == nil) != tt.want {
				t.Errorf("validatePositiveTime(%d) error = %v, want ok=%v", tt.in, err, tt.want)
			}
		})
	}
}

func TestValidateNonNegativeTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		want bool
	}{
		{"positive", 1, true},
		{"large positive", 10000, true},
		{"zero", 0, true},
		{"negative", -5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateNonNegativeTime(tt.in)
			if (err == nil) != tt.want {
				t.Errorf("validateNonNegativeTime(%d) error = %v, want ok=%v", tt.in, err, tt.want)
			}
		})
	}
}
