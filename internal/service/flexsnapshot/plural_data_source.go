package flexsnapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	// TODO: Schema and model must be defined in plural_data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(DataSourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "name", "snapshot_id"},
	})
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// var tfModel TFFlexSnapshotsDSModel
	// resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // TODO: make get request to obtain list of results

	// // connV2 := r.Client.AtlasV2
	// //if err != nil {
	// //	resp.Diagnostics.AddError("error fetching results", err.Error())
	// //	return
	// //}

	// // TODO: process response into new terraform state
	// newFlexSnapshotsModel, diags := NewTFModelPluralDS(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newFlexSnapshotsModel)...)
}
