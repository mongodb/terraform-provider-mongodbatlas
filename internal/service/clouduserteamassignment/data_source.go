package clouduserteamassignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

var _ datasource.DataSource = &cloudUserTeamAssignmentDS{}
var _ datasource.DataSourceWithConfigure = &cloudUserTeamAssignmentDS{}

func DataSource() datasource.DataSource {
	return &cloudUserTeamAssignmentDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type cloudUserTeamAssignmentDS struct {
	config.DSCommon
}

func (d *cloudUserTeamAssignmentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *cloudUserTeamAssignmentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFUserTeamAssignmentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	orgID := state.OrgId.ValueString()
	teamID := state.TeamId.ValueString()
	userID := state.UserId.ValueString()
	username := state.Username.ValueString()

	if username == "" && userID == "" {
		resp.Diagnostics.AddError("invalid configuration", "either username or user_id must be provided")
		return
	}

	var userListResp *admin.PaginatedOrgUser
	var userResp *admin.OrgUserResponse
	var err error

	if userID != "" {
		params := &admin.ListTeamUsersApiParams{
			UserId: &userID,
			OrgId:  orgID,
			TeamId: teamID,
		}
		userListResp, _, err = connV2.MongoDBCloudUsersApi.ListTeamUsersWithParams(ctx, params).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error retrieving resource by user_id: %s", userID), err.Error())
			return
		}

		if userListResp == nil || len(userListResp.GetResults()) == 0 {
			resp.Diagnostics.AddError("resource not found", "no user found with the specified user_id")
			return
		}
		userResp = &(userListResp.GetResults())[0]
	} else if username != "" {
		params := &admin.ListTeamUsersApiParams{
			OrgId:    orgID,
			TeamId:   teamID,
			Username: &username,
		}
		userListResp, _, err = connV2.MongoDBCloudUsersApi.ListTeamUsersWithParams(ctx, params).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error retrieving resource by username: %s", username), err.Error())
			return
		}

		if userListResp == nil || len(userListResp.GetResults()) == 0 {
			resp.Diagnostics.AddError("resource not found", "no user found with the specified username")
			return
		}

		userResp = &(userListResp.GetResults())[0]
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
