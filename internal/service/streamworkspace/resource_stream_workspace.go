package streamworkspace

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

var _ resource.ResourceWithConfigure = &streamWorkspaceRS{}
var _ resource.ResourceWithImportState = &streamWorkspaceRS{}

const streamWorkspaceName = "stream_workspace"

func Resource() resource.Resource {
	return &streamWorkspaceRS{
		RSCommon: config.RSCommon{
			ResourceName: streamWorkspaceName,
		},
	}
}

type streamWorkspaceRS struct {
	config.RSCommon
}

func (r *streamWorkspaceRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *streamWorkspaceRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamWorkspacePlan TFStreamWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamWorkspacePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamWorkspacePlan.ProjectID.ValueString()
	streamWorkspaceReq, diags := NewStreamWorkspaceCreateReq(ctx, &streamWorkspacePlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.CreateStreamInstance(ctx, projectID, streamWorkspaceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newStreamWorkspaceModel, diags := NewTFStreamWorkspace(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamWorkspaceModel)...)
}

func (r *streamWorkspaceRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamWorkspaceState TFStreamWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamWorkspaceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamWorkspaceState.ProjectID.ValueString()
	workspaceName := streamWorkspaceState.WorkspaceName.ValueString()
	apiResp, getResp, err := connV2.StreamsApi.GetStreamInstance(ctx, projectID, workspaceName).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamWorkspaceModel, diags := NewTFStreamWorkspace(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamWorkspaceModel)...)
}

func (r *streamWorkspaceRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var streamWorkspacePlan TFStreamWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamWorkspacePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamWorkspacePlan.ProjectID.ValueString()
	workspaceName := streamWorkspacePlan.WorkspaceName.ValueString()
	streamWorkspaceReq, diags := NewStreamWorkspaceUpdateReq(ctx, &streamWorkspacePlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.UpdateStreamInstance(ctx, projectID, workspaceName, streamWorkspaceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newStreamWorkspaceModel, diags := NewTFStreamWorkspace(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamWorkspaceModel)...)
}

func (r *streamWorkspaceRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamWorkspaceState *TFStreamWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamWorkspaceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamWorkspaceState.ProjectID.ValueString()
	workspaceName := streamWorkspaceState.WorkspaceName.ValueString()
	if _, err := connV2.StreamsApi.DeleteStreamInstance(ctx, projectID, workspaceName).Execute(); err != nil {
		resp.Diagnostics.AddError("error during resource delete", err.Error())
		return
	}
}

func (r *streamWorkspaceRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, workspaceName, err := splitStreamWorkspaceImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search deployment import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_name"), workspaceName)...)
}

func splitStreamWorkspaceImportID(id string) (projectID, workspaceName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{workspace_name}")
		return
	}

	projectID = parts[1]
	workspaceName = parts[2]
	return
}
