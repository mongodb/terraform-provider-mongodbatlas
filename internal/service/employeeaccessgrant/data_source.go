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
	resp.Schema = DataSourceSchema(ctx)
}

func (d *employeeAccessGrantDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
