package acc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var errOutOfCapacity = errors.New("OUT_OF_CAPACITY")

func apiOutOfCapacityError() error {
	err := &admin.GenericOpenAPIError{}
	err.SetError("out of capacity")
	err.SetModel(*admin.NewApiError(400, "OUT_OF_CAPACITY"))
	return err
}

func unexpectedStateError() error {
	return &retry.UnexpectedStateError{State: "FAILED", ExpectedState: []string{"IDLE"}}
}

// attemptResult defines the stubbed outcome of one region attempt in
// CreateClusterWithRegionFallback: the error returned by create, then by wait if create
// succeeded, and by delete if the provisioning failure triggers a cleanup.
type attemptResult struct {
	createErr error
	waitErr   error
	deleteErr error
}

type clusterCreateCallTracker struct {
	createdRegions []string
	deletes        int
}

func (c *clusterCreateCallTracker) funcs(t *testing.T, attempts []attemptResult) acc.ClusterCreateFuncs {
	t.Helper()
	return acc.NewClusterCreateFuncs(
		func(region string) error {
			c.createdRegions = append(c.createdRegions, region)
			attempt := len(c.createdRegions) - 1
			require.Less(t, attempt, len(attempts), "create called more times than expected")
			return attempts[attempt].createErr
		},
		func() error {
			attempt := len(c.createdRegions) - 1
			require.NoError(t, attempts[attempt].createErr, "wait called although create failed")
			return attempts[attempt].waitErr
		},
		func() error {
			c.deletes++
			attempt := len(c.createdRegions) - 1
			require.Error(t, attempts[attempt].waitErr, "delete called although wait did not fail")
			return attempts[attempt].deleteErr
		},
	)
}

func TestCreateClusterWithRegionFallback(t *testing.T) {
	regions := []string{"R1", "R2", "R3"}
	testCases := map[string]struct {
		attempts    []attemptResult
		wantErr     string
		wantCreated []string
		wantDeletes int
	}{
		"success in first region": {
			attempts:    []attemptResult{{}},
			wantCreated: []string{"R1"},
		},
		"create OUT_OF_CAPACITY falls back to next region": {
			attempts: []attemptResult{
				{createErr: errOutOfCapacity},
				{},
			},
			wantCreated: []string{"R1", "R2"},
		},
		"create OUT_OF_CAPACITY api error code falls back to next region": {
			attempts: []attemptResult{
				{createErr: apiOutOfCapacityError()},
				{},
			},
			wantCreated: []string{"R1", "R2"},
		},
		"create OUT_OF_CAPACITY in all regions fails with combined error": {
			attempts: []attemptResult{
				{createErr: errOutOfCapacity},
				{createErr: errOutOfCapacity},
				{createErr: errOutOfCapacity},
			},
			wantErr:     "cluster creation failed in all candidate regions (R1,R2,R3), last error: OUT_OF_CAPACITY",
			wantCreated: regions,
		},
		"provisioning OUT_OF_CAPACITY in all regions deletes each and fails with combined error": {
			attempts: []attemptResult{
				{waitErr: errOutOfCapacity},
				{waitErr: errOutOfCapacity},
				{waitErr: errOutOfCapacity},
			},
			wantErr:     "cluster creation failed in all candidate regions (R1,R2,R3), last error: OUT_OF_CAPACITY",
			wantCreated: regions,
			wantDeletes: 3,
		},
		"unrelated create error fails fast": {
			attempts:    []attemptResult{{createErr: errors.New("boom")}},
			wantErr:     "boom",
			wantCreated: []string{"R1"},
		},
		"provisioning OUT_OF_CAPACITY deletes cluster and falls back": {
			attempts: []attemptResult{
				{waitErr: errOutOfCapacity},
				{},
			},
			wantCreated: []string{"R1", "R2"},
			wantDeletes: 1,
		},
		"provisioning unexpected state deletes cluster and falls back": {
			attempts: []attemptResult{
				{waitErr: unexpectedStateError()},
				{},
			},
			wantCreated: []string{"R1", "R2"},
			wantDeletes: 1,
		},
		"unrelated provisioning error fails fast": {
			attempts:    []attemptResult{{waitErr: errors.New("boom")}},
			wantErr:     "boom",
			wantCreated: []string{"R1"},
		},
		"delete failure fails fast without trying next region": {
			attempts:    []attemptResult{{waitErr: errOutOfCapacity, deleteErr: errors.New("delete boom")}},
			wantErr:     "cluster deletion failed: delete boom (after provisioning failure: OUT_OF_CAPACITY)",
			wantCreated: []string{"R1"},
			wantDeletes: 1,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tracker := &clusterCreateCallTracker{}
			err := acc.CreateClusterWithRegionFallback(t, regions, tracker.funcs(t, tc.attempts))
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
			assert.Equal(t, tc.wantCreated, tracker.createdRegions)
			assert.Equal(t, tc.wantDeletes, tracker.deletes)
		})
	}
}

func TestIsFailedProvisioning(t *testing.T) {
	assert.False(t, acc.IsFailedProvisioning(nil))
	assert.False(t, acc.IsFailedProvisioning(errors.New("boom")))
	assert.False(t, acc.IsFailedProvisioning(fmt.Errorf("wrapped: %w", errors.New("boom"))))
	assert.True(t, acc.IsFailedProvisioning(errOutOfCapacity))
	assert.True(t, acc.IsFailedProvisioning(apiOutOfCapacityError()))
	assert.True(t, acc.IsFailedProvisioning(unexpectedStateError()))
	assert.True(t, acc.IsFailedProvisioning(fmt.Errorf("wrapped: %w", unexpectedStateError())))
}
