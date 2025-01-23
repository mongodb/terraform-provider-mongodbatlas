package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	tpf "github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func AwaitChanges(ctx context.Context, client *config.MongoDBClient, waitParams *tpf.ClusterWaitParams, lastOperation string) sdkv2diag.Diagnostics {
	diags := new(diag.Diagnostics)
	_ = tpf.AwaitChanges(ctx, client, waitParams, lastOperation, diags)
	return conversion.FromTPFDiagsToSDKV2Diags(*diags)
}

func CreateCluster(ctx context.Context, diags *sdkv2diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *tpf.ClusterWaitParams, usingOldShardingConfiguration bool) *admin.ClusterDescription20240805 {
	diagsTPF := new(diag.Diagnostics)
	cluster := tpf.CreateCluster(ctx, diagsTPF, client, req, waitParams, usingOldShardingConfiguration)
	localDiags := conversion.FromTPFDiagsToSDKV2Diags(*diagsTPF)
	*diags = append(*diags, localDiags...)
	return cluster
}

func UpdateAdvancedConfiguration(ctx context.Context, diags *sdkv2diag.Diagnostics, client *config.MongoDBClient, reqLegacy *admin20240530.ClusterDescriptionProcessArgs, reqNew *admin.ClusterDescriptionProcessArgs20240805, waitParams *tpf.ClusterWaitParams) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	diagsTPF := new(diag.Diagnostics)
	legacy, latest, changed = tpf.UpdateAdvancedConfiguration(ctx, diagsTPF, client, reqLegacy, reqNew, waitParams)
	localDiags := conversion.FromTPFDiagsToSDKV2Diags(*diagsTPF)
	*diags = append(*diags, localDiags...)
	return legacy, latest, changed
}
