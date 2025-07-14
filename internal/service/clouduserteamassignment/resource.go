package clouduserteamassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const resourceName = "cloud_user_team_assignment"

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
	// TODO: Schema and model must be defined in resource_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var userTeamAssignment TFUserTeamAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &userTeamAssignment)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: make POST request to Atlas API and handle error in response

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error creating resource", err.Error())
	//	return
	//}
	connV2 := r.Client.AtlasV2
	orgID := userTeamAssignment.OrgID.ValueString()
	teamID := userTeamAssignment.TeamID.ValueString()
	cloudUserTeamAssignmentReq, diags := NewCloudUserTeamAssignmentReq(ctx, &userTeamAssignment)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResp, _, err := connV2.MongoDBCloudUsersApi.AddUserToTeam(ctx, orgID, teamID, cloudUserTeamAssignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	// TODO: process response into new terraform state
	newCloudUserTeamAssignmentModel, diags := NewTFUserTeamAssignmentModel(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserTeamAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var cloudUserTeamAssignmentState TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &cloudUserTeamAssignmentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	if validate.StatusNotFound(apiResp) {
	//		resp.State.RemoveResource(ctx)
	//		return
	//	}
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	newCloudUserTeamAssignmentModel, diags := NewTFModel(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserTeamAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudUserTeamAssignmentReq, diags := NewAtlasReq(ctx, &tfModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// TODO: make PATCH request to Atlas API and handle error in response
	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error updating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state

	newCloudUserTeamAssignmentModel, diags := NewTFModel(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserTeamAssignmentModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var cloudUserTeamAssignmentState *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &cloudUserTeamAssignmentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make Delete request to Atlas API

	// connV2 := r.Client.AtlasV2
	// if _, _, err := connV2.Api.Delete().Execute(); err != nil {
	// 	 resp.Diagnostics.AddError("error deleting resource", err.Error())
	// 	 return
	// }
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: parse req.ID string taking into account documented format. Example:

	// projectID, other, err := splitCloudUserTeamAssignmentImportID(req.ID)
	// if err != nil {
	//	resp.Diagnostics.AddError("error splitting import ID", err.Error())
	//	return
	//}

	// TODO: define attributes that are required for read operation to work correctly. Example:

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
