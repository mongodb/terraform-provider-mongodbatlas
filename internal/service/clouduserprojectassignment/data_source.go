package clouduserprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
	userResp, err := fetchProjectUser(ctx, connV2, projectID, userID, username)
	if err != nil {
		resp.Diagnostics.AddError(errorReadingUser, err.Error())
		return
	}
	if userResp == nil {
		resp.Diagnostics.AddError("resource not found", "no user found with the specified identifier")
		return
	}

	newCloudUserProjectAssignmentModel, diags := NewTFModel(ctx, projectID, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserProjectAssignmentModel)...)
}
