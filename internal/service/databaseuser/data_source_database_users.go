package databaseuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	databaseUsersDSName = "database_users"
)

type DatabaseUsersDS struct {
	config.DSCommon
}

func PluralDataSource() datasource.DataSource {
	return &DatabaseUsersDS{
		DSCommon: config.DSCommon{
			DataSourceName: databaseUsersDSName,
		},
	}
}

var _ datasource.DataSource = &DatabaseUsersDS{}
var _ datasource.DataSourceWithConfigure = &DatabaseUsersDS{}

type TfDatabaseUsersDSModel struct {
	ID        types.String             `tfsdk:"id"`
	ProjectID types.String             `tfsdk:"project_id"`
	Results   []*TfDatabaseUserDSModel `tfsdk:"results"`
}

func (d *DatabaseUsersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},

			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
						"auth_database_name": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
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
						"roles": schema.SetNestedAttribute{
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
						"labels": schema.SetNestedAttribute{
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
						"scopes": schema.SetNestedAttribute{
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
				},
			},
		},
	}
}

func (d *DatabaseUsersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseUsersModel *TfDatabaseUsersDSModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseUsersModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := databaseUsersModel.ProjectID.ValueString()
	connV2 := d.Client.AtlasV2
	dbUser, _, err := connV2.DatabaseUsersApi.ListDatabaseUsers(ctx, projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := NewTFDatabaseUsersModel(ctx, projectID, dbUser.GetResults())
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
