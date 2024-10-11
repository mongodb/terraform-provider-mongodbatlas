package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const advancedClusterName = "advanced_cluster" // TODO: if resource exists this can be deleted

var _ datasource.DataSource = &advancedClusterDS{}
var _ datasource.DataSourceWithConfigure = &advancedClusterDS{}

func DataSource() datasource.DataSource {
	return &advancedClusterDS{
		DSCommon: config.DSCommon{
			DataSourceName: advancedClusterName,
		},
	}
}

type advancedClusterDS struct {
	config.DSCommon
}


func (d *advancedClusterDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: Schema and model must be defined in data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = DataSourceSchema(ctx)
}

func (d *advancedClusterDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var advancedClusterConfig TFAdvancedClusterModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &advancedClusterConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := d.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	newAdvancedClusterModel, diags := NewTFAdvancedCluster(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newAdvancedClusterModel)...)
}
