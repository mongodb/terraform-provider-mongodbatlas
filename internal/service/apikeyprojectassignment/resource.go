package apikeyprojectassignment

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
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
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	projectID := tfModel.ProjectId.ValueString()
	apiKeyID := tfModel.ApiKeyId.ValueString()
	connV2 := r.Client.AtlasV2
	_, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, assignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	// Once CLOUDP-328946 is done, we would use the single GET API to fetch the specific API key project assignment
	apiKeys, err := ListAllProjectAPIKeys(ctx, connV2, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newAPIKeyProjectAssignmentModel, diags := NewTFModel(ctx, apiKeys, projectID, apiKeyID)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newAPIKeyProjectAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var assignmentState TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &assignmentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := assignmentState.ProjectId.ValueString()
	// Once CLOUDP-328946 is done, we would use the single GET API to fetch the specific API key project assignment
	apiKeys, err := ListAllProjectAPIKeys(ctx, connV2, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	apiKeyID := assignmentState.ApiKeyId.ValueString()
	newAPIKeyProjectAssignmentModel, diags := NewTFModel(ctx, apiKeys, projectID, apiKeyID)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
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
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectId.ValueString()
	apiKeyID := tfModel.ApiKeyId.ValueString()

	_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, assignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	// Once CLOUDP-328946 is done, we would use the single GET API to fetch the specific API key project assignment
	apiKeys, err := ListAllProjectAPIKeys(ctx, connV2, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newAssignmentModel, diags := NewTFModel(ctx, apiKeys, projectID, apiKeyID)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
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

func ListAllProjectAPIKeys(ctx context.Context, connV2 *admin.APIClient, projectID string) ([]admin.ApiKeyUserDetails, error) {
	apiKeys, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ApiKeyUserDetails], *http.Response, error) {
		request := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeysWithParams(ctx, &admin.ListProjectApiKeysApiParams{
			GroupId: projectID,
		})
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		return nil, err
	}
	return apiKeys, nil
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
	ok, parts := conversion.ImportSplit(id, 2)
	if ok {
		projectID, apiKeyID = parts[0], parts[1]
		err = conversion.ValidateProjectID(projectID)
		return
	}
	err = errors.New("import format error: to import use the format {project_id}/{api_key_id}")
	return
}
