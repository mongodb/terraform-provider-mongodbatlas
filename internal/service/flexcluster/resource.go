package flexcluster

import (
	"context"
	"errors"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const resourceName = "flex_cluster"
const ErrorUpdateNotAllowed = "update not allowed"

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

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

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flexClusterReq, diags := NewAtlasCreateReq(ctx, &tfModel)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	projectID := tfModel.ProjectId.ValueString()
	clusterName := tfModel.Name.ValueString()

	connV2 := r.Client.AtlasV2
	_, _, err := connV2.FlexClustersApi.CreateFlexCluster(ctx, projectID, flexClusterReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: projectID,
		Name:    clusterName,
	}

	flexClusterResp, err := WaitStateTransition(ctx, flexClusterParams, connV2.FlexClustersApi, []string{retrystrategy.RetryStrategyCreatingState}, []string{retrystrategy.RetryStrategyIdleState})
	if err != nil {
		resp.Diagnostics.AddError("error waiting for resource to be created", err.Error())
		return
	}

	newFlexClusterModel, diags := NewTFModel(ctx, flexClusterResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if conversion.UseNilForEmpty(tfModel.Tags, newFlexClusterModel.Tags) {
		newFlexClusterModel.Tags = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClusterModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var flexClusterState TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &flexClusterState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	flexCluster, apiResp, err := connV2.FlexClustersApi.GetFlexCluster(ctx, flexClusterState.ProjectId.ValueString(), flexClusterState.Name.ValueString()).Execute()
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newFlexClusterModel, diags := NewTFModel(ctx, flexCluster)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if conversion.UseNilForEmpty(flexClusterState.Tags, newFlexClusterModel.Tags) {
		newFlexClusterModel.Tags = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClusterModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flexClusterReq, diags := NewAtlasUpdateReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	projectID := plan.ProjectId.ValueString()
	clusterName := plan.Name.ValueString()

	connV2 := r.Client.AtlasV2
	_, _, err := connV2.FlexClustersApi.UpdateFlexCluster(ctx, projectID, plan.Name.ValueString(), flexClusterReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: projectID,
		Name:    clusterName,
	}

	flexClusterResp, err := WaitStateTransition(ctx, flexClusterParams, connV2.FlexClustersApi, []string{retrystrategy.RetryStrategyUpdatingState}, []string{retrystrategy.RetryStrategyIdleState})
	if err != nil {
		resp.Diagnostics.AddError("error waiting for resource to be updated", err.Error())
		return
	}

	newFlexClusterModel, diags := NewTFModel(ctx, flexClusterResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if conversion.UseNilForEmpty(plan.Tags, newFlexClusterModel.Tags) {
		newFlexClusterModel.Tags = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClusterModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var flexClusterState *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &flexClusterState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	if _, _, err := connV2.FlexClustersApi.DeleteFlexCluster(ctx, flexClusterState.ProjectId.ValueString(), flexClusterState.Name.ValueString()).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: flexClusterState.ProjectId.ValueString(),
		Name:    flexClusterState.Name.ValueString(),
	}

	if err := WaitStateTransitionDelete(ctx, flexClusterParams, connV2.FlexClustersApi); err != nil {
		resp.Diagnostics.AddError("error waiting for resource to be deleted", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, name, err := splitFlexClusterImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
}

func splitFlexClusterImportID(id string) (projectID, clusterName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a flex cluster, use the format {project_id}-{cluster_name}")
		return
	}

	projectID = &parts[1]
	clusterName = &parts[2]

	return
}
