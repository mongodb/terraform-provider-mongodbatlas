//nolint:gocritic
package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &encryptionAtRestPrivateEndpointDS{}
var _ datasource.DataSourceWithConfigure = &encryptionAtRestPrivateEndpointDS{}

func DataSource() datasource.DataSource {
	return &encryptionAtRestPrivateEndpointDS{
		DSCommon: config.DSCommon{
			DataSourceName: encryptionAtRestPrivateEndpointName,
		},
	}
}

type encryptionAtRestPrivateEndpointDS struct {
	config.DSCommon
}

func (d *encryptionAtRestPrivateEndpointDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *encryptionAtRestPrivateEndpointDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var earPrivateEndpointConfig TFEarPrivateEndpointModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &earPrivateEndpointConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := d.Client.AtlasV2
	// if err != nil {
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	//  resp.Diagnostics.Append(resp.State.Set(ctx, NewTFEarPrivateEndpoint(apiResp))...)
}
