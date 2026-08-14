package config

import "testing"

func TestValidLogLevel(t *testing.T) {
	t.Parallel()

	for _, lvl := range []string{LogDebug, LogInfo, LogWarn, LogError} {
		if !validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = false, want true", lvl)
		}
	}
	for _, lvl := range []string{"DEBUG", "Info", "WARN", "Error"} {
		if !validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = false, want true (case-insensitive)", lvl)
		}
	}
	for _, lvl := range []string{"", "verbose", "trace", "fatal", "info ", "debug\n"} {
		if validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = true, want false", lvl)
		}
	}
}
