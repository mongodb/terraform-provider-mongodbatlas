package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241023002/admin"
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
)

var (
	DeprecationMsgOldSchema = fmt.Sprintf("%s %s", constant.DeprecationParam, DeprecationOldSchemaAction)
	pauseRequest            = admin20240805.ClusterDescription20240805{Paused: conversion.Pointer(true)}
	resumeRequest           = admin20240805.ClusterDescription20240805{Paused: conversion.Pointer(false)}
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
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	diags := &resp.Diagnostics
	if diags.HasError() {
		return
	}
	summary, detail := r.createCluster(ctx, &plan, diags, resp)
	if summary != "" || detail != "" {
		diags.AddError(summary, detail)
	}
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	r.readCluster(ctx, &state, &resp.State, diags)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	api := r.Client.AtlasV220240805.ClustersApi
	patchReq := update.PatchPayloadTpf(ctx, diags, &state, &plan, NewAtlasReq)
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	var cluster *admin20240805.ClusterDescription20240805
	var advConfig *admin20240805.ClusterDescriptionProcessArgs20240805
	var legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
	var err error
	if patchReq != nil {
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
				cluster = r.updateAndWait(ctx, &resumeRequest, diags, &plan)
				if cluster == nil {
					return
				}
			}
		}
		cluster = r.updateAndWait(ctx, patchReq, diags, &plan)
		if cluster == nil {
			return
		}
		if pauseAfter {
			cluster = r.updateAndWait(ctx, &pauseRequest, diags, &plan)
		}
	}
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if patchReqProcessArgs != nil {
		advConfig, _, err = api.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgs).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfig", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
		}
	}
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if patchReqProcessArgsLegacy != nil {
		legacyAdvConfig, _, err = r.Client.AtlasV220240530.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgsLegacy).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfigLegacy", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
		}
	}
	if cluster != nil {
		r.convertClusterAddAdvConfig(ctx, legacyAdvConfig, advConfig, cluster, plan.Timeouts, diags, &resp.State)
	} else {
		r.readCluster(ctx, &plan, &resp.State, diags)
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
	if clusterName == "" {
		resp.Diagnostics.AddError("errorDelete", fmt.Sprintf(errorDelete, clusterName, "clusterName is empty"))
		return
	}
	api := r.Client.AtlasV2.ClustersApi
	_, err := api.DeleteCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		diags.AddError("errorDelete", fmt.Sprintf(errorDelete, clusterName, err.Error()))
	}
	_, summary, detail := r.awaitChanges(ctx, &state.Timeouts, diags, projectID, clusterName, "delete")
	if summary != "" || detail != "" {
		diags.AddError(summary, detail)
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}

func (r *rs) createCluster(ctx context.Context, plan *TFModel, diags *diag.Diagnostics, resp *resource.CreateResponse) (summary, detail string) {
	sdkReq := NewAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return "", ""
	}
	api := r.Client.AtlasV220240805.ClustersApi
	apiLegacy := r.Client.AtlasV220240530.ClustersApi
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	_, _, err := api.CreateCluster(ctx, projectID, sdkReq).Execute()
	if err != nil {
		return "errorCreate", fmt.Sprintf(errorCreate, err.Error())
	}
	cluster, summary, detail := r.awaitChanges(ctx, &plan.Timeouts, diags, projectID, clusterName, "create")
	if summary != "" || detail != "" {
		return summary, detail
	}

	var legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
	legacyAdvConfigUpdate := NewAtlasReqAdvancedConfigurationLegacy(ctx, &plan.AdvancedConfiguration, diags)
	if legacyAdvConfigUpdate != nil {
		legacyAdvConfig, _, err = apiLegacy.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, legacyAdvConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			return "errorUpdateeAdvConfigLegacy", fmt.Sprintf(errorCreate, err.Error())
		}
	}

	advConfigUpdate := NewAtlasReqAdvancedConfiguration(ctx, &plan.AdvancedConfiguration, diags)
	var advConfig *admin20240805.ClusterDescriptionProcessArgs20240805
	if advConfigUpdate != nil {
		advConfig, _, err = api.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, advConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			return "errorUpdateAdvConfig", fmt.Sprintf(errorCreate, err.Error())
		}
	}
	return r.convertClusterAddAdvConfig(ctx, legacyAdvConfig, advConfig, cluster, plan.Timeouts, diags, &resp.State)
}

func (r *rs) readCluster(ctx context.Context, model *TFModel, state *tfsdk.State, diags *diag.Diagnostics) {
	clusterName := model.Name.ValueString()
	projectID := model.ProjectID.ValueString()
	api := r.Client.AtlasV220240805.ClustersApi
	readResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, constant.ErrorCodeClusterNotFound) {
			state.RemoveResource(ctx)
			return
		}
		diags.AddError("errorRead", fmt.Sprintf(errorRead, clusterName, err.Error()))
		return
	}
	summary, detail := r.convertClusterAddAdvConfig(ctx, nil, nil, readResp, model.Timeouts, diags, state)
	if summary != "" || detail != "" {
		diags.AddError(summary, detail)
	}
}
func (r *rs) updateAndWait(ctx context.Context, patchReq *admin20240805.ClusterDescription20240805, diags *diag.Diagnostics, tfModel *TFModel) *admin20240805.ClusterDescription20240805 {
	api := r.Client.AtlasV220240805.ClustersApi
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.Name.ValueString()
	_, _, err := api.UpdateCluster(ctx, projectID, clusterName, patchReq).Execute()
	if err != nil {
		diags.AddError("errorUpdate", fmt.Sprintf(errorUpdate, clusterName, err.Error()))
		return nil
	}
	cluster, summary, detail := r.awaitChanges(ctx, &tfModel.Timeouts, diags, projectID, clusterName, "update")
	if summary != "" || detail != "" {
		diags.AddError(summary, detail)
		return nil
	}
	return cluster
}

func (r *rs) awaitChanges(ctx context.Context, t *timeouts.Value, diags *diag.Diagnostics, projectID, clusterName, changeReason string) (cluster *admin20240805.ClusterDescription20240805, summary, detail string) {
	var timeoutDuration time.Duration
	var localDiags diag.Diagnostics
	var targetState = retrystrategy.RetryStrategyIdleState
	var extraPending = []string{}
	switch changeReason {
	case "create":
		timeoutDuration, localDiags = t.Create(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case "update":
		timeoutDuration, localDiags = t.Update(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case "delete":
		timeoutDuration, localDiags = t.Delete(ctx, defaultTimeout)
		diags.Append(localDiags...)
		targetState = retrystrategy.RetryStrategyDeletedState
		extraPending = append(extraPending, retrystrategy.RetryStrategyIdleState)
	default:
		return nil, "errorReadingTimeout", "unknown change reason" + changeReason
	}
	if diags.HasError() {
		return nil, "errorReadingTimeout", ""
	}
	stateConf := CreateStateChangeConfig(ctx, r.Client.AtlasV220240805, projectID, clusterName, targetState, timeoutDuration, extraPending...)
	clusterAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if admin.IsErrorCode(err, constant.ErrorCodeClusterNotFound) && changeReason == "delete" {
			return nil, "", ""
		}
		return nil, "errorAwaitingCluster", fmt.Sprintf(errorCreate, err)
	}
	if targetState == retrystrategy.RetryStrategyDeletedState {
		return nil, "", ""
	}
	cluster, ok := clusterAny.(*admin20240805.ClusterDescription20240805)
	if !ok {
		return nil, "errorAwaitingCluster", fmt.Sprintf(errorCreate, "unexpected type from WaitForStateContext")
	}
	return cluster, "", ""
}

func (r *rs) convertClusterAddAdvConfig(ctx context.Context, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin20240805.ClusterDescriptionProcessArgs20240805, cluster *admin20240805.ClusterDescription20240805, resourceTimeouts timeouts.Value, diags *diag.Diagnostics, state *tfsdk.State) (summary, detail string) {
	api := r.Client.AtlasV220240805.ClustersApi
	apiLegacy := r.Client.AtlasV220240530.ClustersApi
	projectID := cluster.GetGroupId()
	clusterName := cluster.GetName()
	var err error
	if legacyAdvConfig == nil {
		legacyAdvConfig, _, err = apiLegacy.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			return "errorReadAdvConfigLegacy", fmt.Sprintf(errorRead, clusterName, err.Error())
		}
	}
	if advConfig == nil {
		advConfig, _, err = api.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			return "errorReadAdvConfig", fmt.Sprintf(errorRead, clusterName, err.Error())
		}
	}
	model := NewTFModel(ctx, cluster, resourceTimeouts, diags)
	if diags.HasError() {
		return "", ""
	}
	AddAdvancedConfig(ctx, model, advConfig, legacyAdvConfig, diags)
	if diags.HasError() {
		return "", ""
	}
	diags.Append(state.Set(ctx, model)...)
	return "", ""
}
