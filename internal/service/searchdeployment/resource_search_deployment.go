package searchdeployment

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &searchDeploymentRS{}
var _ resource.ResourceWithImportState = &searchDeploymentRS{}

const searchDeploymentName = "search_deployment"

func Resource() resource.Resource {
	return &searchDeploymentRS{
		RSCommon: config.RSCommon{
			ResourceName: searchDeploymentName,
		},
	}
}

type searchDeploymentRS struct {
	config.RSCommon
}

func (r *searchDeploymentRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

const defaultSearchNodeTimeout time.Duration = 3 * time.Hour
const minTimeoutCreateUpdate time.Duration = 1 * time.Minute
const minTimeoutDelete time.Duration = 30 * time.Second

func retryTimeConfig(configuredTimeout, minTimeout time.Duration) retrystrategy.TimeConfig {
	return retrystrategy.TimeConfig{
		Timeout:    configuredTimeout,
		MinTimeout: minTimeout,
		Delay:      1 * time.Minute,
	}
}

func (r *searchDeploymentRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var searchDeploymentPlan TFSearchDeploymentRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	searchDeploymentReq := NewSearchDeploymentReq(ctx, &searchDeploymentPlan)
	if _, _, err := connV2.AtlasSearchApi.CreateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment creation", err.Error())
		return
	}

	createTimeout, diags := searchDeploymentPlan.Timeouts.Create(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	deploymentResp, err := WaitSearchNodeStateTransition(ctx, projectID, clusterName, connV2.AtlasSearchApi,
		retryTimeConfig(createTimeout, minTimeoutCreateUpdate))
	if err != nil {
		resp.Diagnostics.AddError("error during search deployment creation", err.Error())
		return
	}
	newSearchNodeModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *searchDeploymentRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var searchDeploymentPlan TFSearchDeploymentRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	deploymentResp, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting search deployment information", err.Error())
		return
	}

	newSearchNodeModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *searchDeploymentRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var searchDeploymentPlan TFSearchDeploymentRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	searchDeploymentReq := NewSearchDeploymentReq(ctx, &searchDeploymentPlan)
	if _, _, err := connV2.AtlasSearchApi.UpdateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment update", err.Error())
		return
	}

	updateTimeout, diags := searchDeploymentPlan.Timeouts.Update(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	deploymentResp, err := WaitSearchNodeStateTransition(ctx, projectID, clusterName, connV2.AtlasSearchApi,
		retryTimeConfig(updateTimeout, minTimeoutCreateUpdate))
	if err != nil {
		resp.Diagnostics.AddError("error during search deployment update", err.Error())
		return
	}
	newSearchNodeModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *searchDeploymentRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var searchDeploymentState *TFSearchDeploymentRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchDeploymentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentState.ProjectID.ValueString()
	clusterName := searchDeploymentState.ClusterName.ValueString()
	if _, err := connV2.AtlasSearchApi.DeleteAtlasSearchDeployment(ctx, projectID, clusterName).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment delete", err.Error())
		return
	}

	deleteTimeout, diags := searchDeploymentState.Timeouts.Delete(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := WaitSearchNodeDelete(ctx, projectID, clusterName, connV2.AtlasSearchApi, retryTimeConfig(deleteTimeout, minTimeoutDelete)); err != nil {
		resp.Diagnostics.AddError("error during search deployment delete", err.Error())
		return
	}
}

func (r *searchDeploymentRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
