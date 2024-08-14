package streamprocessor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &StreamProccesorDS{}
var _ datasource.DataSourceWithConfigure = &StreamProccesorDS{}

func PluralDataSource() datasource.DataSource {
	return &streamProcessorsDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", StreamProccesorName),
		},
	}
}

type streamProcessorsDS struct {
	config.DSCommon
}

func (d *streamProcessorsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.PaginatedDSSchema(
		map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"instance_name": schema.StringAttribute{
				Required: true,
			},
		},
		DSAttributes(false))
}

type TFStreamProcessorsDSModel struct {
	ID           types.String               `tfsdk:"id"`
	ProjectID    types.String               `tfsdk:"project_id"`
	InstanceName types.String               `tfsdk:"instance_name"`
	Results      []TFStreamProcessorDSModel `tfsdk:"results"`
	PageNum      types.Int64                `tfsdk:"page_num"`
	ItemsPerPage types.Int64                `tfsdk:"items_per_page"`
	TotalCount   types.Int64                `tfsdk:"total_count"`
}
