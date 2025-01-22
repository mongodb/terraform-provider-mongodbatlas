package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithMoveState = &rs{}
var _ resource.ResourceWithUpgradeState = &rs{}

const (
	resourceName                    = "advanced_cluster"
	errorSchemaDowngrade            = "error operation not permitted, nums_shards from 1 -> > 1"
	errorPatchPayload               = "error creating patch payload"
	errorCreate                     = "error creating advanced cluster"
	operationCreateLegacy           = "create API 20240805"
	errorCreateLegacy20240805       = "error creating advanced cluster API 20240805"
	errorDetailDefault              = "cluster name %s. API error detail %s"
	errorUpdateAdvancedConfigLegacy = "error updating advanced cluster advanced configuration options with legacy API"
	errorSchemaUpgradeReadIDs       = "error reading IDs from API when upgrading schema"
	errorReadResource               = "error reading advanced cluster"
	errorAdvancedConfRead           = "error reading Advanced Configuration"
	errorAdvancedConfReadLegacy     = "error reading Advanced Configuration from legacy API"
	errorDelete                     = "error deleting advanced cluster"
	errorUpdate                     = "error updating advanced cluster"
	errorUpdateLegacy20240805       = "error updating advanced cluster legacy API 20240805"
	errorUpdateLegacy20240530       = "error updating advanced cluster legacy API 20240530"
	errorList                       = "error reading  advanced cluster list"
	errorListDetail                 = "project ID %s. Error %s"
	errorTenantUpgrade              = "error upgrading tenant cluster"
	errorReadLegacy20240530         = "error reading cluster with legacy API 20240530"
	errorResolveContainerIDs        = "error resolving container IDs"
	errorRegionPriorities           = "priority values in region_configs must be in descending order"
	errorUnknownChangeReason        = "unknown change reason"
	errorAwaitState                 = "error awaiting cluster to reach desired state"
	errorAwaitStateResultType       = "the result of awaiting cluster wasn't of the expected type"
	errorAdvancedConfUpdate         = "error updating Advanced Configuration"
	errorAdvancedConfUpdateLegacy   = "error updating Advanced Configuration from legacy API"
	errorPinningFCV                 = "error pinning FCV"
	errorUnpinningFCV               = "error unpinning FCV"

	DeprecationOldSchemaAction                   = "Please refer to our examples, documentation, and 1.18.0 migration guide for more details at https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide.html.markdown"
	defaultTimeout                               = 3 * time.Hour
	ErrorCodeClusterNotFound                     = "CLUSTER_NOT_FOUND"
	operationUpdate                              = "update"
	operationCreate                              = "create"
	operationPauseAfterCreate                    = "pause after create"
	operationDelete                              = "delete"
	operationAdvancedConfigurationUpdate20240530 = "update advanced configuration 20240530"
	operationAdvancedConfigurationUpdate         = "update advanced configuration"
	operationTenantUpgrade                       = "tenant upgrade"
	operationPauseAfterUpdate                    = "pause after update"
	operationResumeBeforeUpdate                  = "resume before update"
)

func defaultAPIErrorDetails(clusterName string, err error) string {
	return fmt.Sprintf(errorDetailDefault, clusterName, err.Error())
}

func deprecationMsgOldSchema(name string) string {
	return fmt.Sprintf("%s Name=%s. %s", constant.DeprecationParam, name, DeprecationOldSchemaAction)
}

func resolveClusterReader(ctx context.Context, model *TFModel, diags *diag.Diagnostics, changeReason string) *ClusterReader {
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()

	operationTimeout := resolveTimeout(ctx, &model.Timeouts, changeReason, diags)
	if diags.HasError() {
		return nil
	}
	return &ClusterReader{
		ProjectID:   projectID,
		ClusterName: clusterName,
		Timeout:     operationTimeout,
	}
}

var (
	pauseRequest               = admin.ClusterDescription20240805{Paused: conversion.Pointer(true)}
	resumeRequest              = admin.ClusterDescription20240805{Paused: conversion.Pointer(false)}
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
	model := r.createCluster(ctx, &plan, diags)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model := r.readCluster(ctx, diags, &state, &resp.State)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	reader := resolveClusterReader(ctx, &plan, diags, operationUpdate)
	if diags.HasError() {
		return
	}
	var clusterResp *admin.ClusterDescription20240805

	// FCV update is intentionally handled before any other cluster updates, and will wait for cluster to reach IDLE state before continuing
	clusterResp = r.applyPinnedFCVChanges(ctx, diags, &state, &plan)
	if diags.HasError() {
		return
	}

	stateUsingLegacy := usingLegacySchema(ctx, state.ReplicationSpecs, diags)
	planUsingLegacy := usingLegacySchema(ctx, plan.ReplicationSpecs, diags)
	if planUsingLegacy && !stateUsingLegacy {
		diags.AddError(errorSchemaDowngrade, fmt.Sprintf(errorSchemaDowngradeDetail, plan.Name.ValueString()))
		return
	}
	isSchemaUpgrade := stateUsingLegacy && !planUsingLegacy
	stateReq := normalizeFromTFModel(ctx, &state, diags, false)
	planReq := normalizeFromTFModel(ctx, &plan, diags, isSchemaUpgrade)
	if diags.HasError() {
		return
	}

	patchOptions := update.PatchOptions{
		IgnoreInStatePrefix: []string{"regionConfigs"},
		IgnoreInStateSuffix: []string{"id", "zoneId"}, // replication_spec.*.zone_id|id doesn't have to be included, the API will do its best to create a minimal change
	}
	if findNumShardsUpdates(ctx, &state, &plan, diags) != nil {
		// force update the replicationSpecs when update.PatchPayload will not detect changes by default:
		// `num_shards` updates is only in the legacy ClusterDescription
		patchOptions.ForceUpdateAttr = append(patchOptions.ForceUpdateAttr, "replicationSpecs")
	}
	patchReq, err := update.PatchPayload(stateReq, planReq, patchOptions)
	if err != nil {
		diags.AddError(errorPatchPayload, err.Error())
		return
	}
	if !update.IsZeroValues(patchReq) {
		upgradeRequest := getTenantUpgradeRequest(stateReq, patchReq)
		if upgradeRequest != nil {
			clusterResp = tenantUpgrade(ctx, diags, r.Client, reader, upgradeRequest)
		} else {
			if isSchemaUpgrade {
				specs, err := populateIDValuesUsingNewAPI(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), r.Client.AtlasV2.ClustersApi, patchReq.ReplicationSpecs)
				if err != nil {
					diags.AddError(errorSchemaUpgradeReadIDs, defaultAPIErrorDetails(plan.Name.ValueString(), err))
					return
				}
				patchReq.ReplicationSpecs = specs
			}
			clusterResp = r.applyClusterChanges(ctx, diags, &state, &plan, patchReq, reader)
		}
		if diags.HasError() {
			return
		}
	}
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if diags.HasError() {
		return
	}
	legacyAdvConfig, advConfig, advConfigChanged := updateAdvancedConfiguration(ctx, diags, r.Client, patchReqProcessArgsLegacy, patchReqProcessArgs, reader)
	if diags.HasError() {
		return
	}
	modelOut := &state
	if clusterResp != nil {
		modelOut, _ = getBasicClusterModel(ctx, diags, r.Client, clusterResp, &plan, false)
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
	clusterName := state.Name.ValueString()
	projectID := state.ProjectID.ValueString()
	api := r.Client.AtlasV2.ClustersApi
	params := &admin.DeleteClusterApiParams{
		GroupId:     projectID,
		ClusterName: clusterName,
	}
	if retainBackups := conversion.NilForUnknown(state.RetainBackupsEnabled, state.RetainBackupsEnabled.ValueBoolPointer()); retainBackups != nil {
		params.RetainBackups = retainBackups
	}
	_, err := api.DeleteClusterWithParams(ctx, params).Execute()
	if err != nil {
		diags.AddError(errorDelete, defaultAPIErrorDetails(clusterName, err))
		return
	}
	_ = awaitChanges(ctx, r.Client, &state.Timeouts, diags, projectID, clusterName, operationDelete, "")
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}

func (r *rs) createCluster(ctx context.Context, plan *TFModel, diags *diag.Diagnostics) *TFModel {
	latestReq := normalizeFromTFModel(ctx, plan, diags, true)
	if diags.HasError() {
		return nil
	}
	ids := resolveClusterReader(ctx, plan, diags, operationCreate)
	if diags.HasError() {
		return nil
	}
	var (
		pauseAfter  = latestReq.GetPaused()
		clusterResp *admin.ClusterDescription20240805
	)
	if pauseAfter {
		latestReq.Paused = nil
	}
	if usingLegacySchema(ctx, plan.ReplicationSpecs, diags) {
		legacyReq := newLegacyModel(latestReq)
		clusterResp = createCluster20240805(ctx, diags, r.Client, legacyReq, ids)
	} else {
		clusterResp = createCluster(ctx, diags, r.Client, latestReq, ids)
	}
	if diags.HasError() {
		return nil
	}
	if pauseAfter {
		clusterResp = updateCluster(ctx, diags, r.Client, &pauseRequest, ids, operationCreate)
	}
	emptyAdvancedConfiguration := types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &emptyAdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if diags.HasError() {
		return nil
	}
	legacyAdvConfig, advConfig, _ := updateAdvancedConfiguration(ctx, diags, r.Client, patchReqProcessArgsLegacy, patchReqProcessArgs, ids)
	if diags.HasError() {
		return nil
	}
	if changedCluster := r.applyPinnedFCVChanges(ctx, diags, &TFModel{}, plan); changedCluster != nil {
		clusterResp = changedCluster
	}
	if diags.HasError() {
		return nil
	}

	modelOut, _ := getBasicClusterModel(ctx, diags, r.Client, clusterResp, plan, false)
	if diags.HasError() {
		return nil
	}
	legacyAdvConfig, advConfig = readIfUnsetAdvancedConfiguration(ctx, diags, r.Client, ids, legacyAdvConfig, advConfig)
	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, legacyAdvConfig, advConfig)
	if diags.HasError() {
		return nil
	}
	AddAdvancedConfig(ctx, modelOut, advConfig, legacyAdvConfig, diags)
	return modelOut
}

func (r *rs) readCluster(ctx context.Context, diags *diag.Diagnostics, state *TFModel, respState *tfsdk.State) *TFModel {
	clusterName := state.Name.ValueString()
	projectID := state.ProjectID.ValueString()
	api := r.Client.AtlasV2.ClustersApi
	readResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			respState.RemoveResource(ctx)
			return nil
		}
		diags.AddError(errorReadResource, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	warningIfFCVExpiredOrUnpinnedExternally(diags, state, readResp)
	modelOut, _ := getBasicClusterModel(ctx, diags, r.Client, readResp, state, false)
	if diags.HasError() {
		return nil
	}
	updateModelAdvancedConfig(ctx, diags, r.Client, modelOut, nil, nil)
	if diags.HasError() {
		return nil
	}
	return modelOut
}

func (r *rs) applyPinnedFCVChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) *admin.ClusterDescription20240805 {
	var (
		api         = r.Client.AtlasV2.ClustersApi
		projectID   = plan.ProjectID.ValueString()
		clusterName = plan.Name.ValueString()
	)
	if !state.PinnedFCV.Equal(plan.PinnedFCV) {
		isFCVPresentInConfig := !plan.PinnedFCV.IsNull()
		if isFCVPresentInConfig {
			fcvModel := &TFPinnedFCVModel{}
			// pinned_fcv has been defined or updated expiration date
			if localDiags := plan.PinnedFCV.As(ctx, fcvModel, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
				diags.Append(localDiags...)
				return nil
			}
			if err := PinFCV(ctx, api, projectID, clusterName, fcvModel.ExpirationDate.ValueString()); err != nil {
				diags.AddError(errorUnpinningFCV, defaultAPIErrorDetails(clusterName, err))
				return nil
			}
		} else {
			// pinned_fcv has been removed from the config so unpin method is called
			if _, _, err := api.UnpinFeatureCompatibilityVersion(ctx, projectID, clusterName).Execute(); err != nil {
				diags.AddError(errorUnpinningFCV, defaultAPIErrorDetails(clusterName, err))
				return nil
			}
		}
		// ensures cluster is in IDLE state before continuing with other changes
		return awaitChanges(ctx, r.Client, &plan.Timeouts, diags, projectID, clusterName, operationUpdate, "FCV pinning")
	}
	return nil
}

func (r *rs) applyClusterChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, patchReq *admin.ClusterDescription20240805, reader *ClusterReader) *admin.ClusterDescription20240805 {
	var cluster *admin.ClusterDescription20240805
	if usingLegacySchema(ctx, plan.ReplicationSpecs, diags) {
		// Only updates of replication specs will be done with legacy API
		legacySpecsChanged := r.updateLegacyReplicationSpecs(ctx, state, plan, diags, patchReq.ReplicationSpecs)
		if diags.HasError() {
			return nil
		}
		patchReq.ReplicationSpecs = nil // Already updated by legacy API
		if legacySpecsChanged && update.IsZeroValues(patchReq) {
			return awaitChanges(ctx, r.Client, &plan.Timeouts, diags, plan.ProjectID.ValueString(), plan.Name.ValueString(), operationUpdate, "Legacy Replication Specs Update")
		}
	}
	if update.IsZeroValues(patchReq) {
		return cluster
	}
	pauseAfter := false
	if patchReq.Paused != nil && patchReq.GetPaused() {
		// More changes than pause, need to pause after
		if !reflect.DeepEqual(pauseRequest, *patchReq) {
			pauseAfter = true
			patchReq.Paused = nil
		}
	} else if patchReq.Paused != nil && !patchReq.GetPaused() {
		// More changes than pause, need to resume before applying changes
		if !reflect.DeepEqual(resumeRequest, *patchReq) {
			patchReq.Paused = nil
			_ = updateCluster(ctx, diags, r.Client, &resumeRequest, reader, operationResumeBeforeUpdate)
		}
	}
	if diags.HasError() {
		return nil
	}
	replicationSpecsUpdated := patchReq.ReplicationSpecs != nil
	if replicationSpecsUpdated {
		// Cannot call latest API (2024-10-23 or newer) as it can enable ISS autoscaling
		legacyPatch := newLegacyModel(patchReq)
		cluster = r.updateAndWaitLegacy(ctx, legacyPatch, diags, plan)
	} else {
		cluster = updateCluster(ctx, diags, r.Client, patchReq, reader, operationUpdate)
	}
	if pauseAfter && cluster != nil {
		cluster = updateCluster(ctx, diags, r.Client, &pauseRequest, reader, operationPauseAfterUpdate)
	}
	return cluster
}

func (r *rs) updateLegacyReplicationSpecs(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics, specChanges *[]admin.ReplicationSpec20240805) bool {
	numShardsUpdates := findNumShardsUpdates(ctx, state, plan, diags)
	if diags.HasError() {
		return false
	}
	if specChanges == nil && numShardsUpdates == nil { // No changes to replication specs
		return false
	}
	if specChanges == nil {
		// Use state replication specs as there are no changes in plan except for numShards updates
		specChanges = newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
		if diags.HasError() {
			return false
		}
	}
	numShardsPlan := numShardsMap(ctx, plan.ReplicationSpecs, diags)
	legacyIDs := externalIDToLegacyID(ctx, state.ReplicationSpecs, diags)
	if diags.HasError() {
		return false
	}
	legacyPatch := newLegacyModel20240530ReplicationSpecsAndDiskGBOnly(specChanges, numShardsPlan, state.DiskSizeGB.ValueFloat64Pointer(), legacyIDs)
	if diags.HasError() {
		return false
	}
	api20240530 := r.Client.AtlasV220240530.ClustersApi
	_, _, err := api20240530.UpdateCluster(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), legacyPatch).Execute()
	if err != nil {
		diags.AddError(errorUpdateLegacy20240530, defaultAPIErrorDetails(plan.Name.ValueString(), err))
		return false
	}
	return true
}

func (r *rs) updateAndWaitLegacy(ctx context.Context, patchReq *admin20240805.ClusterDescription20240805, diags *diag.Diagnostics, plan *TFModel) *admin.ClusterDescription20240805 {
	api20240805 := r.Client.AtlasV220240805.ClustersApi
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	_, _, err := api20240805.UpdateCluster(ctx, projectID, clusterName, patchReq).Execute()
	if err != nil {
		diags.AddError(errorUpdateLegacy20240805, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	return awaitChanges(ctx, r.Client, &plan.Timeouts, diags, projectID, clusterName, operationUpdate, "Update Cluster Legacy")
}

func (r *rs) applyTenantUpgrade(ctx context.Context, plan *TFModel, upgradeRequest *admin.LegacyAtlasTenantClusterUpgradeRequest, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	api := r.Client.AtlasV2.ClustersApi
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	upgradeRequest.Name = clusterName
	_, _, err := api.UpgradeSharedCluster(ctx, projectID, upgradeRequest).Execute()
	if err != nil {
		diags.AddError(errorTenantUpgrade, defaultAPIErrorDetails(clusterName, err))
		return nil
	}
	return awaitChanges(ctx, r.Client, &plan.Timeouts, diags, projectID, clusterName, operationUpdate, "Tenant Upgrade")
}

func getBasicClusterModel(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805, modelIn *TFModel, forceLegacySchema bool) (*TFModel, *ExtraAPIInfo) {
	extraInfo := resolveAPIInfo(ctx, diags, client, modelIn, clusterResp, forceLegacySchema)
	if diags.HasError() {
		return nil, nil
	}
	if extraInfo.ForceLegacySchemaFailed { // can't create a model if legacy is forced but cluster does not support it
		return nil, extraInfo
	}
	modelOut := NewTFModel(ctx, clusterResp, modelIn.Timeouts, diags, *extraInfo)
	if diags.HasError() {
		return nil, nil
	}
	overrideAttributesWithPrevStateValue(modelIn, modelOut)
	return modelOut, extraInfo
}

func updateModelAdvancedConfig(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, model *TFModel, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin.ClusterDescriptionProcessArgs20240805) {
	api := client.AtlasV2.ClustersApi
	api20240530 := client.AtlasV220240530.ClustersApi
	projectID := model.ProjectID.ValueString()
	clusterName := model.Name.ValueString()
	var err error
	if legacyAdvConfig == nil {
		legacyAdvConfig, _, err = api20240530.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfReadLegacy, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	if advConfig == nil {
		advConfig, _, err = api.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	AddAdvancedConfig(ctx, model, advConfig, legacyAdvConfig, diags)
}

func warningIfFCVExpiredOrUnpinnedExternally(diags *diag.Diagnostics, state *TFModel, clusterResp *admin.ClusterDescription20240805) {
	fcvPresentInState := !state.PinnedFCV.IsNull()
	newWarnings := GenerateFCVPinningWarningForRead(fcvPresentInState, clusterResp.FeatureCompatibilityVersionExpirationDate)
	diags.Append(newWarnings...)
}

func awaitChanges(ctx context.Context, client *config.MongoDBClient, t *timeouts.Value, diags *diag.Diagnostics, projectID, clusterName, changeReason, lastOperation string) (cluster *admin.ClusterDescription20240805) {
	timeoutDuration := resolveTimeout(ctx, t, changeReason, diags)
	if diags.HasError() {
		return nil
	}
	if lastOperation == "" {
		lastOperation = changeReason
	}
	ids := ClusterReader{
		ProjectID:   projectID,
		ClusterName: clusterName,
		Timeout:     timeoutDuration,
	}
	return AwaitChanges(ctx, client, &ids, lastOperation, diags)
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
		diags.AddError(errorUnknownChangeReason, "unknown change reason "+operationName)
	}
	return timeoutDuration
}
