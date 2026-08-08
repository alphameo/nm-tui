package infra

import "context"

// CaptivePortalOpener opens the captive portal login page
// in the default browser.
type CaptivePortalOpener interface {
	OpenCaptivePortal(ctx context.Context) error
}
