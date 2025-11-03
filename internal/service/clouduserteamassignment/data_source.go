package clouduserteamassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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

	userResp, err := fetchTeamUser(ctx, connV2, orgID, teamID, &userID, &username)
	if err != nil {
		resp.Diagnostics.AddError("error retrieving user", err.Error())
		return
	}
	if userResp == nil {
		resp.Diagnostics.AddError("resource not found", "no user found with the specified identifier")
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
