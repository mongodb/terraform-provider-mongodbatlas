package clean

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312020/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// ErrUnauthorized signals a transient HTTP 401 while accessing a project's resources.
// It happens for projects that are being created or torn down by a concurrent test run.
// Callers should skip the project for this run and let the next run retry it instead of failing.
var ErrUnauthorized = errors.New("unauthorized accessing project resources, skipping for this run")

// SkipUnauthorizedErr returns ErrUnauthorized when the response is a transient 401, so the caller
// can skip the project. Any other error is returned unchanged.
func SkipUnauthorizedErr(resp *http.Response, err error) error {
	if err != nil && validate.StatusUnauthorized(resp) {
		return ErrUnauthorized
	}
	return err
}

// RemoveStreamInstances deletes all stream instances in the project.
// It will also remove all stream processors associated with the stream instance.
func RemoveStreamInstances(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	streamInstances, resp, err := client.StreamsApi.ListStreamWorkspaces(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}

	for _, instance := range streamInstances.GetResults() {
		instanceName := *instance.Name

		if !dryRun {
			_, err = client.StreamsApi.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
			if err != nil && admin.IsErrorCode(err, "STREAM_TENANT_HAS_STREAM_PROCESSORS") {
				streamProcessors, _, spErr := client.StreamsApi.GetStreamProcessors(ctx, projectID, instanceName).Execute()
				if spErr != nil {
					return 0, spErr
				}

				for _, processor := range streamProcessors.GetResults() {
					_, err = client.StreamsApi.DeleteStreamProcessor(ctx, projectID, instanceName, processor.Name).Execute()
					if err != nil {
						return 0, err
					}
				}

				_, err = client.StreamsApi.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
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

// RemovePrivateLinkConnections deletes all Stream Processing Private Link connections in the project.
// Left behind, they block project deletion with CANNOT_CLOSE_GROUP_ACTIVE_STREAMS_RESOURCE.
func RemovePrivateLinkConnections(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	connections, resp, err := client.StreamsApi.ListPrivateLinkConnections(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	results := connections.GetResults()
	for i := range results {
		if !dryRun {
			if _, err = client.StreamsApi.DeletePrivateLinkConnection(ctx, projectID, results[i].GetId()).Execute(); err != nil {
				return 0, err
			}
		}
	}
	return len(results), nil
}
