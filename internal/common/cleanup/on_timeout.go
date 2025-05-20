package cleanup

import (
	"context"
	"errors"
	"time"
)

type AddWarning interface {
	AddWarning(string, string)
}

func OnTimeout(ctx context.Context, timeout time.Duration, warnDiags AddWarning, warningDetail string, cleanup func(context.Context) error) (outCtx context.Context, deferCall func()) {
	outCtx, cancel := context.WithTimeout(ctx, timeout)
	return outCtx, func() {
		cancel()
		if !errors.Is(outCtx.Err(), context.DeadlineExceeded) {
			return
		}
		cleanupWarning := "Failed to create, will perform cleanup due to timeout reached"
		warnDiags.AddWarning(cleanupWarning, warningDetail)
		newContext := context.Background() // Create a new context for cleanup
		if err := cleanup(newContext); err != nil {
			warnDiags.AddWarning("Error during cleanup", warningDetail+" error="+err.Error())
		}
	}
}
