package cleanup_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timeoutDuration = 10 * time.Millisecond
)

func TestCleanupOnErrorSkippedWhenNoTimeout(t *testing.T) {
	cleanupCalled := false
	cleanupFunc := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}
	diags := diag.Diagnostics{}
	_, call := cleanup.OnTimeout(t.Context(), timeoutDuration, diags.AddWarning, "warning detail", cleanupFunc)
	call()
	assert.False(t, cleanupCalled, "cleanup should not be called when there are no timeouts")
}

func TestCleanupOnErrorCalledForATimeout(t *testing.T) {
	cleanupCalled := false
	finalContext := t.Context()
	cleanupFunc := func(callbackCtx context.Context) error {
		cleanupCalled = true
		finalContext = callbackCtx
		return errors.New("cleanup error")
	}
	diags := diag.Diagnostics{}
	diags.AddError("error", "timeout") // diags entry 1
	ctx, call := cleanup.OnTimeout(t.Context(), timeoutDuration, diags.AddWarning, "warning detail", cleanupFunc)
	time.Sleep(2 * timeoutDuration) // Sleep to ensure the timeout is reached
	call()
	assert.True(t, cleanupCalled, "cleanup should be called")
	assert.NotEqual(t, finalContext, ctx, "cleanup should be called with a new context")
	require.NoError(t, finalContext.Err(), "cleanup should be called with a new context that hasn't been cancelled")
	assert.Len(t, diags, 3) // diags entry 2 & 3 are added in the cleanup
	assert.Equal(t, cleanup.CleanupWarning, diags[1].Summary())
	assert.Equal(t, "warning detail", diags[1].Detail())

	assert.Equal(t, "Error during cleanup", diags[2].Summary())
	assert.Equal(t, "warning detail error=cleanup error", diags[2].Detail())
}

func TestReplaceContextDeadlineExceededDiags(t *testing.T) {
	diags := diag.Diagnostics{
		diag.NewErrorDiagnostic(
			"Error creating resource",
			"Error waiting for state to be IDLE: context deadline exceeded",
		),
		diag.NewErrorDiagnostic(
			"Another error",
			"This is a different error",
		),
		diag.NewWarningDiagnostic(
			"Warning with deadline",
			"Warning with context deadline exceeded mentioned",
		),
	}

	expectedSummaries := []string{
		"Error creating resource",
		"Another error",
		"Warning with deadline",
	}
	expectedDetails := []string{
		"Error waiting for state to be IDLE: Timeout reached after 2m0s",
		"This is a different error",
		"Warning with Timeout reached after 2m0s mentioned",
	}

	duration := 2 * time.Minute
	cleanup.ReplaceContextDeadlineExceededDiags(&diags, duration)

	assert.Len(t, diags, 3, "Expected same number of diagnostics")
	for i, diag := range diags {
		assert.Equal(t, expectedSummaries[i], diag.Summary(), "Summary at index %d should match", i)
		assert.Equal(t, expectedDetails[i], diag.Detail(), "Detail at index %d should match", i)
	}
}

type mockResource struct {
	values    map[string]interface{}
	isChanged bool
}

func (m *mockResource) GetOkExists(key string) (interface{}, bool) {
	value, exists := m.values[key]
	return value, exists
}

func (m *mockResource) HasChange(key string) bool {
	return m.isChanged
}

func TestDeleteOnCreateTimeoutInvalidUpdate(t *testing.T) {
	t.Run("No change in delete_on_create_timeout returns empty string", func(t *testing.T) {
		resource := &mockResource{
			values:    map[string]interface{}{},
			isChanged: false,
		}

		result := cleanup.DeleteOnCreateTimeoutInvalidUpdate(resource)
		assert.Empty(t, result)
	})

	t.Run("Change detected but field doesn't exist returns empty string", func(t *testing.T) {
		resource := &mockResource{
			values:    map[string]interface{}{},
			isChanged: true,
		}

		result := cleanup.DeleteOnCreateTimeoutInvalidUpdate(resource)
		assert.Empty(t, result)
	})

	t.Run("Change detected and field exists returns error message", func(t *testing.T) {
		resource := &mockResource{
			values:    map[string]interface{}{"delete_on_create_timeout": true},
			isChanged: true,
		}

		result := cleanup.DeleteOnCreateTimeoutInvalidUpdate(resource)
		expectedMessage := "delete_on_create_timeout cannot be updated or set after import, remove it from the configuration"
		assert.Equal(t, expectedMessage, result)
	})

	t.Run("Change detected with false value still returns error message", func(t *testing.T) {
		resource := &mockResource{
			values:    map[string]interface{}{"delete_on_create_timeout": false},
			isChanged: true,
		}

		result := cleanup.DeleteOnCreateTimeoutInvalidUpdate(resource)
		expectedMessage := "delete_on_create_timeout cannot be updated or set after import, remove it from the configuration"
		assert.Equal(t, expectedMessage, result)
	})
}
