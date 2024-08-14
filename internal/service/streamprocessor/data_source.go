package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const StreamProccesorName = "streamprocessor"

var _ datasource.DataSource = &StreamProccesorDS{}
var _ datasource.DataSourceWithConfigure = &StreamProccesorDS{}

func DataSource() datasource.DataSource {
	return &StreamProccesorDS{
		DSCommon: config.DSCommon{
			DataSourceName: StreamProccesorName,
		},
	}
}

type StreamProccesorDS struct {
	config.DSCommon
}

func (d *StreamProccesorDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: Schema and model must be defined in data_source_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StreamProccesorDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamProccesorConfig TFStreamProcessorDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamProccesorConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamProccesorConfig.ProjectID.ValueString()
	instanceName := streamProccesorConfig.InstanceName.ValueString()
	processorName := streamProccesorConfig.ProcessorName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamProcessor(ctx, projectID, instanceName, processorName).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamTFStreamprocessorDSModelModel, diags := NewTFStreamprocessorDSModel(ctx, projectID, instanceName, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamTFStreamprocessorDSModelModel)...)
}
