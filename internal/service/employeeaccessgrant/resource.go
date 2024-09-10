package employeeaccessgrant

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName     = "employee_access_grant"
	fullResourceName = "mongodbatlas_" + resourceName
	errorCreate      = "Error creating resource " + fullResourceName
	errorRead        = "Error retrieving info for resource " + fullResourceName
	errorDelete      = "Error deleting resource " + fullResourceName
)

var _ resource.ResourceWithConfigure = &employeeAccessGrantRS{}
var _ resource.ResourceWithImportState = &employeeAccessGrantRS{}

func Resource() resource.Resource {
	return &employeeAccessGrantRS{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type employeeAccessGrantRS struct {
	config.RSCommon
}

func (r *employeeAccessGrantRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *employeeAccessGrantRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	atlasReq, err := NewAtlasReq(&tfModel)
	if err != nil {
		resp.Diagnostics.AddError(errorCreate, err.Error())
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	if _, _, err := connV2.ClustersApi.GrantMongoDBEmployeeAccess(ctx, projectID, clusterName, atlasReq).Execute(); err != nil {
		resp.Diagnostics.AddError(errorCreate, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
}

func (r *employeeAccessGrantRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	cluster, httpResp, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if httpResp.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(errorCreate, err.Error())
		return
	}
	atlasResp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	if atlasResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(projectID, clusterName, atlasResp))...)
}

func (r *employeeAccessGrantRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var tfPlan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfPlan)...)
}

func (r *employeeAccessGrantRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := r.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	_, httpResp, err := connV2.ClustersApi.RevokeMongoDBEmployeeAccess(ctx, projectID, clusterName).Execute()
	if err != nil && httpResp.StatusCode != http.StatusNotFound {
		resp.Diagnostics.AddError(errorDelete, err.Error())
		return
	}
}

func (r *employeeAccessGrantRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
