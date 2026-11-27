package networkcontainer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312029/admin"
)

// NewDeleteRefreshFuncForTest exposes resourceDeleteRefreshFunc to package networkcontainer_test.
func NewDeleteRefreshFuncForTest(ctx context.Context, d *schema.ResourceData, client *admin.APIClient) retry.StateRefreshFunc {
	return resourceDeleteRefreshFunc(ctx, d, client)
}
