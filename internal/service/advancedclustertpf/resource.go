package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
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
)

var DeprecationMsgOldSchema = fmt.Sprintf("%s %s", constant.DeprecationParam, DeprecationOldSchemaAction)

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

func (r *rs) createCluster(ctx context.Context, plan *TFModel, diags *diag.Diagnostics, resp *resource.CreateResponse) (summary, detail string) {
	sdkReq := NewAtlasReq(ctx, plan, diags)
	if diags.HasError() {
		return "", ""
	}
	api := r.Client.AtlasV2.ClustersApi
	apiLegacy := r.Client.AtlasV220240530.ClustersApi
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	cluster, _, err := api.CreateCluster(ctx, projectID, sdkReq).Execute()
	if err != nil {
		return "errorCreate", fmt.Sprintf(errorCreate, err.Error())
	}
	var legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
	legacyAdvConfigUpdate := NewAtlasReqAdvancedConfigurationLegacy(ctx, &plan.AdvancedConfiguration, diags)
	if legacyAdvConfigUpdate != nil {
		legacyAdvConfig, _, err = apiLegacy.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, legacyAdvConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			return "errorCreateAdvConfigLegacy", fmt.Sprintf(errorCreate, err.Error())
		}
	}

	advConfigUpdate := NewAtlasReqAdvancedConfiguration(ctx, &plan.AdvancedConfiguration, diags)
	var advConfig *admin.ClusterDescriptionProcessArgs20240805
	if advConfigUpdate != nil {
		advConfig, _, err = api.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, advConfigUpdate).Execute()
		if err != nil {
			// Maybe should be warning instead of error to avoid having to re-create the cluster
			return "errorCreateAdvConfig", fmt.Sprintf(errorCreate, err.Error())
		}
	}
	return r.convertClusterAddAdvConfig(ctx, legacyAdvConfig, advConfig, cluster, plan.Timeouts, diags, &resp.State)
}

func (r *rs) convertClusterAddAdvConfig(ctx context.Context, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin.ClusterDescriptionProcessArgs20240805, cluster *admin.ClusterDescription20240805, resourceTimeouts timeouts.Value, diags *diag.Diagnostics, state *tfsdk.State) (summary, detail string) {
	api := r.Client.AtlasV2.ClustersApi
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
	if !diags.HasError() {
		diags.Append(state.Set(ctx, model)...)
	}
	return "", ""
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

func (r *rs) readCluster(ctx context.Context, model *TFModel, state *tfsdk.State, diags *diag.Diagnostics) {
	clusterName := model.Name.ValueString()
	projectID := model.ProjectID.ValueString()
	api := r.Client.AtlasV2.ClustersApi
	readResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "CLUSTER_NOT_FOUND") {
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

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	api := r.Client.AtlasV2.ClustersApi
	patchReq := update.PatchPayloadTpf(ctx, diags, &state, &plan, NewAtlasReq)
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.Name.ValueString()
	if patchReq != nil {
		_, _, err := api.UpdateCluster(ctx, projectID, clusterName, patchReq).Execute()
		if err != nil {
			diags.AddError("errorUpdate", fmt.Sprintf(errorUpdate, clusterName, err.Error()))
			return
		}
	}
	patchReqProcessArgs := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if patchReqProcessArgs != nil {
		_, _, err := api.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgs).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfig", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
		}
	}
	patchReqProcessArgsLegacy := update.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if patchReqProcessArgsLegacy != nil {
		_, _, err := r.Client.AtlasV220240530.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, patchReqProcessArgsLegacy).Execute()
		if err != nil {
			diags.AddError("errorUpdateAdvancedConfigLegacy", fmt.Sprintf(errorConfigUpdate, clusterName, err.Error()))
		}
	}
	r.readCluster(ctx, &plan, &resp.State, diags)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	name := state.Name.ValueString()
	clusterName := state.Name.ValueString()
	if clusterName == "" {
		resp.Diagnostics.AddError("errorDelete", fmt.Sprintf(errorDelete, name, "clusterName is empty"))
		return
	}
	api := r.Client.AtlasV2.ClustersApi
	_, err := api.DeleteCluster(ctx, state.ProjectID.ValueString(), name).Execute()
	if err != nil {
		diags.AddError("errorDelete", fmt.Sprintf(errorDelete, name, err.Error()))
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}
