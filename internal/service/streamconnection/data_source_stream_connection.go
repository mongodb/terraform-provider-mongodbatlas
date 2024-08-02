package streamconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamConnectionDS{}
var _ datasource.DataSourceWithConfigure = &streamConnectionDS{}

func DataSource() datasource.DataSource {
	return &streamConnectionDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamConnectionName,
		},
	}
}

type streamConnectionDS struct {
	config.DSCommon
}

func (d *streamConnectionDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: DSAttributes(true),
	}
}

// DSAttributes returns the attribute definitions for a single stream connection.
// `withArguments` marks certain attributes as required (for singular data source) or as computed (for plural data source)
func DSAttributes(withArguments bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"project_id": schema.StringAttribute{
			Required: withArguments,
			Computed: !withArguments,
		},
		"instance_name": schema.StringAttribute{
			Required: withArguments,
			Computed: !withArguments,
		},
		"connection_name": schema.StringAttribute{
			Required: withArguments,
			Computed: !withArguments,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},

		// cluster type specific
		"cluster_name": schema.StringAttribute{
			Computed: true,
		},
		"db_role_to_execute": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"role": schema.StringAttribute{
					Computed: true,
				},
				"type": schema.StringAttribute{
					Computed: true,
				},
			},
		},

		// kafka type specific
		"authentication": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"mechanism": schema.StringAttribute{
					Computed: true,
				},
				"password": schema.StringAttribute{
					Computed:  true,
					Sensitive: true,
				},
				"username": schema.StringAttribute{
					Computed: true,
				},
			},
		},
		"bootstrap_servers": schema.StringAttribute{
			Computed: true,
		},
		"config": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
		"security": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"broker_public_certificate": schema.StringAttribute{
					Computed: true,
				},
				"protocol": schema.StringAttribute{
					Computed: true,
				},
			},
		},
		"networking": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"access": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *streamConnectionDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionConfig TFStreamConnectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionConfig.ProjectID.ValueString()
	instanceName := streamConnectionConfig.InstanceName.ValueString()
	connectionName := streamConnectionConfig.ConnectionName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamConnection(ctx, projectID, instanceName, connectionName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, nil, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}
