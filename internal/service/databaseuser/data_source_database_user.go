package databaseuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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

type TfDatabaseUserDSModel struct {
	ID               types.String   `tfsdk:"id"`
	ProjectID        types.String   `tfsdk:"project_id"`
	AuthDatabaseName types.String   `tfsdk:"auth_database_name"`
	Username         types.String   `tfsdk:"username"`
	Description      types.String   `tfsdk:"description"`
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
			"description": schema.StringAttribute{
				Computed: true,
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
	var databaseDSUserModel *TfDatabaseUserDSModel
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

	dbUserModel, diagnostic := NewTFDatabaseDSUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
