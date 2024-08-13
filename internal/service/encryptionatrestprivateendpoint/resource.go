//nolint:gocritic
package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	encryptionAtRestPrivateEndpointName = "encryption_at_rest_private_endpoint"
	warnUnsupportedOperation            = "Operation not supported"
)

var _ resource.ResourceWithConfigure = &encryptionAtRestPrivateEndpointRS{}
var _ resource.ResourceWithImportState = &encryptionAtRestPrivateEndpointRS{}

func Resource() resource.Resource {
	return &encryptionAtRestPrivateEndpointRS{
		RSCommon: config.RSCommon{
			ResourceName: encryptionAtRestPrivateEndpointName,
		},
	}
}

type encryptionAtRestPrivateEndpointRS struct {
	config.RSCommon
}

func (r *encryptionAtRestPrivateEndpointRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *encryptionAtRestPrivateEndpointRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var earPrivateEndpointPlan TFEarPrivateEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &earPrivateEndpointPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//  encryptionAtRestPrivateEndpointReq := NewEarPrivateEndpointReq(&earPrivateEndpointPlan)

	// TODO: make POST request to Atlas API and handle error in response
	// connV2 := r.Client.AtlasV2
	// if err != nil {
	//	resp.Diagnostics.AddError("error creating resource", err.Error())
	//	return
	//}

	// TODO: process response into new terraform state
	// newencryptionAtRestPrivateEndpointModel := NewTFencryptionAtRestPrivateEndpoint(apiResp)
	// resp.Diagnostics.Append(resp.State.Set(ctx, newencryptionAtRestPrivateEndpointModel)...)
}

func (r *encryptionAtRestPrivateEndpointRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var earPrivateEndpointState TFEarPrivateEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &earPrivateEndpointState)...)
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
	//}

	// TODO: process response into new terraform state
	//  resp.Diagnostics.Append(resp.State.Set(ctx, NewTFEarPrivateEndpoint(apiResp))...)
}

func (r *encryptionAtRestPrivateEndpointRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(warnUnsupportedOperation, "Updating the private endpoint for encryption at rest is not supported. To modify your infrastructure, please delete the existing mongodbatlas_encryption_at_rest_private_endpoint resource and create a new one with the necessary updates")
}

func (r *encryptionAtRestPrivateEndpointRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var earPrivateEndpointState *TFEarPrivateEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &earPrivateEndpointState)...)
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

func (r *encryptionAtRestPrivateEndpointRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: parse req.ID string taking into account documented format. Example:

	// projectID, other, err := splitencryptionAtRestPrivateEndpointImportID(req.ID)
	// if err != nil {
	//	resp.Diagnostics.AddError("error splitting import ID", err.Error())
	//	return
	//}

	// TODO: define attributes that are required for read operation to work correctly. Example:

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
