package clouduserorgassignment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &cloudUserOrgAssignmentDS{}
var _ datasource.DataSourceWithConfigure = &cloudUserOrgAssignmentDS{}

func DataSource() datasource.DataSource {
	return &cloudUserOrgAssignmentDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type cloudUserOrgAssignmentDS struct {
	config.DSCommon
}

func (d *cloudUserOrgAssignmentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *cloudUserOrgAssignmentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cloudUserOrgAssignmentConfig TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cloudUserOrgAssignmentConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	orgID := cloudUserOrgAssignmentConfig.OrgId.ValueString()
	username := cloudUserOrgAssignmentConfig.Username.ValueString()
	userID := cloudUserOrgAssignmentConfig.UserId.ValueString()

	if username == "" && userID == "" {
		resp.Diagnostics.AddError("invalid configuration", "either username or user_id must be provided")
		return
	}

	var orgUser *admin.OrgUserResponse
	var err error

	if userID != "" {
		orgUser, _, err = connV2.MongoDBCloudUsersApi.GetOrganizationUser(ctx, orgID, userID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error retrieving resource by user_id: %s", userID), err.Error())
			return
		}
	} else {
		params := &admin.ListOrganizationUsersApiParams{
			OrgId:    orgID,
			Username: &username,
		}
		usersResp, _, err := connV2.MongoDBCloudUsersApi.ListOrganizationUsersWithParams(ctx, params).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("error retrieving resource by username: %s", username), err.Error())
			return
		}

		if usersResp == nil || usersResp.Results == nil || len(*usersResp.Results) == 0 {
			resp.Diagnostics.AddError("resource not found", "no user found with the specified username")
			return
		}

		orgUser = &(*usersResp.Results)[0]
	}

	tfModel, diags := NewTFModel(ctx, orgUser, cloudUserOrgAssignmentConfig.OrgId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
}
