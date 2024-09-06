package employeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &employeeAccessDS{}
var _ datasource.DataSourceWithConfigure = &employeeAccessDS{}

func DataSource() datasource.DataSource {
	return &employeeAccessDS{
		DSCommon: config.DSCommon{
			DataSourceName: employeeAccessName,
		},
	}
}

type employeeAccessDS struct {
	config.DSCommon
}

func (d *employeeAccessDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: Schema and model must be defined in data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = DataSourceSchema(ctx)
}

func (d *employeeAccessDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var employeeAccessConfig TFEmployeeAccessDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &employeeAccessConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := d.Client.AtlasV2
	// if err != nil {
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	// }

	// TODO: process response into new terraform state
	// newEmployeeAccessModel, diags := NewTFEmployeeAccess(ctx, apiResp)
	// if diags.HasError() {
	//	resp.Diagnostics.Append(diags...)
	//	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newEmployeeAccessModel)...)
}
