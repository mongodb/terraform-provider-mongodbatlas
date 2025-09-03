package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	v2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

func GetClusterDetails(ctx context.Context, client *config.MongoDBClient, projectID, clusterName string) (cluster *admin.ClusterDescription20240805, flexCluster *admin.FlexClusterDescription20241113, diags v2diag.Diagnostics) {
	fwdiags := fwdiag.Diagnostics{}
	cluster, flexCluster = advancedclustertpf.GetClusterDetails(ctx, &fwdiags, projectID, clusterName, client, false)
	return cluster, flexCluster, conversion.FromTPFDiagsToSDKV2Diags(fwdiags)
}
