package cleanup

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	CleanupWarning = "Failed to create resource. Will run cleanup due to the operation timing out"
)

// HandleCreateTimeout helps to implement Create in long-running operations.
// It will delete the resource if the creation times out and delete_on_create_timeout is enabled.
// It returns an error with additional information which should be used instead of the original error.
func HandleCreateTimeout(deleteOnCreateTimeout bool, errWait error, cleanup func() error) error {
	if _, isTimeoutErr := errWait.(*retry.TimeoutError); !isTimeoutErr {
		return errWait
	}
	if !deleteOnCreateTimeout {
		return errors.Join(errWait, errors.New("cleanup won't be run because delete_on_create_timeout is false"))
	}
	errWait = errors.Join(errWait, errors.New("will run cleanup because delete_on_create_timeout is true, if you think this error is transient, please try creating the resource again in a few minutes"))
	if errCleanup := cleanup(); errCleanup != nil {
		errWait = errors.Join(errWait, errors.New("cleanup failed: "+errCleanup.Error()))
	}
	return errWait
}

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
		if strings.Contains(d.Detail(), contextDeadlineExceeded) {
			(*diags)[i] = diag.NewErrorDiagnostic(
				d.Summary(),
				strings.ReplaceAll(d.Detail(), contextDeadlineExceeded, TimeoutReachedPrefix+duration.String()),
			)
		}
	}
}
