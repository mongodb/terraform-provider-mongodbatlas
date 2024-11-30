package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithMoveState = &rs{}

const (
	errorCreate                    = "error creating advanced cluster: %s"
	errorRead                      = "error reading  advanced cluster (%s): %s"
	errorDelete                    = "error deleting advanced cluster (%s): %s"
	errorUpdate                    = "error updating advanced cluster (%s): %s"
	errorConfigUpdate              = "error updating advanced cluster configuration options (%s): %s"
	errorConfigRead                = "error reading advanced cluster configuration options (%s): %s"
	ErrorClusterSetting            = "error setting `%s` for MongoDB Cluster (%s): %s"
	ErrorAdvancedConfRead          = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
	ErrorClusterAdvancedSetting    = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	ErrorAdvancedClusterListStatus = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
	ErrorOperationNotPermitted     = "error operation not permitted"
	ignoreLabel                    = "Infrastructure Tool"
	DeprecationOldSchemaAction     = "Please refer to our examples, documentation, and 1.18.0 migration guide for more details at https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide.html.markdown"
	defaultTimeout                 = 3 * time.Hour
	ErrorCodeClusterNotFound       = "CLUSTER_NOT_FOUND"
	changeReasonUpdate             = "update"
	changeReasonCreate             = "create"
	changeReasonDelete             = "delete"
)

var (
	DeprecationMsgOldSchema = fmt.Sprintf("%s %s", constant.DeprecationParam, DeprecationOldSchemaAction)
	pauseRequest            = admin.ClusterDescription20240805{Paused: conversion.Pointer(true)}
	resumeRequest           = admin.ClusterDescription20240805{Paused: conversion.Pointer(false)}
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
	resp.Schema = ResourceSchema(ctx)
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
	model := r.readCluster(ctx, &state, &resp.State, diags, true)
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
	cluster := r.applyClusterChanges(ctx, diags, &state, &plan)
	if diags.HasError() {
		return
	}
	legacyAdvConfig, advConfig := r.applyAdvancedConfigurationChanges(ctx, diags, &state, &plan)
	if diags.HasError() {
		return
	}
	var model *TFModel
	if cluster == nil { // no cluster updates, TODO: See if we can only update the advanced config if no cluster changes
		model = r.readCluster(ctx, &plan, &resp.State, diags, false)
	} else {
		model = r.convertClusterAddAdvConfig(ctx, legacyAdvConfig, advConfig, cluster, &plan, diags)
	}
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
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
	_, err := api.DeleteCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		diags.AddError("errorDelete", fmt.Sprintf(errorDelete, clusterName, err.Error()))
		return
	}
	_ = AwaitChanges(ctx, r.Client.AtlasV2.ClustersApi, &state.Timeouts, diags, projectID, clusterName, changeReasonDelete)
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}

func (r *rs) createCluster(ctx context.Context, plan *TFModel, diags *diag.Diagnostics) *TFModel {
	latestReq := normalizeFromTFModel(ctx, plan, diags, true)
	if diags.HasError() {
		return nil
	}
	var (
		projectID         = plan.ProjectID.ValueString()
		clusterName       = plan.Name.ValueString()
		apiLegacy20240530 = r.Client.AtlasV220240530.ClustersApi
		apiLatest         = r.Client.AtlasV2.ClustersApi
		err               error
		pauseAfter        = latestReq.GetPaused()
	)
	_, _, err = apiLatest.CreateCluster(ctx, projectID, latestReq).Execute()
	if err != nil {
		diags.AddError("errorCreate", fmt.Sprintf(errorCreate, err.Error()))
		return nil
	}
	cluster := AwaitChanges(ctx, r.Client.AtlasV2.ClustersApi, &plan.Timeouts, diags, projectID, clusterName, changeReasonCreate)
	if diags.HasError() {
		return nil
	}
	if pauseAfter {
		cluster = r.updateAndWait(ctx, &pauseRequest, diags, plan)
	}
	var legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
	legacyAdvConfigUpdate := NewAtlasReqAdvancedConfigurationLegacy(ctx, &plan.AdvancedConfiguration, diags)
	if !update.IsEmpty(legacyAdvConfigUpdate) {
		legacyAdvConfig, _, err = apiLegacy20240530.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, legacyAdvConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			diags.AddError("errorUpdateeAdvConfigLegacy", fmt.Sprintf(errorCreate, err.Error()))
			return nil
		}
	}

	advConfigUpdate := NewAtlasReqAdvancedConfiguration(ctx, &plan.AdvancedConfiguration, diags)
	var advConfig *admin.ClusterDescriptionProcessArgs20240805
	if !update.IsEmpty(advConfigUpdate) {
		advConfig, _, err = apiLatest.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, advConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			diags.AddError("errorUpdateAdvConfig", fmt.Sprintf(errorCreate, err.Error()))
			return nil
		}
	}
	return r.convertClusterAddAdvConfig(ctx, legacyAdvConfig, advConfig, cluster, plan, diags)
}

func (r *rs) readCluster(ctx context.Context, model *TFModel, state *tfsdk.State, diags *diag.Diagnostics, allowNotFound bool) *TFModel {
	clusterName := model.Name.ValueString()
	projectID := model.ProjectID.ValueString()
	api := r.Client.AtlasV2.ClustersApi
	readResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) && allowNotFound {
			state.RemoveResource(ctx)
			return nil
		}
		diags.AddError("errorRead", fmt.Sprintf(errorRead, clusterName, err.Error()))
		return nil
	}
	return r.convertClusterAddAdvConfig(ctx, nil, nil, readResp, model, diags)
}
func (r *rs) applyAdvancedConfigurationChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) (*admin20240530.ClusterDescriptionProcessArgs, *admin.ClusterDescriptionProcessArgs20240805) {
	var (
		apiLatest       = r.Client.AtlasV2.ClustersApi
		projectID       = plan.ProjectID.ValueString()
		clusterName     = plan.Name.ValueString()
		err             error
		advConfig       *admin.ClusterDescriptionProcessArgs20240805
		legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
	)
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if patchReqProcessArgs != nil {
		advConfig, _, err = apiLatest.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgs).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfig", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
			return legacyAdvConfig, advConfig
		}
	}
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if patchReqProcessArgsLegacy != nil {
		legacyAdvConfig, _, err = r.Client.AtlasV220240530.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgsLegacy).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfigLegacy", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
		}
	}
	return legacyAdvConfig, advConfig
}

func (r *rs) applyClusterChanges(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) *admin.ClusterDescription20240805 {
	stateReq := normalizeFromTFModel(ctx, state, diags, false)
	planReq := normalizeFromTFModel(ctx, plan, diags, false)
	if diags.HasError() {
		return nil
	}
	normalizePatchState(stateReq)
	patchReq, err := update.PatchPayload(stateReq, planReq)
	if err != nil {
		diags.AddError("errorPatchPayload", err.Error())
		return nil
	}
	var cluster *admin.ClusterDescription20240805
	if usingLegacySchema(ctx, plan.ReplicationSpecs, diags) {
		// Only updates of replication specs will be done with legacy API
		legacySpecsChanged := r.updateLegacyReplicationSpecs(ctx, state, plan, diags, patchReq.ReplicationSpecs)
		if diags.HasError() {
			return nil
		}
		patchReq.ReplicationSpecs = nil // already updated by legacy API
		if legacySpecsChanged && update.IsEmpty(patchReq) {
			return AwaitChanges(ctx, r.Client.AtlasV2.ClustersApi, &plan.Timeouts, diags, plan.ProjectID.ValueString(), plan.Name.ValueString(), changeReasonUpdate)
		}
	}
	if update.IsEmpty(patchReq) {
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
			_ = r.updateAndWait(ctx, &resumeRequest, diags, plan)
		}
	}
	if diags.HasError() {
		return nil
	}
	replicationSpecsUpdated := patchReq.ReplicationSpecs != nil
	if replicationSpecsUpdated {
		// cannot call latest API (2024-10-23 or newer) as it can enable ISS autoscaling
		legacyPatch := newLegacyModel(patchReq)
		cluster = r.updateAndWaitLegacy(ctx, legacyPatch, diags, plan)
	} else {
		cluster = r.updateAndWait(ctx, patchReq, diags, plan)
	}
	if pauseAfter && cluster != nil {
		cluster = r.updateAndWait(ctx, &pauseRequest, diags, plan)
	}
	return cluster
}

func (r *rs) updateLegacyReplicationSpecs(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics, specChanges *[]admin.ReplicationSpec20240805) bool {
	numShardsUpdates := findNumShardsUpdates(ctx, state, plan, diags)
	if diags.HasError() {
		return false
	}
	if specChanges == nil && numShardsUpdates == nil { // no changes to replication specs
		return false
	}
	if specChanges == nil {
		// use state replication specs as there are no changes in paln except for numShardsupdates
		specChanges = newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
		if diags.HasError() {
			return false
		}
	}
	numShardsPlan := numShardsMap(ctx, plan.ReplicationSpecs, diags)
	if diags.HasError() {
		return false
	}
	legacyPatch := newLegacyModel20240530OnlyReplicationSpecsAndDiskGB(specChanges, numShardsPlan, state.DiskSizeGB.ValueFloat64Pointer())
	api20240530 := r.Client.AtlasV220240530.ClustersApi
	api20240530.UpdateCluster(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), legacyPatch)
	_, _, err := api20240530.UpdateCluster(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), legacyPatch).Execute()
	if err != nil {
		diags.AddError("errorUpdateLegacy", fmt.Sprintf(errorUpdate, plan.Name.ValueString(), err.Error()))
		return false
	}
	return true
}

func (r *rs) updateAndWait(ctx context.Context, patchReq *admin.ClusterDescription20240805, diags *diag.Diagnostics, tfModel *TFModel) *admin.ClusterDescription20240805 {
	api := r.Client.AtlasV2.ClustersApi
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.Name.ValueString()
	_, _, err := api.UpdateCluster(ctx, projectID, clusterName, patchReq).Execute()
	if err != nil {
		diags.AddError("errorUpdate", fmt.Sprintf(errorUpdate, clusterName, err.Error()))
		return nil
	}
	return AwaitChanges(ctx, r.Client.AtlasV2.ClustersApi, &tfModel.Timeouts, diags, projectID, clusterName, changeReasonUpdate)
}
func (r *rs) updateAndWaitLegacy(ctx context.Context, patchReq *admin20240805.ClusterDescription20240805, diags *diag.Diagnostics, tfModel *TFModel) *admin.ClusterDescription20240805 {
	api := r.Client.AtlasV220240805.ClustersApi
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.Name.ValueString()
	_, _, err := api.UpdateCluster(ctx, projectID, clusterName, patchReq).Execute()
	if err != nil {
		diags.AddError("errorUpdateLegacy", fmt.Sprintf(errorUpdate, clusterName, err.Error()))
		return nil
	}
	return AwaitChanges(ctx, r.Client.AtlasV2.ClustersApi, &tfModel.Timeouts, diags, projectID, clusterName, changeReasonUpdate)
}

func (r *rs) convertClusterAddAdvConfig(ctx context.Context, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin.ClusterDescriptionProcessArgs20240805, cluster *admin.ClusterDescription20240805, modelIn *TFModel, diags *diag.Diagnostics) *TFModel {
	apiLatest := r.Client.AtlasV2.ClustersApi
	apiLegacy := r.Client.AtlasV220240530.ClustersApi
	projectID := cluster.GetGroupId()
	clusterName := cluster.GetName()
	var err error
	if legacyAdvConfig == nil {
		legacyAdvConfig, _, err = apiLegacy.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError("errorReadAdvConfigLegacy", fmt.Sprintf(errorRead, clusterName, err.Error()))
			return nil
		}
	}
	if advConfig == nil {
		advConfig, _, err = apiLatest.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError("errorReadAdvConfig", fmt.Sprintf(errorRead, clusterName, err.Error()))
			return nil
		}
	}
	legacyInfo := resolveLegacyInfo(ctx, modelIn, diags, cluster, apiLegacy)
	modelOut := NewTFModel(ctx, cluster, modelIn.Timeouts, diags, legacyInfo)
	if diags.HasError() {
		return nil
	}
	AddAdvancedConfig(ctx, modelOut, advConfig, legacyAdvConfig, diags)
	if diags.HasError() {
		return nil
	}
	overrideKnowTPFIssueFields(modelIn, modelOut)
	return modelOut
}
