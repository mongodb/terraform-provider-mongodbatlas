package pushbasedlogexport

import (
	"context"
	"log"
	"slices"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	pushBasedLogExportName               = "push_based_log_export"
	defaultTimeout         time.Duration = 15 * time.Minute
	minTimeoutCreateUpdate time.Duration = 1 * time.Minute
	minTimeoutDelete       time.Duration = 30 * time.Second
	retryTimeDelay         time.Duration = 10 * time.Second
)

var _ resource.ResourceWithConfigure = &pushBasedLogExportRS{}
var _ resource.ResourceWithImportState = &pushBasedLogExportRS{}

func Resource() resource.Resource {
	return &pushBasedLogExportRS{
		RSCommon: config.RSCommon{
			ResourceName: pushBasedLogExportName,
		},
	}
}

type pushBasedLogExportRS struct {
	config.RSCommon
}

func (r *pushBasedLogExportRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *pushBasedLogExportRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfPlan TFPushBasedLogExportRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logExportConfigReq := NewPushBasedLogExportCreateReq(&tfPlan)

	connV2 := r.Client.AtlasV2
	projectID := tfPlan.ProjectID.ValueString()
	if _, err := connV2.PushBasedLogExportApi.CreatePushBasedLogConfiguration(ctx, projectID, logExportConfigReq).Execute(); err != nil {
		resp.Diagnostics.AddError("Error when creating push-based log export configuration", err.Error())

		if err := unconfigureFailedPushBasedLog(ctx, connV2, projectID); err != nil {
			resp.Diagnostics.AddError("Error when unconfiguring push-based log export configuration", err.Error())
			return
		}
		return
	}

	timeout, diags := tfPlan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	logExportConfigResp, err := WaitStateTransition(ctx, projectID, connV2.PushBasedLogExportApi,
		retryTimeConfig(timeout, minTimeoutCreateUpdate))
	if err != nil {
		resp.Diagnostics.AddError("Error when creating push-based log export configuration", err.Error())

		if err := unconfigureFailedPushBasedLog(ctx, connV2, projectID); err != nil {
			resp.Diagnostics.AddError("Error when unconfiguring push-based log export configuration", err.Error())
			return
		}
		return
	}

	newTFModel, diags := NewTFPushBasedLogExport(ctx, projectID, logExportConfigResp, &tfPlan.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTFModel)...)
}

func (r *pushBasedLogExportRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfState TFPushBasedLogExportRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := tfState.ProjectID.ValueString()
	logConfig, getResp, err := connV2.PushBasedLogExportApi.GetPushBasedLogConfiguration(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error when getting push-based log export configuration", err.Error())
		return
	}

	newTFModel, diags := NewTFPushBasedLogExport(ctx, projectID, logConfig, &tfState.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTFModel)...)
}

func (r *pushBasedLogExportRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var tfPlan TFPushBasedLogExportRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logExportConfigReq := NewPushBasedLogExportUpdateReq(&tfPlan)

	connV2 := r.Client.AtlasV2
	projectID := tfPlan.ProjectID.ValueString()
	if _, err := connV2.PushBasedLogExportApi.UpdatePushBasedLogConfiguration(ctx, projectID, logExportConfigReq).Execute(); err != nil {
		resp.Diagnostics.AddError("Error when updating push-based log export configuration", err.Error())
		return
	}

	timeout, diags := tfPlan.Timeouts.Update(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	logExportConfigResp, err := WaitStateTransition(ctx, projectID, connV2.PushBasedLogExportApi,
		retryTimeConfig(timeout, minTimeoutCreateUpdate))
	if err != nil {
		resp.Diagnostics.AddError("Error when updating push-based log export configuration", err.Error())
		return
	}

	newTFModel, diags := NewTFPushBasedLogExport(ctx, projectID, logExportConfigResp, &tfPlan.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newTFModel)...)
}

func (r *pushBasedLogExportRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var tfState *TFPushBasedLogExportRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := tfState.ProjectID.ValueString()
	if _, err := connV2.PushBasedLogExportApi.DeletePushBasedLogConfiguration(ctx, projectID).Execute(); err != nil {
		resp.Diagnostics.AddError("Error when deleting push-based log export configuration", err.Error())
		return
	}

	deleteTimeout, diags := tfState.Timeouts.Delete(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := WaitResourceDelete(ctx, projectID, connV2.PushBasedLogExportApi, retryTimeConfig(deleteTimeout, minTimeoutDelete)); err != nil {
		resp.Diagnostics.AddError("Error when deleting push-based log export configuration", err.Error())
		return
	}
}

func (r *pushBasedLogExportRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), req.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func retryTimeConfig(configuredTimeout, minTimeout time.Duration) retrystrategy.TimeConfig {
	return retrystrategy.TimeConfig{
		Timeout:    configuredTimeout,
		MinTimeout: minTimeout,
		Delay:      retryTimeDelay,
	}
}

func unconfigureFailedPushBasedLog(ctx context.Context, connV2 *admin.APIClient, projectID string) error {
	logConfig, _, _ := connV2.PushBasedLogExportApi.GetPushBasedLogConfiguration(ctx, projectID).Execute()
	if logConfig != nil && slices.Contains(failureStates, *logConfig.State) {
		log.Printf("[INFO] Unconfiguring push-based log export for project due to create failure: %s", projectID)
		if _, err := connV2.PushBasedLogExportApi.DeletePushBasedLogConfiguration(ctx, projectID).Execute(); err != nil {
			return err
		}
	}
	return nil
}
