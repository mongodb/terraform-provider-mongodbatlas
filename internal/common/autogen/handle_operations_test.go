package autogen_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitFailedStateError_withErrorDescription(t *testing.T) {
	wait := &autogen.WaitReq{
		FailedStates:             []string{"FAILED", "CANCELED"},
		ErrorDescriptionProperty: "errorMessage",
	}
	objJSON := map[string]any{
		"state":        "FAILED",
		"errorMessage": "The restore could not complete because 1 collection was not found",
	}

	err := autogen.WaitFailedStateError(wait, "FAILED", objJSON)

	require.Error(t, err)
	assert.Equal(t, "The restore could not complete because 1 collection was not found", err.Error())
}

func TestWaitFailedStateError_withoutErrorDescription(t *testing.T) {
	wait := &autogen.WaitReq{
		FailedStates:             []string{"FAILED"},
		ErrorDescriptionProperty: "errorMessage",
	}

	t.Run("property missing", func(t *testing.T) {
		err := autogen.WaitFailedStateError(wait, "FAILED", map[string]any{"state": "FAILED"})
		require.Error(t, err)
		assert.Equal(t, `operation failed with state "FAILED"`, err.Error())
	})

	t.Run("property empty", func(t *testing.T) {
		err := autogen.WaitFailedStateError(wait, "FAILED", map[string]any{"errorMessage": ""})
		require.Error(t, err)
		assert.Equal(t, `operation failed with state "FAILED"`, err.Error())
	})

	t.Run("property not a string", func(t *testing.T) {
		err := autogen.WaitFailedStateError(wait, "FAILED", map[string]any{"errorMessage": 12})
		require.Error(t, err)
		assert.Equal(t, `operation failed with state "FAILED"`, err.Error())
	})
}

func TestWaitFailedStateError_pendingReturnsNil(t *testing.T) {
	wait := &autogen.WaitReq{
		FailedStates:             []string{"FAILED", "CANCELED"},
		ErrorDescriptionProperty: "errorMessage",
	}
	objJSON := map[string]any{
		"state":        "IN_PROGRESS",
		"errorMessage": "should be ignored",
	}

	assert.NoError(t, autogen.WaitFailedStateError(wait, "IN_PROGRESS", objJSON))
	assert.NoError(t, autogen.WaitFailedStateError(wait, "SUCCESSFUL", objJSON))
	assert.NoError(t, autogen.WaitFailedStateError(&autogen.WaitReq{}, "FAILED", objJSON))
}
