package streamconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamConnectionDS{}
var _ datasource.DataSourceWithConfigure = &streamConnectionDS{}

func DataSource() datasource.DataSource {
	return &streamConnectionDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamConnectionName,
		},
	}
}

type streamConnectionDS struct {
	config.DSCommon
}

func (d *streamConnectionDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "connection_name"},
	})
}

func (d *streamConnectionDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionConfig TFStreamConnectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionConfig.ProjectID.ValueString()
	effectiveWorkspaceName := getEffectiveWorkspaceName(&streamConnectionConfig)
	if effectiveWorkspaceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	connectionName := streamConnectionConfig.ConnectionName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamConnection(ctx, projectID, effectiveWorkspaceName, connectionName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	instanceName := streamConnectionConfig.InstanceName.ValueString()
	workspaceName := streamConnectionConfig.WorkspaceName.ValueString()
	newStreamConnectionModel, diags := NewTFStreamConnectionWithInstanceName(ctx, projectID, instanceName, workspaceName, nil, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}
