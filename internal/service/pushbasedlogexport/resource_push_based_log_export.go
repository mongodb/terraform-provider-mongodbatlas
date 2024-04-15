package pushbasedlogexport

import (
    "context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const pushBasedLogExportName = "push_based_log_export"

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
    // TODO: Schema and model must be defined in resource_push_based_log_export_schema.go. Details on scaffolding this file found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
	resp.Schema = ResourceSchema(ctx)
}

func (r *pushBasedLogExportRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var pushBasedLogExportPlan TFPushBasedLogExportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &pushBasedLogExportPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

    pushBasedLogExportReq, diags := NewPushBasedLogExportReq(ctx, &pushBasedLogExportPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	
	// TODO: make POST request to Atlas API and handle error in response

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error creating resource", err.Error())
	//	return
	//}
	
    // TODO: process response into new terraform state
	newPushBasedLogExportModel, diags := NewTFPushBasedLogExport(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newPushBasedLogExportModel)...)
}

func (r *pushBasedLogExportRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var pushBasedLogExportState TFPushBasedLogExportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &pushBasedLogExportState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: make get request to resource

	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error fetching resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	newPushBasedLogExportModel, diags := NewTFPushBasedLogExport(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newPushBasedLogExportModel)...)
}

func (r *pushBasedLogExportRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var pushBasedLogExportPlan TFPushBasedLogExportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &pushBasedLogExportPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pushBasedLogExportReq, diags := NewPushBasedLogExportReq(ctx, &pushBasedLogExportPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// TODO: make PATCH request to Atlas API and handle error in response
	// connV2 := r.Client.AtlasV2
	//if err != nil {
	//	resp.Diagnostics.AddError("error updating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state

	newPushBasedLogExportModel, diags := NewTFPushBasedLogExport(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newPushBasedLogExportModel)...)
}

func (r *pushBasedLogExportRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var pushBasedLogExportState *TFPushBasedLogExportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &pushBasedLogExportState)...)
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

func (r *pushBasedLogExportRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: parse req.ID string taking into account documented format. Example:
	
	// projectID, other, err := splitPushBasedLogExportImportID(req.ID)
	// if err != nil {
	//	resp.Diagnostics.AddError("error splitting import ID", err.Error())
	//	return
	//}

	// TODO: define attributes that are required for read operation to work correctly. Example:

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
