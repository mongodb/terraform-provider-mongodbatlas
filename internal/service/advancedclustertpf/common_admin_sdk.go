package advancedclustertpf

import (
	"context"
	"fmt"
	"net/http"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

func CreateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams, usingNewShardingConfig bool) *admin.ClusterDescription20240805 {
	var (
		pauseAfter  = req.GetPaused()
		clusterResp *admin.ClusterDescription20240805
	)
	if pauseAfter {
		req.Paused = nil
	}
	if usingNewShardingConfig {
		clusterResp = createClusterLatest(ctx, diags, client, req, waitParams)
	} else {
		oldReq := ConvertClusterDescription20241023to20240805(req)
		clusterResp = createCluster20240805(ctx, diags, client, oldReq, waitParams)
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

// ProcessArgs.ClusterAdvancedConfig is managed through create/updateCluster APIs instead of /processArgs APIs but since corresponding TF attributes
// belong in the advanced_configuration attribute we still need to check for any changes
func UpdateAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, p *ProcessArgs, waitParams *ClusterWaitParams) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	var (
		err             error
		advConfig       *admin.ClusterDescriptionProcessArgs20240805
		legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
		projectID       = waitParams.ProjectID
		clusterName     = waitParams.ClusterName
	)
	if !update.IsZeroValues(p.ArgsDefault) {
		changed = true
		advConfig, _, err = client.AtlasV2.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, p.ArgsDefault).Execute()
		if err != nil {
			addErrorDiag(diags, operationAdvancedConfigurationUpdate, defaultAPIErrorDetails(clusterName, err))
			return nil, nil, false
		}
		_ = AwaitChanges(ctx, client, waitParams, operationAdvancedConfigurationUpdate, diags)
		if diags.HasError() {
			return nil, nil, false
		}
	}
	if !update.IsZeroValues(p.ClusterAdvancedConfig) {
		changed = true
	}
	return legacyAdvConfig, advConfig, changed
}

func ReadIfUnsetAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, configNew *admin.ClusterDescriptionProcessArgs20240805) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805) {
	var err error
	if configNew == nil {
		configNew, _, err = client.AtlasV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	return nil, configNew
}

func UpgradeTenant(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpgradeSharedCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationTenantUpgrade, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChangesUpgrade(ctx, client, waitParams, operationTenantUpgrade, diags)
}

func UpgradeFlexToDedicated(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.AtlasTenantClusterUpgradeRequest20240805) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.FlexClustersApi.UpgradeFlexCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationFlexUpgrade, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChangesUpgrade(ctx, client, waitParams, operationFlexUpgrade, diags)
}

func PinFCV(ctx context.Context, api admin.ClustersApi, projectID, clusterName, expirationDateStr string) error {
	expirationTime, ok := conversion.StringToTime(expirationDateStr)
	if !ok {
		return fmt.Errorf("expiration_date format is incorrect: %s", expirationDateStr)
	}
	req := admin.PinFCV{
		ExpirationDate: &expirationTime,
	}
	if _, err := api.PinFeatureCompatibilityVersion(ctx, projectID, clusterName, &req).Execute(); err != nil {
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
		if !admin.IsErrorCode(err, "CANNOT_USE_FLEX_CLUSTER_IN_CLUSTER_API") {
			addErrorDiag(diags, operationDelete, defaultAPIErrorDetails(waitParams.ClusterName, err))
			return
		}
		err := flexcluster.DeleteFlexCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, client.AtlasV2.FlexClustersApi, waitParams.Timeout)
		if err != nil {
			addErrorDiag(diags, operationDeleteFlex, defaultAPIErrorDetails(waitParams.ClusterName, err))
			return
		}
	}
	AwaitChanges(ctx, client, waitParams, operationDelete, diags)
}

func DeleteClusterNoWait(client *config.MongoDBClient, projectID, clusterName string, isFlex bool) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var cleanResp *http.Response
		var cleanErr error
		if isFlex {
			cleanResp, cleanErr = client.AtlasV2.FlexClustersApi.DeleteFlexCluster(ctx, projectID, clusterName).Execute()
		} else {
			cleanResp, cleanErr = client.AtlasV2.ClustersApi.DeleteCluster(ctx, projectID, clusterName).Execute()
		}
		if validate.StatusNotFound(cleanResp) {
			return nil
		}
		return cleanErr
	}
}

func GetClusterDetails(ctx context.Context, diags *diag.Diagnostics, projectID, clusterName string, client *config.MongoDBClient, fcvPresentInState bool) (clusterDesc *admin.ClusterDescription20240805, flexClusterResp *admin.FlexClusterDescription20241113) {
	isFlex := false
	clusterDesc, resp, err := client.AtlasV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) || admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			return nil, nil
		}
		if isFlex = admin.IsErrorCode(err, "CANNOT_USE_FLEX_CLUSTER_IN_CLUSTER_API"); !isFlex {
			diags.AddError(errorReadResource, defaultAPIErrorDetails(clusterName, err))
			return nil, nil
		}
	}

	if !isFlex && fcvPresentInState && clusterDesc != nil {
		newWarnings := GenerateFCVPinningWarningForRead(fcvPresentInState, clusterDesc.FeatureCompatibilityVersionExpirationDate)
		diags.Append(newWarnings...)
	}

	if isFlex {
		flexClusterResp, err = flexcluster.GetFlexCluster(ctx, projectID, clusterName, client.AtlasV2.FlexClustersApi)
		if err != nil {
			diags.AddError(fmt.Sprintf(flexcluster.ErrorReadFlex, clusterName, err), defaultAPIErrorDetails(clusterName, err))
			return nil, nil
		}
	}
	return clusterDesc, flexClusterResp
}
