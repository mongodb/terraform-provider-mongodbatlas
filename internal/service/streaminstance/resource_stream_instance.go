package streaminstance

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &streamInstanceRS{}
var _ resource.ResourceWithImportState = &streamInstanceRS{}

const streamInstanceName = "stream_instance"

func Resource() resource.Resource {
	return &streamInstanceRS{
		RSCommon: config.RSCommon{
			ResourceName: streamInstanceName,
		},
	}
}

type streamInstanceRS struct {
	config.RSCommon
}

func (r *streamInstanceRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *streamInstanceRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamInstancePlan TFStreamInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamInstancePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamInstancePlan.ProjectID.ValueString()
	streamInstanceReq, diags := NewStreamInstanceCreateReq(ctx, &streamInstancePlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.CreateStreamInstance(ctx, projectID, streamInstanceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newStreamInstanceModel, diags := NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstanceModel)...)
}

func (r *streamInstanceRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamInstanceState TFStreamInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamInstanceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamInstanceState.ProjectID.ValueString()
	instanceName := streamInstanceState.InstanceName.ValueString()
	apiResp, getResp, err := connV2.StreamsApi.GetStreamInstance(ctx, projectID, instanceName).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamInstanceModel, diags := NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstanceModel)...)
}

func (r *streamInstanceRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var streamInstancePlan TFStreamInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamInstancePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamInstancePlan.ProjectID.ValueString()
	instanceName := streamInstancePlan.InstanceName.ValueString()
	streamInstanceReq, diags := NewStreamInstanceUpdateReq(ctx, &streamInstancePlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.UpdateStreamInstance(ctx, projectID, instanceName, streamInstanceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newStreamInstanceModel, diags := NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstanceModel)...)
}

func (r *streamInstanceRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamInstanceState *TFStreamInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamInstanceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamInstanceState.ProjectID.ValueString()
	instanceName := streamInstanceState.InstanceName.ValueString()
	if _, err := connV2.StreamsApi.DeleteStreamInstance(ctx, projectID, instanceName).Execute(); err != nil {
		resp.Diagnostics.AddError("error during resource delete", err.Error())
		return
	}
}

func (r *streamInstanceRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, instanceName, err := splitStreamInstanceImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search deployment import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), instanceName)...)
}

func splitStreamInstanceImportID(id string) (projectID, instanceName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{instance_name}")
		return
	}

	projectID = parts[1]
	instanceName = parts[2]
	return
}
