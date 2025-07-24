package flexcluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName                = "flex_cluster"
	ErrorUpdateNotAllowed       = "update not allowed"
	FlexClusterType             = "FLEX"
	ErrorCreateFlex             = "error creating flex cluster: %s"
	ErrorReadFlex               = "error reading flex cluster (%s): %s"
	ErrorUpdateFlex             = "error updating flex cluster: %s"
	ErrorUpgradeFlex            = "error upgrading to a flex cluster: %s"
	ErrorDeleteFlex             = "error deleting a flex cluster (%s): %s"
	ErrorNonUpdatableAttributes = "flex cluster update is not supported except for tags and termination_protection_enabled fields"
)

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
	flexClusterResp, err := CreateFlexClusterNew(ctx, projectID, clusterName, flexClusterReq, connV2.FlexClustersApi)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(ErrorCreateFlex, err.Error()), fmt.Sprintf("Name: %s, Project ID: %s", clusterName, projectID))
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
	projectID := flexClusterState.ProjectId.ValueString()
	clusterName := flexClusterState.Name.ValueString()
	flexCluster, apiResp, err := connV2.FlexClustersApi.GetFlexCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf(ErrorReadFlex, projectID, clusterName), err.Error())
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

	flexClusterResp, err := UpdateFlexCluster(ctx, projectID, clusterName, flexClusterReq, connV2.FlexClustersApi)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(ErrorUpdateFlex, clusterName), err.Error())
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

	projectID := flexClusterState.ProjectId.ValueString()
	clusterName := flexClusterState.Name.ValueString()
	err := DeleteFlexCluster(ctx, projectID, clusterName, connV2.FlexClustersApi)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(ErrorDeleteFlex, projectID, clusterName), err.Error())
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

// TODO: using CreateFlexClusterNew to avoid changes in adv_cluster and running their acc tests while doing this PR, rename to CreateFlexCluster and update adv_cluster call params before merging.
func CreateFlexClusterNew(ctx context.Context, projectID, clusterName string, flexClusterReq *admin.FlexClusterDescriptionCreate20241113, client admin.FlexClustersApi) (*admin.FlexClusterDescription20241113, error) {
	_, _, err := client.CreateFlexCluster(ctx, projectID, flexClusterReq).Execute()
	if err != nil {
		return nil, err
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: projectID,
		Name:    clusterName,
	}

	flexClusterResp, err := WaitStateTransition(ctx, flexClusterParams, client, []string{retrystrategy.RetryStrategyCreatingState, retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyRepairingState}, []string{retrystrategy.RetryStrategyIdleState}, false, nil)
	if err != nil {
		return nil, err
	}
	return flexClusterResp, nil
}

// TODO: keeping CreateFlexCluster to avoid changes in adv_cluster and running their acc tests while doing this PR, remove before merging.
func CreateFlexCluster(ctx context.Context, projectID, clusterName string, flexClusterReq *admin.FlexClusterDescriptionCreate20241113, client admin.FlexClustersApi) (*admin.FlexClusterDescription20241113, error) {
	return CreateFlexClusterNew(ctx, projectID, clusterName, flexClusterReq, client)
}

func GetFlexCluster(ctx context.Context, projectID, clusterName string, client admin.FlexClustersApi) (*admin.FlexClusterDescription20241113, error) {
	flexCluster, _, err := client.GetFlexCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, err
	}
	return flexCluster, nil
}

func UpdateFlexCluster(ctx context.Context, projectID, clusterName string, flexClusterReq *admin.FlexClusterDescriptionUpdate20241113, client admin.FlexClustersApi) (*admin.FlexClusterDescription20241113, error) {
	_, _, err := client.UpdateFlexCluster(ctx, projectID, clusterName, flexClusterReq).Execute()
	if err != nil {
		return nil, err
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: projectID,
		Name:    clusterName,
	}

	flexClusterResp, err := WaitStateTransition(ctx, flexClusterParams, client, []string{retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyRepairingState}, []string{retrystrategy.RetryStrategyIdleState}, false, nil)
	if err != nil {
		return nil, err
	}
	return flexClusterResp, nil
}

func DeleteFlexCluster(ctx context.Context, projectID, clusterName string, client admin.FlexClustersApi) error {
	if _, err := client.DeleteFlexCluster(ctx, projectID, clusterName).Execute(); err != nil {
		return err
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: projectID,
		Name:    clusterName,
	}

	return WaitStateTransitionDelete(ctx, flexClusterParams, client)
}

func ListFlexClusters(ctx context.Context, projectID string, client admin.FlexClustersApi) (*[]admin.FlexClusterDescription20241113, error) {
	params := admin.ListFlexClustersApiParams{
		GroupId: projectID,
	}

	flexClusters, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.FlexClusterDescription20241113], *http.Response, error) {
		request := client.ListFlexClustersWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})

	if err != nil {
		return nil, err
	}

	return &flexClusters, nil
}
