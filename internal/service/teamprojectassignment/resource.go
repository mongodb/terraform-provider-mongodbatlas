package teamprojectassignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName          = "team_project_assignment"
	errorFetchingResource = "error fetching resource"
	invalidImportID       = "invalid import ID format"
	errorAssigment        = "error assigning Team to ProjectID (%s):"
	errorUpdate           = "error updating TeamID(%s) in ProjectID(%s):"
	errorDelete           = "error deleting TeamID(%s) from ProjectID(%s):"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

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
	resp.Schema = resourceSchema()
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectId.ValueString()
	teamID := plan.TeamId.ValueString()
	teamProjectReq, diags := NewAtlasReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, _, err := connV2.TeamsApi.AddAllTeamsToProject(ctx, projectID, teamProjectReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorAssigment, projectID), err.Error())
		return
	}

	apiResp, _, err := connV2.TeamsApi.GetProjectTeam(ctx, projectID, teamID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorAssigment, projectID), err.Error())
		return
	}
	newTeamProjectAssignmentModel, diags := NewTFModel(ctx, apiResp, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTeamProjectAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	teamID := state.TeamId.ValueString()

	apiResp, httpResp, err := connV2.TeamsApi.GetProjectTeam(ctx, projectID, teamID).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errorFetchingResource, err.Error())
		return
	}

	newTeamProjectAssignmentModel, diags := NewTFModel(ctx, apiResp, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTeamProjectAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectId.ValueString()
	teamID := plan.TeamId.ValueString()

	updateReq, diags := NewAtlasUpdateReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, _, err := connV2.TeamsApi.UpdateTeamRoles(ctx, projectID, teamID, updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorUpdate, teamID, projectID), "API response is nil")
		return
	}

	apiResp, _, err := connV2.TeamsApi.GetProjectTeam(ctx, projectID, teamID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorUpdate, teamID, projectID), "API response is nil")
		return
	}

	newTeamProjectAssignmentModel, diags := NewTFModel(ctx, apiResp, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTeamProjectAssignmentModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	teamID := state.TeamId.ValueString()

	httpResp, err := connV2.TeamsApi.RemoveProjectTeam(ctx, projectID, teamID).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf(errorDelete, teamID, projectID), err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	ok, parts := conversion.ImportSplit(req.ID, 2)
	if !ok {
		resp.Diagnostics.AddError(invalidImportID, "expected 'project_id/team_id', got: "+importID)
		return
	}
	projectID, teamID := parts[0], parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("team_id"), teamID)...)
}
