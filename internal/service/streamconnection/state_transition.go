package streamconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
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
		if validate.StatusNotFound(resp) {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}
