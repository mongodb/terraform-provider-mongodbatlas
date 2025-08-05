package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamWorkspaceDS{}
var _ datasource.DataSourceWithConfigure = &streamWorkspaceDS{}

func DataSource() datasource.DataSource {
	return &streamWorkspaceDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamWorkspaceName,
		},
	}
}

type streamWorkspaceDS struct {
	config.DSCommon
}

func (d *streamWorkspaceDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "workspace_name"},
	})
}

func (d *streamWorkspaceDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamWorkspaceConfig TFStreamWorkspaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamWorkspaceConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamWorkspaceConfig.ProjectID.ValueString()
	workspaceName := streamWorkspaceConfig.WorkspaceName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamInstance(ctx, projectID, workspaceName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamWorkspaceModel, diags := NewTFStreamWorkspace(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamWorkspaceModel)...)
}
