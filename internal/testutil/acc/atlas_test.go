package acc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"

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

// TestRemainingProjectResources covers the diagnostics printed when a shared project cannot be deleted, so the
// teardown failure names the resource that blocked it. Flex clusters are listed separately because the Clusters API
// does not return clusters created through the Flex Clusters API.
func TestRemainingProjectResources(t *testing.T) {
	const projectID = "6ad95b900c3871ac99ea9ca7"
	cluster := admin.ClusterDescription20240805{Name: new("test-acc-tf-c-1"), StateName: new("DELETING")}
	flexCluster := admin.FlexClusterDescription20241113{Name: new("test-acc-tf-2"), StateName: new("IDLE")}
	clustersErr := errors.New("clusters boom")
	flexClustersErr := errors.New("flex clusters boom")
	testCases := map[string]struct {
		clusters        *admin.PaginatedClusterDescription20240805
		clustersErr     error
		flexClusters    *admin.PaginatedFlexClusters20241113
		flexClustersErr error
		want            string
	}{
		"reports clusters and flex clusters with their state": {
			clusters:     admin.NewPaginatedClusterDescription20240805([]admin.ClusterDescription20240805{cluster}),
			flexClusters: admin.NewPaginatedFlexClusters20241113([]admin.FlexClusterDescription20241113{flexCluster}),
			want:         "cluster test-acc-tf-c-1 (state DELETING), flex cluster test-acc-tf-2 (state IDLE)",
		},
		"reports when no cluster is returned": {
			clusters:     admin.NewPaginatedClusterDescription20240805([]admin.ClusterDescription20240805{}),
			flexClusters: admin.NewPaginatedFlexClusters20241113([]admin.FlexClusterDescription20241113{}),
			want:         "no clusters reported by the Clusters and Flex Clusters APIs",
		},
		"reports listing errors instead of failing": {
			clustersErr:     clustersErr,
			flexClustersErr: flexClustersErr,
			want:            "clusters could not be listed: clusters boom, flex clusters could not be listed: flex clusters boom",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			clustersAPI := mockadmin.NewClustersAPI(t)
			clustersAPI.EXPECT().ListClusters(mock.Anything, projectID).Return(admin.ListClustersApiRequest{ApiService: clustersAPI}).Once()
			clustersAPI.EXPECT().ListClustersExecute(mock.Anything).Return(tc.clusters, nil, tc.clustersErr).Once()
			flexClustersAPI := mockadmin.NewFlexClustersAPI(t)
			flexClustersAPI.EXPECT().ListFlexClusters(mock.Anything, projectID).Return(admin.ListFlexClustersApiRequest{ApiService: flexClustersAPI}).Once()
			flexClustersAPI.EXPECT().ListFlexClustersExecute(mock.Anything).Return(tc.flexClusters, nil, tc.flexClustersErr).Once()

			client := &admin.APIClient{ClustersAPI: clustersAPI, FlexClustersAPI: flexClustersAPI}
			assert.Equal(t, tc.want, acc.RemainingProjectResources(client, projectID))
		})
	}
}
