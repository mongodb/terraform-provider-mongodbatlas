package encryptionatrestprivateendpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

const (
	encryptionAtRestPrivateEndpointName = "encryption_at_rest_private_endpoint"
	warnUnsupportedOperation            = "Operation not supported"
	failedStatusErrorMessage            = "Private endpoint is in a failed status"
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

	privateEndpointReq := NewEarPrivateEndpointReq(&earPrivateEndpointPlan)
	connV2 := r.Client.AtlasV2
	projectID := earPrivateEndpointPlan.ProjectID.ValueString()
	cloudProvider := earPrivateEndpointPlan.CloudProvider.ValueString()
	createResp, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.CreateEncryptionAtRestPrivateEndpoint(ctx, projectID, cloudProvider, privateEndpointReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	finalResp, err := waitStateTransition(ctx, projectID, cloudProvider, createResp.GetId(), connV2.EncryptionAtRestUsingCustomerKeyManagementApi)
	if err != nil {
		resp.Diagnostics.AddError("error when waiting for status transition in creation", err.Error())
		return
	}

	privateEndpointModel := NewTFEarPrivateEndpoint(*finalResp, projectID)
	resp.Diagnostics.Append(resp.State.Set(ctx, privateEndpointModel)...)
	if err := getErrorMsgForFailedStatus(finalResp); err != nil {
		resp.Diagnostics.AddError(failedStatusErrorMessage, err.Error())
	}
}

func (r *encryptionAtRestPrivateEndpointRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var earPrivateEndpointState TFEarPrivateEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &earPrivateEndpointState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := earPrivateEndpointState.ProjectID.ValueString()
	cloudProvider := earPrivateEndpointState.CloudProvider.ValueString()
	endpointID := earPrivateEndpointState.ID.ValueString()

	endpointModel, apiResp, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRestPrivateEndpoint(ctx, projectID, cloudProvider, endpointID).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFEarPrivateEndpoint(*endpointModel, projectID))...)
	if err := getErrorMsgForFailedStatus(endpointModel); err != nil {
		resp.Diagnostics.AddError(failedStatusErrorMessage, err.Error())
	}
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

	connV2 := r.Client.AtlasV2
	projectID := earPrivateEndpointState.ProjectID.ValueString()
	cloudProvider := earPrivateEndpointState.CloudProvider.ValueString()
	endpointID := earPrivateEndpointState.ID.ValueString()
	if _, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.RequestEncryptionAtRestPrivateEndpointDeletion(ctx, projectID, cloudProvider, endpointID).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}

	model, err := WaitDeleteStateTransition(ctx, projectID, cloudProvider, endpointID, connV2.EncryptionAtRestUsingCustomerKeyManagementApi)
	if err != nil {
		resp.Diagnostics.AddError("error when waiting for status transition in delete", err.Error())
		return
	}
	if err := getErrorMsgForFailedStatus(model); err != nil {
		resp.Diagnostics.AddError(failedStatusErrorMessage, err.Error())
	}
}

func (r *encryptionAtRestPrivateEndpointRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, cloudProvider, privateEndpointID, err := splitEncryptionAtRestPrivateEndpointImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cloud_provider"), cloudProvider)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), privateEndpointID)...)
}

func splitEncryptionAtRestPrivateEndpointImportID(id string) (projectID, cloudProvider, privateEndpointID string, err error) {
	re := regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("use the format {project_id}-{cloud_provider}-{private_endpoint_id}")
		return
	}

	projectID = parts[1]
	cloudProvider = parts[2]
	privateEndpointID = parts[3]
	return
}

func getErrorMsgForFailedStatus(model *admin.EARPrivateEndpoint) error {
	if model.GetStatus() != retrystrategy.RetryStrategyFailedState {
		return nil
	}
	msg := model.GetErrorMessage()
	return fmt.Errorf("detail of error message: %s", msg)
}
