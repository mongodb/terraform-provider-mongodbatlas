package controlplaneipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const controlPlaneIPAddressesName = "control_plane_ip_addresses"

var _ datasource.DataSource = &controlPlaneIPAddressesDS{}
var _ datasource.DataSourceWithConfigure = &controlPlaneIPAddressesDS{}

func DataSource() datasource.DataSource {
	return &controlPlaneIPAddressesDS{
		DSCommon: config.DSCommon{
			DataSourceName: controlPlaneIPAddressesName,
		},
	}
}

type controlPlaneIPAddressesDS struct {
	config.DSCommon
}

func (d *controlPlaneIPAddressesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *controlPlaneIPAddressesDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	connV2 := d.Client.AtlasV2
	apiResp, _, err := connV2.RootApi.ListControlPlaneAddresses(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching control plane ip addresses", err.Error())
		return
	}
	newControlPlaneIPAddressesModel, diags := NewTFControlPlaneIPAddresses(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newControlPlaneIPAddressesModel)...)
}
