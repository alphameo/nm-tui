package infra

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// OpenCaptivePortal opens captive portal web-page in browser
func OpenCaptivePortal(ctx context.Context) error {
	ip, err := getGatewayIP(ctx)
	if err != nil {
		slog.Error(err.Error(),
			"err", err)
		return err
	}

	url := fmt.Sprintf("http://%s", ip.String())

	err = openURL(url)
	if err != nil {
		stderr := ExtractStderr(err)
		slog.Error(ErrOpenCaptivePortal.Error(),
			"err", err,
			"stderr", stderr)
		return err
	}
	return nil
}

// Equivalent to xdg-open "http://$(ip --oneline route get 1.1.1.1 | awk '{print $3}')"
func getGatewayIP(ctx context.Context) (net.IP, error) {
	ipargs := []string{"--oneline", "route", "get", "1.1.1.1"}
	route, err := exec.CommandContext(ctx, "ip", ipargs...).Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetGatewayIP, err)
	}
	out := strings.Split(string(route), " ")
	if len(out) < 3 {
		return nil, fmt.Errorf("%w: unexpected format", ErrGetGatewayIP)
	}
	return net.ParseIP(out[2]), nil
}

// openURL opens the URL in the default browser (cross-platform)
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedPlarform, runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
