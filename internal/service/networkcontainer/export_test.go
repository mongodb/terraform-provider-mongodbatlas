package networkcontainer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
)

// NewRefreshFuncForTest exposes resourceRefreshFunc to package networkcontainer_test.
func NewRefreshFuncForTest(ctx context.Context, d *schema.ResourceData, client *admin.APIClient) retry.StateRefreshFunc {
	return resourceRefreshFunc(ctx, d, client)
}
