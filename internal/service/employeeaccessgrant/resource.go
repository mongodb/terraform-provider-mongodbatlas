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
	resp.Schema = ResourceSchema(ctx)
}

func (r *employeeAccessGrantRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfPlan TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfPlan)...)
}

func (r *employeeAccessGrantRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *employeeAccessGrantRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var tfPlan TFEmployeeAccessGrantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tfPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfPlan)...)
}

func (r *employeeAccessGrantRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *employeeAccessGrantRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
