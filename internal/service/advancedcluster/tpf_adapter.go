package advancedcluster

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

func AwaitChanges(ctx context.Context, client *config.MongoDBClient, projectID, clusterName, lastOperation string, timeoutDuration time.Duration) sdkv2diag.Diagnostics {
	diags := &diag.Diagnostics{}
	ids := &advancedclustertpf.ClusterReader{
		ProjectID:   projectID,
		ClusterName: clusterName,
		Timeout:     timeoutDuration,
	}
	_ = advancedclustertpf.AwaitChanges(ctx, client, ids, lastOperation, diags)
	return conversion.FromTPFDiagsToSDKV2Diags(*diags)
}
