package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
)

var _ datasource.DataSource = &streamsWorkspaceDS{}
var _ datasource.DataSourceWithConfigure = &streamsWorkspaceDS{}

func DataSource() datasource.DataSource {
	return &streamsWorkspaceDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamsWorkspaceName,
		},
	}
}

type streamsWorkspaceDS struct {
	config.DSCommon
}

func (d *streamsWorkspaceDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "workspace_name"},
	})
}

func (d *streamsWorkspaceDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamsWorkspaceConfig TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamsWorkspaceConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamsWorkspaceConfig.ProjectID.ValueString()
	workspaceName := streamsWorkspaceConfig.WorkspaceName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamWorkspace(ctx, projectID, workspaceName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newInstanceModel, diags := streaminstance.NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var newWorkspaceModel TFModel
	newWorkspaceModel.FromInstanceModel(newInstanceModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, newWorkspaceModel)...)
}
