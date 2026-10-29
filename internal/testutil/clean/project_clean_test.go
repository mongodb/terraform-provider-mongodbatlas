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

func TestRemoveProjectResourcesFailure(t *testing.T) {
	clustersAPI := mockadmin.NewClustersAPI(t)
	clustersAPI.EXPECT().ListClusters(mock.Anything, projectID).Return(admin.ListClustersApiRequest{ApiService: clustersAPI}).Once()
	clustersAPI.EXPECT().ListClustersExecute(mock.Anything).
		Return(nil, &http.Response{StatusCode: http.StatusUnauthorized}, errors.New("unauthorized")).Once()

	// Leaving later APIs unset verifies that cleanup stops at the failed step.
	client := &admin.APIClient{ClustersAPI: clustersAPI}
	changes, err := clean.RemoveProjectResources(t.Context(), t.Logf, false, client, projectID)
	require.ErrorContains(t, err, "clusters:")
	require.ErrorIs(t, err, clean.ErrUnauthorized)
	assert.Empty(t, changes)
}

// TestRemoveFlexClusters covers the flex clusters step, needed because the Clusters API list does not return them.
func TestRemoveFlexClusters(t *testing.T) {
	testCases := map[string]struct {
		deleteErr   error
		deleteResp  *http.Response
		wantErr     string
		firstState  string
		wantDeletes int
		dryRun      bool
	}{
		"deletes every flex cluster": {
			firstState:  retrystrategy.RetryStrategyIdleState,
			wantDeletes: 2,
		},
		"a flex cluster already gone is not an error": {
			deleteResp:  &http.Response{StatusCode: http.StatusNotFound},
			deleteErr:   errors.New("not found"),
			firstState:  retrystrategy.RetryStrategyIdleState,
			wantDeletes: 2,
		},
		"returns other delete errors": {
			deleteResp:  &http.Response{StatusCode: http.StatusInternalServerError},
			deleteErr:   errors.New("boom"),
			wantErr:     "boom",
			firstState:  retrystrategy.RetryStrategyIdleState,
			wantDeletes: 1,
		},
		"a flex cluster already being deleted is left alone": {
			firstState:  retrystrategy.RetryStrategyDeletingState,
			wantDeletes: 1,
		},
		"dry run only counts them": {
			firstState: retrystrategy.RetryStrategyIdleState,
			dryRun:     true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			flexClusters := admin.NewPaginatedFlexClusters20241113([]admin.FlexClusterDescription20241113{
				{Name: new("flex-1"), StateName: new(tc.firstState)},
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
			count, err := clean.RemoveFlexClusters(t.Context(), t.Logf, tc.dryRun, client, projectID)
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	}
}

func TestRemoveFederatedDatabasesIncompleteDeletion(t *testing.T) {
	apiErr := apiError(http.StatusNotFound, "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_ID")
	dataFederationAPI := mockadmin.NewDataFederationAPI(t)
	dataFederationAPI.EXPECT().ListDataFederation(mock.Anything, projectID).Return(admin.ListDataFederationApiRequest{ApiService: dataFederationAPI}).Once()
	dataFederationAPI.EXPECT().ListDataFederationExecute(mock.Anything).Return(nil, nil, apiErr).Once()

	client := &admin.APIClient{DataFederationAPI: dataFederationAPI}
	count, err := clean.RemoveFederatedDatabases(t.Context(), t.Logf, false, client, projectID)
	require.ErrorContains(t, err, "requires manual cleanup")
	require.ErrorIs(t, err, apiErr)
	assert.Zero(t, count)
}

func TestRemoveClustersDeletionRace(t *testing.T) {
	testCases := map[string]struct {
		deleteErr  error
		deleteResp *http.Response
		wantErr    bool
	}{
		"already deleted HTTP response": {
			deleteErr:  errors.New("not found"),
			deleteResp: &http.Response{StatusCode: http.StatusNotFound},
		},
		"already deleted API error": {
			deleteErr: apiError(http.StatusNotFound, "CLUSTER_NOT_FOUND"),
		},
		"other deletion errors still fail": {
			deleteErr: errors.New("boom"),
			wantErr:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			clustersAPI := mockadmin.NewClustersAPI(t)
			clusters := &admin.PaginatedClusterDescription20240805{Results: []admin.ClusterDescription20240805{
				{Name: new("finishing-deletion"), StateName: new(retrystrategy.RetryStrategyDeletingState)},
				{Name: new("remaining-cluster"), StateName: new(retrystrategy.RetryStrategyIdleState)},
			}}
			clustersAPI.EXPECT().ListClusters(mock.Anything, projectID).Return(admin.ListClustersApiRequest{ApiService: clustersAPI}).Once()
			clustersAPI.EXPECT().ListClustersExecute(mock.Anything).Return(clusters, nil, nil).Once()
			clustersAPI.EXPECT().DeleteCluster(mock.Anything, projectID, "finishing-deletion").Return(admin.DeleteClusterApiRequest{ApiService: clustersAPI}).Once()
			clustersAPI.EXPECT().DeleteClusterExecute(mock.Anything).Return(tc.deleteResp, tc.deleteErr).Once()
			if !tc.wantErr {
				clustersAPI.EXPECT().DeleteCluster(mock.Anything, projectID, "remaining-cluster").Return(admin.DeleteClusterApiRequest{ApiService: clustersAPI}).Once()
				clustersAPI.EXPECT().DeleteClusterExecute(mock.Anything).Return(nil, nil).Once()
			}

			client := &admin.APIClient{ClustersAPI: clustersAPI}
			count, err := clean.RemoveClusters(t.Context(), t.Logf, false, client, projectID)
			if tc.wantErr {
				require.ErrorIs(t, err, tc.deleteErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	}
}

func apiError(status int, code string) error {
	err := &admin.GenericOpenAPIError{}
	err.SetError(code)
	err.SetModel(*admin.NewApiError(status, code))
	return err
}
