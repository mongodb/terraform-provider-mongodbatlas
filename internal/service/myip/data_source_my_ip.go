package myip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	myIP = "my_ip"
)

type ds struct {
	config.DSCommon
}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: myIP,
		},
	}
}

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

type Model struct {
	IPAddress types.String `tfsdk:"ip_address"`
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "My IP",
		Attributes: map[string]schema.Attribute{
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "The IP.",
				Computed:            true,
			},
		},
	}
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseDSUserConfig *Model
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseDSUserConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, _, err := d.Client.Atlas.IPInfo.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error getting access list entry", err.Error())
		return
	}

	accessListEntry := &Model{
		IPAddress: types.StringValue(info.CurrentIPv4Address),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &accessListEntry)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
