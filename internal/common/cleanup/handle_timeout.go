package cleanup

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

const (
	CleanupWarning                           = "Failed to create resource. Will run cleanup due to the operation timing out"
	DeleteOnCreateTimeoutInvalidErrorMessage = "delete_on_create_timeout cannot be updated or set after import, remove it from the configuration"
)

// HandleCreateTimeout helps to implement Create in long-running operations.
// It deletes the resource if the creation times out and `delete_on_create_timeout` is enabled.
// It returns an error with additional information which should be used instead of the original error.
func HandleCreateTimeout(deleteOnCreateTimeout bool, errWait error, cleanup func(context.Context) error) error {
	if _, isTimeoutErr := errWait.(*retry.TimeoutError); !isTimeoutErr {
		return errWait
	}
	if !deleteOnCreateTimeout {
		return errors.Join(errWait, errors.New("cleanup won't be run because delete_on_create_timeout is false"))
	}
	errWait = errors.Join(errWait, errors.New("will run cleanup because delete_on_create_timeout is true. If you suspect a transient error, wait before retrying to allow resource deletion to finish"))
	// cleanup uses a new context as existing one is expired.
	if errCleanup := cleanup(context.Background()); errCleanup != nil {
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

const (
	OperationCreate = "create"
	OperationUpdate = "update"
	OperationDelete = "delete"
)

// ResolveTimeout extracts the appropriate timeout duration from the model for the given operation
func ResolveTimeout(ctx context.Context, t *timeouts.Value, operationName string, diags *diag.Diagnostics) time.Duration {
	var (
		timeoutDuration time.Duration
		localDiags      diag.Diagnostics
	)
	switch operationName {
	case OperationCreate:
		timeoutDuration, localDiags = t.Create(ctx, constant.DefaultTimeout)
		diags.Append(localDiags...)
	case OperationUpdate:
		timeoutDuration, localDiags = t.Update(ctx, constant.DefaultTimeout)
		diags.Append(localDiags...)
	case OperationDelete:
		timeoutDuration, localDiags = t.Delete(ctx, constant.DefaultTimeout)
		diags.Append(localDiags...)
	default:
		timeoutDuration = constant.DefaultTimeout
	}
	return timeoutDuration
}

// ResolveDeleteOnCreateTimeout returns true if delete_on_create_timeout should be enabled.
// Default behavior is true when not explicitly set to false.
func ResolveDeleteOnCreateTimeout(deleteOnCreateTimeout types.Bool) bool {
	// If null or unknown, default to true
	if deleteOnCreateTimeout.IsNull() || deleteOnCreateTimeout.IsUnknown() {
		return true
	}
	// Otherwise use the explicit value
	return deleteOnCreateTimeout.ValueBool()
}

type resourceInterface interface {
	GetOkExists(key string) (any, bool)
	HasChange(key string) bool
}

// DeleteOnCreateTimeoutInvalidUpdate returns an error if the `delete_on_create_timeout` attribute has been updated to true/false
// This use case differs slightly from the behavior of TPF customplanmodifier.CreateOnlyBoolWithDefault:
// - from a given value (true/false) --> `null`.
// While the TPF implementation keeps the state value (UseStateForUnknown behavior),
// The SDKv2 implementation will set the state value to null (Optional-only attribute).
func DeleteOnCreateTimeoutInvalidUpdate(resource resourceInterface) string {
	if !resource.HasChange("delete_on_create_timeout") {
		return ""
	}
	if _, exists := resource.GetOkExists("delete_on_create_timeout"); exists {
		return DeleteOnCreateTimeoutInvalidErrorMessage
	}
	return ""
}
