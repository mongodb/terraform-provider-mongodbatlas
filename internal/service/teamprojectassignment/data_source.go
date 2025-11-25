package teamprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &teamProjectAssignmentDS{}
var _ datasource.DataSourceWithConfigure = &teamProjectAssignmentDS{}

func DataSource() datasource.DataSource {
	return &teamProjectAssignmentDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type teamProjectAssignmentDS struct {
	config.DSCommon
}

func (d *teamProjectAssignmentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *teamProjectAssignmentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	teamID := state.TeamId.ValueString()

	apiResp, httpResp, err := connV2.TeamsApi.GetGroupTeam(ctx, projectID, teamID).Execute()
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
