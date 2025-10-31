package streamworkspace

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
)

var _ resource.ResourceWithConfigure = &streamsWorkspaceRS{}
var _ resource.ResourceWithImportState = &streamsWorkspaceRS{}

const streamsWorkspaceName = "stream_workspace"

func Resource() resource.Resource {
	return &streamsWorkspaceRS{
		RSCommon: config.RSCommon{
			ResourceName: streamsWorkspaceName,
		},
	}
}

type streamsWorkspaceRS struct {
	config.RSCommon
}

func (r *streamsWorkspaceRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *streamsWorkspaceRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamsWorkspacePlan TFStreamsWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamsWorkspacePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert workspace model to instance model
	instanceModel := streamsWorkspacePlan.AsInstanceModel()

	connV2 := r.Client.AtlasV2
	projectID := instanceModel.ProjectID.ValueString()
	streamInstanceReq, diags := streaminstance.NewStreamInstanceCreateReq(ctx, instanceModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.CreateStreamWorkspace(ctx, projectID, streamInstanceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newInstanceModel, diags := streaminstance.NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Convert back to workspace model
	var newWorkspaceModel TFStreamsWorkspaceModel
	newWorkspaceModel.FromInstanceModel(newInstanceModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, newWorkspaceModel)...)
}

func (r *streamsWorkspaceRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamsWorkspaceState TFStreamsWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamsWorkspaceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamsWorkspaceState.ProjectID.ValueString()
	workspaceName := streamsWorkspaceState.WorkspaceName.ValueString()
	apiResp, getResp, err := connV2.StreamsApi.GetStreamWorkspace(ctx, projectID, workspaceName).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newInstanceModel, diags := streaminstance.NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Convert back to workspace model
	var newWorkspaceModel TFStreamsWorkspaceModel
	newWorkspaceModel.FromInstanceModel(newInstanceModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, newWorkspaceModel)...)
}

func (r *streamsWorkspaceRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var streamsWorkspacePlan TFStreamsWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamsWorkspacePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert workspace model to instance model
	instanceModel := streamsWorkspacePlan.AsInstanceModel()

	connV2 := r.Client.AtlasV2
	projectID := instanceModel.ProjectID.ValueString()
	workspaceName := instanceModel.InstanceName.ValueString()
	streamInstanceReq, diags := streaminstance.NewStreamInstanceUpdateReq(ctx, instanceModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.UpdateStreamWorkspace(ctx, projectID, workspaceName, streamInstanceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newInstanceModel, diags := streaminstance.NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Convert back to workspace model
	var newWorkspaceModel TFStreamsWorkspaceModel
	newWorkspaceModel.FromInstanceModel(newInstanceModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, newWorkspaceModel)...)
}

func (r *streamsWorkspaceRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamsWorkspaceState *TFStreamsWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamsWorkspaceState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamsWorkspaceState.ProjectID.ValueString()
	workspaceName := streamsWorkspaceState.WorkspaceName.ValueString()
	if _, err := connV2.StreamsApi.DeleteStreamWorkspace(ctx, projectID, workspaceName).Execute(); err != nil {
		resp.Diagnostics.AddError("error during resource delete", err.Error())
		return
	}
}

func (r *streamsWorkspaceRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, workspaceName, err := splitStreamsWorkspaceImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting streams workspace import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_name"), workspaceName)...)
}

func splitStreamsWorkspaceImportID(id string) (projectID, workspaceName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{workspace_name}")
		return
	}

	projectID, workspaceName = parts[1], parts[2]
	return
}

func (r *streamsWorkspaceRS) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configModel TFStreamsWorkspaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !configModel.WorkspaceName.IsNull() && !configModel.WorkspaceName.IsUnknown() {
		workspaceName := configModel.WorkspaceName.ValueString()
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9-]*$`, workspaceName); !matched {
			resp.Diagnostics.AddAttributeError(
				path.Root("workspace_name"),
				"Invalid workspace name",
				"Workspace name must contain only alphanumeric characters and hyphens, and cannot start with a hyphen.",
			)
		}
	}
}
