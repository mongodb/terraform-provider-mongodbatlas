package advancedcluster

import (
	"context"
	"fmt"
	"net/http"

	// "go.mongodb.org/atlas-sdk/v20250312009/admin" TODO: don't use normal SDK while hidden tls1.3 field
	"github.com/mongodb/atlas-sdk-go/admin" // TODO: change to SDK before merging to master

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

func CreateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	var (
		pauseAfter  = req.GetPaused()
		clusterResp *admin.ClusterDescription20240805
	)
	if pauseAfter {
		req.Paused = nil
	}
	clusterResp = createClusterLatest(ctx, diags, client, req, waitParams)
	if diags.HasError() {
		return nil
	}
	if pauseAfter {
		clusterResp = updateCluster(ctx, diags, client, &pauseRequest, waitParams, operationPauseAfterCreate)
	}
	return clusterResp
}

func createClusterLatest(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasPreview.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationCreate, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationCreate, diags)
}

func updateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams, operationName string) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasPreview.ClustersApi.UpdateCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationName, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationName, diags)
}

// ProcessArgs.ClusterAdvancedConfig is managed through create/updateCluster APIs instead of /processArgs APIs but since corresponding TF attributes
// belong in the advanced_configuration attribute we still need to check for any changes
func UpdateAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, p *ProcessArgs, waitParams *ClusterWaitParams) (latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	var (
		err         error
		advConfig   *admin.ClusterDescriptionProcessArgs20240805
		projectID   = waitParams.ProjectID
		clusterName = waitParams.ClusterName
	)
	if !update.IsZeroValues(p.ArgsDefault) {
		changed = true
		advConfig, _, err = client.AtlasPreview.ClustersApi.UpdateProcessArgs(ctx, projectID, clusterName, p.ArgsDefault).Execute()
		if err != nil {
			addErrorDiag(diags, operationAdvancedConfigurationUpdate, defaultAPIErrorDetails(clusterName, err))
			return nil, false
		}
		_ = AwaitChanges(ctx, client, waitParams, operationAdvancedConfigurationUpdate, diags)
		if diags.HasError() {
			return nil, false
		}
	}
	if !update.IsZeroValues(p.ClusterAdvancedConfig) {
		changed = true
	}
	return advConfig, changed
}

func ReadIfUnsetAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, configNew *admin.ClusterDescriptionProcessArgs20240805) (latest *admin.ClusterDescriptionProcessArgs20240805) {
	var err error
	if configNew == nil {
		configNew, _, err = client.AtlasPreview.ClustersApi.GetProcessArgs(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	return configNew
}

func UpgradeTenant(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasPreview.ClustersApi.UpgradeTenantUpgrade(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationTenantUpgrade, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChangesUpgrade(ctx, client, waitParams, operationTenantUpgrade, diags)
}

func UpgradeFlexToDedicated(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.AtlasTenantClusterUpgradeRequest20240805) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasPreview.FlexClustersApi.TenantUpgrade(ctx, waitParams.ProjectID, req).Execute()
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
	_, err := client.AtlasPreview.ClustersApi.DeleteClusterWithParams(ctx, params).Execute()
	if err != nil {
		if !admin.IsErrorCode(err, "CANNOT_USE_FLEX_CLUSTER_IN_CLUSTER_API") {
			addErrorDiag(diags, operationDelete, defaultAPIErrorDetails(waitParams.ClusterName, err))
			return
		}
		err := flexcluster.DeleteFlexCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, client.AtlasPreview.FlexClustersApi, waitParams.Timeout)
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
			cleanResp, cleanErr = client.AtlasPreview.FlexClustersApi.DeleteFlexCluster(ctx, projectID, clusterName).Execute()
		} else {
			cleanResp, cleanErr = client.AtlasPreview.ClustersApi.DeleteCluster(ctx, projectID, clusterName).Execute()
		}
		if validate.StatusNotFound(cleanResp) {
			return nil
		}
		return cleanErr
	}
}

func GetClusterDetails(ctx context.Context, diags *diag.Diagnostics, projectID, clusterName string, client *config.MongoDBClient, fcvPresentInState bool) (clusterDesc *admin.ClusterDescription20240805, flexClusterResp *admin.FlexClusterDescription20241113) {
	isFlex := false
	clusterDesc, resp, err := client.AtlasPreview.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
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
		flexClusterResp, err = flexcluster.GetFlexCluster(ctx, projectID, clusterName, client.AtlasPreview.FlexClustersApi)
		if err != nil {
			diags.AddError(fmt.Sprintf(flexcluster.ErrorReadFlex, clusterName, err), defaultAPIErrorDetails(clusterName, err))
			return nil, nil
		}
	}
	return clusterDesc, flexClusterResp
}
