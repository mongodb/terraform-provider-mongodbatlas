package streamprocessor

import (
	"context"
	"fmt"

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
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields:  []string{"project_id", "processor_name"},
		OverridenFields: dataSourceOverridenFields(),
	})
}

func dataSourceOverridenFields() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"instance_name": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Label that identifies the stream processing workspace.",
			DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("workspace_name")),
			},
		},
		"workspace_name": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Label that identifies the stream processing workspace. Conflicts with `instance_name`.",
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("instance_name")),
			},
		},
	}
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
	apiResp, _, err := connV2.StreamsApi.GetStreamProcessor(ctx, projectID, workspaceOrInstanceName, processorName).Execute()

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
