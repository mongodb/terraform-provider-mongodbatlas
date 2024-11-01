package flexcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", resourceName),
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = PluralDataSourceSchema(ctx)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModelDSP
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	apiResp, _, err := connV2.FlexClustersApi.ListFlexClusters(ctx, tfModel.ProjectId.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error reading plural data source", err.Error())
		return
	}

	newFlexClustersModel, diags := NewTFModelDSP(ctx, tfModel.ProjectId.ValueString(), apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClustersModel)...)
}
