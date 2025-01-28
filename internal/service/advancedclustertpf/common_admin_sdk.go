package advancedclustertpf

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2Diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func CreateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams, usingOldShardingConfiguration bool) *admin.ClusterDescription20240805 {
	var (
		pauseAfter  = req.GetPaused()
		clusterResp *admin.ClusterDescription20240805
	)
	if pauseAfter {
		req.Paused = nil
	}
	if usingOldShardingConfiguration {
		legacyReq := ConvertClusterDescription20241023to20240805(req)
		clusterResp = createCluster20240805(ctx, diags, client, legacyReq, waitParams)
	} else {
		clusterResp = createClusterLatest(ctx, diags, client, req, waitParams)
	}
	if diags.HasError() {
		return nil
	}
	if pauseAfter {
		clusterResp = updateCluster(ctx, diags, client, &pauseRequest, waitParams, operationPauseAfterCreate)
	}
	return clusterResp
}

func createCluster20240805(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin20240805.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV220240805.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationCreate20240805, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationCreate20240805, diags)
}

func createClusterLatest(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationCreate, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationCreate, diags)
}

func updateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams, operationName string) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpdateCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationName, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationName, diags)
}

func UpdateAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, reqLegacy *admin20240530.ClusterDescriptionProcessArgs, reqNew *admin.ClusterDescriptionProcessArgs20240805, waitParams *ClusterWaitParams) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	var (
		err             error
		advConfig       *admin.ClusterDescriptionProcessArgs20240805
		legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
		projectID       = waitParams.ProjectID
		clusterName     = waitParams.ClusterName
	)
	if !update.IsZeroValues(reqNew) {
		changed = true
		advConfig, _, err = client.AtlasV2.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, reqNew).Execute()
		if err != nil {
			addErrorDiag(diags, operationAdvancedConfigurationUpdate, defaultAPIErrorDetails(clusterName, err))
			return nil, nil, false
		}
		_ = AwaitChanges(ctx, client, waitParams, operationAdvancedConfigurationUpdate, diags)
		if diags.HasError() {
			return nil, nil, false
		}
	}
	if !update.IsZeroValues(reqLegacy) {
		changed = true
		legacyAdvConfig, _, err = client.AtlasV220240530.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, reqLegacy).Execute()
		if err != nil {
			addErrorDiag(diags, operationAdvancedConfigurationUpdate20240530, defaultAPIErrorDetails(clusterName, err))
			diags.AddError(errorAdvancedConfUpdateLegacy, defaultAPIErrorDetails(clusterName, err))
			return nil, nil, false
		}
		_ = AwaitChanges(ctx, client, waitParams, operationAdvancedConfigurationUpdate20240530, diags)
		if diags.HasError() {
			return nil, nil, false
		}
	}
	return legacyAdvConfig, advConfig, changed
}

func ReadIfUnsetAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, configLegacy *admin20240530.ClusterDescriptionProcessArgs, configNew *admin.ClusterDescriptionProcessArgs20240805) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805) {
	var err error
	if configLegacy == nil {
		configLegacy, _, err = client.AtlasV220240530.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfReadLegacy, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	if configNew == nil {
		configNew, _, err = client.AtlasV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	return configLegacy, configNew
}

func TenantUpgrade(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpgradeSharedCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationTenantUpgrade, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationTenantUpgrade, diags)
}

func PinFCV(ctx context.Context, api admin.ClustersApi, projectID, clusterName, expirationDateStr string) error {
	expirationTime, ok := conversion.StringToTime(expirationDateStr)
	if !ok {
		return fmt.Errorf("expiration_date format is incorrect: %s", expirationDateStr)
	}
	req := admin.PinFCV{
		ExpirationDate: &expirationTime,
	}
	if _, _, err := api.PinFeatureCompatibilityVersion(ctx, projectID, clusterName, &req).Execute(); err != nil {
		return err
	}
	return nil
}

func DeleteCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, retainBackups *bool) {
	params := &admin.DeleteClusterApiParams{
		GroupId:       waitParams.ProjectID,
		ClusterName:   waitParams.ClusterName,
		RetainBackups: retainBackups,
	}
	_, err := client.AtlasV2.ClustersApi.DeleteClusterWithParams(ctx, params).Execute()
	if err != nil {
		addErrorDiag(diags, operationDelete, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return
	}
	AwaitChanges(ctx, client, waitParams, operationDelete, diags)
}

func ReadCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, fcvPresentInState bool) *admin.ClusterDescription20240805 {
	readResp, _, err := client.AtlasV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			return nil
		}
		diags.AddError(errorReadResource, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	if fcvPresentInState {
		newWarnings := GenerateFCVPinningWarningForRead(fcvPresentInState, readResp.FeatureCompatibilityVersionExpirationDate)
		diags.Append(newWarnings...)
	}
	return readResp
}

func GetClusterDetails(ctx context.Context, d *schema.ResourceData, projectID, clusterName string, client *admin.APIClient) (isFlex bool, clusterDesc *admin.ClusterDescription20240805, flexClusterResp *admin.FlexClusterDescription20241113, diags sdkv2Diag.Diagnostics, err error) {
	clusterDesc, resp, err := client.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil, nil, nil, nil
		}
		if isFlex = admin.IsErrorCode(err, "CANNOT_USE_FLEX_CLUSTER_IN_CLUSTER_API"); !isFlex {
			return false, nil, nil, sdkv2Diag.FromErr(fmt.Errorf("error reading  advanced cluster (%s): %s", clusterName, err)), err
		}
	}
	if isFlex {
		clusterName := d.Get("name").(string)
		flexClusterResp, err = flexcluster.GetFlexCluster(ctx, projectID, clusterName, client.FlexClustersApi)
		if err != nil {
			return true, nil, nil, sdkv2Diag.FromErr(fmt.Errorf(flexcluster.ErrorReadFlex, clusterName, err)), err
		}
	}
	return isFlex, clusterDesc, flexClusterResp, nil, nil
}
