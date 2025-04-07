// Code generated by terraform-provider-mongodbatlas using `make generate-resource`. DO NOT EDIT.

package databaseuserapi

import (
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

const apiVersionHeader = "application/vnd.atlas.2023-01-01+json"

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: "database_user_api",
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

	reqBody, err := autogeneration.Marshal(&plan, false)
	if err != nil {
		resp.Diagnostics.AddError("error during create operation", err.Error())
		return
	}

	pathParams := map[string]string{
		"groupId": plan.GroupId.ValueString(),
	}
	apiResp, err := r.Client.UntypedAPICall(ctx, &config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  "/api/atlas/v2/groups/{groupId}/databaseUsers",
		PathParams:    pathParams,
		Method:        http.MethodPost,
		Body:          reqBody,
	})

	if err != nil {
		resp.Diagnostics.AddError("error during create operation", err.Error())
		return
	}

	respBody, err := io.ReadAll(apiResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("error during create operation", err.Error())
		return
	}

	// Use the plan as the base model to set the response state
	if err := autogeneration.Unmarshal(respBody, &plan); err != nil {
		resp.Diagnostics.AddError("error during create operation", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathParams := map[string]string{
		"groupId":      state.GroupId.ValueString(),
		"databaseName": state.DatabaseName.ValueString(),
		"username":     state.Username.ValueString(),
	}
	apiResp, err := r.Client.UntypedAPICall(ctx, &config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  "/api/atlas/v2/groups/{groupId}/databaseUsers/{databaseName}/{username}",
		PathParams:    pathParams,
		Method:        http.MethodGet,
	})

	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during get operation", err.Error())
		return
	}

	respBody, err := io.ReadAll(apiResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("error during get operation", err.Error())
		return
	}

	// Use the current state as the base model to set the response state
	if err := autogeneration.Unmarshal(respBody, &state); err != nil {
		resp.Diagnostics.AddError("error during get operation", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: code generation logic for update will be handled in milestone 2
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathParams := map[string]string{
		"groupId":      state.GroupId.ValueString(),
		"databaseName": state.DatabaseName.ValueString(),
		"username":     state.Username.ValueString(),
	}
	if _, err := r.Client.UntypedAPICall(ctx, &config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  "/api/atlas/v2/groups/{groupId}/databaseUsers/{databaseName}/{username}",
		PathParams:    pathParams,
		Method:        http.MethodDelete,
	}); err != nil {
		resp.Diagnostics.AddError("error during delete", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: code generation logic for import will be handled in milestone 2
}
