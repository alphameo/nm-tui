package infra

import (
	"context"
	"errors"
)

var (
	ErrOpenCaptivePortal   error = errors.New("failed to open captive portal")
	ErrGetGatewayIP        error = errors.New("failed to get gateway ip")
	ErrUnsupportedPlarform error = errors.New("unsupported platform")
)

// CaptivePortalOpener opens the captive portal login page in the default browser.
type CaptivePortalOpener interface {
	OpenCaptivePortal(ctx context.Context) error
}
