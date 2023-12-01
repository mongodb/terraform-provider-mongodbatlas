package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamInstanceDS{}
var _ datasource.DataSourceWithConfigure = &streamInstanceDS{}

func DataSource() datasource.DataSource {
	return &streamInstanceDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamInstanceName,
		},
	}
}

type streamInstanceDS struct {
	config.DSCommon
}

func (d *streamInstanceDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: DSAttributes(true),
	}
}

func DSAttributes(definingArguments bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"instance_name": schema.StringAttribute{
			Required: definingArguments,
			Computed: !definingArguments,
		},
		"project_id": schema.StringAttribute{
			Required: definingArguments,
			Computed: !definingArguments,
		},
		"data_process_region": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"cloud_provider": schema.StringAttribute{
					Computed: true,
				},
				"region": schema.StringAttribute{
					Computed: true,
				},
			},
		},
		"hostnames": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
	}
}

func (d *streamInstanceDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamInstanceConfig TFStreamInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamInstanceConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamInstanceConfig.ProjectID.ValueString()
	instanceName := streamInstanceConfig.InstanceName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamInstance(ctx, projectID, instanceName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamInstanceModel, diags := NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstanceModel)...)
}
