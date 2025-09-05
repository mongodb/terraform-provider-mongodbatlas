package projectipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const projectIPAddressesName = "project_ip_addresses"

var _ datasource.DataSource = &projectIPAddressesDS{}
var _ datasource.DataSourceWithConfigure = &projectIPAddressesDS{}

func DataSource() datasource.DataSource {
	return &projectIPAddressesDS{
		DSCommon: config.DSCommon{
			DataSourceName: projectIPAddressesName,
		},
	}
}

type projectIPAddressesDS struct {
	config.DSCommon
}

func (d *projectIPAddressesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *projectIPAddressesDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	connV2 := d.Client.AtlasV2
	var databaseDSUserConfig *TFProjectIpAddressesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseDSUserConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectIPAddresses, _, err := connV2.ProjectsApi.GetGroupIpAddresses(ctx, databaseDSUserConfig.ProjectId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting project's IP addresses", err.Error())
		return
	}

	newProjectIPAddresses, diags := NewTFProjectIPAddresses(ctx, projectIPAddresses)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newProjectIPAddresses)...)
}
