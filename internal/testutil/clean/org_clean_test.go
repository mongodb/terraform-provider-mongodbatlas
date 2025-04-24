package clean_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	itemsPerPage                   = 100
	keepProjectsCreatedWithinHours = 5
	// Resource cleanup for a project can be slow, especially when there are active clusters, that can take more than 10 minutes to delete
	// Once 5 minutes are passed, we give up deleting and hope for the project to be deleted within the next run
	retryInterval = 60 * time.Second
	retryAttempts = 5
)

var (
	botProjectPrefixes = []string{
		"cfn-test-bot-",
		"test-acc-tf-p-",
	}
	// keptPrefixes has the prefix of the projects that we want to delete their resources but keep the projects themselves.
	// Useful when a feature flag or cloud provider is configured outside of the test
	keptPrefixes = []string{
		"test-acc-tf-p-keep",
	}
	projectRetryDeleteErrors = []string{
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS",
		"CANNOT_CLOSE_GROUP_ACTIVE_PEERING_CONNECTIONS",
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_LAKES",
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_FEDERATION_PRIVATE_ENDPOINTS",
		"CANNOT_CLOSE_GROUP_ACTIVE_STREAMS_RESOURCE",
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_PRIVATE_ENDPOINT_SERVICES",
	}
)

func TestSingleProjectRemoval(t *testing.T) {
	projectToClean := os.Getenv("MONGODB_ATLAS_CLEAN_PROJECT_ID")
	if projectToClean == "" {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_PROJECT_ID=project-id to run")
	}
	client := acc.ConnV2()
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	changes := removeProjectResources(t.Context(), t, dryRun, client, projectToClean)
	if changes != "" {
		t.Logf("project %s %s", projectToClean, changes)
	}
	err := deleteProject(t.Context(), client, projectToClean)
	require.NoError(t, err)
}

// Using a test to simplify logging and parallelization
func TestCleanProjectAndClusters(t *testing.T) {
	cleanOrg, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ORG"))
	if !cleanOrg {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_ORG=true to run")
	}
	client := acc.ConnV2()
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	onlyZeroClusters, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ONLY_WHEN_NO_CLUSTERS"))
	skipProjectsAfter := time.Now().Add(-keepProjectsCreatedWithinHours * time.Hour)
	retryAttemptsStr := os.Getenv("MONGODB_ATLAS_CLEAN_RETRY_ATTEMPTS")
	runRetries := retryAttempts
	if retryAttemptsStr != "" {
		attempts, err := strconv.Atoi(retryAttemptsStr)
		require.NoError(t, err)
		runRetries = attempts
	}
	projects := readAllProjects(t.Context(), t, client)
	projectsBefore := len(projects)
	t.Logf("found %d projects (DRY_RUN=%t)", projectsBefore, dryRun)
	projectsToClean := map[string]string{}
	projectInfos := []string{}
	for _, p := range projects {
		skipReason := projectSkipReason(&p, skipProjectsAfter, onlyZeroClusters)
		projectName := p.GetName()
		projectID := p.GetId()
		if skipReason != "" {
			t.Logf("skip project %s (%s), reason: %s", projectName, projectID, skipReason)
			continue
		}
		projectInfos = append(projectInfos, fmt.Sprintf("Project created at %s name %s (%s)", p.GetCreated().Format(time.RFC3339), projectName, p.GetId()))
		projectsToClean[projectName] = projectID
	}
	t.Logf("deleting project resources and optionally delete projects for %d projects", len(projectsToClean))
	slices.Sort(projectInfos)
	t.Log(strings.Join(projectInfos, "\n"))
	var deleteErrors int
	var emptyProjectCount int
	for name, projectID := range projectsToClean {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			changes := removeProjectResources(t.Context(), t, dryRun, client, projectID)
			if changes != "" {
				t.Logf("project %s %s", name, changes)
			}
			if skipProjectDelete(name) {
				t.Logf("keep project empty, but no delete %s (%s)", name, projectID)
				emptyProjectCount++
				return
			}
			var err error
			for i := range runRetries {
				attempt := i + 1
				if attempt > 1 {
					time.Sleep(retryInterval)
				}
				t.Logf("attempt %d to delete project %s", attempt, name)
				if dryRun {
					return
				}
				err = deleteProject(t.Context(), client, projectID)
				if err == nil {
					return
				}
				retryCode := findRetryErrorCode(err)
				if retryCode != "" {
					t.Logf("attempt %d, project %s has active resources, waiting, error: %s", attempt, projectID, retryCode)
					continue
				}
			}
			t.Logf("failed to delete project %s: %s", name, err)
			deleteErrors++
		})
	}
	t.Cleanup(func() {
		projectsAfter := readAllProjects(context.Background(), t, client) //nolint:usetesting // reason: using context.Background() here intentionally because t.Context() is canceled at cleanup
		t.Logf("SUMMARY\nProjects changed from %d to %d\ndelete_errors=%d\nempty_project_count=%d\nDRY_RUN=%t", projectsBefore, len(projectsAfter), deleteErrors, emptyProjectCount, dryRun)
	})
}

func readAllProjects(ctx context.Context, t *testing.T, client *admin.APIClient) []admin.Group {
	t.Helper()
	projects, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.Group], *http.Response, error) {
		return client.ProjectsApi.ListProjects(ctx).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	})
	require.NoError(t, err)
	return projects
}

func findRetryErrorCode(err error) string {
	if err == nil {
		return ""
	}
	for _, retryErr := range projectRetryDeleteErrors {
		if admin.IsErrorCode(err, retryErr) {
			return retryErr
		}
	}
	return ""
}
func deleteProject(ctx context.Context, client *admin.APIClient, projectID string) error {
	_, err := client.ProjectsApi.DeleteProject(ctx, projectID).Execute()
	if err == nil || admin.IsErrorCode(err, "PROJECT_NOT_FOUND") {
		return nil
	}
	return err
}

func removeProjectResources(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) string {
	t.Helper()
	changes := []string{}
	clustersRemoved := removeClusters(ctx, t, dryRun, client, projectID)
	if clustersRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d clusters", clustersRemoved))
	}
	serverlessClustersRemoved := removeServerlessClusters(ctx, t, dryRun, client, projectID)
	if serverlessClustersRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d serverless clusters", serverlessClustersRemoved))
	}
	peeringsRemoved := removeNetworkPeering(ctx, t, dryRun, client, projectID)
	if peeringsRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d peerings", peeringsRemoved))
	}
	datalakesRemoved := removeDataLakePipelines(ctx, t, dryRun, client, projectID)
	if datalakesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d datalake pipelines", datalakesRemoved))
	}
	federatedEndpointsRemoved := removeFederatedDatabasePrivateEndpoints(ctx, t, dryRun, client, projectID)
	if federatedEndpointsRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d federated private endpoints", federatedEndpointsRemoved))
	}
	federatedDatabasesRemoved := removeFederatedDatabases(ctx, t, dryRun, client, projectID)
	if federatedDatabasesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d federated databases", federatedDatabasesRemoved))
	}
	streamInstancesRemoved := removeStreamInstances(ctx, t, dryRun, client, projectID)
	if streamInstancesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d stream instances", streamInstancesRemoved))
	}
	privateEndpointServicesRemoved := removePrivateEndpointServices(ctx, t, dryRun, client, projectID)
	if privateEndpointServicesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d private endpoint services", privateEndpointServicesRemoved))
	}
	return strings.Join(changes, ", ")
}

func projectSkipReason(p *admin.Group, skipProjectsAfter time.Time, onlyEmpty bool) string {
	usesBotPrefix := false
	for _, botPrefix := range botProjectPrefixes {
		if strings.HasPrefix(p.GetName(), botPrefix) {
			usesBotPrefix = true
			break
		}
	}
	if !usesBotPrefix {
		return "not bot project"
	}
	if p.GetCreated().After(skipProjectsAfter) {
		return "created after " + skipProjectsAfter.Format("2006-01-02T15:04")
	}
	if onlyEmpty && p.GetClusterCount() > 0 {
		return "has clusters"
	}
	return ""
}

func skipProjectDelete(name string) bool {
	for _, keepPrefix := range keptPrefixes {
		if strings.HasPrefix(name, keepPrefix) {
			return true
		}
	}
	return false
}

func removeClusters(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	require.NoError(t, err)
	clustersResults := clusters.GetResults()

	for i := range clustersResults {
		c := clustersResults[i]
		cName := c.GetName()
		t.Logf("delete cluster %s", cName)
		if !dryRun {
			_, err = client.ClustersApi.DeleteCluster(ctx, projectID, cName).Execute()
			if admin.IsErrorCode(err, "CLUSTER_ALREADY_REQUESTED_DELETION") {
				t.Logf("cluster %s already requested deletion", cName)
				continue
			}
			require.NoError(t, err)
		}
	}
	return len(clustersResults)
}

func removeServerlessClusters(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	clusters, _, err := client.ServerlessInstancesApi.ListServerlessInstances(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	require.NoError(t, err)
	clustersResults := clusters.GetResults()
	for i := range clustersResults {
		c := clustersResults[i]
		cName := c.GetName()
		t.Logf("delete serverless cluster %s", cName)
		if !dryRun {
			_, _, err = client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, cName).Execute()
			if admin.IsErrorCode(err, "SERVERLESS_INSTANCE_ALREADY_REQUESTED_DELETION") {
				t.Logf("serverless cluster %s already requested deletion", cName)
				continue
			}
			require.NoError(t, err)
		}
	}
	return len(clustersResults)
}

func removeNetworkPeering(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	peeringIDs := []string{}
	for _, providerName := range []string{constant.AWS, constant.AZURE, constant.GCP} {
		peering, _, err := client.NetworkPeeringApi.ListPeeringConnectionsWithParams(ctx, &admin.ListPeeringConnectionsApiParams{
			ProviderName: &providerName,
			GroupId:      projectID,
		}).ItemsPerPage(itemsPerPage).Execute()
		require.NoError(t, err)
		peeringResults := peering.GetResults()
		for i := range peeringResults {
			p := peeringResults[i]
			peerID := p.GetId()
			peeringIDs = append(peeringIDs, peerID)
		}
	}
	for _, peerID := range peeringIDs {
		t.Logf("delete peering %s", peerID)
		if !dryRun {
			_, _, err := client.NetworkPeeringApi.DeletePeeringConnection(ctx, projectID, peerID).Execute()
			if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") {
				t.Logf("peering %s already requested deletion", peerID)
				continue
			}
			require.NoError(t, err)
		}
	}
	return len(peeringIDs)
}

func removeDataLakePipelines(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	datalakeResults, _, err := client.DataLakePipelinesApi.ListPipelines(ctx, projectID).Execute()
	require.NoError(t, err)
	for _, p := range datalakeResults {
		pipelineID := p.GetId()
		t.Logf("delete pipeline %s", pipelineID)
		if !dryRun {
			_, err = client.DataLakePipelinesApi.DeletePipeline(ctx, projectID, pipelineID).Execute()
			require.NoError(t, err)
		}
	}
	return len(datalakeResults)
}

func removeStreamInstances(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	streamInstances, _, err := client.StreamsApi.ListStreamInstances(ctx, projectID).Execute()
	require.NoError(t, err)

	for _, instance := range *streamInstances.Results {
		instanceName := *instance.Name
		id := instance.GetId()
		t.Logf("delete stream instance %s", id)

		if !dryRun {
			_, err = client.StreamsApi.DeleteStreamInstance(ctx, projectID, instanceName).Execute()
			if err != nil && admin.IsErrorCode(err, "STREAM_TENANT_HAS_STREAM_PROCESSORS") {
				t.Logf("stream instance %s has stream processors, attempting to delete", id)
				streamProcessors, _, spErr := client.StreamsApi.ListStreamProcessors(ctx, projectID, instanceName).Execute()
				require.NoError(t, spErr)

				for _, processor := range *streamProcessors.Results {
					t.Logf("delete stream processor %s", processor.Id)
					_, err = client.StreamsApi.DeleteStreamProcessor(ctx, projectID, instanceName, processor.Name).Execute()
					require.NoError(t, err)
				}
				t.Logf("retry delete stream instance %s after removing stream processors", id)
				_, err = client.StreamsApi.DeleteStreamInstance(ctx, projectID, instanceName).Execute()
				require.NoError(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
	return len(*streamInstances.Results)
}

func removePrivateEndpointServices(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	totalCount := 0
	cloudProviders := []string{"AWS", "AZURE", "GCP"}

	for _, provider := range cloudProviders {
		endpointServices, _, err := client.PrivateEndpointServicesApi.ListPrivateEndpointServices(ctx, projectID, provider).Execute()
		if err != nil {
			t.Errorf("failed to list private endpoint services for %s: %v", provider, err)
			continue
		}

		for _, service := range endpointServices {
			id := service.GetId()
			t.Logf("delete private endpoint service %s for provider %s", id, provider)
			if !dryRun {
				_, err := client.PrivateEndpointServicesApi.DeletePrivateEndpointService(ctx, projectID, provider, id).Execute()
				require.NoError(t, err)
			}
		}
		totalCount += len(endpointServices)
	}

	return totalCount
}

func removeFederatedDatabasePrivateEndpoints(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	paginatedResults, _, err := client.DataFederationApi.ListDataFederationPrivateEndpoints(ctx, projectID).Execute()
	require.NoError(t, err)
	endpoints := paginatedResults.GetResults()
	for _, f := range endpoints {
		endpointID := f.GetEndpointId()
		t.Logf("delete federated private endpoint %s", endpointID)
		if !dryRun {
			_, err = client.DataFederationApi.DeleteDataFederationPrivateEndpoint(ctx, projectID, endpointID).Execute()
			require.NoError(t, err)
		}
	}
	return len(endpoints)
}

func removeFederatedDatabases(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	federatedResults, _, err := client.DataFederationApi.ListFederatedDatabases(ctx, projectID).Execute()
	if admin.IsErrorCode(err, "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_ID") {
		t.Logf("no federated databases found for project %s, must delete this manually from the UI", projectID) // Deletion task was only partially successful - deleted the storage config but not the tenant config (internal slack thread)
		return 0
	}
	require.NoError(t, err)
	for _, f := range federatedResults {
		federatedName := f.GetName()
		t.Logf("delete federated %s", federatedName)
		if !dryRun {
			_, err = client.DataFederationApi.DeleteFederatedDatabase(ctx, projectID, federatedName).Execute()
			require.NoError(t, err)
		}
	}
	return len(federatedResults)
}
