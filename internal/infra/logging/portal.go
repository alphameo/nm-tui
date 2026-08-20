package logging

import (
	"context"
	"log/slog"

	"github.com/alphameo/nm-tui/internal/infra"
)

// PortalMiddleware implements infra.CaptivePortalOpener by delegating to the
// wrapped implementation. Successes are logged at Debug level, failures at
// Error level along with the exit code of the failed command when the error
// is an [*exec.ExitError].
type PortalMiddleware struct {
	middleware

	portal infra.CaptivePortalOpener
}

// NewPortal returns a *PortalMiddleware wrapping the given opener.
func NewPortal(logger *slog.Logger, portal infra.CaptivePortalOpener) *PortalMiddleware {
	return &PortalMiddleware{
		middleware: middleware{logger: logger, prefix: "portal"},
		portal:     portal,
	}
}

func (m *PortalMiddleware) OpenCaptivePortal(ctx context.Context) error {
	return m.call("open_captive_portal", func() error {
		return m.portal.OpenCaptivePortal(ctx)
	})
}
