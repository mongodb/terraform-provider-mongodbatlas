package searchdeployment

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

const resourceName = "search_deployment"

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type rs struct {
	config.RSCommon
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

const defaultSearchNodeTimeout time.Duration = 3 * time.Hour
const minTimeoutCreateUpdate time.Duration = 1 * time.Minute
const minTimeoutDelete time.Duration = 30 * time.Second

func RetryTimeConfig(configuredTimeout, minTimeout time.Duration) retrystrategy.TimeConfig {
	return retrystrategy.TimeConfig{
		Timeout:    configuredTimeout,
		MinTimeout: minTimeout,
		Delay:      1 * time.Minute,
	}
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFSearchDeploymentRSModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.ClusterName.ValueString()
	createReq := NewSearchDeploymentReq(ctx, &plan)
	createTimeout, localDiags := plan.Timeouts.Create(ctx, defaultSearchNodeTimeout)
	diags.Append(localDiags...)
	if diags.HasError() {
		return
	}
	if plan.DeleteOnCreateTimeout.ValueBool() {
		var deferCall func()
		deleteOnTimeout := func(newCtx context.Context) error {
			cleanup.ReplaceContextDeadlineExceededDiags(diags, createTimeout)
			_, err := connV2.AtlasSearchApi.DeleteAtlasSearchDeployment(newCtx, projectID, clusterName).Execute()
			return err
		}
		ctx, deferCall = cleanup.OnTimeout(
			ctx, createTimeout, diags.AddWarning, fmt.Sprintf("Search Deployment %s, (%s)", clusterName, projectID), deleteOnTimeout,
		)
		defer deferCall()
	}
	if _, _, err := connV2.AtlasSearchApi.CreateAtlasSearchDeployment(ctx, projectID, clusterName, &createReq).Execute(); err != nil {
		diags.AddError("error during search deployment creation", err.Error())
		return
	}

	deploymentResp, err := WaitSearchNodeStateTransition(ctx, projectID, clusterName, connV2.AtlasSearchApi,
		RetryTimeConfig(createTimeout, minTimeoutCreateUpdate))
	if err != nil {
		diags.AddError("error during search deployment creation", err.Error())
		return
	}
	outModel, localDiags := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &plan.Timeouts, false)
	diags.Append(localDiags...)
	if diags.HasError() {
		return
	}
	outModel.SkipWaitOnUpdate = plan.SkipWaitOnUpdate
	outModel.DeleteOnCreateTimeout = plan.DeleteOnCreateTimeout
	diags.Append(resp.State.Set(ctx, outModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFSearchDeploymentRSModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectID.ValueString()
	clusterName := state.ClusterName.ValueString()
	deploymentResp, getResp, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		diags.AddError("error getting search deployment information", err.Error())
		return
	}

	if IsNotFoundDeploymentResponse(deploymentResp) {
		resp.State.RemoveResource(ctx)
		return
	}

	outModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &state.Timeouts, false)
	diags.Append(diagnostics...)
	if diags.HasError() {
		return
	}
	outModel.SkipWaitOnUpdate = state.SkipWaitOnUpdate
	outModel.DeleteOnCreateTimeout = state.DeleteOnCreateTimeout
	diags.Append(resp.State.Set(ctx, outModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFSearchDeploymentRSModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectID.ValueString()
	clusterName := plan.ClusterName.ValueString()
	updateReq := NewSearchDeploymentReq(ctx, &plan)
	deploymentResp, _, err := connV2.AtlasSearchApi.UpdateAtlasSearchDeployment(ctx, projectID, clusterName, &updateReq).Execute()
	if err != nil {
		diags.AddError("error during search deployment update", err.Error())
		return
	}

	updateTimeout, localDiags := plan.Timeouts.Update(ctx, defaultSearchNodeTimeout)
	diags.Append(localDiags...)
	if diags.HasError() {
		return
	}
	if !plan.SkipWaitOnUpdate.ValueBool() {
		deploymentResp, err = WaitSearchNodeStateTransition(ctx, projectID, clusterName, connV2.AtlasSearchApi,
			RetryTimeConfig(updateTimeout, minTimeoutCreateUpdate))
		if err != nil {
			diags.AddError("error during search deployment update", err.Error())
			return
		}
	}

	outModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &plan.Timeouts, false)
	diags.Append(diagnostics...)
	if diags.HasError() {
		return
	}
	outModel.SkipWaitOnUpdate = plan.SkipWaitOnUpdate
	outModel.DeleteOnCreateTimeout = plan.DeleteOnCreateTimeout
	diags.Append(resp.State.Set(ctx, outModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TFSearchDeploymentRSModel
	diags := &resp.Diagnostics
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectID.ValueString()
	clusterName := state.ClusterName.ValueString()
	if _, err := connV2.AtlasSearchApi.DeleteAtlasSearchDeployment(ctx, projectID, clusterName).Execute(); err != nil {
		diags.AddError("error during search deployment delete", err.Error())
		return
	}

	deleteTimeout, localDiags := state.Timeouts.Delete(ctx, defaultSearchNodeTimeout)
	diags.Append(localDiags...)
	if diags.HasError() {
		return
	}
	if err := WaitSearchNodeDelete(ctx, projectID, clusterName, connV2.AtlasSearchApi, RetryTimeConfig(deleteTimeout, minTimeoutDelete)); err != nil {
		diags.AddError("error during search deployment delete", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, clusterName, err := splitSearchNodeImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search deployment import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_name"), clusterName)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func splitSearchNodeImportID(id string) (projectID, clusterName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{cluster_name}")
		return
	}

	projectID = parts[1]
	clusterName = parts[2]
	return
}

func IsNotFoundDeploymentResponse(deploymentResp *admin.ApiSearchDeploymentResponse) bool {
	return deploymentResp == nil || deploymentResp.Id == nil
}
