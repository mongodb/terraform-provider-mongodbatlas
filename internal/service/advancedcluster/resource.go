package advancedcluster

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

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
	resourceName                = "advanced_cluster"
	errorSchemaDowngrade        = "error operation not permitted, nums_shards from 1 -> > 1"
	errorPatchPayload           = "error creating patch payload"
	errorDetailDefault          = "cluster name: %s, API error details: %s"
	errorSchemaUpgradeReadIDs   = "error reading IDs from API when upgrading schema"
	errorReadResource           = "error reading advanced cluster"
	errorAdvancedConfRead       = "error reading Advanced Configuration"
	errorAdvancedConfReadLegacy = "error reading Advanced Configuration from legacy API"
	errorUpdateLegacy20240530   = "error updating advanced cluster legacy API 20240530"
	errorList                   = "error reading  advanced cluster list"
	errorListDetail             = "project ID %s. Error %s"
	errorReadLegacy20240530     = "error reading cluster with legacy API 20240530"
	errorResolveContainerIDs    = "error resolving container IDs"
	errorRegionPriorities       = "priority values in region_configs must be in descending order"

	ErrorCodeClusterNotFound                     = "CLUSTER_NOT_FOUND"
	operationUpdate                              = "update"
	operationCreate                              = "create"
	operationCreate20240805                      = "create (legacy)"
	operationPauseAfterCreate                    = "pause after create"
	operationDelete                              = "delete"
	operationDeleteFlex                          = "flex delete"
	operationAdvancedConfigurationUpdate20240530 = "update advanced configuration (legacy)"
	operationAdvancedConfigurationUpdate         = "update advanced configuration"
	operationTenantUpgrade                       = "tenant upgrade"
	operationPauseAfterUpdate                    = "pause after update"
	operationResumeBeforeUpdate                  = "resume before update"
	operationReplicationSpecsUpdateLegacy        = "update replication specs legacy"
	operationFCVPinning                          = "FCV pinning"
	operationFCVUnpinning                        = "FCV unpinning"
	operationFlexUpgrade                         = "flex upgrade"
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
	// useStateForUnknowns will try to convert the field to `Target Type: []advancedcluster.TFReplicationSpecsModel`.
	// But since the field is unknown the user gets an error: `Error: Value Conversion Error`.
	if plan.ReplicationSpecs.IsUnknown() {
		return
	}

	useStateForUnknowns(ctx, diags, &state, &plan)
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
	isFlex := IsFlex(latestReq.ReplicationSpecs)
	projectID, clusterName := waitParams.ProjectID, waitParams.ClusterName
	clusterDetailStr := fmt.Sprintf("Cluster name %s (project_id=%s).", clusterName, projectID)
	if cleanup.ResolveDeleteOnCreateTimeout(plan.DeleteOnCreateTimeout) {
		var deferCall func()
		ctx, deferCall = cleanup.OnTimeout(
			ctx, waitParams.Timeout, diags.AddWarning, clusterDetailStr, DeleteClusterNoWait(r.Client, projectID, clusterName, isFlex),
		)
		defer deferCall()
	}
	if isFlex {
		flexClusterReq := NewFlexCreateReq(latestReq.GetName(), latestReq.GetTerminationProtectionEnabled(), latestReq.Tags, latestReq.ReplicationSpecs)
		flexClusterResp, err := flexcluster.CreateFlexCluster(ctx, plan.ProjectID.ValueString(), latestReq.GetName(), flexClusterReq, r.Client.AtlasV2.FlexClustersApi, &waitParams.Timeout)
		if err != nil {
			diags.AddError(fmt.Sprintf(flexcluster.ErrorCreateFlex, clusterDetailStr), err.Error())
			return
		}
		newFlexClusterModel := NewTFModelFlexResource(ctx, diags, flexClusterResp, GetPriorityOfFlexReplicationSpecs(latestReq.ReplicationSpecs), &plan)
		if diags.HasError() {
			return
		}
		diags.Append(resp.State.Set(ctx, newFlexClusterModel)...)
		return
	}
	clusterResp := CreateCluster(ctx, diags, r.Client, latestReq, waitParams)

	emptyAdvancedConfiguration := types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	patchReqProcessArgs := update.PatchPayloadCluster(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if diags.HasError() {
		return
	}
	p := &ProcessArgs{
		ArgsDefault:           patchReqProcessArgs,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	}
	advConfig, _ := UpdateAdvancedConfiguration(ctx, diags, r.Client, p, waitParams)
	if diags.HasError() {
		return
	}
	if changedCluster := r.applyPinnedFCVChanges(ctx, diags, &TFModel{}, &plan, waitParams); changedCluster != nil {
		clusterResp = changedCluster
	}
	if diags.HasError() {
		return
	}

	modelOut := getBasicClusterModelResource(ctx, diags, r.Client, clusterResp, &plan)
	if diags.HasError() {
		return
	}
	advConfig = ReadIfUnsetAdvancedConfiguration(ctx, diags, r.Client, waitParams.ProjectID, waitParams.ClusterName, advConfig)

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
	cluster, flexCluster := GetClusterDetails(ctx, diags, projectID, clusterName, r.Client, !state.PinnedFCV.IsNull())
	if diags.HasError() {
		return
	}
	if cluster == nil && flexCluster == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if flexCluster != nil {
		newFlexClusterModel := NewTFModelFlexResource(ctx, diags, flexCluster, GetPriorityOfFlexReplicationSpecs(newAtlasReq(ctx, &state, diags).ReplicationSpecs), &state)
		if diags.HasError() {
			return
		}
		diags.Append(resp.State.Set(ctx, newFlexClusterModel)...)
		return
	}
	modelOut := getBasicClusterModelResource(ctx, diags, r.Client, cluster, &state)
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
			clusterResp = UpgradeFlexToDedicated(ctx, diags, r.Client, waitParams, diff.upgradeFlexToDedicatedReq)
		case diff.isUpgradeTenant():
			clusterResp = UpgradeTenant(ctx, diags, r.Client, waitParams, diff.upgradeTenantReq)
		case diff.isClusterPatchOnly():
			clusterResp = r.applyClusterChanges(ctx, diags, &state, &plan, diff.clusterPatchOnlyReq, waitParams)
		}
		if diags.HasError() {
			return
		}
	}
	// clusterResp can be nil if there are no changes to the cluster, for example when `delete_on_create_timeout` is changed or only advanced configuration is changed
	if clusterResp == nil {
		var flexResp *admin.FlexClusterDescription20241113
		clusterResp, flexResp = GetClusterDetails(ctx, diags, waitParams.ProjectID, waitParams.ClusterName, r.Client, false)
		// This should never happen since the switch case should handle the two flex cases (update/upgrade) and return, but keeping it here for safety.
		if flexResp != nil {
			flexPriority := GetPriorityOfFlexReplicationSpecs(newAtlasReq(ctx, &plan, diags).ReplicationSpecs)
			if flexOut := NewTFModelFlexResource(ctx, diags, flexResp, flexPriority, &plan); flexOut != nil {
				diags.Append(resp.State.Set(ctx, flexOut)...)
			}
			return
		}
	}
	modelOut := getBasicClusterModelResource(ctx, diags, r.Client, clusterResp, &plan)
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
	advConfig, advConfigChanged := UpdateAdvancedConfiguration(ctx, diags, r.Client, p, waitParams)
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
	DeleteCluster(ctx, diags, r.Client, waitParams, retainBackups)
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

func (r *rs) applyClusterChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, patchReq *admin.ClusterDescription20240805, waitParams *ClusterWaitParams) *admin.ClusterDescription20240805 {
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

func getBasicClusterModelResource(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805, modelIn *TFModel) *TFModel {
	if diags.HasError() {
		return nil
	}
	modelOut := getBasicClusterModel(ctx, diags, client, clusterResp)
	if modelOut != nil {
		modelOut.Timeouts = modelIn.Timeouts
		overrideAttributesWithPrevStateValue(modelIn, modelOut)
	}
	return modelOut
}

func getBasicClusterModel(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805) *TFModel {
	var (
		projectID   = clusterResp.GetGroupId()
		clusterName = clusterResp.GetName()
	)
	containerIDs, err := resolveContainerIDs(ctx, projectID, clusterResp, client.AtlasV2.NetworkPeeringApi)
	if err != nil {
		diags.AddError(errorResolveContainerIDs, fmt.Sprintf("cluster name = %s, error details: %s", clusterName, err.Error()))
		return nil
	}

	modelOut := NewTFModel(ctx, clusterResp, diags, containerIDs)
	if diags.HasError() {
		return nil
	}
	return modelOut
}

func updateModelAdvancedConfig(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, model *TFModel,
	p *ProcessArgs) {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	advConfig := ReadIfUnsetAdvancedConfiguration(ctx, diags, client, projectID, clusterName, p.ArgsDefault)
	if diags.HasError() {
		return
	}
	p.ArgsDefault = advConfig

	AddAdvancedConfig(ctx, model, p, diags)
}

func resolveClusterWaitParams(ctx context.Context, model *TFModel, diags *diag.Diagnostics, operation string) *ClusterWaitParams {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	operationTimeout := cleanup.ResolveTimeout(ctx, &model.Timeouts, operation, diags)
	if diags.HasError() {
		return nil
	}
	return &ClusterWaitParams{
		ProjectID:   projectID,
		ClusterName: clusterName,
		Timeout:     operationTimeout,
		IsDelete:    operation == operationDelete,
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

	if IsFlex(planReq.ReplicationSpecs) {
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
		IgnoreInStatePrefix: []string{"replicationSpecs"}, // only use config values for replicationSpecs, state values might come from the UseStateForUnknowns and shouldn't be used, `id` is added in updateLegacyReplicationSpecs
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
	flexCluster := FlexUpgrade(ctx, diags, client, waitParams, GetUpgradeToFlexClusterRequest(configReq))
	if diags.HasError() {
		return nil
	}
	return NewTFModelFlexResource(ctx, diags, flexCluster, GetPriorityOfFlexReplicationSpecs(configReq.ReplicationSpecs), plan)
}

func handleFlexUpdate(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, plan *TFModel) *TFModel {
	configReq := newAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return nil
	}
	clusterName := plan.Name.ValueString()
	flexCluster, err := flexcluster.UpdateFlexCluster(ctx, plan.ProjectID.ValueString(), clusterName,
		GetFlexClusterUpdateRequest(configReq.Tags, configReq.TerminationProtectionEnabled),
		client.AtlasV2.FlexClustersApi, waitParams.Timeout)
	if err != nil {
		diags.AddError(fmt.Sprintf(flexcluster.ErrorUpdateFlex, clusterName), err.Error())
		return nil
	}
	return NewTFModelFlexResource(ctx, diags, flexCluster, GetPriorityOfFlexReplicationSpecs(configReq.ReplicationSpecs), plan)
}
