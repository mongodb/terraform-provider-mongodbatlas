package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
)

var _ datasource.DataSource = &advancedClustersDS{}
var _ datasource.DataSourceWithConfigure = &advancedClustersDS{}

func PluralDataSource() datasource.DataSource {
	return &advancedClustersDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", advancedClusterName),
		},
	}
}

type advancedClustersDS struct {
	config.DSCommon
}

func (d *advancedClustersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: Schema and model must be defined in plural_data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = PluralDataSourceSchema(ctx)
}

func (d *advancedClustersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var advancedClustersConfig TFAdvancedClustersDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &advancedClustersConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to obtain list of results

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error fetching results", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	newAdvancedClustersModel, diags := NewTFAdvancedClusters(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newAdvancedClustersModel)...)
}
