package atlasuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	AtlasUserDataSourceName = "atlas_user"
	errorUserRead           = "error getting atlas users(%s): %s"
)

var _ datasource.DataSource = &atlasUserDS{}
var _ datasource.DataSourceWithConfigure = &atlasUserDS{}

type tfAtlasUserDSModel struct {
	ID           types.String           `tfsdk:"id"`
	UserID       types.String           `tfsdk:"user_id"`
	Username     types.String           `tfsdk:"username"`
	Country      types.String           `tfsdk:"country"`
	CreatedAt    types.String           `tfsdk:"created_at"`
	EmailAddress types.String           `tfsdk:"email_address"`
	FirstName    types.String           `tfsdk:"first_name"`
	LastAuth     types.String           `tfsdk:"last_auth"`
	LastName     types.String           `tfsdk:"last_name"`
	MobileNumber types.String           `tfsdk:"mobile_number"`
	TeamIDs      []string               `tfsdk:"team_ids"`
	Links        []tfLinkModel          `tfsdk:"links"`
	Roles        []tfAtlasUserRoleModel `tfsdk:"roles"`
}

type tfLinkModel struct {
	Href types.String `tfsdk:"href"`
	Rel  types.String `tfsdk:"rel"`
}

type tfAtlasUserRoleModel struct {
	GroupID  types.String `tfsdk:"group_id"`
	OrgID    types.String `tfsdk:"org_id"`
	RoleName types.String `tfsdk:"role_name"`
}

func DataSource() datasource.DataSource {
	return &atlasUserDS{
		DSCommon: config.DSCommon{
			DataSourceName: AtlasUserDataSourceName,
		},
	}
}

type atlasUserDS struct {
	config.DSCommon
}

func (d *atlasUserDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
				DeprecationMessage: "Please use user_id id attribute instead",
				Computed:           true,
			},
			"user_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
				},
			},
			"username": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("user_id")),
				},
			},
			"country": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"email_address": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Computed: true,
			},
			"last_auth": schema.StringAttribute{
				Computed: true,
			},
			"last_name": schema.StringAttribute{
				Computed: true,
			},
			"mobile_number": schema.StringAttribute{
				Computed: true,
			},
			"team_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"links": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Computed: true,
						},
						"rel": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"roles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Computed: true,
						},
						"org_id": schema.StringAttribute{
							Computed: true,
						},
						"role_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *atlasUserDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	connV2 := d.Client.AtlasV2

	var atlasUserConfig tfAtlasUserDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &atlasUserConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if atlasUserConfig.UserID.IsNull() && atlasUserConfig.Username.IsNull() {
		resp.Diagnostics.AddError(errorMissingAttributesSummary, "either user_id or username must be configured")
		return
	}

	var (
		err  error
		user *admin.CloudAppUser
	)
	if !atlasUserConfig.UserID.IsNull() {
		userID := atlasUserConfig.UserID.ValueString()
		user, _, err = connV2.MongoDBCloudUsersApi.GetUser(ctx, userID).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting User from Atlas", fmt.Sprintf(errorUserRead, userID, err.Error()))
			return
		}
	} else {
		username := atlasUserConfig.Username.ValueString()
		user, _, err = connV2.MongoDBCloudUsersApi.GetUserByUsername(ctx, username).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting User from Atlas", fmt.Sprintf(errorUserRead, username, err.Error()))
			return
		}
	}

	userResultState := newTFAtlasUserDSModel(user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &userResultState)...)
}

func newTFAtlasUserDSModel(user *admin.CloudAppUser) tfAtlasUserDSModel {
	return tfAtlasUserDSModel{
		ID:           types.StringPointerValue(user.Id),
		UserID:       types.StringPointerValue(user.Id),
		Username:     types.StringValue(user.Username),
		Country:      types.StringValue(user.Country),
		CreatedAt:    types.StringPointerValue(conversion.TimePtrToStringPtr(user.CreatedAt)),
		EmailAddress: types.StringValue(user.EmailAddress),
		FirstName:    types.StringValue(user.FirstName),
		LastAuth:     types.StringPointerValue(conversion.TimePtrToStringPtr(user.LastAuth)),
		LastName:     types.StringValue(user.LastName),
		MobileNumber: types.StringValue(user.MobileNumber),
		TeamIDs:      user.GetTeamIds(),
		Links:        newTFLinksList(user.GetLinks()),
		Roles:        newTFRolesList(user.GetRoles()),
	}
}

func newTFLinksList(links []admin.Link) []tfLinkModel {
	if links == nil {
		return nil
	}
	resLinks := make([]tfLinkModel, len(links))
	for i, value := range links {
		resLink := tfLinkModel{
			Href: types.StringPointerValue(value.Href),
			Rel:  types.StringPointerValue(value.Rel),
		}
		resLinks[i] = resLink
	}
	return resLinks
}

func newTFRolesList(roles []admin.CloudAccessRoleAssignment) []tfAtlasUserRoleModel {
	if roles == nil {
		return nil
	}
	resRoles := make([]tfAtlasUserRoleModel, len(roles))
	for i, value := range roles {
		resRole := tfAtlasUserRoleModel{
			GroupID:  types.StringPointerValue(value.GroupId),
			OrgID:    types.StringPointerValue(value.OrgId),
			RoleName: types.StringPointerValue(value.RoleName),
		}
		resRoles[i] = resRole
	}
	return resRoles
}
