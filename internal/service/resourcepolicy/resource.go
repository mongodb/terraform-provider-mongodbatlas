package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const resourcePolicyName = "resource_policy"

var _ resource.ResourceWithConfigure = &resourcePolicyRS{}
var _ resource.ResourceWithImportState = &resourcePolicyRS{}

func Resource() resource.Resource {
	return &resourcePolicyRS{
		RSCommon: config.RSCommon{
			ResourceName: resourcePolicyName,
		},
	}
}

type resourcePolicyRS struct {
	config.RSCommon
}

func (r *resourcePolicyRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *resourcePolicyRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resourcePolicyPlan TFResourcePolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resourcePolicyPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// resourcePolicyReq, diags := NewResourcePolicyReq(ctx, &resourcePolicyPlan)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// TODO: make POST request to Atlas API and handle error in response

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error creating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	// newResourcePolicyModel, diags := NewTFResourcePolicy(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var resourcePolicyState TFResourcePolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resourcePolicyState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
	//		resp.State.RemoveResource(ctx)
	//		return
	//	}
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	// newResourcePolicyModel, diags := NewTFResourcePolicy(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resourcePolicyPlan TFResourcePolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resourcePolicyPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// resourcePolicyReq, diags := NewResourcePolicyReq(ctx, &resourcePolicyPlan)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// TODO: make PATCH request to Atlas API and handle error in response
	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error updating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state

	// newResourcePolicyModel, diags := NewTFResourcePolicy(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourcePolicyState *TFResourcePolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resourcePolicyState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make Delete request to Atlas API

	// connV2 := r.Client.AtlasV2
	// if _, _, err := connV2.Api.Delete().Execute(); err != nil {
	// 	 resp.Diagnostics.AddError("error deleting resource", err.Error())
	// 	 return
	// }
}

func (r *resourcePolicyRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: parse req.ID string taking into account documented format. Example:

	// projectID, other, err := splitResourcePolicyImportID(req.ID)
	// if err != nil {
	//	resp.Diagnostics.AddError("error splitting import ID", err.Error())
	//	return
	//}

	// TODO: define attributes that are required for read operation to work correctly. Example:

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
