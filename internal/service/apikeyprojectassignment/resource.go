package apikeyprojectassignment

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

const resourceName = "api_key_project_assignment"

var (
	_ resource.ResourceWithConfigure   = &rs{}
	_ resource.ResourceWithImportState = &rs{}
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
	var tfModel TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignmentReq, diags := NewAtlasCreateReq(ctx, &tfModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	connV2 := r.Client.AtlasV2
	_, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, tfModel.ProjectId.ValueString(), tfModel.ApiKeyId.ValueString(), assignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var assignmentState TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &assignmentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := assignmentState.ProjectId.ValueString()
	apiKeys, apiResp, err := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeys(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	apiKeyID := assignmentState.ApiKeyId.ValueString()
	newAPIKeyProjectAssignmentModel, diags := NewTFModel(ctx, apiKeys, apiKeyID, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newAPIKeyProjectAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignmentReq, diags := NewAtlasUpdateReq(ctx, &tfModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectId.ValueString()
	apiKeyID := tfModel.ApiKeyId.ValueString()

	apiResp, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, assignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newAssignmentModel, diags := NewTFModelPatch(ctx, apiResp, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newAssignmentModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var assignmentState *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &assignmentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2

	projectID := assignmentState.ProjectId.ValueString()
	apiKeyID := assignmentState.ApiKeyId.ValueString()
	if _, err := connV2.ProgrammaticAPIKeysApi.RemoveProjectApiKey(ctx, projectID, apiKeyID).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, apiKeyID, err := splitAPIKeyProjectAssignmentImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("api_key_id"), apiKeyID)...)
}

func splitAPIKeyProjectAssignmentImportID(id string) (projectID, apiKeyID string, err error) {
	re := regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{api_key_id}")
		return
	}

	projectID = parts[1]
	apiKeyID = parts[2]
	return
}
