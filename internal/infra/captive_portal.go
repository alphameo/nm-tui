package infra

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
)

const XDGOpen = "xdg-open"

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetOutboundIP, err)
	}
	defer func() { _ = conn.Close() }()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// OpenCaptivePortal opens captive portal web-page in browser
// Equivalent to xdg-open "http://$(ip --oneline route get 1.1.1.1 | awk '{print $3}')"
func OpenCaptivePortal(ctx context.Context) error {
	ip, err := getOutboundIP()
	if err != nil {
		slog.Error(err.Error(),
			"err", err)
		return err
	}

	url := fmt.Sprintf("http://%s", ip.String())

	cmd := exec.CommandContext(ctx, XDGOpen, url)
	if err := cmd.Start(); err != nil {
		stderr := ExtractStderr(err)
		slog.Error(ErrOpenCaptivePortal.Error(),
			"err", err,
			"stderr", stderr)
		return err
	}
	return nil
}
