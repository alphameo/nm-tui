// Package portal opens the captive portal login page in the default browser.
package portal

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
)

// Opener opens captive portal login pages using the default browser.
type Opener struct{}

func New() *Opener {
	return &Opener{}
}

// OpenCaptivePortal opens the captive portal web-page in the default browser.
func (p *Opener) OpenCaptivePortal(ctx context.Context) error {
	ip, err := p.getGatewayIP(ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s", ip.String())

	return p.openURL(url)
}

// getGatewayIP is equivalent to
// xdg-open "http://$(ip --oneline route get 1.1.1.1 | awk '{print $3}')"
func (p *Opener) getGatewayIP(ctx context.Context) (net.IP, error) {
	ipargs := []string{"--oneline", "route", "get", "1.1.1.1"}
	route, err := exec.CommandContext(ctx, "ip", ipargs...).Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", infra.ErrGetGatewayIP, err)
	}
	out := strings.Split(string(route), " ")
	if len(out) < 3 {
		return nil, fmt.Errorf("%w: unexpected format", infra.ErrGetGatewayIP)
	}
	return net.ParseIP(out[2]), nil
}

// openURL opens the URL in the default browser (cross-platform).
func (p *Opener) openURL(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("%w: %s", infra.ErrUnsupportedPlarform, runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
