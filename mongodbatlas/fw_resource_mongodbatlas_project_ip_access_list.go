package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	cstmvalidator "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/validator"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorAccessListCreate          = "error creating Project IP Access List information: %s"
	errorAccessListRead            = "error getting Project IP Access List information: %s"
	errorAccessListDelete          = "error deleting Project IP Access List information: %s"
	projectIPAccessListTimeout     = 45 * time.Minute
	projectIPAccessListTimeoutRead = 2 * time.Minute
	projectIPAccessListMinTimeout  = 2 * time.Second
	projectIPAccessListDelay       = 4 * time.Second
	projectIPAccessListRetry       = 2 * time.Minute
)

type tfProjectIPAccessListModel struct {
	ID               types.String   `tfsdk:"id"`
	ProjectID        types.String   `tfsdk:"project_id"`
	CIDRBlock        types.String   `tfsdk:"cidr_block"`
	IPAddress        types.String   `tfsdk:"ip_address"`
	AWSSecurityGroup types.String   `tfsdk:"aws_security_group"`
	Comment          types.String   `tfsdk:"comment"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type ProjectIPAccessListRS struct {
	RSCommon
}

func NewProjectIPAccessListRS() resource.Resource {
	return &ProjectIPAccessListRS{
		RSCommon: RSCommon{
			resourceName: projectIPAccessList,
		},
	}
}

var _ resource.ResourceWithConfigure = &ProjectIPAccessListRS{}
var _ resource.ResourceWithImportState = &ProjectIPAccessListRS{}

func (r *ProjectIPAccessListRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr_block": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cstmvalidator.ValidCIDR(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("ip_address"),
					}...),
				},
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					cstmvalidator.ValidIP(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("aws_security_group"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
			},
			"aws_security_group": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("ip_address"),
						path.MatchRelative().AtParent().AtName("cidr_block"),
					}...),
				},
			},
			"comment": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Delete: true,
				Read:   true,
			}),
		},
	}
}

func (r *ProjectIPAccessListRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var projectIPAccessListModel *tfProjectIPAccessListModel

	diags := req.Plan.Get(ctx, &projectIPAccessListModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if projectIPAccessListModel.CIDRBlock.IsNull() && projectIPAccessListModel.IPAddress.IsNull() && projectIPAccessListModel.AWSSecurityGroup.IsNull() {
		resp.Diagnostics.AddError("validation error", "cidr_block, ip_address or aws_security_group needs to contain a value")
		return
	}

	conn := r.client.Atlas
	projectID := projectIPAccessListModel.ProjectID.ValueString()
	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created", "failed"},
		Refresh: func() (interface{}, string, error) {
			_, _, err := conn.ProjectIPAccessList.Create(ctx, projectID, newMongoDBProjectIPAccessList(projectIPAccessListModel))
			if err != nil {
				if strings.Contains(err.Error(), "Unexpected error") ||
					strings.Contains(err.Error(), "UNEXPECTED_ERROR") ||
					strings.Contains(err.Error(), "500") {
					return nil, "pending", nil
				}
				return nil, "failed", fmt.Errorf(errorAccessListCreate, err)
			}

			accessListEntry := projectIPAccessListModel.IPAddress.ValueString()
			if projectIPAccessListModel.CIDRBlock.ValueString() != "" {
				accessListEntry = projectIPAccessListModel.CIDRBlock.ValueString()
			} else if projectIPAccessListModel.AWSSecurityGroup.ValueString() != "" {
				accessListEntry = projectIPAccessListModel.AWSSecurityGroup.ValueString()
			}

			entry, exists, err := isEntryInProjectAccessList(ctx, conn, projectID, accessListEntry)
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
		Timeout:    projectIPAccessListTimeout,
		Delay:      projectIPAccessListDelay,
		MinTimeout: projectIPAccessListMinTimeout,
	}

	// Wait, catching any errors
	accessList, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error while waiting for resource creation", err.Error())
		return
	}

	entry, ok := accessList.(*matlas.ProjectIPAccessList)
	if !ok {
		resp.Diagnostics.AddError("error", errorAccessListCreate)
		return
	}

	projectIPAccessListNewModel := newTFProjectIPAccessListModel(projectIPAccessListModel, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectIPAccessListNewModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFProjectIPAccessListModel(projectIPAccessListModel *tfProjectIPAccessListModel, projectIPAccessList *matlas.ProjectIPAccessList) *tfProjectIPAccessListModel {
	entry := projectIPAccessList.IPAddress
	if projectIPAccessList.CIDRBlock != "" {
		entry = projectIPAccessList.CIDRBlock
	} else if projectIPAccessList.AwsSecurityGroup != "" {
		entry = projectIPAccessList.AwsSecurityGroup
	}

	id := encodeStateID(map[string]string{
		"entry":      entry,
		"project_id": projectIPAccessList.GroupID,
	})

	return &tfProjectIPAccessListModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(projectIPAccessList.GroupID),
		CIDRBlock:        types.StringValue(projectIPAccessList.CIDRBlock),
		IPAddress:        types.StringValue(projectIPAccessList.IPAddress),
		AWSSecurityGroup: types.StringValue(projectIPAccessList.AwsSecurityGroup),
		Comment:          types.StringValue(projectIPAccessList.Comment),
		Timeouts:         projectIPAccessListModel.Timeouts,
	}
}

func newMongoDBProjectIPAccessList(projectIPAccessListModel *tfProjectIPAccessListModel) []*matlas.ProjectIPAccessList {
	return []*matlas.ProjectIPAccessList{
		{
			AwsSecurityGroup: projectIPAccessListModel.AWSSecurityGroup.ValueString(),
			CIDRBlock:        projectIPAccessListModel.CIDRBlock.ValueString(),
			IPAddress:        projectIPAccessListModel.IPAddress.ValueString(),
			Comment:          projectIPAccessListModel.Comment.ValueString(),
		},
	}
}

func (r *ProjectIPAccessListRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var projectIPAccessListModelState *tfProjectIPAccessListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &projectIPAccessListModelState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	decodedIDMap := decodeStateID(projectIPAccessListModelState.ID.ValueString())
	if len(decodedIDMap) != 2 {
		resp.Diagnostics.AddError("error during the reading operation", "the provided resource ID is not correct")
		return
	}

	timeout, diags := projectIPAccessListModelState.Timeouts.Read(ctx, projectIPAccessListTimeoutRead)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		accessList, httpResponse, err := conn.ProjectIPAccessList.Get(ctx, decodedIDMap["project_id"], decodedIDMap["entry"])
		if err != nil {
			// case 404
			// deleted in the backend case
			if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				resp.Diagnostics.AddError("resource not found", err.Error())
				return nil
			}

			if httpResponse != nil && httpResponse.StatusCode == http.StatusInternalServerError {
				return retry.RetryableError(err)
			}

			resp.Diagnostics.AddError("error getting project ip access list information", err.Error())
			return nil
		}

		projectIPAccessListNewModel := newTFProjectIPAccessListModel(projectIPAccessListModelState, accessList)
		resp.Diagnostics.Append(resp.State.Set(ctx, &projectIPAccessListNewModel)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error during the read operation", err.Error())
	}
}

func (r *ProjectIPAccessListRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var projectIPAccessListModelState *tfProjectIPAccessListModel

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

	conn := r.client.Atlas
	projectID := projectIPAccessListModelState.ProjectID.ValueString()

	timeout, diags := projectIPAccessListModelState.Timeouts.Delete(ctx, projectIPAccessListTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		httpResponse, err := conn.ProjectIPAccessList.Delete(ctx, projectID, entry)
		if err != nil {
			if httpResponse != nil && httpResponse.StatusCode == http.StatusInternalServerError {
				return retry.RetryableError(err)
			}

			if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
				return nil
			}

			resp.Diagnostics.AddError("error deleting the entry", fmt.Sprintf(errorAccessListDelete, err.Error()))
			return retry.NonRetryableError(fmt.Errorf(errorAccessListDelete, err))
		}

		entry, httpResponse, err := conn.ProjectIPAccessList.Get(ctx, projectID, entry)
		if err != nil {
			if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
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

func (r *ProjectIPAccessListRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "-", 2)

	if len(parts) != 2 {
		resp.Diagnostics.AddError("import format error", "to import a projectIP Access List, use the format {project_id}-{entry}")
	}

	projectID := parts[0]
	entry := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), encodeStateID(map[string]string{
		"entry":      entry,
		"project_id": projectID,
	}))...)
}

func isEntryInProjectAccessList(ctx context.Context, conn *matlas.Client, projectID, entry string) (*matlas.ProjectIPAccessList, bool, error) {
	var out matlas.ProjectIPAccessList
	err := retry.RetryContext(ctx, projectIPAccessListRetry, func() *retry.RetryError {
		accessList, httpResponse, err := conn.ProjectIPAccessList.Get(ctx, projectID, entry)
		if err != nil {
			switch {
			case httpResponse != nil && httpResponse.StatusCode == http.StatusInternalServerError:
				return retry.RetryableError(err)
			case httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound:
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

// Update is not supported
func (r *ProjectIPAccessListRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
