package clouduserprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

var _ datasource.DataSource = &cloudUserProjectAssignmentDS{}
var _ datasource.DataSourceWithConfigure = &cloudUserProjectAssignmentDS{}

func DataSource() datasource.DataSource {
	return &cloudUserProjectAssignmentDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type cloudUserProjectAssignmentDS struct {
	config.DSCommon
}

func (d *cloudUserProjectAssignmentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *cloudUserProjectAssignmentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	userID := state.UserId.ValueString()
	username := state.Username.ValueString()

	if username == "" && userID == "" {
		resp.Diagnostics.AddError("invalid configuration", "either username or user_id must be provided")
		return
	}

	var userResp *admin.GroupUserResponse
	var err error

	if userID != "" {
		userResp, _, err = connV2.MongoDBCloudUsersApi.GetProjectUser(ctx, projectID, userID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(errorReadingByUserID, err.Error())
			return
		}
	} else if username != "" {
		params := &admin.ListProjectUsersApiParams{
			GroupId:  projectID,
			Username: &username,
		}
		usersResp, _, err := connV2.MongoDBCloudUsersApi.ListProjectUsersWithParams(ctx, params).Execute()
		if err != nil {
			resp.Diagnostics.AddError(errorReadingByUsername, err.Error())
			return
		}
		if usersResp == nil || len(usersResp.GetResults()) == 0 {
			resp.Diagnostics.AddError("resource not found", "no user found with username "+username)
			return
		}
		userResp = &usersResp.GetResults()[0]
	}
	if userResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newCloudUserProjectAssignmentModel, diags := NewTFModel(ctx, projectID, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserProjectAssignmentModel)...)
}
