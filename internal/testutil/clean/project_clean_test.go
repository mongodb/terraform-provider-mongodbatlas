package clean_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/clean"
)

const projectID = "6ad95b900c3871ac99ea9ca7"

func noLog(string, ...any) {}

// TestRemoveProjectResourcesFailure covers how a failed step is reported: the error names the step, so the teardown
// failure says what blocked the cleanup, and a transient 401 still matches ErrUnauthorized so the org cleanup skips
// the project. No later step runs after a failure.
func TestRemoveProjectResourcesFailure(t *testing.T) {
	testCases := map[string]struct {
		listErr          error
		resp             *http.Response
		wantErr          string
		wantUnauthorized bool
	}{
		"names the failed step": {
			listErr: errors.New("boom"),
			wantErr: "clusters: boom",
		},
		"keeps unauthorized detectable": {
			resp:             &http.Response{StatusCode: http.StatusUnauthorized},
			listErr:          errors.New("unauthorized"),
			wantErr:          "clusters: " + clean.ErrUnauthorized.Error(),
			wantUnauthorized: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			clustersAPI := mockadmin.NewClustersAPI(t)
			clustersAPI.EXPECT().ListClusters(mock.Anything, projectID).Return(admin.ListClustersApiRequest{ApiService: clustersAPI}).Once()
			clustersAPI.EXPECT().ListClustersExecute(mock.Anything).Return(nil, tc.resp, tc.listErr).Once()

			client := &admin.APIClient{ClustersAPI: clustersAPI}
			changes, err := clean.RemoveProjectResources(t.Context(), noLog, false, client, projectID)
			require.EqualError(t, err, tc.wantErr)
			assert.Empty(t, changes)
			assert.Equal(t, tc.wantUnauthorized, errors.Is(err, clean.ErrUnauthorized))
		})
	}
}

// TestRemoveFlexClusters covers the flex clusters step, needed because the Clusters API list does not return them.
func TestRemoveFlexClusters(t *testing.T) {
	testCases := map[string]struct {
		deleteErr   error
		deleteResp  *http.Response
		wantErr     string
		wantCount   int
		wantDeletes int
		deleting    bool // when true, the first flex cluster is already being deleted
		dryRun      bool
	}{
		"deletes every flex cluster": {
			wantCount:   2,
			wantDeletes: 2,
		},
		"a flex cluster already gone is not an error": {
			deleteResp:  &http.Response{StatusCode: http.StatusNotFound},
			deleteErr:   errors.New("not found"),
			wantCount:   2,
			wantDeletes: 2,
		},
		"returns other delete errors": {
			deleteResp:  &http.Response{StatusCode: http.StatusInternalServerError},
			deleteErr:   errors.New("boom"),
			wantErr:     "boom",
			wantDeletes: 1,
		},
		"a flex cluster already being deleted is left alone": {
			deleting:    true,
			wantCount:   2,
			wantDeletes: 1,
		},
		"dry run only counts them": {
			dryRun:    true,
			wantCount: 2,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			flexClusters := admin.NewPaginatedFlexClusters20241113([]admin.FlexClusterDescription20241113{
				{Name: new("flex-1"), StateName: new(stateName(tc.deleting))},
				{Name: new("flex-2"), StateName: new(retrystrategy.RetryStrategyIdleState)},
			})
			flexClustersAPI := mockadmin.NewFlexClustersAPI(t)
			flexClustersAPI.EXPECT().ListFlexClusters(mock.Anything, projectID).Return(admin.ListFlexClustersApiRequest{ApiService: flexClustersAPI}).Once()
			flexClustersAPI.EXPECT().ListFlexClustersExecute(mock.Anything).Return(flexClusters, nil, nil).Once()
			if !tc.dryRun {
				flexClustersAPI.EXPECT().DeleteFlexCluster(mock.Anything, projectID, mock.Anything).Return(admin.DeleteFlexClusterApiRequest{ApiService: flexClustersAPI}).Times(tc.wantDeletes)
				flexClustersAPI.EXPECT().DeleteFlexClusterExecute(mock.Anything).Return(tc.deleteResp, tc.deleteErr).Times(tc.wantDeletes)
			}

			client := &admin.APIClient{FlexClustersAPI: flexClustersAPI}
			count, err := clean.RemoveFlexClusters(t.Context(), noLog, tc.dryRun, client, projectID)
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func stateName(deleting bool) string {
	if deleting {
		return retrystrategy.RetryStrategyDeletingState
	}
	return retrystrategy.RetryStrategyIdleState
}
