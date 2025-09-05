package clean

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

// RemoveStreamInstances deletes all stream instances in the project.
// It will also remove all stream processors associated with the stream instance.
func RemoveStreamInstances(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	streamInstances, _, err := client.StreamsApi.ListStreamWorkspaces(ctx, projectID).Execute()
	if err != nil {
		return 0, err
	}

	for _, instance := range *streamInstances.Results {
		instanceName := *instance.Name

		if !dryRun {
			_, err = client.StreamsApi.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
			if err != nil && admin.IsErrorCode(err, "STREAM_TENANT_HAS_STREAM_PROCESSORS") {
				streamProcessors, _, spErr := client.StreamsApi.GetStreamProcessors(ctx, projectID, instanceName).Execute()
				if spErr != nil {
					return 0, spErr
				}

				for _, processor := range *streamProcessors.Results {
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
	return len(*streamInstances.Results), nil
}
