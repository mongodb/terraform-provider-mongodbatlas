package employeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &employeeAccessGrantDS{}
var _ datasource.DataSourceWithConfigure = &employeeAccessGrantDS{}

func DataSource() datasource.DataSource {
	return &employeeAccessGrantDS{
		DSCommon: config.DSCommon{
			DataSourceName: employeeAccessGrantName,
		},
	}
}

type employeeAccessGrantDS struct {
	config.DSCommon
}

func (d *employeeAccessGrantDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: Schema and model must be defined in data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = DataSourceSchema(ctx)
}

func (d *employeeAccessGrantDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var employeeAccessGrantConfig TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &employeeAccessGrantConfig)...)
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
	// newEmployeeAccessGrantModel, diags := NewTFEmployeeAccessGrant(ctx, apiResp)
	// if diags.HasError() {
	//	resp.Diagnostics.Append(diags...)
	//	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newEmployeeAccessGrantModel)...)
}
