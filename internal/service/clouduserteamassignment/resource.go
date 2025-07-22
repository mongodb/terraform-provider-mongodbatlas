package clouduserteamassignment

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
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
	resp.Schema = resourceSchema()
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFUserTeamAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	orgID := plan.OrgId.ValueString()
	teamID := plan.TeamId.ValueString()
	cloudUserTeamAssignmentReq, diags := NewUserTeamAssignmentReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResp, _, err := connV2.MongoDBCloudUsersApi.AddUserToTeam(ctx, orgID, teamID, cloudUserTeamAssignmentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error assigning user to TeamID(%s):", teamID), err.Error())
		return
	}

	newUserTeamAssignmentModel, diags := NewTFUserTeamAssignmentModel(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newUserTeamAssignmentModel.OrgId = plan.OrgId
	newUserTeamAssignmentModel.TeamId = plan.TeamId
	resp.Diagnostics.Append(resp.State.Set(ctx, newUserTeamAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFUserTeamAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	orgID := state.OrgId.ValueString()
	teamID := state.TeamId.ValueString()

	var userListResp *admin.PaginatedOrgUser
	var httpResp *http.Response
	var err error

	var userResp *admin.OrgUserResponse
	if !state.UserId.IsNull() && state.UserId.ValueString() != "" {
		userID := state.UserId.ValueString()
		userListResp, httpResp, err = connV2.MongoDBCloudUsersApi.ListTeamUsers(ctx, orgID, teamID).Execute()

		if err != nil {
			if validate.StatusNotFound(httpResp) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Error getting team users by user_id", err.Error())
			return
		}
		if userListResp != nil {
			if len(userListResp.GetResults()) == 0 {
				resp.State.RemoveResource(ctx)
				return
			}
			results := userListResp.GetResults()
			for i := range results {
				if results[i].GetId() == userID {
					userResp = &results[i]
					break
				}
			}
		}
	} else if !state.Username.IsNull() && state.Username.ValueString() != "" { // required for import
		username := state.Username.ValueString()
		params := &admin.ListTeamUsersApiParams{
			Username: &username,
			OrgId:    orgID,
			TeamId:   teamID,
		}
		userListResp, httpResp, err = connV2.MongoDBCloudUsersApi.ListTeamUsersWithParams(ctx, params).Execute()

		if err != nil {
			if validate.StatusNotFound(httpResp) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Error getting team users by username", err.Error())
			return
		}
		if userListResp != nil {
			if len(userListResp.GetResults()) == 0 {
				resp.State.RemoveResource(ctx)
				return
			}
			userResp = &(userListResp.GetResults())[0]
		}
	}

	if userResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newCloudUserTeamAssignmentModel, diags := NewTFUserTeamAssignmentModel(ctx, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newCloudUserTeamAssignmentModel.OrgId = state.OrgId
	newCloudUserTeamAssignmentModel.TeamId = state.TeamId
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserTeamAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TFUserTeamAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	orgID := state.OrgId.ValueString()
	userID := state.UserId.ValueString()
	teamID := state.TeamId.ValueString()

	userInfo := &admin.AddOrRemoveUserFromTeam{
		Id: userID,
	}

	_, httpResp, err := connV2.MongoDBCloudUsersApi.RemoveUserFromTeam(ctx, orgID, teamID, userInfo).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error deleting user(%s) from TeamID(%s):", userID, teamID), err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	ok, parts := conversion.ImportSplit(req.ID, 3)
	if !ok {
		resp.Diagnostics.AddError("invalid import ID format", "expected 'org_id/team_id/user_id' or 'org_id/team_id/username', got: "+importID)
		return
	}
	orgID, teamID, user := parts[0], parts[1], parts[2]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("team_id"), teamID)...)

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

	if emailRegex.MatchString(user) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), user)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), user)...)
	}
}
