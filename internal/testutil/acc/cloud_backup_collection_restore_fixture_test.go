package acc_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAfterSnapshotUTCSeconds_convertsToUnixUTC(t *testing.T) {
	created := time.Date(2026, 8, 20, 15, 4, 5, 0, time.FixedZone("CEST", 2*60*60))
	want := time.Date(2026, 8, 20, 13, 4, 5, 0, time.UTC).Unix()
	assert.Equal(t, want, acc.AfterSnapshotUTCSecondsForTest(created))
}

func TestErrIfCloudBackupOrPITDisabled_failsWhenBackupOff(t *testing.T) {
	err := acc.ErrIfCloudBackupOrPITDisabledForTest("c1", false, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "c1")
	assert.Contains(t, err.Error(), "cloud backup")
}

func TestErrIfCloudBackupOrPITDisabled_failsWhenPITOff(t *testing.T) {
	err := acc.ErrIfCloudBackupOrPITDisabledForTest("c1", true, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "c1")
	assert.Contains(t, err.Error(), "point-in-time")
}

func TestErrIfCloudBackupOrPITDisabled_okWhenBothOn(t *testing.T) {
	require.NoError(t, acc.ErrIfCloudBackupOrPITDisabledForTest("c1", true, true))
}

func TestErrIfSampleDatasetLoadFailed_failedState(t *testing.T) {
	err := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-123", "FAILED")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "job-123")
	assert.Contains(t, err.Error(), "c1")
}

func TestErrIfSampleDatasetLoadFailed_workingOrCompleted(t *testing.T) {
	for _, state := range []string{"WORKING", "COMPLETED"} {
		t.Run(state, func(t *testing.T) {
			require.NoError(t, acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-123", state))
		})
	}
}
