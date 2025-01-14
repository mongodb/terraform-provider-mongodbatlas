package clean_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

const (
	itemsPerPage                   = 100
	keepProjectsCreatedWithinHours = 5
	retryInterval                  = 60 * time.Second
	retryAttempts                  = 10
)

var (
	botProjectPrefixes = []string{
		"cfn-test-bot-",
		"test-acc-tf-p-",
	}
	keptPrefixes = []string{
		"test-acc-tf-p-keep",
	}
	projectRetryDeleteErrors = []string{
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS",
		"CANNOT_CLOSE_GROUP_ACTIVE_PEERING_CONNECTIONS",
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_LAKES",
		"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_FEDERATION_PRIVATE_ENDPOINTS",
	}
)

func TestCleanProjectAndClusters(t *testing.T) {
	client := acc.ConnV2()
	ctx := context.Background()
	cleanOrg, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ORG"))
	if !cleanOrg {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_ORG=true to run")
	}
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	onlyZeroClusters, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ONLY_0_CLUSTERS"))
	skipProjectsAfter := time.Now().Add(-keepProjectsCreatedWithinHours * time.Hour)
	maxDeleteCount := 250
	projects, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.Group], *http.Response, error) {
		return client.ProjectsApi.ListProjects(ctx).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	})
	require.NoError(t, err)
	t.Logf("found %d projects (DRY_RUN=%t)", len(projects), dryRun)
	projectsToDelete := map[string]string{}
	for _, p := range projects {
		if len(projectsToDelete) > maxDeleteCount {
			t.Logf("reached max delete count %d", maxDeleteCount)
			break
		}
		skipReason := projectSkipReason(&p, skipProjectsAfter, onlyZeroClusters)
		projectName := p.GetName()
		if skipReason != "" {
			t.Logf("skip project %s, reason: %s", projectName, skipReason)
			continue
		}
		projectID := p.GetId()
		projectsToDelete[projectName] = projectID
	}
	for name, projectID := range projectsToDelete {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			changes := removeProjectResources(ctx, t, dryRun, client, projectID)
			if changes != "" {
				t.Logf("project %s %s", name, changes)
			}
			for i := range retryAttempts {
				attempt := i + 1
				if attempt > 1 {
					time.Sleep(retryInterval)
				}
				t.Logf("attempt %d to delete project %s", attempt, name)
				if dryRun {
					return
				}
				err = deleteProject(ctx, client, projectID)
				if err == nil {
					return
				}
				retryCode := findRetryErrorCode(err)
				if retryCode != "" {
					t.Logf("attempt %d, project %s has active resources, waiting, error: %s", attempt, projectID, retryCode)
					continue
				}
			}
			require.False(t, true, "failed to delete project %s: %s", name, err)
		})
	}
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
	_, _, err := client.ProjectsApi.DeleteProject(ctx, projectID).Execute()
	if err == nil || admin.IsErrorCode(err, "PROJECT_NOT_FOUND") {
		return nil
	}
	return err
}

func removeProjectResources(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) string {
	t.Helper()
	clustersRemoved := removeClusters(ctx, t, dryRun, client, projectID)
	changes := []string{}
	if clustersRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d clusters", clustersRemoved))
	}
	peeringsRemoved := removeNetworkPeering(ctx, t, dryRun, client, projectID)
	if peeringsRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d peerings", peeringsRemoved))
	}
	datalakesRemoved := removeDataLakePipelines(ctx, t, dryRun, client, projectID)
	if datalakesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d datalake pipelines", datalakesRemoved))
	}
	federatedDatabasesRemoved := removeFederatedDatabases(ctx, t, dryRun, client, projectID)
	if federatedDatabasesRemoved > 0 {
		changes = append(changes, fmt.Sprintf("removed %d federated databases", federatedDatabasesRemoved))
	}
	return strings.Join(changes, ", ")
}

func projectSkipReason(p *admin.Group, skipProjectsAfter time.Time, onlyEmpty bool) string {
	for _, blessedPrefix := range keptPrefixes {
		if strings.HasPrefix(p.GetName(), blessedPrefix) {
			return "blessed prefix: " + blessedPrefix
		}
	}
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

func removeClusters(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	require.NoError(t, err)
	clustersResults := clusters.GetResults()

	for i := range clustersResults {
		c := clustersResults[i]
		cName := c.GetName()
		t.Logf("delete cluster %s", cName)
		if !dryRun || c.GetStateName() != "DELETING" {
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

func removeNetworkPeering(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	peering, _, err := client.NetworkPeeringApi.ListPeeringConnections(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	require.NoError(t, err)
	peeringResults := peering.GetResults()
	for i := range peeringResults {
		p := peeringResults[i]
		peerID := p.GetId()
		t.Logf("delete peering %s", peerID)
		if !dryRun {
			_, _, err = client.NetworkPeeringApi.DeletePeeringConnection(ctx, projectID, peerID).Execute()
			if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") {
				t.Logf("peering %s already requested deletion", peerID)
				continue
			}
			require.NoError(t, err)
		}
	}
	return len(peeringResults)
}

func removeDataLakePipelines(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	datalakeResults, _, err := client.DataLakePipelinesApi.ListPipelines(ctx, projectID).Execute()
	require.NoError(t, err)
	for _, p := range datalakeResults {
		pipelineID := p.GetId()
		t.Logf("delete pipeline %s", pipelineID)
		if !dryRun {
			_, _, err = client.DataLakePipelinesApi.DeletePipeline(ctx, projectID, pipelineID).Execute()
			require.NoError(t, err)
		}
	}
	return len(datalakeResults)
}

func removeFederatedDatabases(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, projectID string) int {
	t.Helper()
	federatedResults, _, err := client.DataFederationApi.ListFederatedDatabases(ctx, projectID).Execute()
	if admin.IsErrorCode(err, "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_ID") {
		t.Logf("no federated databases found for project %s", projectID)
		return 0
	}
	require.NoError(t, err)
	for _, f := range federatedResults {
		federatedName := f.GetName()
		t.Logf("delete federated %s", federatedName)
		if !dryRun {
			_, _, err = client.DataFederationApi.DeleteFederatedDatabase(ctx, projectID, federatedName).Execute()
			require.NoError(t, err)
		}
	}
	return len(federatedResults)
}
