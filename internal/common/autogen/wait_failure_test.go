package autogen_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWaitModel struct {
	ProjectID   types.String `tfsdk:"project_id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	JobID       types.String `tfsdk:"job_id"`
	Count       types.Int64  `tfsdk:"count"`
}

func testWaitModel() *fakeWaitModel {
	return &fakeWaitModel{
		ProjectID:   types.StringValue("proj"),
		ClusterName: types.StringValue("cluster"),
		JobID:       types.StringValue("job1"),
		Count:       types.Int64Value(3),
	}
}

func TestDefaultFormatWaitFailure_withErrorDescription(t *testing.T) {
	wait := &autogen.WaitReq{
		ErrorDescriptionProperty: "errorMessage",
		IDAttributes:             []string{"project_id", "cluster_name", "job_id"},
	}
	req := autogen.WaitFailure{
		LastJSON: map[string]any{
			"state":        "FAILED",
			"errorMessage": "The restore could not complete because 1 collection was not found",
		},
		Model:     testWaitModel(),
		LastState: "FAILED",
	}

	err := autogen.DefaultFormatWaitFailure(wait, req)

	require.Error(t, err)
	assert.Equal(t, `project_id="proj", cluster_name="cluster", job_id="job1": The restore could not complete because 1 collection was not found`, err.Error())
}

func TestDefaultFormatWaitFailure_withoutErrorDescription(t *testing.T) {
	wait := &autogen.WaitReq{
		ErrorDescriptionProperty: "errorMessage",
		IDAttributes:             []string{"project_id", "cluster_name", "job_id"},
	}

	t.Run("property missing", func(t *testing.T) {
		err := autogen.DefaultFormatWaitFailure(wait, autogen.WaitFailure{
			LastJSON:  map[string]any{"state": "FAILED"},
			Model:     testWaitModel(),
			LastState: "FAILED",
		})
		require.Error(t, err)
		assert.Equal(t, `project_id="proj", cluster_name="cluster", job_id="job1": operation failed with state "FAILED"`, err.Error())
	})

	t.Run("property empty", func(t *testing.T) {
		err := autogen.DefaultFormatWaitFailure(wait, autogen.WaitFailure{
			LastJSON:  map[string]any{"errorMessage": ""},
			Model:     testWaitModel(),
			LastState: "FAILED",
		})
		require.Error(t, err)
		assert.Equal(t, `project_id="proj", cluster_name="cluster", job_id="job1": operation failed with state "FAILED"`, err.Error())
	})

	t.Run("property not a string", func(t *testing.T) {
		err := autogen.DefaultFormatWaitFailure(wait, autogen.WaitFailure{
			LastJSON:  map[string]any{"errorMessage": 12},
			Model:     testWaitModel(),
			LastState: "FAILED",
		})
		require.Error(t, err)
		assert.Equal(t, `project_id="proj", cluster_name="cluster", job_id="job1": operation failed with state "FAILED"`, err.Error())
	})
}

func TestDefaultFormatWaitFailure_timeoutWrapsSDKError(t *testing.T) {
	timeoutErr := &retry.TimeoutError{
		LastState:     "INITIALIZING",
		Timeout:       time.Minute,
		ExpectedState: []string{"SUCCESSFUL"},
	}
	wait := &autogen.WaitReq{
		IDAttributes: []string{"project_id", "cluster_name", "job_id"},
	}

	err := autogen.DefaultFormatWaitFailure(wait, autogen.WaitFailure{
		Model:      testWaitModel(),
		TimeoutErr: timeoutErr,
		LastState:  "INITIALIZING",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), `project_id="proj", cluster_name="cluster", job_id="job1":`)
	assert.Contains(t, err.Error(), "timeout while waiting")
	var got *retry.TimeoutError
	require.ErrorAs(t, err, &got)
	assert.Equal(t, timeoutErr, got)
}

func TestIsWaitContinueState(t *testing.T) {
	pending := []string{"INITIALIZING", "IN_PROGRESS"}
	target := []string{"SUCCESSFUL"}

	assert.True(t, autogen.IsWaitContinueState(pending, target, "INITIALIZING"))
	assert.True(t, autogen.IsWaitContinueState(pending, target, "SUCCESSFUL"))
	assert.False(t, autogen.IsWaitContinueState(pending, target, "FAILED"))
	assert.False(t, autogen.IsWaitContinueState(pending, target, "CANCELED"))
}

func TestDefaultFormatWaitFailure_timeoutPrefixWithoutID(t *testing.T) {
	timeoutErr := fmt.Errorf("timeout while waiting")
	err := autogen.DefaultFormatWaitFailure(&autogen.WaitReq{}, autogen.WaitFailure{
		TimeoutErr: timeoutErr,
		LastState:  "INITIALIZING",
	})
	require.Error(t, err)
	assert.Equal(t, "timeout while waiting", err.Error())
}
