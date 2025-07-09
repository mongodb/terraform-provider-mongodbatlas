package clouduserorgassignment

import (
	"context"
	"regexp"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const resourceName = "cloud_user_org_assignment"

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

// var _ resource.ResourceWithMoveState = &rs{} TODO: follow up ticket to implement move state

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
	resp.Schema = resourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	orgID := plan.OrgId.ValueString()
	orgUserRequest, diags := NewOrgUserReq(ctx, &plan)

	apiResp, _, err := connV2.MongoDBCloudUsersApi.CreateOrganizationUser(ctx, orgID, orgUserRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	getUserResp, _, err := connV2.MongoDBCloudUsersApi.GetOrganizationUser(ctx, orgID, apiResp.Id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, getUserResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newCloudUserOrgAssignmentModel.OrgId = plan.OrgId

	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserOrgAssignmentModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	orgID := state.OrgId.ValueString()
	var userResp *admin.OrgUserResponse
	var err error

	if !state.UserId.IsNull() && state.UserId.ValueString() != "" {
		userID := state.UserId.ValueString()
		userResp, _, err = connV2.MongoDBCloudUsersApi.GetOrganizationUser(ctx, orgID, userID).Execute()
	} else if !state.Username.IsNull() && state.Username.ValueString() != "" {
		username := state.Username.ValueString()
		params := &admin.ListOrganizationUsersApiParams{
			OrgId:    orgID,
			Username: &username,
		}
		usersResp, _, err := connV2.MongoDBCloudUsersApi.ListOrganizationUsersWithParams(ctx, params).Execute()
		if err == nil && usersResp != nil && usersResp.Results != nil && len(*usersResp.Results) > 0 {
			userResp = &(*usersResp.Results)[0]
		} else if err == nil {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	if err != nil {
		// Note: validate.StatusNotFound expects *http.Response, but we don't have it here. This is a limitation for username path.
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, userResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newCloudUserOrgAssignmentModel.OrgId = state.OrgId

	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserOrgAssignmentModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	orgID := plan.OrgId.ValueString()
	userID := plan.UserId.ValueString()
	if userID == "" {
		resp.Diagnostics.AddError("missing user_id", "user_id (id) must be set in state for update operation")
		return
	}

	updateReq, diags := NewAtlasUpdateReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, _, err := connV2.MongoDBCloudUsersApi.UpdateOrganizationUser(ctx, orgID, userID, updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	getUserResp, _, err := connV2.MongoDBCloudUsersApi.GetOrganizationUser(ctx, orgID, userID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching updated resource", err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, getUserResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newCloudUserOrgAssignmentModel.OrgId = plan.OrgId
	resp.Diagnostics.Append(resp.State.Set(ctx, newCloudUserOrgAssignmentModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	orgID := state.OrgId.ValueString()
	userID := state.UserId.ValueString()
	if userID == "" {
		resp.Diagnostics.AddError("missing user_id", "user_id (id) must be set in state for delete operation")
		return
	}

	httpResp, err := connV2.MongoDBCloudUsersApi.RemoveOrganizationUser(ctx, orgID, userID).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	var orgID, userID string
	parts := strings.Split(importID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("invalid import ID format", "expected 'org_id/user_id' or 'org_id/username', got: "+importID)
		return
	}
	orgID, userID = parts[0], parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

	if emailRegex.MatchString(userID) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), userID)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userID)...)
	}
}
