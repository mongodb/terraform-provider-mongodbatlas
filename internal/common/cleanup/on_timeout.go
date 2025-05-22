package cleanup

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	CleanupWarning = "Failed to create resource. Will run cleanup due to the operation timing out"
)

// OnTimeout creates a new context with a timeout and a deferred function that will run `cleanup` when the context hit the timeout (no timeout=no-op).
// Remember to always call the returned `deferCall` function: `defer deferCall()`.
// `warningDetail` should have resource identifiable information, for example cluster name and project ID.
// warnDiags(summary, detail) are called:
// 1. Before the cleanup call.
// 2. (Only if the cleanup fails) Details of the cleanup error.
func OnTimeout(ctx context.Context, timeout time.Duration, warnDiags func(string, string), warningDetail string, cleanup func(context.Context) error) (outCtx context.Context, deferCall func()) {
	outCtx, cancel := context.WithTimeout(ctx, timeout)
	return outCtx, func() {
		cancel()
		if !errors.Is(outCtx.Err(), context.DeadlineExceeded) {
			return
		}
		warnDiags(CleanupWarning, warningDetail)
		newContext := context.Background() // Create a new context for cleanup as the old context is expired
		if err := cleanup(newContext); err != nil {
			warnDiags("Error during cleanup", warningDetail+" error="+err.Error())
		}
	}
}

const (
	contextDeadlineExceeded = "context deadline exceeded"
	TimeoutReachedPrefix    = "Timeout reached after "
)

func ReplaceContextDeadlineExceededDiags(diags *diag.Diagnostics, duration time.Duration) {
	for i := range len(*diags) {
		d := (*diags)[i]
		if d.Severity() != diag.SeverityError {
			continue
		}
		if strings.Contains(d.Detail(), contextDeadlineExceeded) {
			(*diags)[i] = diag.NewErrorDiagnostic(
				d.Summary(),
				strings.ReplaceAll(d.Detail(), contextDeadlineExceeded, TimeoutReachedPrefix+duration.String()),
			)
		}
	}
}
