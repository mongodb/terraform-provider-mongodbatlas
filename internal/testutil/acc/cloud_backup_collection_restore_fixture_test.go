package acc_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"
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
	err := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-123", "FAILED", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "job-123")
	assert.Contains(t, err.Error(), "c1")
}

func TestErrIfSampleDatasetLoadFailed_failedStateReportsReason(t *testing.T) {
	err := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-123", "FAILED", "Failed to mongorestore")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to mongorestore")
}

func TestErrIfSampleDatasetLoadFailed_workingOrCompleted(t *testing.T) {
	for _, state := range []string{"WORKING", "COMPLETED"} {
		t.Run(state, func(t *testing.T) {
			require.NoError(t, acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-123", state, ""))
		})
	}
}

type stubSampleDatasetAttempt struct {
	requestErr    error
	jobID         string
	polls         []stubSampleDatasetPoll
	nilRequestJob bool
}

type stubSampleDatasetPoll struct {
	err          error
	state        string
	errorMessage string
	nilJob       bool
}

func TestPopulateWithSampleData(t *testing.T) {
	const (
		projectID   = "6b03890e7f7520b62b046f90"
		clusterName = "test-acc-tf-c-1845508866424288436"
		maxAttempts = 3
	)
	fastTiming := acc.NewSampleDatasetLoadTiming(time.Second, 500*time.Millisecond, time.Millisecond, time.Millisecond, time.Millisecond, maxAttempts)
	noPrimaryMsg := "error connecting to host: failed to connect to mongodb://host:27017: server selection error: server selection timeout, current topology: { Type: ReplicaSetNoPrimary }"
	dupKeyMsg := "continuing through error: E11000 duplicate key error collection: sample_training.grades index: _id_ dup key: { _id: ObjectId('56d5f7ef604eb380b0d9cee5') }"
	conflictMsg := "Target cluster has conflicting non-empty namespace with the dataset"

	testCases := map[string]struct {
		wantErr     string
		attempts    []stubSampleDatasetAttempt
		canceled    bool
		wantTainted bool
	}{
		"completes after WORKING poll": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{state: "WORKING"},
						{state: "COMPLETED"},
					},
				},
			},
		},
		"retries RequestSampleDatasetLoad after transient ReplicaSetNoPrimary failure": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{state: "FAILED", errorMessage: noPrimaryMsg},
					},
				},
				{
					jobID: "job-2",
					polls: []stubSampleDatasetPoll{
						{state: "COMPLETED"},
					},
				},
			},
		},
		"fails fast without same-cluster retry when mongorestore aborts with E11000 partial data": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{state: "WORKING"},
						{state: "FAILED", errorMessage: dupKeyMsg},
					},
				},
			},
			wantErr:     dupKeyMsg,
			wantTainted: true,
		},
		"fails after exhausting maxAttempts on repeated ReplicaSetNoPrimary": {
			attempts: []stubSampleDatasetAttempt{
				{jobID: "job-1", polls: []stubSampleDatasetPoll{{state: "FAILED", errorMessage: noPrimaryMsg}}},
				{jobID: "job-2", polls: []stubSampleDatasetPoll{{state: "FAILED", errorMessage: noPrimaryMsg}}},
				{jobID: "job-3", polls: []stubSampleDatasetPoll{{state: "FAILED", errorMessage: noPrimaryMsg}}},
			},
			wantErr: "sample dataset load job-3 failed for cluster " + projectID + ":" + clusterName + ": " + noPrimaryMsg,
		},
		"does not retry non-retryable FAILED reason": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{state: "FAILED", errorMessage: conflictMsg},
					},
				},
			},
			wantErr:     conflictMsg,
			wantTainted: true,
		},
		"stops retrying if second attempt hits conflicting namespace": {
			attempts: []stubSampleDatasetAttempt{
				{jobID: "job-1", polls: []stubSampleDatasetPoll{{state: "FAILED", errorMessage: noPrimaryMsg}}},
				{jobID: "job-2", polls: []stubSampleDatasetPoll{{state: "FAILED", errorMessage: conflictMsg}}},
			},
			wantErr:     conflictMsg,
			wantTainted: true,
		},
		"does not retry GetSampleDatasetLoad transport error containing ReplicaSetNoPrimary": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{err: errors.New("transport error: ReplicaSetNoPrimary")},
					},
				},
			},
			wantErr: "transport error: ReplicaSetNoPrimary",
		},
		"does not retry RequestSampleDatasetLoad API error": {
			attempts: []stubSampleDatasetAttempt{
				{requestErr: errors.New("request boom")},
			},
			wantErr: "loading sample data set error: request boom",
		},
		"fails fast when RequestSampleDatasetLoad returns nil job": {
			attempts: []stubSampleDatasetAttempt{
				{nilRequestJob: true},
			},
			wantErr: "loading sample data set returned nil job",
		},
		"fails fast when GetSampleDatasetLoad returns nil job": {
			attempts: []stubSampleDatasetAttempt{
				{
					jobID: "job-1",
					polls: []stubSampleDatasetPoll{
						{nilJob: true},
					},
				},
			},
			wantErr: "sample dataset load job-1 returned nil job",
		},
		"stops during initial delay when context is already canceled": {
			canceled: true,
			wantErr:  "context canceled",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tc.canceled {
				cancel()
			}

			clustersAPI := mockadmin.NewClustersAPI(t)
			for _, attempt := range tc.attempts {
				clustersAPI.EXPECT().
					RequestSampleDatasetLoad(mock.Anything, projectID, clusterName).
					Return(admin.RequestSampleDatasetLoadApiRequest{ApiService: clustersAPI}).
					Once()
				if attempt.requestErr != nil {
					clustersAPI.EXPECT().
						RequestSampleDatasetLoadExecute(mock.Anything).
						Return(nil, nil, attempt.requestErr).
						Once()
					continue
				}
				if attempt.nilRequestJob {
					clustersAPI.EXPECT().
						RequestSampleDatasetLoadExecute(mock.Anything).
						Return(nil, nil, nil).
						Once()
					continue
				}
				reqStatus := admin.NewSampleDatasetStatus()
				reqStatus.SetId(attempt.jobID)
				clustersAPI.EXPECT().
					RequestSampleDatasetLoadExecute(mock.Anything).
					Return(reqStatus, nil, nil).
					Once()

				for _, poll := range attempt.polls {
					clustersAPI.EXPECT().
						GetSampleDatasetLoad(mock.Anything, projectID, attempt.jobID).
						Return(admin.GetSampleDatasetLoadApiRequest{ApiService: clustersAPI}).
						Once()
					if poll.err != nil {
						clustersAPI.EXPECT().
							GetSampleDatasetLoadExecute(mock.Anything).
							Return(nil, nil, poll.err).
							Once()
						continue
					}
					if poll.nilJob {
						clustersAPI.EXPECT().
							GetSampleDatasetLoadExecute(mock.Anything).
							Return(nil, nil, nil).
							Once()
						continue
					}
					pollStatus := admin.NewSampleDatasetStatus()
					pollStatus.SetId(attempt.jobID)
					pollStatus.SetState(poll.state)
					if poll.errorMessage != "" {
						pollStatus.SetErrorMessage(poll.errorMessage)
					}
					clustersAPI.EXPECT().
						GetSampleDatasetLoadExecute(mock.Anything).
						Return(pollStatus, nil, nil).
						Once()
				}
			}

			err := acc.PopulateWithSampleDataWithTiming(ctx, clustersAPI, projectID, clusterName, fastTiming)
			assert.Equal(t, tc.wantTainted, acc.IsTaintedSampleDatasetClusterErr(err))
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}

	t.Run("stops during retry delay when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		polled := make(chan struct{})
		go func() {
			<-polled
			time.Sleep(200 * time.Millisecond)
			cancel()
		}()

		longRetryTiming := acc.NewSampleDatasetLoadTiming(15*time.Second, time.Second, time.Millisecond, 0, 10*time.Second, maxAttempts)
		clustersAPI := mockadmin.NewClustersAPI(t)
		clustersAPI.EXPECT().
			RequestSampleDatasetLoad(mock.Anything, projectID, clusterName).
			Return(admin.RequestSampleDatasetLoadApiRequest{ApiService: clustersAPI}).
			Once()
		reqStatus := admin.NewSampleDatasetStatus()
		reqStatus.SetId("job-1")
		clustersAPI.EXPECT().
			RequestSampleDatasetLoadExecute(mock.Anything).
			Return(reqStatus, nil, nil).
			Once()
		clustersAPI.EXPECT().
			GetSampleDatasetLoad(mock.Anything, projectID, "job-1").
			Return(admin.GetSampleDatasetLoadApiRequest{ApiService: clustersAPI}).
			Once()
		pollStatus := admin.NewSampleDatasetStatus()
		pollStatus.SetId("job-1")
		pollStatus.SetState("FAILED")
		pollStatus.SetErrorMessage(noPrimaryMsg)
		clustersAPI.EXPECT().
			GetSampleDatasetLoadExecute(mock.Anything).
			Run(func(admin.GetSampleDatasetLoadApiRequest) { close(polled) }).
			Return(pollStatus, nil, nil).
			Once()

		err := acc.PopulateWithSampleDataWithTiming(ctx, clustersAPI, projectID, clusterName, longRetryTiming)
		require.ErrorIs(t, err, context.Canceled)
		require.ErrorContains(t, err, "ReplicaSetNoPrimary")
	})
}

func TestCreateSeededClusterWithRetry(t *testing.T) {
	dupKeyErr := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-1", "FAILED", "continuing through error: E11000 duplicate key error collection: sample_training.grades")
	conflictErr := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-2", "FAILED", "Target cluster has conflicting non-empty namespace with the dataset")
	noPrimaryErr := acc.ErrIfSampleDatasetLoadFailedForTest("p1", "c1", "job-3", "FAILED", "server selection error: current topology: { Type: ReplicaSetNoPrimary }")

	assert.True(t, acc.IsTaintedSampleDatasetClusterErr(dupKeyErr))
	assert.True(t, acc.IsTaintedSampleDatasetClusterErr(conflictErr))
	assert.False(t, acc.IsTaintedSampleDatasetClusterErr(noPrimaryErr))
	assert.False(t, acc.IsTaintedSampleDatasetClusterErr(errors.New("E11000 duplicate key error")))

	type clusterAttempt struct {
		createErr   error
		populateErr error
		deleteErr   error
		name        string
	}

	testCases := map[string]struct {
		wantName      string
		wantErr       string
		attempts      []clusterAttempt
		wantPopulated []string
		wantDeleted   []string
	}{
		"succeeds on first cluster without deleting": {
			attempts: []clusterAttempt{
				{name: "c1"},
			},
			wantName:      "c1",
			wantPopulated: []string{"c1"},
		},
		"deletes tainted cluster and retries on fresh cluster after E11000": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: dupKeyErr},
				{name: "c2"},
			},
			wantName:      "c2",
			wantPopulated: []string{"c1", "c2"},
			wantDeleted:   []string{"c1"},
		},
		"deletes tainted cluster and retries on fresh cluster after conflicting namespace": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: conflictErr},
				{name: "c2"},
			},
			wantName:      "c2",
			wantPopulated: []string{"c1", "c2"},
			wantDeleted:   []string{"c1"},
		},
		"does not delete or recreate cluster on non-tainted populate failure": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: noPrimaryErr},
			},
			wantName:      "c1",
			wantErr:       "ReplicaSetNoPrimary",
			wantPopulated: []string{"c1"},
		},
		"returns second cluster and error when fresh cluster also fails with E11000": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: dupKeyErr},
				{name: "c2", populateErr: dupKeyErr},
			},
			wantName:      "c2",
			wantErr:       "E11000 duplicate key error",
			wantPopulated: []string{"c1", "c2"},
			wantDeleted:   []string{"c1"},
		},
		"preserves first cluster name and fails fast when deleting tainted cluster fails": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: dupKeyErr, deleteErr: errors.New("cluster deletion failed: delete boom")},
			},
			wantName:      "c1",
			wantErr:       "cluster deletion failed: delete boom",
			wantPopulated: []string{"c1"},
			wantDeleted:   []string{"c1"},
		},
		"fails fast when initial cluster creation fails": {
			attempts: []clusterAttempt{
				{createErr: errors.New("create boom")},
			},
			wantErr: "create boom",
		},
		"returns empty cluster name when second cluster creation fails after deleting first": {
			attempts: []clusterAttempt{
				{name: "c1", populateErr: dupKeyErr},
				{createErr: errors.New("second create boom")},
			},
			wantName:      "",
			wantErr:       "second create boom",
			wantPopulated: []string{"c1"},
			wantDeleted:   []string{"c1"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var (
				idx       int
				populated []string
				deleted   []string
			)
			fns := acc.NewSeededClusterFuncs(
				func() (string, error) {
					require.Less(t, idx, len(tc.attempts), "create called more times than expected")
					cur := tc.attempts[idx]
					return cur.name, cur.createErr
				},
				func(clusterName string) error {
					populated = append(populated, clusterName)
					return tc.attempts[idx].populateErr
				},
				func(clusterName string) error {
					deleted = append(deleted, clusterName)
					err := tc.attempts[idx].deleteErr
					idx++
					return err
				},
			)

			gotName, err := acc.CreateSeededClusterWithRetry(2, fns)
			assert.Equal(t, tc.wantName, gotName)
			assert.Equal(t, tc.wantPopulated, populated)
			assert.Equal(t, tc.wantDeleted, deleted)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestRequireSnapshotSeedCollections(t *testing.T) {
	const (
		projectID   = "6b03890e7f7520b62b046f90"
		clusterName = "test-acc-tf-c-1845508866424288436"
		snapshotID  = "6b038d00c461001f46454800"
	)

	t.Run("succeeds when both required seed collections exist", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(admin.NewDiskBackupCollectionResponse(acc.CollectionRestoreSeedCollectionName), &http.Response{StatusCode: http.StatusOK}, nil).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedRestaurantsDB, acc.CollectionRestoreSeedRestaurantsCol).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(admin.NewDiskBackupCollectionResponse(acc.CollectionRestoreSeedRestaurantsCol), &http.Response{StatusCode: http.StatusOK}, nil).
			Once()

		require.NoError(t, acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID))
	})

	t.Run("fails with missing message when first seed collection returns 404", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(nil, &http.Response{StatusCode: http.StatusNotFound}, errors.New("404 COLLECTION_NOT_FOUND")).
			Once()

		err := acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID)
		require.ErrorContains(t, err, "is missing required seed collection sample_mflix.movies")
		require.ErrorContains(t, err, "404 COLLECTION_NOT_FOUND")
	})

	t.Run("fails with missing message when second seed collection returns 404", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(admin.NewDiskBackupCollectionResponse(acc.CollectionRestoreSeedCollectionName), &http.Response{StatusCode: http.StatusOK}, nil).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedRestaurantsDB, acc.CollectionRestoreSeedRestaurantsCol).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(nil, &http.Response{StatusCode: http.StatusNotFound}, errors.New("404 COLLECTION_NOT_FOUND")).
			Once()

		err := acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID)
		require.ErrorContains(t, err, "is missing required seed collection sample_restaurants.restaurants")
	})

	t.Run("distinguishes non-404 API error from missing seed collection", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("500 INTERNAL_SERVER_ERROR")).
			Once()

		err := acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID)
		require.ErrorContains(t, err, "verify seed collection sample_mflix.movies")
		require.NotContains(t, err.Error(), "is missing required seed collection")
	})

	t.Run("handles nil http response on transport error", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(nil, nil, errors.New("dial tcp: i/o timeout")).
			Once()

		err := acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID)
		require.ErrorContains(t, err, "verify seed collection sample_mflix.movies")
		require.NotContains(t, err.Error(), "is missing required seed collection")
	})

	t.Run("fails when API returns nil collection without error", func(t *testing.T) {
		backupsAPI := mockadmin.NewCloudBackupsAPI(t)
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollection(mock.Anything, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, acc.CollectionRestoreSeedCollectionName).
			Return(admin.GetSnapshotDatabaseCollectionApiRequest{ApiService: backupsAPI}).
			Once()
		backupsAPI.EXPECT().
			GetSnapshotDatabaseCollectionExecute(mock.Anything).
			Return(nil, &http.Response{StatusCode: http.StatusOK}, nil).
			Once()

		err := acc.RequireSnapshotSeedCollections(t.Context(), backupsAPI, projectID, clusterName, snapshotID)
		require.ErrorContains(t, err, "returned nil seed collection sample_mflix.movies")
	})
}
