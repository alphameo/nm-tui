package infra

import (
	"context"
	"errors"
)

var (
	ErrOpenCaptivePortal   = errors.New("failed to open captive portal")
	ErrGetGatewayIP        = errors.New("failed to get gateway ip")
	ErrUnsupportedPlarform = errors.New("unsupported platform")
)

// CaptivePortalOpener opens the captive portal login page in the default browser.
type CaptivePortalOpener interface {
	OpenCaptivePortal(ctx context.Context) error
}
