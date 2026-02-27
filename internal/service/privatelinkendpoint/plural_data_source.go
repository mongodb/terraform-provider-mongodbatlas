package privatelinkendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: "privatelink_endpoints",
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = PluralDataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFPrivateLinkEndpointsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateEndpoints, _, err := d.Client.AtlasV2.PrivateEndpointServicesApi.ListPrivateEndpointService(ctx, state.ProjectID.ValueString(), state.ProviderName.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting Private Endpoints", err.Error())
		return
	}

	results, diags := newTFPrivateLinkEndpointResults(ctx, privateEndpoints)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Results = results
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
