package projectipaccesslist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorAccessListCreate = "error creating Project IP Access List information: %s"
	errorAccessListRead   = "error getting Project IP Access List information: %s"
	errorAccessListUpdate = "error updating Project IP Access List information: %s"
	errorAccessListDelete = "error deleting Project IP Access List information: %s"
	timeoutCreateDelete   = 45 * time.Minute
	timeoutRead           = 2 * time.Minute
	timeoutUpdate         = 45 * time.Minute
	timeoutRetryItem      = 2 * time.Minute
	minTimeoutCreate      = 10 * time.Second
	delayCreate           = 10 * time.Second
)

type projectIPAccessListRS struct {
	config.RSCommon
}

func Resource() resource.Resource {
	return &projectIPAccessListRS{
		RSCommon: config.RSCommon{
			ResourceName: projectIPAccessList,
		},
	}
}

var _ resource.ResourceWithConfigure = &projectIPAccessListRS{}
var _ resource.ResourceWithImportState = &projectIPAccessListRS{}

func (r *projectIPAccessListRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *projectIPAccessListRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var projectIPAccessListModel *TfProjectIPAccessListModel

	diags := req.Plan.Get(ctx, &projectIPAccessListModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if projectIPAccessListModel.CIDRBlock.IsNull() && projectIPAccessListModel.IPAddress.IsNull() && projectIPAccessListModel.AWSSecurityGroup.IsNull() {
		resp.Diagnostics.AddError("validation error", "cidr_block, ip_address or aws_security_group needs to contain a value")
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := projectIPAccessListModel.ProjectID.ValueString()
	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created", "failed"},
		Refresh: func() (any, string, error) {
			_, _, err := connV2.ProjectIPAccessListApi.CreateAccessListEntry(ctx, projectID, NewMongoDBProjectIPAccessList(projectIPAccessListModel)).Execute()
			// Atlas Create is called inside refresh because this limitation: This endpoint doesn't support concurrent POST requests. You must submit multiple POST requests synchronously.
			if err != nil {
				if strings.Contains(err.Error(), "Unexpected error") ||
					strings.Contains(err.Error(), "UNEXPECTED_ERROR") ||
					strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorAccessListCreate, err)
			}

			accessListEntry := getAccessListEntry(projectIPAccessListModel)

			entry, exists, err := isEntryInProjectAccessList(ctx, connV2, projectID, accessListEntry)
			if err != nil {
				if strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				if strings.Contains(err.Error(), "404") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorAccessListCreate, err)
			}
			if !exists {
				return nil, "pending", nil
			}

			return entry, "created", nil
		},
		Timeout:    timeoutCreateDelete,
		Delay:      delayCreate,
		MinTimeout: minTimeoutCreate,
	}

	// Wait, catching any errors
	accessList, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error while waiting for resource creation", err.Error())
		return
	}

	entry, ok := accessList.(*admin.NetworkPermissionEntry)
	if !ok {
		resp.Diagnostics.AddError("error", errorAccessListCreate)
		return
	}

	projectIPAccessListNewModel := NewTfProjectIPAccessListModel(projectIPAccessListModel, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectIPAccessListNewModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectIPAccessListRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var projectIPAccessListModelState *TfProjectIPAccessListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &projectIPAccessListModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	decodedIDMap := conversion.DecodeStateID(projectIPAccessListModelState.ID.ValueString())
	if len(decodedIDMap) != 2 {
		resp.Diagnostics.AddError("error during the reading operation", "the provided resource ID is not correct")
		return
	}

	timeout, diags := projectIPAccessListModelState.Timeouts.Read(ctx, timeoutRead)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		accessList, httpResponse, err := connV2.ProjectIPAccessListApi.GetAccessListEntry(ctx, decodedIDMap["project_id"], decodedIDMap["entry"]).Execute()
		if err != nil {
			// case 404
			// deleted in the backend case
			if validate.StatusNotFound(httpResponse) {
				resp.State.RemoveResource(ctx)
				return nil
			}

			if validate.StatusInternalServerError(httpResponse) {
				return retry.RetryableError(err)
			}

			resp.Diagnostics.AddError("error getting project ip access list information", err.Error())
			return nil
		}

		projectIPAccessListNewModel := NewTfProjectIPAccessListModel(projectIPAccessListModelState, accessList)
		resp.Diagnostics.Append(resp.State.Set(ctx, &projectIPAccessListNewModel)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error during the read operation", err.Error())
	}
}

func (r *projectIPAccessListRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var projectIPAccessListModelState *TfProjectIPAccessListModel

	resp.Diagnostics.Append(req.State.Get(ctx, &projectIPAccessListModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry := getAccessListEntry(projectIPAccessListModelState)

	connV2 := r.Client.AtlasV2
	projectID := projectIPAccessListModelState.ProjectID.ValueString()

	timeout, diags := projectIPAccessListModelState.Timeouts.Delete(ctx, timeoutCreateDelete)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		httpResponse, err := connV2.ProjectIPAccessListApi.DeleteAccessListEntry(ctx, projectID, entry).Execute()
		if err != nil {
			if validate.StatusInternalServerError(httpResponse) {
				return retry.RetryableError(err)
			}

			if validate.StatusNotFound(httpResponse) {
				return nil
			}

			resp.Diagnostics.AddError("error deleting the entry", fmt.Sprintf(errorAccessListDelete, err.Error()))
			return retry.NonRetryableError(fmt.Errorf(errorAccessListDelete, err))
		}

		entry, httpResponse, err := connV2.ProjectIPAccessListApi.GetAccessListEntry(ctx, projectID, entry).Execute()
		if err != nil {
			if validate.StatusNotFound(httpResponse) {
				return nil
			}

			return retry.RetryableError(err)
		}

		if entry != nil {
			return retry.RetryableError(fmt.Errorf(errorAccessListDelete, "Access list still exists"))
		}

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error during the read operation", err.Error())
	}
}

func (r *projectIPAccessListRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "-", 2)

	if len(parts) != 2 {
		resp.Diagnostics.AddError("import format error", "to import a projectIP Access List, use the format {project_id}-{entry}")
		return
	}

	projectID := parts[0]
	entry := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), conversion.EncodeStateID(map[string]string{
		"entry":      entry,
		"project_id": projectID,
	}))...)
}

func isEntryInProjectAccessList(ctx context.Context, connV2 *admin.APIClient, projectID, entry string) (*admin.NetworkPermissionEntry, bool, error) {
	var out admin.NetworkPermissionEntry
	err := retry.RetryContext(ctx, timeoutRetryItem, func() *retry.RetryError {
		accessList, httpResponse, err := connV2.ProjectIPAccessListApi.GetAccessListEntry(ctx, projectID, entry).Execute()
		if err != nil {
			switch {
			case validate.StatusInternalServerError(httpResponse):
				return retry.RetryableError(err)
			case validate.StatusNotFound(httpResponse):
				return retry.RetryableError(err)
			default:
				return retry.NonRetryableError(fmt.Errorf(errorAccessListRead, err))
			}
		}

		out = *accessList
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return &out, true, nil
}

// Update is supported only for the comment field, the rest of the fields will trigger a replace.
func (r *projectIPAccessListRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var projectIPAccessListState *TfProjectIPAccessListModel
	var projectIPAccessListPlan *TfProjectIPAccessListModel
	connV2 := r.Client.AtlasV2

	resp.Diagnostics.Append(req.State.Get(ctx, &projectIPAccessListState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &projectIPAccessListPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := projectIPAccessListState.ProjectID.ValueString()

	updatedProjectIPAccessList := &TfProjectIPAccessListModel{
		ID:               projectIPAccessListState.ID,
		ProjectID:        projectIPAccessListState.ProjectID,
		CIDRBlock:        projectIPAccessListState.CIDRBlock,
		IPAddress:        projectIPAccessListState.IPAddress,
		AWSSecurityGroup: projectIPAccessListState.AWSSecurityGroup,

		// Only comment and timeouts can be updated without replace.
		Comment:  projectIPAccessListPlan.Comment,
		Timeouts: projectIPAccessListPlan.Timeouts,
	}

	timeout, diags := projectIPAccessListPlan.Timeouts.Update(ctx, timeoutUpdate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"updated", "failed"},
		Refresh: func() (any, string, error) {
			_, _, err := connV2.ProjectIPAccessListApi.CreateAccessListEntry(ctx, projectID, NewMongoDBProjectIPAccessList(updatedProjectIPAccessList)).Execute()
			if err != nil {
				if strings.Contains(err.Error(), "Unexpected error") ||
					strings.Contains(err.Error(), "UNEXPECTED_ERROR") ||
					strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorAccessListUpdate, err)
			}

			accessListEntry := getAccessListEntry(updatedProjectIPAccessList)

			entry, exists, err := isEntryInProjectAccessList(ctx, connV2, projectID, accessListEntry)
			if err != nil {
				if strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				if strings.Contains(err.Error(), "404") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorAccessListUpdate, err)
			}
			if !exists {
				return nil, "pending", nil
			}

			return entry, "updated", nil
		},
		Timeout:    timeout,
		Delay:      delayCreate,
		MinTimeout: minTimeoutCreate,
	}

	// Wait, catching any errors
	accessList, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error while waiting for resource update", err.Error())
		return
	}

	entry, ok := accessList.(*admin.NetworkPermissionEntry)
	if !ok {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unexpected type %T returned from state change, expected *admin.NetworkPermissionEntry", accessList))
		return
	}

	projectIPAccessListNewModel := NewTfProjectIPAccessListModel(updatedProjectIPAccessList, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectIPAccessListNewModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
