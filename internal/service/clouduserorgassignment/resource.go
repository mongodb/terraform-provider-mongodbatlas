package clouduserorgassignment

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

const resourceName = "cloud_user_org_assignment"

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithMoveState = &rs{}

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
	orgID := plan.OrgId.ValueString()
	orgUserRequest, diags := NewOrgUserReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResp, _, err := connV2.MongoDBCloudUsersApi.CreateOrganizationUser(ctx, orgID, orgUserRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error assigning user to OrgID(%s):", orgID), err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, apiResp, orgID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

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
	var httpResp *http.Response
	var err error

	if !state.UserId.IsNull() && state.UserId.ValueString() != "" {
		userID := state.UserId.ValueString()
		userResp, httpResp, err = connV2.MongoDBCloudUsersApi.GetOrganizationUser(ctx, orgID, userID).Execute()
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
	} else if !state.Username.IsNull() && state.Username.ValueString() != "" { // required for import
		username := state.Username.ValueString()
		params := &admin.ListOrganizationUsersApiParams{
			OrgId:    orgID,
			Username: &username,
		}
		usersResp, _, err := connV2.MongoDBCloudUsersApi.ListOrganizationUsersWithParams(ctx, params).Execute()
		if err == nil && usersResp != nil && usersResp.Results != nil {
			if len(*usersResp.Results) == 0 {
				resp.State.RemoveResource(ctx)
				return
			}
			userResp = &(*usersResp.Results)[0]
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error fetching user(%s) from OrgID(%s):", userResp.Username, orgID), err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, userResp, orgID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

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
	username := plan.Username.ValueString()

	updateReq, diags := NewAtlasUpdateReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResp, _, err := connV2.MongoDBCloudUsersApi.UpdateOrganizationUser(ctx, orgID, userID, updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error updating user(%s) in OrgID(%s):", username, orgID), err.Error())
		return
	}

	newCloudUserOrgAssignmentModel, diags := NewTFModel(ctx, apiResp, orgID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
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
	username := state.Username.ValueString()

	httpResp, err := connV2.MongoDBCloudUsersApi.RemoveOrganizationUser(ctx, orgID, userID).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error deleting user(%s) from OrgID(%s):", username, orgID), err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	ok, parts := conversion.ImportSplit(req.ID, 2)
	if !ok {
		resp.Diagnostics.AddError("invalid import ID format", "expected 'org_id/user_id' or 'org_id/username', got: "+importID)
		return
	}
	orgID, userID := parts[0], parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

	if emailRegex.MatchString(userID) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), userID)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userID)...)
	}
}
