package advancedcluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

// ProcessArgs.ClusterAdvancedConfig is managed through create/updateCluster APIs instead of /processArgs APIs but since corresponding TF attributes
// belong in the advanced_configuration attribute we still need to check for any changes
func updateAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, p *ProcessArgs, waitParams *ClusterWaitParams) (latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	var (
		err         error
		advConfig   *admin.ClusterDescriptionProcessArgs20240805
		projectID   = waitParams.ProjectID
		clusterName = waitParams.ClusterName
	)
	if !isEmptyProcessArgs(p.ArgsDefault) {
		// Read current API processArgs and recompute the diff against them.
		// This avoids unnecessary PATCH calls when the API already has the desired values,
		// which happens after a moved block migration where the TF state has all-null fields
		// but the API already has the correct values from the source resource.
		p.ArgsDefault = recalculatePatchProcessArgs(ctx, diags, client, projectID, clusterName, p.ArgsDefault)
		if diags.HasError() {
			return nil, false
		}
	}
	if !isEmptyProcessArgs(p.ArgsDefault) {
		changed = true
		advConfig, _, err = client.AtlasV2.ClustersApi.UpdateProcessArgs(ctx, projectID, clusterName, p.ArgsDefault).Execute()
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

// isEmptyProcessArgs checks if the processArgs request would produce an empty JSON body ("{}").
// This is used instead of update.IsZeroValues to decide whether to call UpdateProcessArgs, because
// IsZeroValues uses reflect.DeepEqual which can return false for structs that still serialize to "{}".
// This happens because Go's json omitempty omits non-nil pointers to empty values (e.g., *[]string
// pointing to []string{} is omitted), while reflect.DeepEqual treats &[]string{} as different from nil.
// Without this check, a PATCH with an empty body "{}" would be sent to the processArgs API unnecessarily.
func isEmptyProcessArgs(req *admin.ClusterDescriptionProcessArgs20240805) bool {
	if req == nil {
		return true
	}
	b, err := json.Marshal(req)
	if err != nil {
		return false
	}
	return string(b) == "{}"
}

// recalculatePatchProcessArgs reads the current API processArgs and recomputes the diff against the proposed patch.
// This ensures we only send a PATCH when there's an actual change. PatchPayload only considers Add and Replace
// operations as changes (not Remove), so comparing the full API response against a partial patch request is safe —
// fields not in the patch (nil → absent from JSON → Remove operation) are ignored.
func recalculatePatchProcessArgs(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, patchReq *admin.ClusterDescriptionProcessArgs20240805) *admin.ClusterDescriptionProcessArgs20240805 {
	currentArgs, _, err := client.AtlasV2.ClustersApi.GetProcessArgs(ctx, projectID, clusterName).Execute()
	if err != nil {
		diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	recalculated, err := update.PatchPayload(currentArgs, patchReq)
	if err != nil {
		diags.AddError("error recalculating processArgs patch", err.Error())
		return nil
	}
	return recalculated
}

// recalculateClusterPatch reads the current cluster from the API and recomputes the diff against the proposed patch.
// This avoids unnecessary PATCH calls when the API already has the desired values, which happens after a state
// upgrade (v1→v3) where the TF state has null values for Optional-only attributes (e.g., backup_enabled) but
// the API already has the correct values. The state/plan diff produces false changes because the state was
// overridden to null by overrideAttributesWithPrevStateValue.
// Similarly, zone_name defaults to "ZoneName managed by Terraform" in the plan when the user doesn't configure it,
// which differs from the API's actual zone_name, producing a false replicationSpecs change.
func recalculateClusterPatch(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, planReq *admin.ClusterDescription20240805) *admin.ClusterDescription20240805 {
	currentCluster, _, err := client.AtlasV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		diags.AddError(errorReadResource, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	patchOptions := update.PatchOptions{
		IgnoreInStatePrefix: []string{"replicationSpecs"},
	}
	recalculated, err := update.PatchPayload(currentCluster, planReq, patchOptions)
	if err != nil {
		diags.AddError("error recalculating cluster patch", err.Error())
		return nil
	}
	return recalculated
}

func readIfUnsetAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, configNew *admin.ClusterDescriptionProcessArgs20240805) (latest *admin.ClusterDescriptionProcessArgs20240805) {
	var err error
	if configNew == nil {
		configNew, _, err = client.AtlasV2.ClustersApi.GetProcessArgs(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	return configNew
}

func upgradeTenant(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpgradeTenantUpgrade(ctx, waitParams.ProjectID, req).Execute()
	if err != nil {
		addErrorDiag(diags, operationTenantUpgrade, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChangesUpgrade(ctx, client, waitParams, operationTenantUpgrade, diags)
}

func upgradeFlexToDedicated(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.AtlasTenantClusterUpgradeRequest20240805) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.FlexClustersApi.TenantUpgrade(ctx, waitParams.ProjectID, req).Execute()
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

func deleteCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, retainBackups *bool) {
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
	_ = AwaitChanges(ctx, client, waitParams, operationDelete, diags)
}

func deleteClusterNoWait(client *config.MongoDBClient, projectID, clusterName string, isFlex bool) func(ctx context.Context) error {
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
	clusterDesc, resp, err := client.AtlasV2.ClustersApi.GetCluster(ctx, projectID, clusterName).UseEffectiveFieldsReplicationSpecs(true).Execute()
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
