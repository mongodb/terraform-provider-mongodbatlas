//nolint:gocritic
package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	// "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/types"
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
	// TODO: Schema and model must be defined in data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = DataSourceSchema(ctx)
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var tfModel TFModel
	// resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// TODO: make get request to resource

	// connV2 := d.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	// newFlexClusterModel, diags := NewTFModel(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClusterModel)...)
}