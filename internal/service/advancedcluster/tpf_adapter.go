package advancedcluster

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func AwaitChanges(ctx context.Context, api *admin.APIClient, projectID, clusterName, lastOperation string, timeoutDuration time.Duration) sdkv2diag.Diagnostics {
	diags := &diag.Diagnostics{}
	_ = advancedclustertpf.AwaitChanges(ctx, api.ClustersApi, projectID, clusterName, lastOperation, timeoutDuration, diags)
	return conversion.FromTPFDiagsToSDKV2Diags(*diags)
}
