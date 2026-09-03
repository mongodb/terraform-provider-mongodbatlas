package streamprocessor

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &StreamProccesorDS{}
var _ datasource.DataSourceWithConfigure = &StreamProccesorDS{}

func DataSource() datasource.DataSource {
	return &StreamProccesorDS{
		DSCommon: config.DSCommon{
			DataSourceName: StreamProcessorName,
		},
	}
}

type StreamProccesorDS struct {
	config.DSCommon
}

func (d *StreamProccesorDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	fields := dataSourceOverridenFields()
	maps.Copy(fields, dataSourceOptionsOverridenField())
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields:  []string{"project_id", "processor_name"},
		OverridenFields: fields,
	})
}

// dataSourceOptionsOverridenField redefines options for the data sources without
// resume_from_checkpoint, which is request-only and never returned by the API.
//
// IMPORTANT: keep in sync with the options attribute in ResourceSchema.
func dataSourceOptionsOverridenField() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"options": dsschema.SingleNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Optional configuration for the stream processor. Empty `options` objects are not supported.",
			Attributes: map[string]dsschema.Attribute{
				"dlq": dsschema.SingleNestedAttribute{
					Computed:            true,
					MarkdownDescription: "Dead letter queue for the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/reference/glossary/#std-term-dead-letter-queue) for more information.",
					Attributes: map[string]dsschema.Attribute{
						"coll": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the collection to use for the DLQ.",
						},
						"connection_name": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the connection to write DLQ messages to. Must be an Atlas connection.",
						},
						"db": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the database to use for the DLQ.",
						},
					},
				},
				"autoscaling": dsschema.SingleNestedAttribute{
					Computed:            true,
					MarkdownDescription: "Vertical autoscaling configuration for the stream processor. When present, the processor automatically scales its tier between `min_tier` and `max_tier` based on load; `tier` is used only as the initial/baseline tier and the running tier is reported by `effective_tier`. To disable autoscaling, remove this block.",
					Attributes: map[string]dsschema.Attribute{
						"min_tier": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Tier floor for autoscaling (scale-down limit). When not set, it defaults to the lower of the processor `tier` and the workspace default tier.",
						},
						"max_tier": dsschema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Tier ceiling for autoscaling (scale-up limit). When not set, it defaults to the workspace maximum tier.",
						},
					},
				},
			},
		},
	}
}

// dataSourceOverridenFields returns the root-level attributes that differ from the resource schema.
// It deliberately excludes options: the plural data source passes this as OverridenRootFields, where
// adding options would create a stray root-level attribute.
func dataSourceOverridenFields() map[string]dsschema.Attribute {
	fields := map[string]dsschema.Attribute{}
	fields["instance_name"] = dsschema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Label that identifies the stream processing workspace.",
		DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
		Validators: []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("workspace_name")),
		},
	}
	fields["workspace_name"] = dsschema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Label that identifies the stream processing workspace. Conflicts with `instance_name`.",
		Validators: []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("instance_name")),
		},
	}
	return fields
}

func (d *StreamProccesorDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamProccesorConfig TFStreamProcessorDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamProccesorConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamProccesorConfig.ProjectID.ValueString()
	workspaceOrInstanceName := GetWorkspaceOrInstanceName(streamProccesorConfig.WorkspaceName, streamProccesorConfig.InstanceName)

	processorName := streamProccesorConfig.ProcessorName.ValueString()
	apiResp, _, err := connV2.StreamsAPI.GetStreamProcessor(ctx, projectID, workspaceOrInstanceName, processorName).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	instanceName := streamProccesorConfig.InstanceName.ValueString()
	workspaceName := streamProccesorConfig.WorkspaceName.ValueString()

	newStreamTFStreamprocessorDSModelModel, diags := NewTFStreamprocessorDSModel(ctx, projectID, instanceName, workspaceName, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamTFStreamprocessorDSModelModel)...)
}
