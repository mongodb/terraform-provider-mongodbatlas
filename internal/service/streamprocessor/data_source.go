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
		RequiredFields: []string{"project_id", "processor_name"},
		OverridenFields: map[string]dsschema.Attribute{
			"instance_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance.",
				DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("workspace_name")),
				},
			},
			"workspace_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance. Conflicts with `instance_name`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("instance_name")),
				},
			},
		},
	})
}

// getWorkspaceOrInstanceNameForDS returns the workspace name from workspace_name or instance_name field for datasource model
func getWorkspaceOrInstanceNameForDS(model *TFStreamProcessorDSModel) string {
	if !model.WorkspaceName.IsNull() && !model.WorkspaceName.IsUnknown() {
		return model.WorkspaceName.ValueString()
	}
	if !model.InstanceName.IsNull() && !model.InstanceName.IsUnknown() {
		return model.InstanceName.ValueString()
	}
	return ""
}

func (d *StreamProccesorDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamProccesorConfig TFStreamProcessorDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamProccesorConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamProccesorConfig.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceNameForDS(&streamProccesorConfig)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

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
