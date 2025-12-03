package advancedcluster

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithMoveState = &rs{}
var _ resource.ResourceWithUpgradeState = &rs{}
var _ resource.ResourceWithModifyPlan = &rs{}

const (
	resourceName             = "advanced_cluster"
	errorPatchPayload        = "error creating patch payload"
	errorDetailDefault       = "cluster name: %s, API error details: %s"
	errorReadResource        = "error reading advanced cluster"
	errorAdvancedConfRead    = "error reading Advanced Configuration"
	errorList                = "error reading advanced cluster list"
	errorListDetail          = "project ID %s. Error %s"
	errorResolveContainerIDs = "error resolving container IDs"
	errorRegionPriorities    = "priority values in region_configs must be in descending order"

	ErrorCodeClusterNotFound             = "CLUSTER_NOT_FOUND"
	operationUpdate                      = "update"
	operationCreate                      = "create"
	operationPauseAfterCreate            = "pause after create"
	operationDelete                      = "delete"
	operationDeleteFlex                  = "flex delete"
	operationAdvancedConfigurationUpdate = "update advanced configuration"
	operationTenantUpgrade               = "tenant upgrade"
	operationPauseAfterUpdate            = "pause after update"
	operationResumeBeforeUpdate          = "resume before update"
	operationFCVPinning                  = "FCV pinning"
	operationFCVUnpinning                = "FCV unpinning"
	operationFlexUpgrade                 = "flex upgrade"
)

func addErrorDiag(diags *diag.Diagnostics, errorLocator, details string) {
	diags.AddError("Error in "+errorLocator, details)
}

func defaultAPIErrorDetails(clusterName string, err error) string {
	return fmt.Sprintf(errorDetailDefault, clusterName, err.Error())
}

var (
	resumeRequest = admin.ClusterDescription20240805{Paused: conversion.Pointer(false)}
	pauseRequest  = admin.ClusterDescription20240805{Paused: conversion.Pointer(true)}
)

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type rs struct {
	config.RSCommon
}

// ModifyPlan is called before plan is shown to the user and right before the plan is applied.
// Why do we need this? Why can't we use planmodifier.UseStateForUnknown in different fields?
// 1. UseStateForUnknown always copies the state for unknown values. However, that leads to `Error: Provider produced inconsistent result after apply` in some cases (see implementation below).
// 2. Adding the different UseStateForUnknown is very verbose.
func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.Plan.Raw.IsFullyKnown() { // Return early unless it is an Update
		return
	}
	var plan, state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	// The replication specs can be unknown if the cluster depends on another resource.
	// handleModifyPlan will try to convert the field to `Target Type: []advancedcluster.TFReplicationSpecsModel`.
	// But since the field is unknown the user gets an error: `Error: Value Conversion Error`.
	if plan.ReplicationSpecs.IsUnknown() {
		return
	}

	handleModifyPlan(ctx, diags, &state, &plan)
	if diags.HasError() {
		return
	}
	diags.Append(resp.Plan.Set(ctx, plan)...)
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}
	latestReq := newAtlasReq(ctx, &plan, diags)
	if diags.HasError() {
		return
	}
	waitParams := resolveClusterWaitParams(ctx, &plan, diags, operationCreate)
	if diags.HasError() {
		return
	}
	isFlex := isFlex(latestReq.ReplicationSpecs)
	projectID, clusterName := waitParams.ProjectID, waitParams.ClusterName
	clusterDetailStr := fmt.Sprintf("Cluster name %s (project_id=%s).", clusterName, projectID)
	if plan.DeleteOnCreateTimeout.ValueBool() {
		var deferCall func()
		ctx, deferCall = cleanup.OnTimeout(
			ctx, waitParams.Timeout, diags.AddWarning, clusterDetailStr, deleteClusterNoWait(r.Client, projectID, clusterName, isFlex),
		)
		defer deferCall()
	}
	if isFlex {
		flexClusterReq := newFlexCreateReq(latestReq.GetName(), latestReq.GetTerminationProtectionEnabled(), latestReq.Tags, latestReq.ReplicationSpecs)
		flexClusterResp, err := flexcluster.CreateFlexCluster(ctx, plan.ProjectID.ValueString(), latestReq.GetName(), flexClusterReq, r.Client.AtlasV2.FlexClustersApi, &waitParams.Timeout)
		if err != nil {
			diags.AddError(fmt.Sprintf(flexcluster.ErrorCreateFlex, clusterDetailStr), err.Error())
			return
		}
		newFlexClusterModel := newTFModelFlex(ctx, diags, flexClusterResp, getPriorityOfFlexReplicationSpecs(latestReq.ReplicationSpecs), &plan)
		if diags.HasError() {
			return
		}
		diags.Append(resp.State.Set(ctx, newFlexClusterModel)...)
		return
	}
	clusterResp := createCluster(ctx, diags, r.Client, latestReq, waitParams)

	emptyAdvancedConfiguration := types.ObjectNull(advancedConfigurationObjType.AttrTypes)
	patchReqProcessArgs := update.PatchPayloadCluster(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if diags.HasError() {
		return
	}
	p := &ProcessArgs{
		ArgsDefault:           patchReqProcessArgs,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	}
	advConfig, _ := updateAdvancedConfiguration(ctx, diags, r.Client, p, waitParams)
	if diags.HasError() {
		return
	}
	if changedCluster := r.applyPinnedFCVChanges(ctx, diags, &TFModel{}, &plan, waitParams); changedCluster != nil {
		clusterResp = changedCluster
	}
	if diags.HasError() {
		return
	}

	modelOut := getBasicClusterModel(ctx, diags, r.Client, clusterResp, &plan)
	if diags.HasError() {
		return
	}
	advConfig = readIfUnsetAdvancedConfiguration(ctx, diags, r.Client, waitParams.ProjectID, waitParams.ClusterName, advConfig)

	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, &ProcessArgs{
		ArgsDefault:           advConfig,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	})
	if diags.HasError() {
		return
	}
	diags.Append(resp.State.Set(ctx, modelOut)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	clusterName := state.Name.ValueString()
	projectID := state.ProjectID.ValueString()
	cluster, flexCluster := GetClusterDetails(ctx, diags, projectID, clusterName, r.Client, !state.PinnedFCV.IsNull(), state.UseEffectiveFields.ValueBool())
	if diags.HasError() {
		return
	}
	if cluster == nil && flexCluster == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if flexCluster != nil {
		newFlexClusterModel := newTFModelFlex(ctx, diags, flexCluster, getPriorityOfFlexReplicationSpecs(newAtlasReq(ctx, &state, diags).ReplicationSpecs), &state)
		if diags.HasError() {
			return
		}
		diags.Append(resp.State.Set(ctx, newFlexClusterModel)...)
		return
	}
	modelOut := getBasicClusterModel(ctx, diags, r.Client, cluster, &state)
	if diags.HasError() {
		return
	}
	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, &ProcessArgs{
		ArgsDefault:           nil,
		ClusterAdvancedConfig: cluster.AdvancedConfiguration,
	})
	if diags.HasError() {
		return
	}
	diags.Append(resp.State.Set(ctx, modelOut)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	waitParams := resolveClusterWaitParams(ctx, &plan, diags, operationUpdate)
	if diags.HasError() {
		return
	}

	// FCV update is intentionally handled before any other cluster updates, and will wait for cluster to reach IDLE state before continuing
	clusterResp := r.applyPinnedFCVChanges(ctx, diags, &state, &plan, waitParams)
	if diags.HasError() {
		return
	}

	{
		diff := findClusterDiff(ctx, &state, &plan, diags)
		if diags.HasError() {
			return
		}
		switch {
		case diff.isUpgradeTenantToFlex:
			if flexOut := handleFlexUpgrade(ctx, diags, r.Client, waitParams, &plan); flexOut != nil {
				diags.Append(resp.State.Set(ctx, flexOut)...)
			}
			return
		case diff.isUpdateOfFlex:
			if flexOut := handleFlexUpdate(ctx, diags, r.Client, waitParams, &plan); flexOut != nil {
				diags.Append(resp.State.Set(ctx, flexOut)...)
			}
			return
		case diff.isUpgradeFlexToDedicated():
			clusterResp = upgradeFlexToDedicated(ctx, diags, r.Client, waitParams, diff.upgradeFlexToDedicatedReq)
		case diff.isUpgradeTenant():
			clusterResp = upgradeTenant(ctx, diags, r.Client, waitParams, diff.upgradeTenantReq)
		case diff.isClusterPatchOnly():
			clusterResp = r.applyClusterChanges(ctx, diags, diff.clusterPatchOnlyReq, waitParams)
		}
		if diags.HasError() {
			return
		}
	}
	// clusterResp can be nil if there are no changes to the cluster, for example when `delete_on_create_timeout` is changed or only advanced configuration is changed
	if clusterResp == nil {
		var flexResp *admin.FlexClusterDescription20241113
		clusterResp, flexResp = GetClusterDetails(ctx, diags, waitParams.ProjectID, waitParams.ClusterName, r.Client, false, waitParams.UseEffectiveFields)
		if diags.HasError() {
			return
		}
		// This should never happen since the switch case should handle the two flex cases (update/upgrade) and return, but keeping it here for safety.
		if flexResp != nil {
			flexPriority := getPriorityOfFlexReplicationSpecs(newAtlasReq(ctx, &plan, diags).ReplicationSpecs)
			if flexOut := newTFModelFlex(ctx, diags, flexResp, flexPriority, &plan); flexOut != nil {
				diags.Append(resp.State.Set(ctx, flexOut)...)
			}
			return
		}
	}
	modelOut := getBasicClusterModel(ctx, diags, r.Client, clusterResp, &plan)
	if diags.HasError() {
		return
	}
	patchReqProcessArgs := update.PatchPayloadCluster(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if diags.HasError() {
		return
	}
	p := &ProcessArgs{
		ArgsDefault:           patchReqProcessArgs,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	}
	advConfig, advConfigChanged := updateAdvancedConfiguration(ctx, diags, r.Client, p, waitParams)
	if diags.HasError() {
		return
	}
	if advConfigChanged {
		updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, &ProcessArgs{
			ArgsDefault:           advConfig,
			ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
		})
		if diags.HasError() {
			return
		}
	} else {
		modelOut.AdvancedConfiguration = state.AdvancedConfiguration
	}
	diags.Append(resp.State.Set(ctx, modelOut)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	waitParams := resolveClusterWaitParams(ctx, &state, diags, operationDelete)
	if diags.HasError() {
		return
	}
	retainBackups := conversion.NilForUnknown(state.RetainBackupsEnabled, state.RetainBackupsEnabled.ValueBoolPointer())
	deleteCluster(ctx, diags, r.Client, waitParams, retainBackups)
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}

func (r *rs) applyPinnedFCVChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	var (
		api         = r.Client.AtlasV2.ClustersApi
		projectID   = waitParams.ProjectID
		clusterName = waitParams.ClusterName
	)
	if state.PinnedFCV.Equal(plan.PinnedFCV) {
		return nil
	}
	isFCVPresentInConfig := !plan.PinnedFCV.IsNull()
	if isFCVPresentInConfig {
		fcvModel := &TFPinnedFCVModel{}
		// pinned_fcv has been defined or updated expiration date
		if localDiags := plan.PinnedFCV.As(ctx, fcvModel, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
			diags.Append(localDiags...)
			return nil
		}
		if err := PinFCV(ctx, api, projectID, clusterName, fcvModel.ExpirationDate.ValueString()); err != nil {
			addErrorDiag(diags, operationFCVPinning, defaultAPIErrorDetails(clusterName, err))
			return nil
		}
		return AwaitChanges(ctx, r.Client, waitParams, operationFCVPinning, diags)
	}
	// pinned_fcv has been removed from the config so unpin method is called
	if _, err := api.UnpinFeatureCompatibilityVersion(ctx, projectID, clusterName).Execute(); err != nil {
		addErrorDiag(diags, operationFCVUnpinning, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	return AwaitChanges(ctx, r.Client, waitParams, operationFCVUnpinning, diags)
}

func (r *rs) applyClusterChanges(ctx context.Context, diags *diag.Diagnostics, patchReq *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	// paused = `false` is sent in an isolated request before other changes to avoid error from API: Cannot update cluster while it is paused or being paused.
	var result *admin.ClusterDescription20240805
	if patchReq.Paused != nil && !patchReq.GetPaused() {
		patchReq.Paused = nil
		_ = updateCluster(ctx, diags, r.Client, &resumeRequest, waitParams, operationResumeBeforeUpdate)
	}

	// paused = `true` is sent in an isolated request after other changes have been applied to avoid error from API: Cannot update and pause cluster at the same time
	var pauseAfterOtherChanges = false
	if patchReq.Paused != nil && patchReq.GetPaused() {
		patchReq.Paused = nil
		pauseAfterOtherChanges = true
	}

	result = updateCluster(ctx, diags, r.Client, patchReq, waitParams, operationUpdate)

	if pauseAfterOtherChanges {
		result = updateCluster(ctx, diags, r.Client, &pauseRequest, waitParams, operationPauseAfterUpdate)
	}
	return result
}

func getBasicClusterModel(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805, modelIn *TFModel) *TFModel {
	containerIDs := resolveContainerIDsOrError(ctx, diags, clusterResp, client.AtlasV2.NetworkPeeringApi)
	if diags.HasError() {
		return nil
	}
	modelOut := newTFModel(ctx, clusterResp, diags, containerIDs)
	if diags.HasError() {
		return nil
	}
	overrideAttributesWithPrevStateValue(modelIn, modelOut)
	return modelOut
}

func resolveContainerIDsOrError(ctx context.Context, diags *diag.Diagnostics, clusterResp *admin.ClusterDescription20240805, api admin.NetworkPeeringApi) map[string]string {
	projectID := clusterResp.GetGroupId()
	clusterName := clusterResp.GetName()
	containerIDs, err := resolveContainerIDs(ctx, projectID, clusterResp, api)
	if err != nil {
		diags.AddError(errorResolveContainerIDs, fmt.Sprintf("cluster name = %s, error details: %s", clusterName, err.Error()))
		return nil
	}
	return containerIDs
}

func readAdvancedConfigIfUnset(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, projectID, clusterName string, p *ProcessArgs) {
	advConfig := readIfUnsetAdvancedConfiguration(ctx, diags, client, projectID, clusterName, p.ArgsDefault)
	if diags.HasError() {
		return
	}
	p.ArgsDefault = advConfig
}

func updateModelAdvancedConfig(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, model *TFModel, p *ProcessArgs) {
	readAdvancedConfigIfUnset(ctx, diags, client, model.ProjectID.ValueString(), model.Name.ValueString(), p)
	if !diags.HasError() {
		model.AdvancedConfiguration = buildAdvancedConfigObjType(ctx, p, diags)
	}
}

func updateModelAdvancedConfigDS(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, model *TFModelDS, p *ProcessArgs) {
	readAdvancedConfigIfUnset(ctx, diags, client, model.ProjectID.ValueString(), model.Name.ValueString(), p)
	if !diags.HasError() {
		model.AdvancedConfiguration = buildAdvancedConfigObjType(ctx, p, diags)
	}
}

func createCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
	pauseAfter := req.GetPaused()
	if pauseAfter {
		req.Paused = nil
	}
	_, _, err := client.AtlasV2.ClustersApi.CreateCluster(ctx, waitParams.ProjectID, req).UseEffectiveInstanceFields(waitParams.UseEffectiveFields).Execute()
	if err != nil {
		addErrorDiag(diags, operationCreate, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	clusterResp := AwaitChanges(ctx, client, waitParams, operationCreate, diags)
	if diags.HasError() {
		return nil
	}
	if pauseAfter {
		clusterResp = updateCluster(ctx, diags, client, &pauseRequest, waitParams, operationPauseAfterCreate)
	}
	return clusterResp
}

func updateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, waitParams *ClusterWaitParams, operationName string) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpdateCluster(ctx, waitParams.ProjectID, waitParams.ClusterName, req).UseEffectiveInstanceFields(waitParams.UseEffectiveFields).Execute()
	if err != nil {
		addErrorDiag(diags, operationName, defaultAPIErrorDetails(waitParams.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, waitParams, operationName, diags)
}

func resolveClusterWaitParams(ctx context.Context, model *TFModel, diags *diag.Diagnostics, operation string) *ClusterWaitParams {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	operationTimeout := cleanup.ResolveTimeout(ctx, &model.Timeouts, operation, diags)
	if diags.HasError() {
		return nil
	}
	return &ClusterWaitParams{
		ProjectID:          projectID,
		ClusterName:        clusterName,
		Timeout:            operationTimeout,
		IsDelete:           operation == operationDelete,
		UseEffectiveFields: model.UseEffectiveFields.ValueBool(),
	}
}

type clusterDiff struct {
	clusterPatchOnlyReq       *admin.ClusterDescription20240805
	upgradeTenantReq          *admin.LegacyAtlasTenantClusterUpgradeRequest
	upgradeFlexToDedicatedReq *admin.AtlasTenantClusterUpgradeRequest20240805
	isUpgradeTenantToFlex     bool
	isUpdateOfFlex            bool
}

func (c *clusterDiff) isClusterPatchOnly() bool {
	return !update.IsZeroValues(c.clusterPatchOnlyReq)
}

func (c *clusterDiff) isUpgradeTenant() bool {
	return c.upgradeTenantReq != nil
}

func (c *clusterDiff) isUpgradeFlexToDedicated() bool {
	return c.upgradeFlexToDedicatedReq != nil
}

func (c *clusterDiff) isAnyUpgrade() bool {
	return c.isUpgradeTenantToFlex || c.isUpgradeTenant() || c.isUpgradeFlexToDedicated()
}

// findClusterDiff should be called only in Update, e.g. it will fail for a flex cluster with no changes.
func findClusterDiff(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics) clusterDiff {
	stateReq := newAtlasReq(ctx, state, diags)
	planReq := newAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return clusterDiff{}
	}

	if isFlex(planReq.ReplicationSpecs) {
		if isValidUpgradeTenantToFlex(stateReq, planReq) {
			return clusterDiff{isUpgradeTenantToFlex: true}
		}
		if isValidUpdateOfFlex(stateReq, planReq) {
			return clusterDiff{isUpdateOfFlex: true}
		}
		diags.AddError(flexcluster.ErrorNonUpdatableAttributes, "")
		return clusterDiff{}
	}

	patchOptions := update.PatchOptions{
		IgnoreInStatePrefix: []string{"replicationSpecs"}, // only use config values for replicationSpecs, state values might come from the UseStateForUnknown and shouldn't be used, `id` is added in updateLegacyReplicationSpecs
	}
	patchReq, err := update.PatchPayload(stateReq, planReq, patchOptions)
	if err != nil {
		diags.AddError(errorPatchPayload, err.Error())
		return clusterDiff{}
	}
	if update.IsZeroValues(patchReq) { // No changes to cluster
		return clusterDiff{}
	}
	upgradeTenantReq := getUpgradeTenantRequest(stateReq, patchReq)
	upgradeFlexToDedicatedReq := getUpgradeFlexToDedicatedRequest(stateReq, patchReq)
	if upgradeTenantReq != nil || upgradeFlexToDedicatedReq != nil {
		return clusterDiff{upgradeTenantReq: upgradeTenantReq, upgradeFlexToDedicatedReq: upgradeFlexToDedicatedReq}
	}
	return clusterDiff{clusterPatchOnlyReq: patchReq}
}

func handleFlexUpgrade(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, plan *TFModel) *TFModel {
	configReq := newAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return nil
	}
	flexCluster := flexUpgrade(ctx, diags, client, waitParams, getUpgradeToFlexClusterRequest(configReq))
	if diags.HasError() {
		return nil
	}
	return newTFModelFlex(ctx, diags, flexCluster, getPriorityOfFlexReplicationSpecs(configReq.ReplicationSpecs), plan)
}

func handleFlexUpdate(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, plan *TFModel) *TFModel {
	configReq := newAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return nil
	}
	clusterName := plan.Name.ValueString()
	flexCluster, err := flexcluster.UpdateFlexCluster(ctx, plan.ProjectID.ValueString(), clusterName,
		getFlexClusterUpdateRequest(configReq.Tags, configReq.TerminationProtectionEnabled),
		client.AtlasV2.FlexClustersApi, waitParams.Timeout)
	if err != nil {
		diags.AddError(fmt.Sprintf(flexcluster.ErrorUpdateFlex, clusterName), err.Error())
		return nil
	}
	return newTFModelFlex(ctx, diags, flexCluster, getPriorityOfFlexReplicationSpecs(configReq.ReplicationSpecs), plan)
}
