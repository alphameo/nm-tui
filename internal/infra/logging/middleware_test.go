package logging_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/infra/logging"
)

type record struct {
	level slog.Level
	msg   string
	attrs map[string]any
}

type captureHandler struct {
	records []record
	level   slog.Level
}

func (h *captureHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level
}

func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	rec := record{level: r.Level, msg: r.Message, attrs: map[string]any{}}
	r.Attrs(func(a slog.Attr) bool {
		rec.attrs[a.Key] = a.Value.Any()
		return true
	})
	h.records = append(h.records, rec)
	return nil
}
func (h *captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *captureHandler) WithGroup(string) slog.Handler      { return h }

func newCapture(level slog.Level) *captureHandler {
	return &captureHandler{level: level}
}

type wifiStub struct {
	infra.NetworksManager

	scanErr    error
	connectErr error
	pass       string
}

func (s *wifiStub) ListNetworksWithRescan(context.Context) ([]infra.AvailableNetwork, error) {
	return nil, s.scanErr
}

func (s *wifiStub) ConnectToNetwork(context.Context, string, string) error {
	return s.connectErr
}

func (s *wifiStub) GetProfilePassword(context.Context, string) (string, error) {
	return s.pass, s.scanErr
}

type networkStub struct {
	infra.DeviceManager
}

type portalStub struct {
	infra.CaptivePortalOpener
}

func exitErr(t *testing.T, code int) error {
	t.Helper()
	// #nosec G204 -- test helper; code is a fixed integer from the test
	err := exec.Command("sh", "-c",
		fmt.Sprintf("exit %d", code)).
		Run()
	if _, ok := errors.AsType[*exec.ExitError](err); !ok {
		t.Fatalf("want *exec.ExitError, got %T", err)
	}
	return err
}

func TestScanWifisFailureLogsErrorAndExitCode(t *testing.T) {
	t.Parallel()

	h := newCapture(slog.LevelDebug)
	m := logging.New(slog.New(h), &wifiStub{scanErr: exitErr(t, 3)}, &networkStub{}, &portalStub{})

	if _, err := m.ListNetworksWithRescan(context.Background()); err == nil {
		t.Fatal("want error, got nil")
	}
	if len(h.records) != 1 {
		t.Fatalf("want 1 record, got %d", len(h.records))
	}
	r := h.records[0]
	if r.level != slog.LevelError {
		t.Errorf("level = %v, want error", r.level)
	}
	if got := r.attrs["exit_code"]; got != int64(3) {
		t.Errorf("exit_code = %v (type %T), want 3", got, got)
	}
	if _, ok := r.attrs["error"]; !ok {
		t.Error("missing error attribute")
	}
}

func TestSecretsNeverLogged(t *testing.T) {
	t.Parallel()

	h := newCapture(slog.LevelDebug)
	pass := "super-secret-pass-1"
	m := logging.New(slog.New(h), &wifiStub{pass: pass}, &networkStub{}, &portalStub{})

	got, err := m.GetProfilePassword(context.Background(), "home")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != pass {
		t.Fatalf("password passthrough broken: got %q", got)
	}

	if err = m.ConnectToNetwork(context.Background(), "home", pass); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range h.records {
		if strings.Contains(r.msg, pass) {
			t.Errorf("record message leaks password: %q", r.msg)
		}
		for k, v := range r.attrs {
			if s, ok := v.(string); ok && strings.Contains(s, pass) {
				t.Errorf("record attribute %s leaks password: %q", k, s)
			}
		}
	}
}
