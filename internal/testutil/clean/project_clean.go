package clean

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312027/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

const itemsPerPage = 100

// projectRetryDeleteErrors are the project deletion errors caused by dependent resources that are still terminating,
// so retrying the deletion later can succeed.
var projectRetryDeleteErrors = []string{
	"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS",
	"CANNOT_CLOSE_GROUP_ACTIVE_PEERING_CONNECTIONS",
	"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_LAKES",
	"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_FEDERATION_PRIVATE_ENDPOINTS",
	"CANNOT_CLOSE_GROUP_ACTIVE_STREAMS_RESOURCE",
	"CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_PRIVATE_ENDPOINT_SERVICES",
}

// ProjectDeleteRetryCode returns the error code when a project deletion failed because of dependent resources that
// are still terminating, and an empty string for nil or any other error.
func ProjectDeleteRetryCode(err error) string {
	if err == nil {
		return ""
	}
	for _, code := range projectRetryDeleteErrors {
		if admin.IsErrorCode(err, code) {
			return code
		}
	}
	return ""
}

// DeleteProject deletes a project, treating a project that no longer exists as deleted.
func DeleteProject(ctx context.Context, client *admin.APIClient, projectID string) error {
	_, err := client.ProjectsAPI.DeleteGroup(ctx, projectID).Execute()
	if err == nil || admin.IsErrorCode(err, "PROJECT_NOT_FOUND") {
		return nil
	}
	return err
}

// RemoveProjectResources deletes the resources that block a project deletion: clusters, flex clusters, network
// peerings, federated databases and private endpoints, stream instances and private link connections, private
// endpoint services, and encryption at rest private endpoints. It stops at the first step that fails and returns what
// it removed. It only requests the deletions, so callers retry the project deletion while the resources terminate.
// Both the org cleanup job and the acceptance-test teardown use it.
// logf receives the per-resource progress lines, so callers pass `t.Logf` in tests and a print function elsewhere.
func RemoveProjectResources(ctx context.Context, logf func(format string, args ...any), dryRun bool, client *admin.APIClient, projectID string) (string, error) {
	steps := []struct {
		run   func() (int, error)
		label string
	}{
		{label: "clusters", run: func() (int, error) { return removeClusters(ctx, logf, dryRun, client, projectID) }},
		{label: "flex clusters", run: func() (int, error) { return removeFlexClusters(ctx, logf, dryRun, client, projectID) }},
		{label: "peerings", run: func() (int, error) { return removeNetworkPeering(ctx, logf, dryRun, client, projectID) }},
		{label: "federated private endpoints", run: func() (int, error) {
			return removeFederatedDatabasePrivateEndpoints(ctx, logf, dryRun, client, projectID)
		}},
		{label: "federated databases", run: func() (int, error) { return removeFederatedDatabases(ctx, logf, dryRun, client, projectID) }},
		{label: "stream instances", run: func() (int, error) { return RemoveStreamInstances(ctx, dryRun, client, projectID) }},
		{label: "stream private link connections", run: func() (int, error) { return RemovePrivateLinkConnections(ctx, dryRun, client, projectID) }},
		{label: "private endpoint services", run: func() (int, error) { return removePrivateEndpointServices(ctx, logf, dryRun, client, projectID) }},
		{label: "encryption at rest private endpoints", run: func() (int, error) {
			return removeEncryptionAtRestPrivateEndpoints(ctx, logf, dryRun, client, projectID)
		}},
	}
	changes := []string{}
	for _, step := range steps {
		count, err := step.run()
		if err != nil {
			return "", fmt.Errorf("%s: %w", step.label, err)
		}
		if count > 0 {
			changes = append(changes, fmt.Sprintf("removed %d %s", count, step.label))
		}
	}
	return strings.Join(changes, ", "), nil
}

func removeClusters(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	clusters, resp, err := client.ClustersAPI.ListClusters(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	clustersResults := clusters.GetResults()
	deleteFailures := []string{}

	for i := range clustersResults {
		c := clustersResults[i]
		cName := c.GetName()
		logf("delete cluster %s", cName)
		if !dryRun {
			_, err = client.ClustersAPI.DeleteCluster(ctx, projectID, cName).Execute()
			if admin.IsErrorCode(err, "CANNOT_TERMINATE_CLUSTER_WITH_UNDERGOING_REGIONAL_OUTAGE_SIMULATION") {
				logf("cluster %s has ongoing region outage simulation, deleting it now", cName)
				_, _, err = client.ClusterOutageSimulationAPI.EndOutageSimulation(ctx, projectID, cName).Execute()
				if err != nil {
					deleteFailures = append(deleteFailures, fmt.Sprintf("Unable to delete %s (cluster simulation stuck), might require manual action: %s", cName, err))
					continue
				}
				_, err = client.ClustersAPI.DeleteCluster(ctx, projectID, cName).Execute() // Retry deletion after ending outage simulation
			}
			if admin.IsErrorCode(err, "CLUSTER_ALREADY_REQUESTED_DELETION") {
				logf("cluster %s already requested deletion", cName)
				continue
			}
			if err != nil {
				return 0, err
			}
		}
	}
	if len(deleteFailures) > 0 {
		return 0, errors.New(strings.Join(deleteFailures, "\n"))
	}
	return len(clustersResults), nil
}

// removeFlexClusters deletes the flex clusters, which the Clusters API list does not return when they were created
// through the Flex Clusters API. A flex cluster that is already gone is not an error.
func removeFlexClusters(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	flexClusters, resp, err := client.FlexClustersAPI.ListFlexClusters(ctx, projectID).ItemsPerPage(itemsPerPage).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	results := flexClusters.GetResults()
	for i := range results {
		name := results[i].GetName()
		logf("delete flex cluster %s", name)
		if !dryRun {
			if resp, err := client.FlexClustersAPI.DeleteFlexCluster(ctx, projectID, name).Execute(); err != nil && !validate.StatusNotFound(resp) {
				return 0, err
			}
		}
	}
	return len(results), nil
}

func removeNetworkPeering(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	peeringIDs := []string{}
	for _, providerName := range []string{constant.AWS, constant.AZURE, constant.GCP} {
		peering, resp, err := client.NetworkPeeringAPI.ListGroupPeersWithParams(ctx, &admin.ListGroupPeersApiParams{
			ProviderName: &providerName,
			GroupId:      projectID,
		}).ItemsPerPage(itemsPerPage).Execute()
		if err != nil {
			return 0, SkipUnauthorizedErr(resp, err)
		}
		peeringResults := peering.GetResults()
		for i := range peeringResults {
			p := peeringResults[i]
			peerID := p.GetId()
			peeringIDs = append(peeringIDs, peerID)
		}
	}
	for _, peerID := range peeringIDs {
		logf("delete peering %s", peerID)
		if !dryRun {
			_, _, err := client.NetworkPeeringAPI.DeleteGroupPeer(ctx, projectID, peerID).Execute()
			if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") {
				logf("peering %s already requested deletion", peerID)
				continue
			}
			if err != nil {
				return 0, err
			}
		}
	}
	return len(peeringIDs), nil
}

func removePrivateEndpointServices(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	totalCount := 0
	cloudProviders := []string{constant.AWS, constant.AZURE, constant.GCP}

	for _, provider := range cloudProviders {
		endpointServices, resp, err := client.PrivateEndpointServicesAPI.ListPrivateEndpointService(ctx, projectID, provider).Execute()
		if err != nil {
			return 0, SkipUnauthorizedErr(resp, err)
		}

		for i := range endpointServices {
			id := endpointServices[i].GetId()
			logf("delete private endpoint service %s for provider %s", id, provider)
			if !dryRun {
				if _, err := client.PrivateEndpointServicesAPI.DeletePrivateEndpointService(ctx, projectID, provider, id).Execute(); err != nil {
					return 0, err
				}
			}
		}
		totalCount += len(endpointServices)
	}

	return totalCount, nil
}

func removeFederatedDatabasePrivateEndpoints(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	paginatedResults, resp, err := client.DataFederationAPI.ListPrivateEndpointIds(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	endpoints := paginatedResults.GetResults()
	for _, f := range endpoints {
		endpointID := f.GetEndpointId()
		logf("delete federated private endpoint %s", endpointID)
		if !dryRun {
			if _, err = client.DataFederationAPI.DeletePrivateEndpointId(ctx, projectID, endpointID).Execute(); err != nil {
				return 0, err
			}
		}
	}
	return len(endpoints), nil
}

func removeFederatedDatabases(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	federatedResults, resp, err := client.DataFederationAPI.ListDataFederation(ctx, projectID).Execute()
	if admin.IsErrorCode(err, "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_ID") {
		logf("no federated databases found for project %s, must delete this manually from the UI", projectID) // Deletion task was only partially successful - deleted the storage config but not the tenant config (internal slack thread)
		return 0, nil
	}
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	for _, f := range federatedResults {
		federatedName := f.GetName()
		logf("delete federated %s", federatedName)
		if !dryRun {
			if _, err = client.DataFederationAPI.DeleteDataFederation(ctx, projectID, federatedName).Execute(); err != nil {
				return 0, err
			}
		}
	}
	return len(federatedResults), nil
}

func removeEncryptionAtRestPrivateEndpoints(ctx context.Context, logf func(string, ...any), dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	endpointsCount := 0
	for _, cloudProvider := range []string{constant.AWS, constant.AZURE} {
		privateEndpoints, resp, err := client.EncryptionAtRestUsingCustomerKeyManagementAPI.ListRestPrivateEndpoints(ctx, projectID, cloudProvider).Execute()
		if err != nil {
			return 0, SkipUnauthorizedErr(resp, err)
		}
		endpoints := privateEndpoints.GetResults()
		endpointsCount += len(endpoints)
		for _, endpoint := range endpoints {
			endpointID := endpoint.GetId()
			logf("delete encryption at rest private endpoint %s", endpointID)
			if !dryRun {
				if _, err = client.EncryptionAtRestUsingCustomerKeyManagementAPI.RequestPrivateEndpointDeletion(ctx, projectID, cloudProvider, endpointID).Execute(); err != nil {
					return 0, err
				}
				time.Sleep(30 * time.Second)
			}
		}
	}
	return endpointsCount, nil
}
