package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamInstanceDS{}
var _ datasource.DataSourceWithConfigure = &streamInstanceDS{}

func DataSource() datasource.DataSource {
	return &streamInstanceDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamInstanceName,
		},
	}
}

type streamInstanceDS struct {
	config.DSCommon
}

func (d *streamInstanceDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "instance_name"},
	})
}

func (d *streamInstanceDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamInstanceConfig TFStreamInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamInstanceConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamInstanceConfig.ProjectID.ValueString()
	instanceName := streamInstanceConfig.InstanceName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamWorkspace(ctx, projectID, instanceName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamInstanceModel, diags := NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstanceModel)...)
}
