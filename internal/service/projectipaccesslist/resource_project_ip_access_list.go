package projectipaccesslist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/concurrency"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

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
	errorAccessListDelete = "error deleting Project IP Access List information: %s"
	timeoutCreateDelete   = 45 * time.Minute
	timeoutRead           = 2 * time.Minute
	timeoutRetryItem      = 2 * time.Minute
	delayCreate           = 10 * time.Second
	minTimeoutCreate      = 10 * time.Second
)

var createAccessListEntryMutex = concurrency.NewMutexKV()

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
	timeout, diags := projectIPAccessListModel.Timeouts.Create(ctx, timeoutCreateDelete)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := createEntry(ctx, connV2, projectIPAccessListModel, timeout, errorAccessListCreate)
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}
	if entry == nil {
		resp.Diagnostics.AddError("error", fmt.Errorf(errorAccessListCreate, "entry is nil").Error())
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

	entry := projectIPAccessListModelState.CIDRBlock.ValueString()
	if projectIPAccessListModelState.IPAddress.ValueString() != "" {
		entry = projectIPAccessListModelState.IPAddress.ValueString()
	} else if projectIPAccessListModelState.AWSSecurityGroup.ValueString() != "" {
		entry = projectIPAccessListModelState.AWSSecurityGroup.ValueString()
	}

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

func createEntry(ctx context.Context, connV2 *admin.APIClient, projectIPAccessListModel *TfProjectIPAccessListModel, timeout time.Duration, errorMsg string) (*admin.NetworkPermissionEntry, error) {
	projectID := projectIPAccessListModel.ProjectID.ValueString()
	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created", "failed"},
		Refresh: func() (any, string, error) {
			// Each access list entry is its own resource, which leads to concurrent calls within a single apply unless explicitly handled by the user.
			// From API docs: "This endpoint doesn't support concurrent POST requests. You must submit multiple POST requests synchronously."
			// Locking on a project level to avoid race conditions within a single apply. Still, we verify that the entry was added to the access list and retry otherwise in case of an external update.
			createAccessListEntryMutex.Lock(projectID)
			_, httpResponse, err := connV2.ProjectIPAccessListApi.CreateAccessListEntry(ctx, projectID, NewMongoDBProjectIPAccessList(projectIPAccessListModel)).Execute()
			// Unlock immediately after create to allow parallel reads (intentionally not deferring).
			createAccessListEntryMutex.Unlock(projectID)
			if err != nil {
				if validate.StatusInternalServerError(httpResponse) {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorMsg, err)
			}

			accessListEntry := projectIPAccessListModel.IPAddress.ValueString()
			if projectIPAccessListModel.CIDRBlock.ValueString() != "" {
				accessListEntry = projectIPAccessListModel.CIDRBlock.ValueString()
			} else if projectIPAccessListModel.AWSSecurityGroup.ValueString() != "" {
				accessListEntry = projectIPAccessListModel.AWSSecurityGroup.ValueString()
			}

			entry, exists, err := isEntryInProjectAccessList(ctx, connV2, projectID, accessListEntry)
			if err != nil {
				if strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				if strings.Contains(err.Error(), "404") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorMsg, err)
			}
			if !exists {
				return nil, "pending", nil
			}

			return entry, "created", nil
		},
		Timeout:    timeout,
		Delay:      delayCreate,
		MinTimeout: minTimeoutCreate,
	}

	// Wait, catching any errors
	accessList, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	entry, ok := accessList.(*admin.NetworkPermissionEntry)
	if !ok {
		return nil, fmt.Errorf(errorMsg, "invalid result type")
	}

	return entry, nil
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
