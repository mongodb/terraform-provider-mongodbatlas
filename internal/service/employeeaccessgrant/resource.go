package employeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const employeeAccessGrantName = "employee_access_grant"

var _ resource.ResourceWithConfigure = &employeeAccessGrantRS{}
var _ resource.ResourceWithImportState = &employeeAccessGrantRS{}

func Resource() resource.Resource {
	return &employeeAccessGrantRS{
		RSCommon: config.RSCommon{
			ResourceName: employeeAccessGrantName,
		},
	}
}

type employeeAccessGrantRS struct {
	config.RSCommon
}

func (r *employeeAccessGrantRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// TODO: Schema and model must be defined in resource_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = ResourceSchema(ctx)
}

func (r *employeeAccessGrantRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var employeeAccessGrantPlan TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &employeeAccessGrantPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// employeeAccessGrantReq, diags := NewEmployeeAccessReq(ctx, &employeeAccessGrantPlan)
	// if diags.HasError() {
	//		resp.Diagnostics.Append(diags...)
	//		return
	//}

	// TODO: make POST request to Atlas API and handle error in response

	// connV2 := r.Client.AtlasV2
	// if err != nil {
	//	resp.Diagnostics.AddError("error creating resource", err.Error())
	//	return
	// }

	// TODO: process response into new terraform state
	// newEmployeeAccessModel, diags := NewTFEmployeeAccess(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	//	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newEmployeeAccessModel)...)
}

func (r *employeeAccessGrantRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var employeeAccessGrantState TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &employeeAccessGrantState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := r.Client.AtlasV2
	// if err != nil {
	//	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
	//		resp.State.RemoveResource(ctx)
	//		return
	//	}
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	// }

	// TODO: process response into new terraform state
	// newEmployeeAccessModel, diags := NewTFEmployeeAccess(ctx, apiResp)
	// if diags.HasError() {
	//	resp.Diagnostics.Append(diags...)
	//	return
	// }
	// resp.Diagnostics.Append(resp.State.Set(ctx, newEmployeeAccessModel)...)
}

func (r *employeeAccessGrantRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var employeeAccessGrantPlan TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &employeeAccessGrantPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// employeeAccessGrantReq, diags := NewEmployeeAccessReq(ctx, &employeeAccessGrantPlan)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	//	return
	// }

	// TODO: make PATCH request to Atlas API and handle error in response
	// connV2 := r.Client.AtlasV2
	// if err != nil {
	//	resp.Diagnostics.AddError("error updating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state

	// newEmployeeAccessModel, diags := NewTFEmployeeAccess(ctx, apiResp)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	//	return
	//}
	// resp.Diagnostics.Append(resp.State.Set(ctx, newEmployeeAccessModel)...)
}

func (r *employeeAccessGrantRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var employeeAccessGrantState *TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &employeeAccessGrantState)...)
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

func (r *employeeAccessGrantRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: parse req.ID string taking into account documented format. Example:

	// projectID, other, err := splitEmployeeAccessGrantImportID(req.ID)
	// if err != nil {
	//	resp.Diagnostics.AddError("error splitting import ID", err.Error())
	//	return
	//}

	// TODO: define attributes that are required for read operation to work correctly. Example:

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
