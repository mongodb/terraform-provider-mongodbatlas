package advancedclustertpf

import (
	"context"
	"fmt"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithMoveState = &rs{}
var _ resource.ResourceWithUpgradeState = &rs{}
var _ resource.ResourceWithModifyPlan = &rs{}

const (
	resourceName                  = "advanced_cluster"
	errorSchemaDowngrade          = "error operation not permitted, nums_shards from 1 -> > 1"
	errorPatchPayload             = "error creating patch payload"
	errorDetailDefault            = "cluster name %s. API error detail %s"
	errorSchemaUpgradeReadIDs     = "error reading IDs from API when upgrading schema"
	errorReadResource             = "error reading advanced cluster"
	errorAdvancedConfRead         = "error reading Advanced Configuration"
	errorAdvancedConfReadLegacy   = "error reading Advanced Configuration from legacy API"
	errorUpdateLegacy20240530     = "error updating advanced cluster legacy API 20240530"
	errorList                     = "error reading  advanced cluster list"
	errorListDetail               = "project ID %s. Error %s"
	errorReadLegacy20240530       = "error reading cluster with legacy API 20240530"
	errorResolveContainerIDs      = "error resolving container IDs"
	errorRegionPriorities         = "priority values in region_configs must be in descending order"
	errorAdvancedConfUpdateLegacy = "error updating Advanced Configuration from legacy API"

	DeprecationOldSchemaAction                   = "Please refer to our examples, documentation, and 1.18.0 migration guide for more details at https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide.html.markdown"
	defaultTimeout                               = 3 * time.Hour
	ErrorCodeClusterNotFound                     = "CLUSTER_NOT_FOUND"
	operationUpdate                              = "update"
	operationCreate                              = "create"
	operationCreate20240805                      = "create (legacy)"
	operationPauseAfterCreate                    = "pause after create"
	operationDelete                              = "delete"
	operationAdvancedConfigurationUpdate20240530 = "update advanced configuration (legacy)"
	operationAdvancedConfigurationUpdate         = "update advanced configuration"
	operationTenantUpgrade                       = "tenant upgrade"
	operationPauseAfterUpdate                    = "pause after update"
	operationResumeBeforeUpdate                  = "resume before update"
	operationReplicationSpecsUpdateLegacy        = "update replication specs legacy"
	operationFCVPinning                          = "FCV pinning"
	operationFCVUnpinning                        = "FCV unpinning"
)

func addErrorDiag(diags *diag.Diagnostics, errorLocator, details string) {
	diags.AddError("Error in "+errorLocator, details)
}

func defaultAPIErrorDetails(clusterName string, err error) string {
	return fmt.Sprintf(errorDetailDefault, clusterName, err.Error())
}

func deprecationMsgOldSchema(name string) string {
	return fmt.Sprintf("%s Name=%s. %s", constant.DeprecationParam, name, DeprecationOldSchemaAction)
}

var (
	resumeRequest              = admin.ClusterDescription20240805{Paused: conversion.Pointer(false)}
	pauseRequest               = admin.ClusterDescription20240805{Paused: conversion.Pointer(true)}
	errorSchemaDowngradeDetail = "Cluster name %s. " + fmt.Sprintf("cannot increase num_shards to > 1 under the current configuration. New shards can be defined by adding new replication spec objects; %s", DeprecationOldSchemaAction)
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
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() { // Can be null in case of destroy
		return
	}
	var plan, state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
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
	latestReq := normalizeFromTFModel(ctx, &plan, diags, true)
	if diags.HasError() {
		return
	}
	waitParams := resolveClusterWaitParams(ctx, &plan, diags, operationCreate)
	if diags.HasError() {
		return
	}
	clusterResp := CreateCluster(ctx, diags, r.Client, latestReq, waitParams, usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags))
	emptyAdvancedConfiguration := types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if diags.HasError() {
		return
	}
	legacyAdvConfig, advConfig, _ := UpdateAdvancedConfiguration(ctx, diags, r.Client, patchReqProcessArgsLegacy, patchReqProcessArgs, waitParams)
	if diags.HasError() {
		return
	}
	if changedCluster := r.applyPinnedFCVChanges(ctx, diags, &TFModel{}, &plan, waitParams); changedCluster != nil {
		clusterResp = changedCluster
	}
	if diags.HasError() {
		return
	}

	modelOut, _ := getBasicClusterModelResource(ctx, diags, r.Client, clusterResp, &plan)
	if diags.HasError() {
		return
	}
	legacyAdvConfig, advConfig = ReadIfUnsetAdvancedConfiguration(ctx, diags, r.Client, waitParams.ProjectID, waitParams.ClusterName, legacyAdvConfig, advConfig)
	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, legacyAdvConfig, advConfig)
	if diags.HasError() {
		return
	}
	AddAdvancedConfig(ctx, modelOut, advConfig, legacyAdvConfig, diags)
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
	readResp := ReadCluster(ctx, diags, r.Client, projectID, clusterName, !state.PinnedFCV.IsNull())
	if diags.HasError() {
		return
	}
	if readResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	modelOut, _ := getBasicClusterModelResource(ctx, diags, r.Client, readResp, &state)
	if diags.HasError() {
		return
	}
	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, nil, nil)
	if diags.HasError() {
		return
	}
	diags.Append(resp.State.Set(ctx, modelOut)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, config TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &config)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	waitParams := resolveClusterWaitParams(ctx, &config, diags, operationUpdate)
	if diags.HasError() {
		return
	}
	var clusterResp *admin.ClusterDescription20240805

	// FCV update is intentionally handled before any other cluster updates, and will wait for cluster to reach IDLE state before continuing
	clusterResp = r.applyPinnedFCVChanges(ctx, diags, &state, &config, waitParams)
	if diags.HasError() {
		return
	}
	patchOptions := update.PatchOptions{
		IgnoreInStatePrefix: []string{"replicationSpecs"},
		IgnoreInStateSuffix: []string{"zoneId"}, // replication_spec.*.zone_id doesn't have to be included, the API will do its best to create a minimal change
	}
	if usingNewShardingConfig(ctx, config.ReplicationSpecs, diags) {
		patchOptions.IgnoreInStateSuffix = append(patchOptions.IgnoreInStateSuffix, "id") // Not safe to send replication_spec.*.id when using the new schema: replicationSpecs.java.util.ArrayList[0].id attribute does not match expected format
	}
	patchReq, upgradeReq := findClusterDiff(ctx, &state, &config, diags, &patchOptions)
	if diags.HasError() {
		return
	}
	if upgradeReq != nil {
		clusterResp = TenantUpgrade(ctx, diags, r.Client, waitParams, upgradeReq)
		if diags.HasError() {
			return
		}
	}
	if !update.IsZeroValues(patchReq) {
		clusterResp = r.applyClusterChanges(ctx, diags, &state, &config, patchReq, waitParams)
		if diags.HasError() {
			return
		}
	}
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &config.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &config.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if diags.HasError() {
		return
	}
	legacyAdvConfig, advConfig, advConfigChanged := UpdateAdvancedConfiguration(ctx, diags, r.Client, patchReqProcessArgsLegacy, patchReqProcessArgs, waitParams)
	if diags.HasError() {
		return
	}
	var modelOut *TFModel
	if clusterResp == nil { // no Atlas updates needed but override is still needed (e.g. tags going from nil to [] or vice versa)
		modelOut = &state
		overrideAttributesWithPrevStateValue(&config, modelOut)
	} else {
		modelOut, _ = getBasicClusterModelResource(ctx, diags, r.Client, clusterResp, &config)
		if diags.HasError() {
			return
		}
	}
	if advConfigChanged {
		updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, legacyAdvConfig, advConfig)
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
	if _, _, err := api.UnpinFeatureCompatibilityVersion(ctx, projectID, clusterName).Execute(); err != nil {
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

	if !usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) {
		// With old sharding config we call older API (2023-02-01) for updating replication specs to avoid cluster having asymmetric autoscaling mode. Old sharding config can only represent symmetric clusters.
		r.updateLegacyReplicationSpecs(ctx, state, plan, diags, patchReq.ReplicationSpecs)
		if diags.HasError() {
			return nil
		}
		patchReq.ReplicationSpecs = nil // Already updated by 2023-02-01 API
		if update.IsZeroValues(patchReq) && !pauseAfterOtherChanges {
			return AwaitChanges(ctx, r.Client, waitParams, operationReplicationSpecsUpdateLegacy, diags)
		}
	}

	// latest API can be used safely because if old sharding config is used replication specs will not be included in this request
	result = updateCluster(ctx, diags, r.Client, patchReq, waitParams, operationUpdate)

	if pauseAfterOtherChanges {
		result = updateCluster(ctx, diags, r.Client, &pauseRequest, waitParams, operationPauseAfterUpdate)
	}
	return result
}

func (r *rs) updateLegacyReplicationSpecs(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics, specChanges *[]admin.ReplicationSpec20240805) {
	numShardsUpdates := findNumShardsUpdates(ctx, state, plan, diags)
	if diags.HasError() {
		return
	}
	if specChanges == nil && numShardsUpdates == nil { // No changes to replication specs
		return
	}
	if specChanges == nil {
		// Use state replication specs as there are no changes in plan except for numShards updates
		specChanges = newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
		if diags.HasError() {
			return
		}
	}
	numShardsPlan := numShardsMap(ctx, plan.ReplicationSpecs, diags)
	legacyIDs := externalIDToLegacyID(ctx, state.ReplicationSpecs, diags)
	if diags.HasError() {
		return
	}
	legacyPatch := newLegacyModel20240530ReplicationSpecsAndDiskGBOnly(specChanges, numShardsPlan, state.DiskSizeGB.ValueFloat64Pointer(), legacyIDs)
	if diags.HasError() {
		return
	}
	api20240530 := r.Client.AtlasV220240530.ClustersApi
	_, _, err := api20240530.UpdateCluster(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), legacyPatch).Execute()
	if err != nil {
		diags.AddError(errorUpdateLegacy20240530, defaultAPIErrorDetails(plan.Name.ValueString(), err))
	}
}

func getBasicClusterModelResource(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805, modelIn *TFModel) (*TFModel, *ExtraAPIInfo) {
	useReplicationSpecPerShard := usingNewShardingConfig(ctx, modelIn.ReplicationSpecs, diags)
	if diags.HasError() {
		return nil, nil
	}
	modelOut, apiInfo := getBasicClusterModel(ctx, diags, client, clusterResp, useReplicationSpecPerShard)
	if modelOut != nil {
		modelOut.Timeouts = modelIn.Timeouts
		overrideAttributesWithPrevStateValue(modelIn, modelOut)
	}
	return modelOut, apiInfo
}

func getBasicClusterModel(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805, useReplicationSpecPerShard bool) (*TFModel, *ExtraAPIInfo) {
	extraInfo := resolveAPIInfo(ctx, diags, client, clusterResp, useReplicationSpecPerShard)
	if diags.HasError() {
		return nil, nil
	}
	if extraInfo.UseOldShardingConfigFailed { // can't create a model if the cluster does not support old sharding config
		return nil, extraInfo
	}
	modelOut := NewTFModel(ctx, clusterResp, diags, *extraInfo)
	if diags.HasError() {
		return nil, nil
	}
	return modelOut, extraInfo
}

func updateModelAdvancedConfig(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, model *TFModel, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin.ClusterDescriptionProcessArgs20240805) {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	legacyAdvConfig, advConfig = ReadIfUnsetAdvancedConfiguration(ctx, diags, client, projectID, clusterName, legacyAdvConfig, advConfig)
	if diags.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, advConfig, legacyAdvConfig, diags)
}

func resolveClusterWaitParams(ctx context.Context, model *TFModel, diags *diag.Diagnostics, operation string) *ClusterWaitParams {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	operationTimeout := resolveTimeout(ctx, &model.Timeouts, operation, diags)
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

func resolveTimeout(ctx context.Context, t *timeouts.Value, operationName string, diags *diag.Diagnostics) time.Duration {
	var (
		timeoutDuration time.Duration
		localDiags      diag.Diagnostics
	)
	switch operationName {
	case operationCreate:
		timeoutDuration, localDiags = t.Create(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case operationUpdate:
		timeoutDuration, localDiags = t.Update(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case operationDelete:
		timeoutDuration, localDiags = t.Delete(ctx, defaultTimeout)
		diags.Append(localDiags...)
	default:
		timeoutDuration = defaultTimeout
	}
	return timeoutDuration
}

func findClusterDiff(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics, options *update.PatchOptions) (*admin.ClusterDescription20240805, *admin.LegacyAtlasTenantClusterUpgradeRequest) {
	isShardingUpgrade := isShardingConfigUpgrade(ctx, state, plan, diags)
	if diags.HasError() {
		return nil, nil
	} // old sharding config  (num_shards > 1) to new one
	stateReq := normalizeFromTFModel(ctx, state, diags, false)
	planReq := normalizeFromTFModel(ctx, plan, diags, isShardingUpgrade)
	if diags.HasError() {
		return nil, nil
	}
	if findNumShardsUpdates(ctx, state, plan, diags) != nil {
		// force update the replicationSpecs when update.PatchPayload will not detect changes by default:
		// `num_shards` updates is only in the legacy ClusterDescription
		options.ForceUpdateAttr = append(options.ForceUpdateAttr, "replicationSpecs")
	}
	patchReq, err := update.PatchPayload(stateReq, planReq, *options)
	if err != nil {
		diags.AddError(errorPatchPayload, err.Error())
		return nil, nil
	}
	if update.IsZeroValues(patchReq) { // No changes to cluster
		return nil, nil
	}
	upgradeRequest := getTenantUpgradeRequest(stateReq, patchReq)
	if upgradeRequest != nil {
		return nil, upgradeRequest
	}
	return patchReq, nil
}

func isShardingConfigUpgrade(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics) bool {
	stateUsingNewSharding := usingNewShardingConfig(ctx, state.ReplicationSpecs, diags)
	planUsingNewSharding := usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags)
	if stateUsingNewSharding && !planUsingNewSharding {
		diags.AddError(errorSchemaDowngrade, fmt.Sprintf(errorSchemaDowngradeDetail, plan.Name.ValueString()))
		return false
	}
	return !stateUsingNewSharding && planUsingNewSharding
}
