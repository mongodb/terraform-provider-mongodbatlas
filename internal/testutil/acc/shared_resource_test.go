package acc_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"
)

func Test_NextProjectIDClusterName(t *testing.T) {
	projectIDReturner := func(name string) string {
		return fmt.Sprintf("%s-name", name)
	}
	projectIDs := map[string]int{}
	clusterNames := map[string]int{}

	addProjectIDClusterName := func(nodeCount int, freeTierClusterCount int) {
		projectID, clusterName := acc.NextProjectIDClusterName(nodeCount, freeTierClusterCount, projectIDReturner)
		projectIDs[projectID]++
		clusterNames[clusterName]++
	}
	for range acc.MaxClusterNodesPerProject {
		addProjectIDClusterName(1, 0)
	}
	assert.Len(t, projectIDs, 1)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject)
	addProjectIDClusterName(1, 0)
	assert.Len(t, projectIDs, 2)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+1)
	addProjectIDClusterName(acc.MaxClusterNodesPerProject, 0)
	assert.Len(t, projectIDs, 3)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+2)
	addProjectIDClusterName(1, 0)
	assert.Len(t, projectIDs, 4)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+3)
	addProjectIDClusterName(0, 1) // adds free tier, shares existing project
	assert.Len(t, projectIDs, 4)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+4)
	addProjectIDClusterName(0, 1) // second free tier cluster creates a new project
	assert.Len(t, projectIDs, 5)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+5)
}

func TestDeleteLocalProjectSharedResources(t *testing.T) {
	const (
		projectID    = "6ad95b900c3871ac99ea9ca7"
		instanceName = "execution-instance"
	)

	t.Run("does not touch execution projects or unowned resources", func(t *testing.T) {
		require.False(t, acc.DeleteLocalProjectSharedResources(t.Context(), nil, "", "cluster", instanceName))
		require.False(t, acc.DeleteLocalProjectSharedResources(t.Context(), nil, projectID, "", ""))
	})

	t.Run("deletes only the run's stream instance and its processors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
		defer cancel()
		streamsAPI := mockadmin.NewStreamsAPI(t)
		streamsAPI.EXPECT().DeleteStreamWorkspace(ctx, projectID, instanceName).
			Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI}).Twice()
		streamsAPI.EXPECT().DeleteStreamWorkspaceExecute(mock.Anything).
			Return(nil, apiError("STREAM_TENANT_HAS_STREAM_PROCESSORS")).Once()
		streamsAPI.EXPECT().GetStreamProcessors(ctx, projectID, instanceName).
			Return(admin.GetStreamProcessorsApiRequest{ApiService: streamsAPI}).Once()
		streamsAPI.EXPECT().GetStreamProcessorsExecute(mock.Anything).
			Return(&admin.PaginatedApiStreamsStreamProcessorWithStats{
				Results: []admin.StreamsProcessorWithStats{{Name: "execution-processor"}},
			}, nil, nil).Once()
		streamsAPI.EXPECT().DeleteStreamProcessor(ctx, projectID, instanceName, "execution-processor").
			Return(admin.DeleteStreamProcessorApiRequest{ApiService: streamsAPI}).Once()
		streamsAPI.EXPECT().DeleteStreamProcessorExecute(mock.Anything).Return(nil, nil).Once()
		streamsAPI.EXPECT().DeleteStreamWorkspaceExecute(mock.Anything).Return(nil, nil).Once()

		require.False(t, acc.DeleteLocalProjectSharedResources(ctx, &admin.APIClient{StreamsAPI: streamsAPI}, projectID, "", instanceName))
	})

	t.Run("a stream instance already gone is not an error", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsAPI(t)
		streamsAPI.EXPECT().DeleteStreamWorkspace(mock.Anything, projectID, instanceName).
			Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI}).Once()
		streamsAPI.EXPECT().DeleteStreamWorkspaceExecute(mock.Anything).
			Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("not found")).Once()

		require.False(t, acc.DeleteLocalProjectSharedResources(t.Context(), &admin.APIClient{StreamsAPI: streamsAPI}, projectID, "", instanceName))
	})

	t.Run("stops the cluster deletion wait when the teardown context ends", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		clustersAPI := mockadmin.NewClustersAPI(t)
		clustersAPI.EXPECT().DeleteCluster(ctx, projectID, "execution-cluster").
			Return(admin.DeleteClusterApiRequest{ApiService: clustersAPI}).Once()
		clustersAPI.EXPECT().DeleteClusterExecute(mock.Anything).
			Run(func(admin.DeleteClusterApiRequest) { cancel() }).Return(nil, nil).Once()

		require.True(t, acc.DeleteLocalProjectSharedResources(ctx, &admin.APIClient{ClustersAPI: clustersAPI}, projectID, "execution-cluster", ""))
	})

	t.Run("reports stream deletion failure with the teardown context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		streamsAPI := mockadmin.NewStreamsAPI(t)
		streamsAPI.EXPECT().DeleteStreamWorkspace(ctx, projectID, instanceName).
			Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI}).Once()
		streamsAPI.EXPECT().DeleteStreamWorkspaceExecute(mock.Anything).Return(nil, context.Canceled).Once()

		require.True(t, acc.DeleteLocalProjectSharedResources(ctx, &admin.APIClient{StreamsAPI: streamsAPI}, projectID, "", instanceName))
	})
}

func TestDeleteExecutionProjectsDoesNotBlockOtherProjects(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	fastProjectStarted := make(chan struct{})
	projectsAPI := mockadmin.NewProjectsAPI(t)
	projectsAPI.EXPECT().DeleteGroup(ctx, "slow-project").
		Run(func(context.Context, string) {
			select {
			case <-fastProjectStarted:
			case <-ctx.Done():
			}
		}).Return(admin.DeleteGroupApiRequest{ApiService: projectsAPI}).Once()
	projectsAPI.EXPECT().DeleteGroup(ctx, "fast-project").
		Run(func(context.Context, string) { close(fastProjectStarted) }).
		Return(admin.DeleteGroupApiRequest{ApiService: projectsAPI}).Once()
	projectsAPI.EXPECT().DeleteGroupExecute(mock.Anything).Return(nil, nil).Twice()

	acc.DeleteExecutionProjects(ctx, &admin.APIClient{ProjectsAPI: projectsAPI}, []string{"slow-project", "fast-project"})
	require.NoError(t, ctx.Err(), "the slow project must not prevent another project's deletion before the deadline")
}
