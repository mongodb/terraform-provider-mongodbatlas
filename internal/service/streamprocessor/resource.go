package streamprocessor

import (
	"context"
	"errors"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const StreamProcessorName = "stream_processor"

var _ resource.ResourceWithConfigure = &streamProcessorRS{}
var _ resource.ResourceWithImportState = &streamProcessorRS{}

const (
	errorCreateStartActions    = "You need to fix the processor and import the resource or delete it manually and re-run terraform apply."
	errorCreateStart           = "Error starting stream processor. " + errorCreateStartActions
	errorCreateStartTransition = "Error changing state of stream processor. " + errorCreateStartActions
)

func Resource() resource.Resource {
	return &streamProcessorRS{
		RSCommon: config.RSCommon{
			ResourceName: StreamProcessorName,
		},
	}
}

type streamProcessorRS struct {
	config.RSCommon
}

func (r *streamProcessorRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *streamProcessorRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	streamProcessorReq, diags := NewStreamProcessorReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var needsStarting bool
	if !plan.State.IsNull() && !plan.State.IsUnknown() {
		switch plan.State.ValueString() {
		case StartedState:
			needsStarting = true
		case CreatedState:
			needsStarting = false
		default:
			resp.Diagnostics.AddError("When creating a stream processor, the only valid states are CREATED and STARTED", "")
			return
		}
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectID.ValueString()
	instanceName := plan.InstanceName.ValueString()
	processorName := plan.ProcessorName.ValueString()
	_, _, err := connV2.StreamsApi.CreateStreamProcessor(ctx, projectID, instanceName, streamProcessorReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	streamProcessorParams := &admin.GetStreamProcessorApiParams{
		GroupId:       projectID,
		TenantName:    instanceName,
		ProcessorName: processorName,
	}

	streamProcessorResp, err := WaitStateTransition(ctx, streamProcessorParams, connV2.StreamsApi, []string{InitiatingState, CreatingState}, []string{CreatedState})
	if err != nil {
		resp.Diagnostics.AddError("Error creating stream processor", err.Error())
		return
	}

	if needsStarting {
		_, err := connV2.StreamsApi.StartStreamProcessorWithParams(ctx,
			&admin.StartStreamProcessorApiParams{
				GroupId:       projectID,
				TenantName:    instanceName,
				ProcessorName: processorName,
			},
		).Execute()
		if err != nil {
			resp.Diagnostics.AddError(errorCreateStart, err.Error())
			return
		}
		streamProcessorResp, err = WaitStateTransition(ctx, streamProcessorParams, connV2.StreamsApi, []string{CreatedState}, []string{StartedState})
		if err != nil {
			resp.Diagnostics.AddError(errorCreateStartTransition, err.Error())
			return
		}
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, projectID, instanceName, streamProcessorResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamProcessorModel)...)
}

func (r *streamProcessorRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFStreamProcessorRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2

	projectID := state.ProjectID.ValueString()
	instanceName := state.InstanceName.ValueString()
	streamProcessor, apiResp, err := connV2.StreamsApi.GetStreamProcessor(ctx, projectID, instanceName, state.ProcessorName.ValueString()).Execute()
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, projectID, instanceName, streamProcessor)
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

	plannedState := plan.State.ValueString()
	if plannedState == "" {
		plannedState = state.State.ValueString()
	}

	projectID := plan.ProjectID.ValueString()
	instanceName := plan.InstanceName.ValueString()
	processorName := plan.ProcessorName.ValueString()
	currentState := state.State.ValueString()
	connV2 := r.Client.AtlasV2
	var streamProcessorResp *admin.StreamsProcessorWithStats

	// requestParams are needed for the state transition via the GET API
	requestParams := &admin.GetStreamProcessorApiParams{
		GroupId:       projectID,
		TenantName:    instanceName,
		ProcessorName: processorName,
	}

	if errMsg, isValidStateTransition := ValidateUpdateStateTransition(currentState, plannedState); !isValidStateTransition {
		resp.Diagnostics.AddError(errMsg, "")
		return
	}

	// we must stop the current stream processor if the current state is started
	if currentState == StartedState {
		_, err := connV2.StreamsApi.StopStreamProcessorWithParams(ctx,
			&admin.StopStreamProcessorApiParams{
				GroupId:       plan.ProjectID.ValueString(),
				TenantName:    plan.InstanceName.ValueString(),
				ProcessorName: plan.ProcessorName.ValueString(),
			},
		).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error stopping stream processor", err.Error())
			return
		}

		// wait for transition from started to stopped
		_, err = WaitStateTransition(ctx, requestParams, r.Client.AtlasV2.StreamsApi, []string{StartedState}, []string{StoppedState})
		if err != nil {
			resp.Diagnostics.AddError("Error changing state of stream processor", err.Error())
			return
		}
	}

	// modify the stream processor
	modifyAPIRequestParams, diags := NewStreamProcessorUpdateReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	streamProcessorResp, _, err := r.Client.AtlasV2.StreamsApi.ModifyStreamProcessorWithParams(ctx, modifyAPIRequestParams).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error modifying stream processor", err.Error())
		return
	}

	// start the stream processor if the desired state is started
	if plannedState == StartedState {
		_, err := r.Client.AtlasV2.StreamsApi.StartStreamProcessorWithParams(ctx,
			&admin.StartStreamProcessorApiParams{
				GroupId:       projectID,
				TenantName:    instanceName,
				ProcessorName: processorName,
			},
		).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error starting stream processor", err.Error())
			return
		}

		// wait for transition from stopped to started
		streamProcessorResp, err = WaitStateTransition(ctx, requestParams, r.Client.AtlasV2.StreamsApi, []string{StoppedState}, []string{StartedState})
		if err != nil {
			resp.Diagnostics.AddError("Error changing state of stream processor", err.Error())
			return
		}
	}

	newStreamProcessorModel, diags := NewStreamProcessorWithStats(ctx, projectID, instanceName, streamProcessorResp)
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
	projectID, instanceName, processorName, err := splitImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), instanceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("processor_name"), processorName)...)
}

func splitImportID(id string) (projectID, instanceName, processorName *string, err error) {
	var re = regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a stream processor, use the format {instance_name}-{project_id}-(processor_name)")
		return
	}

	instanceName = &parts[1]
	projectID = &parts[2]
	processorName = &parts[3]

	return
}
