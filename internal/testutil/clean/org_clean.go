package clean

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312022/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// ErrUnauthorized signals a transient HTTP 401 while accessing a project's resources (e.g. a
// project a concurrent run is creating or tearing down). Callers skip the project and retry next run.
var ErrUnauthorized = errors.New("unauthorized accessing project resources, skipping for this run")

// SkipUnauthorizedErr maps a transient 401 to ErrUnauthorized so the caller can skip the project;
// any other error is returned unchanged.
func SkipUnauthorizedErr(resp *http.Response, err error) error {
	if validate.StatusUnauthorized(resp) {
		return ErrUnauthorized
	}
	return err
}

// RemoveStreamInstances deletes all stream instances in the project.
// It will also remove all stream processors associated with the stream instance.
func RemoveStreamInstances(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	streamInstances, resp, err := client.StreamsAPI.ListStreamWorkspaces(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}

	for _, instance := range streamInstances.GetResults() {
		instanceName := *instance.Name

		if !dryRun {
			_, err = client.StreamsAPI.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
			if err != nil && admin.IsErrorCode(err, "STREAM_TENANT_HAS_STREAM_PROCESSORS") {
				streamProcessors, spResp, spErr := client.StreamsAPI.GetStreamProcessors(ctx, projectID, instanceName).Execute()
				if spErr != nil {
					return 0, SkipUnauthorizedErr(spResp, spErr)
				}

				processors := streamProcessors.GetResults()
				for i := range processors {
					_, err = client.StreamsAPI.DeleteStreamProcessor(ctx, projectID, instanceName, processors[i].Name).Execute()
					if err != nil {
						return 0, err
					}
				}

				_, err = client.StreamsAPI.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
				if err != nil {
					return 0, err
				}
			} else if err != nil {
				return 0, err
			}
		}
	}
	return len(streamInstances.GetResults()), nil
}

// RemovePrivateLinkConnections deletes all Stream Processing Private Link connections in the project;
// left behind, they block project deletion with CANNOT_CLOSE_GROUP_ACTIVE_STREAMS_RESOURCE.
func RemovePrivateLinkConnections(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	connections, resp, err := client.StreamsAPI.ListPrivateLinkConnections(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	results := connections.GetResults()
	for i := range results {
		if !dryRun {
			_, err = client.StreamsAPI.DeletePrivateLinkConnection(ctx, projectID, results[i].GetId()).Execute()
			if admin.IsErrorCode(err, "STREAM_PRIVATE_LINK_IN_USE") {
				continue // still referenced by a stream connection, leave it for the next run
			}
			if err != nil {
				return 0, err
			}
		}
	}
	return len(results), nil
}
