package databaseuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

type databaseUserDS struct {
	config.DSCommon
}

func DataSource() datasource.DataSource {
	return &databaseUserDS{
		DSCommon: config.DSCommon{
			DataSourceName: databaseUserResourceName,
		},
	}
}

type tfDatabaseUserDSModel struct {
	ID               types.String   `tfsdk:"id"`
	ProjectID        types.String   `tfsdk:"project_id"`
	AuthDatabaseName types.String   `tfsdk:"auth_database_name"`
	Username         types.String   `tfsdk:"username"`
	Password         types.String   `tfsdk:"password"`
	X509Type         types.String   `tfsdk:"x509_type"`
	OIDCAuthType     types.String   `tfsdk:"oidc_auth_type"`
	LDAPAuthType     types.String   `tfsdk:"ldap_auth_type"`
	AWSIAMType       types.String   `tfsdk:"aws_iam_type"`
	Roles            []TfRoleModel  `tfsdk:"roles"`
	Labels           []TfLabelModel `tfsdk:"labels"`
	Scopes           []TfScopeModel `tfsdk:"scopes"`
}

var _ datasource.DataSource = &databaseUserDS{}
var _ datasource.DataSourceWithConfigure = &databaseUserDS{}

func (d *databaseUserDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"auth_database_name": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"x509_type": schema.StringAttribute{
				Computed: true,
			},
			"oidc_auth_type": schema.StringAttribute{
				Computed: true,
			},
			"ldap_auth_type": schema.StringAttribute{
				Computed: true,
			},
			"aws_iam_type": schema.StringAttribute{
				Computed: true,
			},
			"roles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"collection_name": schema.StringAttribute{
							Computed: true,
						},
						"database_name": schema.StringAttribute{
							Computed: true,
						},
						"role_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"labels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"scopes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *databaseUserDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseDSUserModel *tfDatabaseUserDSModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseDSUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := databaseDSUserModel.Username.ValueString()
	projectID := databaseDSUserModel.ProjectID.ValueString()
	authDatabaseName := databaseDSUserModel.AuthDatabaseName.ValueString()

	connV2 := d.Client.AtlasV2
	dbUser, _, err := connV2.DatabaseUsersApi.GetDatabaseUser(ctx, projectID, authDatabaseName, username).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseDSUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFDatabaseDSUserModel(ctx context.Context, dbUser *admin.CloudDatabaseUser) (*tfDatabaseUserDSModel, diag.Diagnostics) {
	id := fmt.Sprintf("%s-%s-%s", dbUser.GroupId, dbUser.Username, dbUser.DatabaseName)
	databaseUserModel := &tfDatabaseUserDSModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(dbUser.GroupId),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		Password:         types.StringValue(dbUser.GetPassword()),
		X509Type:         types.StringValue(dbUser.GetX509Type()),
		OIDCAuthType:     types.StringValue(dbUser.GetOidcAuthType()),
		LDAPAuthType:     types.StringValue(dbUser.GetLdapAuthType()),
		AWSIAMType:       types.StringValue(dbUser.GetAwsIAMType()),
		Roles:            NewTFRolesModel(dbUser.Roles),
		Labels:           NewTFLabelsModel(dbUser.Labels),
		Scopes:           NewTFScopesModel(dbUser.Scopes),
	}

	return databaseUserModel, nil
}

func NewTFLabelsModel(labels []admin.ComponentLabel) []TfLabelModel {
	if len(labels) == 0 {
		return nil
	}

	out := make([]TfLabelModel, len(labels))
	for i, v := range labels {
		out[i] = TfLabelModel{
			Key:   types.StringValue(v.GetKey()),
			Value: types.StringValue(v.GetValue()),
		}
	}

	return out
}

func NewTFRolesModel(roles []admin.DatabaseUserRole) []TfRoleModel {
	if len(roles) == 0 {
		return nil
	}

	out := make([]TfRoleModel, len(roles))
	for i, v := range roles {
		out[i] = TfRoleModel{
			RoleName:     types.StringValue(v.RoleName),
			DatabaseName: types.StringValue(v.DatabaseName),
		}

		if v.GetCollectionName() != "" {
			out[i].CollectionName = types.StringValue(v.GetCollectionName())
		}
	}

	return out
}

func NewMongoDBAtlasRoles(roles []*TfRoleModel) []admin.DatabaseUserRole {
	if len(roles) == 0 {
		return []admin.DatabaseUserRole{}
	}

	out := make([]admin.DatabaseUserRole, len(roles))
	for i, v := range roles {
		out[i] = admin.DatabaseUserRole{
			RoleName:       v.RoleName.ValueString(),
			DatabaseName:   v.DatabaseName.ValueString(),
			CollectionName: v.CollectionName.ValueStringPointer(),
		}
	}

	return out
}
