package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
	})
}

// getEffectiveInstanceNameForDS returns the instance name from either instance_name or workspace_name field for datasource model
func getEffectiveInstanceNameForDS(model *TFStreamProcessorDSModel) string {
	if !model.InstanceName.IsNull() && !model.InstanceName.IsUnknown() {
		return model.InstanceName.ValueString()
	}
	if !model.WorkspaceName.IsNull() && !model.WorkspaceName.IsUnknown() {
		return model.WorkspaceName.ValueString()
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
	instanceName := getEffectiveInstanceNameForDS(&streamProccesorConfig)
	if instanceName == "" {
		resp.Diagnostics.AddError("validation error", "either instance_name or workspace_name must be provided")
		return
	}
	processorName := streamProccesorConfig.ProcessorName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamProcessor(ctx, projectID, instanceName, processorName).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamTFStreamprocessorDSModelModel, diags := NewTFStreamprocessorDSModelWithOriginal(ctx, projectID, instanceName, apiResp, &streamProccesorConfig)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamTFStreamprocessorDSModelModel)...)
}
