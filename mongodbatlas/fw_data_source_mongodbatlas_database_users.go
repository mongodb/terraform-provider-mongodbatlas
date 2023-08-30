package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/description"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	databaseUsersDSName = "database_users"
)

type DatabaseUsersDS struct {
	client *MongoDBClient
}

func NewDatabaseUsersDS() datasource.DataSource {
	return &DatabaseUsersDS{}
}

var _ datasource.DataSource = &DatabaseUsersDS{}
var _ datasource.DataSourceWithConfigure = &DatabaseUsersDS{}

func (d *DatabaseUsersDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, databaseUsersDSName)
}

func (d *DatabaseUsersDS) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := ConfigureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	d.client = client
}

type tfDatabaseUsersDSModel struct {
	ID        types.String             `tfsdk:"id"`
	ProjectID types.String             `tfsdk:"project_id"`
	Results   []*tfDatabaseUserDSModel `tfsdk:"results"`
}

func (d *DatabaseUsersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.ID,
				Description:         description.ID,
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.ProjectID,
			},

			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.ID,
							Description:         description.ID,
						},
						"project_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.ProjectID,
						},
						"auth_database_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.ProjectID,
						},
						"username": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.Username,
						},
						"password": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: description.Password,
						},
						"x509_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.X509Type,
						},
						"ldap_auth_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.LDAPAuthYype,
						},
						"aws_iam_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.AWSIAMType,
						},
						"roles": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: description.Roles,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"collection_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.CollectionName,
									},
									"database_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.DatabaseName,
									},
									"role_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.RoleName,
									},
								},
							},
						},
						"labels": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: description.Labels,
							Description:         description.Labels,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.Key,
									},
									"value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.Value,
									},
								},
							},
						},
						"scopes": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: description.Scopes,
							Description:         description.Scopes,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.Name,
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: description.Type,
									},
								},
							},
						},
					},
				},
			},
		},
		MarkdownDescription: description.DatabaseUsersDS,
	}
}

func (d *DatabaseUsersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseUsersModel *tfDatabaseUsersDSModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseUsersModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := databaseUsersModel.ProjectID.ValueString()
	conn := d.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.List(ctx, projectID, nil)
	if err != nil {
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUsersMode(ctx, projectID, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFDatabaseUsersMode(ctx context.Context, projectID string, dbUsers []matlas.DatabaseUser) (*tfDatabaseUsersDSModel, diag.Diagnostics) {
	results := make([]*tfDatabaseUserDSModel, len(dbUsers))
	for i := range dbUsers {
		dbUserModel, d := newTFDatabaseDSUserModel(ctx, &dbUsers[i])
		if d.HasError() {
			return nil, d
		}
		results[i] = dbUserModel
	}

	return &tfDatabaseUsersDSModel{
		ProjectID: types.StringValue(projectID),
		Results:   results,
		ID:        types.StringValue(id.UniqueId()),
	}, nil
}
