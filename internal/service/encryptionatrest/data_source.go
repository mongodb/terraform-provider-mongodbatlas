package encryptionatrest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &encryptionAtRestDS{}
var _ datasource.DataSourceWithConfigure = &encryptionAtRestDS{}

func DataSource() datasource.DataSource {
	return &encryptionAtRestDS{
		DSCommon: config.DSCommon{
			DataSourceName: encryptionAtRestResourceName,
		},
	}
}

type encryptionAtRestDS struct {
	config.DSCommon
}

func (d *encryptionAtRestDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *encryptionAtRestDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var earConfig TFEncryptionAtRestDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &earConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: update before merging to master: connV2 := d.Client.AtlasV2
	connV2 := d.Client.AtlasPreview
	projectID := earConfig.ProjectID.ValueString()

	encryptionResp, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFEncryptionAtRestDSModel(projectID, encryptionResp))...)
}
