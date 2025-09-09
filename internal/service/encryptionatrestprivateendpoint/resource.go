package encryptionatrestprivateendpoint

import (
	"context"
	"errors"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	encryptionAtRestPrivateEndpointName = "encryption_at_rest_private_endpoint"
	warnUnsupportedOperation            = "Operation not supported"
	FailedStatusErrorMessageSummary     = "Private endpoint is in a failed status"
	NonEmptyErrorMessageFieldSummary    = "Something went wrong. Please review the `status` field of this resource"
	PendingAcceptanceWarnMsgSummary     = "Private endpoint may be in PENDING_ACCEPTANCE status"
	PendingAcceptanceWarnMsg            = "Please ensure to approve the private endpoint connection. If recently approved or deleted the endpoint, please ignore this warning & wait for a few minutes for the status to update."
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
	conversion.UpdateSchemaDescription(&resp.Schema)
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
	createResp, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.CreateRestPrivateEndpoint(ctx, projectID, cloudProvider, privateEndpointReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	createTimeout := cleanup.ResolveTimeout(ctx, &earPrivateEndpointPlan.Timeouts, cleanup.OperationCreate, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	finalResp, err := waitStateTransition(ctx, projectID, cloudProvider, createResp.GetId(), connV2.EncryptionAtRestUsingCustomerKeyManagementApi, createTimeout)
	err = cleanup.HandleCreateTimeout(cleanup.ResolveDeleteOnCreateTimeout(earPrivateEndpointPlan.DeleteOnCreateTimeout), err, func(ctxCleanup context.Context) error {
		cleanResp, cleanErr := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.RequestPrivateEndpointDeletion(ctxCleanup, projectID, cloudProvider, createResp.GetId()).Execute()
		if validate.StatusNotFound(cleanResp) {
			return nil
		}
		return cleanErr
	})

	if err != nil {
		resp.Diagnostics.AddError("error when waiting for status transition in creation", err.Error())
		return
	}

	privateEndpointModel := NewTFEarPrivateEndpoint(*finalResp, projectID)
	privateEndpointModel.Timeouts = earPrivateEndpointPlan.Timeouts
	privateEndpointModel.DeleteOnCreateTimeout = earPrivateEndpointPlan.DeleteOnCreateTimeout
	resp.Diagnostics.Append(resp.State.Set(ctx, privateEndpointModel)...)

	diags := CheckErrorMessageAndStatus(finalResp)
	resp.Diagnostics.Append(diags...)
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

	endpointModel, apiResp, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetRestPrivateEndpoint(ctx, projectID, cloudProvider, endpointID).Execute()
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	privateEndpointModel := NewTFEarPrivateEndpoint(*endpointModel, projectID)
	privateEndpointModel.Timeouts = earPrivateEndpointState.Timeouts
	privateEndpointModel.DeleteOnCreateTimeout = earPrivateEndpointState.DeleteOnCreateTimeout
	resp.Diagnostics.Append(resp.State.Set(ctx, privateEndpointModel)...)

	diags := CheckErrorMessageAndStatus(endpointModel)
	resp.Diagnostics.Append(diags...)
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
	if _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.RequestPrivateEndpointDeletion(ctx, projectID, cloudProvider, endpointID).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}

	deleteTimeout := cleanup.ResolveTimeout(ctx, &earPrivateEndpointState.Timeouts, cleanup.OperationDelete, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := WaitDeleteStateTransition(ctx, projectID, cloudProvider, endpointID, connV2.EncryptionAtRestUsingCustomerKeyManagementApi, deleteTimeout)
	if err != nil {
		resp.Diagnostics.AddError("error when waiting for status transition in delete", err.Error())
		return
	}

	diags := CheckErrorMessageAndStatus(model)
	resp.Diagnostics.Append(diags...)
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

func CheckErrorMessageAndStatus(model *admin.EARPrivateEndpoint) diag.Diagnostics {
	var diags diag.Diagnostics

	switch {
	case model.GetStatus() == retrystrategy.RetryStrategyFailedState:
		diags = append(diags, diag.NewErrorDiagnostic(FailedStatusErrorMessageSummary, model.GetErrorMessage()))
	case model.GetErrorMessage() != "":
		diags = append(diags, diag.NewWarningDiagnostic(NonEmptyErrorMessageFieldSummary, model.GetErrorMessage()))
	case model.GetStatus() == retrystrategy.RetryStrategyPendingAcceptanceState:
		diags = append(diags, diag.NewWarningDiagnostic(PendingAcceptanceWarnMsgSummary, PendingAcceptanceWarnMsg))
	}

	return diags
}
