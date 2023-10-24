package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

type DatabaseUserDS struct {
	DSCommon
}

func NewDatabaseUserDS() datasource.DataSource {
	return &DatabaseUserDS{
		DSCommon: DSCommon{
			dataSourceName: databaseUserResourceName,
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
	Roles            []tfRoleModel  `tfsdk:"roles"`
	Labels           []tfLabelModel `tfsdk:"labels"`
	Scopes           []tfScopeModel `tfsdk:"scopes"`
}

type tfRoleModel struct {
	RoleName       types.String `tfsdk:"role_name"`
	CollectionName types.String `tfsdk:"collection_name"`
	DatabaseName   types.String `tfsdk:"database_name"`
}

type tfLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type tfScopeModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

var _ datasource.DataSource = &DatabaseUserDS{}
var _ datasource.DataSourceWithConfigure = &DatabaseUserDS{}

func (d *DatabaseUserDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *DatabaseUserDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseDSUserModel *tfDatabaseUserDSModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseDSUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := databaseDSUserModel.Username.ValueString()
	projectID := databaseDSUserModel.ProjectID.ValueString()
	authDatabaseName := databaseDSUserModel.AuthDatabaseName.ValueString()

	conn := d.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Get(ctx, authDatabaseName, projectID, username)
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

func newTFDatabaseDSUserModel(ctx context.Context, dbUser *matlas.DatabaseUser) (*tfDatabaseUserDSModel, diag.Diagnostics) {
	id := fmt.Sprintf("%s-%s-%s", dbUser.GroupID, dbUser.Username, dbUser.DatabaseName)
	databaseUserModel := &tfDatabaseUserDSModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(dbUser.GroupID),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		Password:         types.StringValue(dbUser.Password),
		X509Type:         types.StringValue(dbUser.X509Type),
		OIDCAuthType:     types.StringValue(dbUser.OIDCAuthType),
		LDAPAuthType:     types.StringValue(dbUser.LDAPAuthType),
		AWSIAMType:       types.StringValue(dbUser.AWSIAMType),
		Roles:            newTFRolesModelLegacy(dbUser.Roles),
		Labels:           newTFLabelsModelLegacy(dbUser.Labels),
		Scopes:           newTFScopesModelLegacy(dbUser.Scopes),
	}

	return databaseUserModel, nil
}

func newTFScopesModelLegacy(scopes []matlas.Scope) []tfScopeModel {
	if len(scopes) == 0 {
		return nil
	}

	out := make([]tfScopeModel, len(scopes))
	for i, v := range scopes {
		out[i] = tfScopeModel{
			Name: types.StringValue(v.Name),
			Type: types.StringValue(v.Type),
		}
	}

	return out
}

func newTFLabelsModelLegacy(labels []matlas.Label) []tfLabelModel {
	if len(labels) == 0 {
		return nil
	}

	out := make([]tfLabelModel, len(labels))
	for i, v := range labels {
		out[i] = tfLabelModel{
			Key:   types.StringValue(v.Key),
			Value: types.StringValue(v.Value),
		}
	}

	return out
}

func newTFRolesModelLegacy(roles []matlas.Role) []tfRoleModel {
	if len(roles) == 0 {
		return nil
	}

	out := make([]tfRoleModel, len(roles))
	for i, v := range roles {
		out[i] = tfRoleModel{
			RoleName:     types.StringValue(v.RoleName),
			DatabaseName: types.StringValue(v.DatabaseName),
		}

		if v.CollectionName != "" {
			out[i].CollectionName = types.StringValue(v.CollectionName)
		}
	}

	return out
}
