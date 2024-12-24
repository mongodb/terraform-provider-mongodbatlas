package streamconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func DeleteStreamConnection(ctx context.Context, api admin.StreamsApi, projectID, instanceName, connectionName string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, resp, err := api.DeleteStreamConnection(ctx, projectID, instanceName, connectionName).Execute()
		if err == nil {
			return nil
		}
		if admin.IsErrorCode(err, "STREAM_KAFKA_CONNECTION_IS_DEPLOYING") {
			return retry.RetryableError(err)
		}
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}
