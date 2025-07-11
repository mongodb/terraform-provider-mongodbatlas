package apikeyprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "api_key_id"},
	})
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := tfModel.ProjectId.ValueString()
	// Once CLOUDP-328946 is done, we would use the single GET API to fetch the specific API key project assignment
	apiKeys, err := ListAllProjectAPIKeys(ctx, connV2, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	apiKeyID := tfModel.ApiKeyId.ValueString()
	newAPIKeyProjectAssignmentModel, diags := NewTFModel(ctx, apiKeys, projectID, apiKeyID)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newAPIKeyProjectAssignmentModel)...)
}
