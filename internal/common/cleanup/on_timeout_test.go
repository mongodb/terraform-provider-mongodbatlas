package cleanup_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanupOnErrorSkippedWhenNoTimeout(t *testing.T) {
	cleanupCalled := false
	cleanupFunc := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}
	diags := diag.Diagnostics{}
	_, call := cleanup.OnTimeout(t.Context(), 1*time.Millisecond, diags.AddWarning, "warning detail", cleanupFunc)
	call()
	assert.False(t, cleanupCalled, "cleanup should not be called when there are no timeouts")
}

func TestCleanupOnErrorCalledForATimeout(t *testing.T) {
	time.Sleep(2 * time.Millisecond)
	cleanupCalled := false
	finalContext := t.Context()
	cleanupFunc := func(callbackCtx context.Context) error {
		cleanupCalled = true
		finalContext = callbackCtx
		return nil
	}
	diags := diag.Diagnostics{}
	diags.AddError("error", "timeout")
	ctx, call := cleanup.OnTimeout(t.Context(), 1*time.Millisecond, diags.AddWarning, "warning detail", cleanupFunc)
	time.Sleep(2 * time.Millisecond)
	call()
	assert.True(t, cleanupCalled, "cleanup should be called")
	assert.NotEqual(t, finalContext, ctx, "cleanup should be called with a new context")
	require.NoError(t, finalContext.Err(), "cleanup should be called with a new context that hasn't been cancelled")
	assert.Len(t, diags, 2)
	assert.Equal(t, "Failed to create, will perform cleanup due to timeout reached", diags[1].Summary())
	assert.Equal(t, "warning detail", diags[1].Detail())
}
