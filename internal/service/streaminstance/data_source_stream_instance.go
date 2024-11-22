package streaminstance

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := schema.Schema{
		Attributes: DSAttributes(true),
	}
	conversion.UpdateSchemaDescription(&ds1)
	requiredFields := []string{"project_id", "instance_name"}
	ds2 := conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil)
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
}

// DSAttributes returns the attribute definitions for a single stream instance.
// `withArguments` marks certain attributes as required (for singular data source) or as computed (for plural data source)
func DSAttributes(withArguments bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"instance_name": schema.StringAttribute{
			Required: withArguments,
			Computed: !withArguments,
		},
		"project_id": schema.StringAttribute{
			Required: withArguments,
			Computed: !withArguments,
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
		"stream_config": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"tier": schema.StringAttribute{
					Computed: true,
				},
			},
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
