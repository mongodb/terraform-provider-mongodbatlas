package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
	if resp.Diagnostics.HasError() {
		return
	}
	sdkReq := NewAtlasReq(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	err := StoreCreatePayload(sdkReq)
	if err != nil {
		resp.Diagnostics.AddError("errorCreate", fmt.Sprintf(errorCreate, err.Error()))
		return
	}
	tfNewModel, shouldReturn := mockedSDK(ctx, &resp.Diagnostics, plan.Timeouts)
	if shouldReturn {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfNewModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ClusterID.IsNull() {
		tfModel, shouldReturn := mockedSDK(ctx, &resp.Diagnostics, state.Timeouts)
		if shouldReturn {
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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
	patchReq := conversion.PatchPayloadTpf(ctx, diags, &state, &plan, NewAtlasReq)
	if patchReq != nil {
		err := StoreUpdatePayload(patchReq)
		if err != nil {
			diags.AddError("error storing update payload", fmt.Sprintf("error storing update payload: %s", err.Error()))
		}
	}
	patchReqProcessArgs := conversion.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfiguration)
	if patchReqProcessArgs != nil {
		err := StoreUpdatePayloadProcessArgs(patchReqProcessArgs)
		if err != nil {
			diags.AddError("error storing update payload advanced config", fmt.Sprintf("error storing update payload: %s", err.Error()))
		}
	}
	patchReqProcessArgsLegacy := conversion.PatchPayloadTpf(ctx, diags, &state.AdvancedConfiguration, &plan.AdvancedConfiguration, NewAtlasReqAdvancedConfigurationLegacy)
	if patchReqProcessArgsLegacy != nil {
		err := StoreUpdatePayloadProcessArgsLegacy(patchReqProcessArgsLegacy)
		if err != nil {
			diags.AddError("error storing update payload advanced config legacy", fmt.Sprintf("error storing update payload: %s", err.Error()))
		}
	}
	if diags.HasError() {
		return
	}

	// TODO: Use requests to do actual updates with Admin API
	tfNewModel, shouldReturn := mockedSDK(ctx, diags, plan.Timeouts)
	// TODO: keep project_id and name from plan to avoid overwriting for move_state tests. We should probably do the same with the rest of attributes
	tfNewModel.Name = plan.Name
	tfNewModel.ProjectID = plan.ProjectID
	if shouldReturn {
		return
	}
	diags.Append(resp.State.Set(ctx, tfNewModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := state.Name.ValueString()
	clusterID := state.ClusterID.ValueString()
	if clusterID == "" {
		resp.Diagnostics.AddError("errorDelete", fmt.Sprintf(errorDelete, name, "clusterID is empty"))
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "name")
}

func mockedSDK(ctx context.Context, diags *diag.Diagnostics, timeout timeouts.Value) (*TFModel, bool) {
	sdkResp, err := ReadClusterResponse()
	if err != nil {
		diags.AddError("errorCreate", fmt.Sprintf(errorCreate, err.Error()))
		return nil, true
	}
	tfNewModel := NewTFModel(ctx, sdkResp, timeout, diags)
	sdkAdvConfig, err := ReadClusterProcessArgsResponse()
	if err != nil {
		diags.AddError("errorCreateAdvConfig", fmt.Sprintf(errorCreate, err.Error()))
		return nil, true
	}
	sdkAdvConfigLegacy, err := ReadClusterProcessArgsResponseLegacy()
	if err != nil {
		diags.AddError("errorCreateAdvConfigLegacy", fmt.Sprintf(errorCreate, err.Error()))
		return nil, true
	}
	AddAdvancedConfig(ctx, tfNewModel, sdkAdvConfig, sdkAdvConfigLegacy, diags)
	if diags.HasError() {
		return nil, true
	}
	return tfNewModel, false
}
