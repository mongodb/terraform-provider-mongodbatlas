package mongodbemployeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName      = "mongodb_employee_access_grant"
	fullResourceName  = "mongodbatlas_" + resourceName
	errorCreateUpdate = "Error setting resource " + fullResourceName
	errorRead         = "Error retrieving info for resource " + fullResourceName
	errorDataSource   = "Error retrieving info for data source " + fullResourceName
	errorDelete       = "Error deleting resource " + fullResourceName
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
	r.createOrUpdate(ctx, req.Plan.Get, &resp.Diagnostics, &resp.State)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	cluster, httpResp, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if validate.StatusNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	apiResp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	if apiResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(projectID, clusterName, apiResp))...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.createOrUpdate(ctx, req.Plan.Get, &resp.Diagnostics, &resp.State)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	httpResp, err := connV2.ClustersApi.RevokeMongoEmployeeAccess(ctx, projectID, clusterName).Execute()
	if err != nil && validate.StatusNotFound(httpResp) {
		resp.Diagnostics.AddError(errorDelete, err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conversion.ImportStateProjectIDClusterName(ctx, req, resp, "project_id", "cluster_name")
}

func (r *rs) createOrUpdate(ctx context.Context, tfModelFunc func(context.Context, any) diag.Diagnostics, diagnostics *diag.Diagnostics, state *tfsdk.State) {
	var tfModel TFModel
	diagnostics.Append(tfModelFunc(ctx, &tfModel)...)
	if diagnostics.HasError() {
		return
	}
	atlasReq, err := NewAtlasReq(&tfModel)
	if err != nil {
		diagnostics.AddError(errorCreateUpdate, err.Error())
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	if _, err := connV2.ClustersApi.GrantMongoEmployeeAccess(ctx, projectID, clusterName, atlasReq).Execute(); err != nil {
		diagnostics.AddError(errorCreateUpdate, err.Error())
		return
	}
	diagnostics.Append(state.Set(ctx, tfModel)...)
}
