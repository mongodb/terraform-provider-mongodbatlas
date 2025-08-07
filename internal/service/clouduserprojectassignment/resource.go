package clouduserprojectassignment

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName           = "cloud_user_project_assignment"
	errorReadingByUserID   = "Error getting project users by user_id"
	errorReadingByUsername = "Error getting project users by username"
	invalidImportID        = "Invalid import ID format"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type rs struct {
	config.RSCommon
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema()
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectId.ValueString()
	projectUserRequest, diags := NewProjectUserReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResp, _, err := connV2.MongoDBCloudUsersApi.AddProjectUser(ctx, projectID, projectUserRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error assigning user to ProjectID(%s):", projectID), err.Error())
		return
	}

	newCloudUserProjectAssignmentModel, diags := NewTFModel(ctx, projectID, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserProjectAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	var userResp *admin.GroupUserResponse
	var httpResp *http.Response
	var err error

	if !state.UserId.IsNull() && state.UserId.ValueString() != "" {
		userID := state.UserId.ValueString()
		userResp, httpResp, err = connV2.MongoDBCloudUsersApi.GetProjectUser(ctx, projectID, userID).Execute()
		if err != nil {
			if validate.StatusNotFound(httpResp) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(errorReadingByUserID, err.Error())
			return
		}
	} else if !state.Username.IsNull() && state.Username.ValueString() != "" { // required for import
		username := state.Username.ValueString()
		params := &admin.ListProjectUsersApiParams{
			GroupId:  projectID,
			Username: &username,
		}
		usersResp, _, err := connV2.MongoDBCloudUsersApi.ListProjectUsersWithParams(ctx, params).Execute()
		if err != nil {
			if validate.StatusNotFound(httpResp) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(errorReadingByUsername, err.Error())
			return
		}
		if usersResp == nil || len(usersResp.GetResults()) == 0 {
			resp.State.RemoveResource(ctx)
			return
		}
		userResp = &usersResp.GetResults()[0]
	}

	if userResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newCloudUserProjectAssignmentModel, diags := NewTFModel(ctx, projectID, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserProjectAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel
	var state TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := plan.ProjectId.ValueString()
	userID := plan.UserId.ValueString()
	username := plan.Username.ValueString()

	addRequests, removeRequests, diags := NewAtlasUpdateReq(ctx, &plan, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, addReq := range addRequests {
		_, _, err := connV2.MongoDBCloudUsersApi.AddProjectRole(ctx, projectID, userID, addReq).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error adding role %s to user(%s) in ProjectID(%s):", addReq.GroupRole, username, projectID),
				err.Error(),
			)
			return
		}
	}

	for _, removeReq := range removeRequests {
		_, _, err := connV2.MongoDBCloudUsersApi.RemoveProjectRole(ctx, projectID, userID, removeReq).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error removing role %s from user(%s) in ProjectID(%s):", removeReq.GroupRole, username, projectID),
				err.Error(),
			)
			return
		}
	}

	var userResp *admin.GroupUserResponse
	var err error
	if !state.UserId.IsNull() && state.UserId.ValueString() != "" {
		userID := state.UserId.ValueString()
		userResp, _, err = connV2.MongoDBCloudUsersApi.GetProjectUser(ctx, projectID, userID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error fetching user(%s) from ProjectID(%s):", username, projectID), err.Error())
			return
		}
	}

	newCloudUserProjectAssignmentModel, diags := NewTFModel(ctx, projectID, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserProjectAssignmentModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := state.ProjectId.ValueString()
	userID := state.UserId.ValueString()
	username := state.Username.ValueString()

	httpResp, err := connV2.MongoDBCloudUsersApi.RemoveProjectUser(ctx, projectID, userID).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error deleting user(%s) from ProjectID(%s):", username, projectID), err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	ok, parts := conversion.ImportSplit(req.ID, 2)
	if !ok {
		resp.Diagnostics.AddError(invalidImportID, "expected 'project_id/user_id' or 'project_id/username', got: "+importID)
		return
	}
	projectID, user := parts[0], parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

	if emailRegex.MatchString(user) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), user)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), user)...)
	}
}
