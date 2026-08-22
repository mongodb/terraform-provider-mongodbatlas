package acc_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectionRestoreEnvFromOS_requiresAllThree(t *testing.T) {
	t.Setenv("MONGODB_ATLAS_PROJECT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLUSTER_NAME", "")
	t.Setenv("MONGODB_ATLAS_SNAPSHOT_ID", "")
	projectID, clusterName, snapshotID, ok := acc.CollectionRestoreEnvFromOSForTest()
	assert.False(t, ok)
	assert.Empty(t, projectID)
	assert.Empty(t, clusterName)
	assert.Empty(t, snapshotID)

	t.Setenv("MONGODB_ATLAS_PROJECT_ID", "p")
	t.Setenv("MONGODB_ATLAS_CLUSTER_NAME", "c")
	projectID, clusterName, snapshotID, ok = acc.CollectionRestoreEnvFromOSForTest()
	assert.False(t, ok)
	assert.Equal(t, "p", projectID)
	assert.Equal(t, "c", clusterName)
	assert.Empty(t, snapshotID)

	t.Setenv("MONGODB_ATLAS_SNAPSHOT_ID", "s")
	projectID, clusterName, snapshotID, ok = acc.CollectionRestoreEnvFromOSForTest()
	require.True(t, ok)
	assert.Equal(t, "p", projectID)
	assert.Equal(t, "c", clusterName)
	assert.Equal(t, "s", snapshotID)
}

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
