package streamprocessor

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

const streamProcessorName = "stream_processor"

var _ resource.ResourceWithConfigure = &streamProcessorRS{}
var _ resource.ResourceWithImportState = &streamProcessorRS{}

func Resource() resource.Resource {
	return &streamProcessorRS{
		RSCommon: config.RSCommon{
			ResourceName: streamProcessorName,
		},
	}
}

type streamProcessorRS struct {
	config.RSCommon
}

func (r *streamProcessorRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *streamProcessorRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamProcessorPlan TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamProcessorPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	streamProcessorReq, diags := NewStreamProcessorReq(ctx, &streamProcessorPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	connV2 := r.Client.AtlasV2
	_, _, err := connV2.StreamsApi.CreateStreamProcessor(ctx, streamProcessorPlan.ProjectID.ValueString(), streamProcessorPlan.InstanceName.ValueString(), streamProcessorReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	streamProcessorParams := &admin.GetStreamProcessorApiParams{
		GroupId:       streamProcessorPlan.ProjectID.ValueString(),
		TenantName:    streamProcessorPlan.InstanceName.ValueString(),
		ProcessorName: streamProcessorPlan.ProcessorName.ValueString(),
	}

	streamProcessorResp, err := WaitStateTransition(ctx, streamProcessorParams, connV2.StreamsApi, []string{InitiatingState, CreatingState}, []string{CreatedState, FailedState})
	if err != nil {
		resp.Diagnostics.AddError("Error creating stream processor", err.Error())
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, streamProcessorResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamProcessorModel)...)
}

func (r *streamProcessorRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamProcessorState TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamProcessorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	streamProcessor, apiResp, err := connV2.StreamsApi.GetStreamProcessor(ctx, streamProcessorState.ProjectID.ValueString(), streamProcessorState.InstanceName.ValueString(), streamProcessorState.ProcessorName.ValueString()).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, streamProcessor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamProcessorModel)...)
}

func (r *streamProcessorRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFStreamProcessorRSModel
	var state TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	pendingStates := []string{CreatedState}
	desiredState := []string{}
	if !plan.State.Equal(state.State) {
		switch plan.State.ValueString() {
		case StartedState:
			desiredState = append(desiredState, StartedState)
			pendingStates = append(pendingStates, StoppedState)
			_, _, err := connV2.StreamsApi.StartStreamProcessorWithParams(ctx,
				&admin.StartStreamProcessorApiParams{
					GroupId:       plan.ProjectID.ValueString(),
					TenantName:    plan.InstanceName.ValueString(),
					ProcessorName: plan.ProcessorName.ValueString(),
				},
			).Execute()
			if err != nil {
				resp.Diagnostics.AddError("Error starting stream processor", err.Error())
			}
		case StoppedState:
			desiredState = append(desiredState, StoppedState)
			pendingStates = append(pendingStates, StartedState)
			_, _, err := connV2.StreamsApi.StopStreamProcessorWithParams(ctx,
				&admin.StopStreamProcessorApiParams{
					GroupId:       plan.ProjectID.ValueString(),
					TenantName:    plan.InstanceName.ValueString(),
					ProcessorName: plan.ProcessorName.ValueString(),
				},
			).Execute()
			if err != nil {
				resp.Diagnostics.AddError("Error stopping stream processor", err.Error())
			}
		default:
			resp.Diagnostics.AddError("transitions to states other than STARTED or STOPPED are not supported", "")
			return
		}
	} else {
		resp.Diagnostics.AddError("updating a Stream Processor is not supported. Please follow this guide to update it", "TODO")
		return
	}
	requestParams := &admin.GetStreamProcessorApiParams{
		GroupId:       plan.ProjectID.ValueString(),
		TenantName:    plan.InstanceName.ValueString(),
		ProcessorName: plan.ProcessorName.ValueString(),
	}

	streamProcessorResp, err := WaitStateTransition(ctx, requestParams, connV2.StreamsApi, pendingStates, desiredState)
	if err != nil {
		resp.Diagnostics.AddError("Error changing state of stream processor", err.Error())
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, streamProcessorResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamProcessorModel)...)
}

func (r *streamProcessorRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamProcessorState *TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamProcessorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	if _, err := connV2.StreamsApi.DeleteStreamProcessor(ctx, streamProcessorState.ProjectID.ValueString(), streamProcessorState.InstanceName.ValueString(), streamProcessorState.ProcessorName.ValueString()).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *streamProcessorRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, instanceName, processorName, err := splitStreamProcessorImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), instanceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("processor_name"), processorName)...)
}

func splitStreamProcessorImportID(id string) (projectID, instanceName, processorName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a stream processor, use the format {project_id}-{instance_name}-{processor_name}")
		return
	}

	projectID = &parts[1]
	instanceName = &parts[2]
	processorName = &parts[3]

	return
}
