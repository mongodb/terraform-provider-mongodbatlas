package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/description"
)

type DatabaseUserDS struct {
	client *MongoDBClient
}

func NewDatabaseUserDS() datasource.DataSource {
	return &DatabaseUserDS{}
}

var _ datasource.DataSource = &DatabaseUserDS{}
var _ datasource.DataSourceWithConfigure = &DatabaseUserDS{}

func (d *DatabaseUserDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, databaseUserResourceName)
}

func (d *DatabaseUserDS) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*MongoDBClient)

	if !ok {
		resp.Diagnostics.AddError(errorConfigureSummary, fmt.Sprintf(errorConfigure, req.ProviderData))
		return
	}
	d.client = client
}

func (d *DatabaseUserDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Description:         description.ProjectID,
			},
			"auth_database_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.ProjectID,
				Description:         description.ProjectID,
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.Username,
				Description:         description.Username,
			},
			"password": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: description.Password,
				Description:         description.Password,
			},
			"x509_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.X509Type,
				Description:         description.X509Type,
			},
			"ldap_auth_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.LDAPAuthYype,
				Description:         description.LDAPAuthYype,
			},
			"aws_iam_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.AWSIAMType,
				Description:         description.AWSIAMType,
			},
		},
		Blocks: map[string]schema.Block{
			"roles": schema.SetNestedBlock{
				MarkdownDescription: description.Roles,
				Description:         description.Roles,
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"collection_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.CollectionName,
							Description:         description.CollectionName,
						},
						"database_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.DatabaseName,
							Description:         description.DatabaseName,
						},
						"role_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.RoleName,
							Description:         description.RoleName,
						},
					},
				},
			},
			"labels": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.Key,
							Description:         description.Key,
						},
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.Value,
							Description:         description.Value,
						},
					},
				},
			},
			"scopes": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.Name,
							Description:         description.Name,
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: description.Type,
							Description:         description.Type,
						},
					},
				},
			},
		},
		MarkdownDescription: description.DatabaseUserDS,
		Description:         description.DatabaseUserDS,
	}
}

func (d *DatabaseUserDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var databaseUserModel *tfDatabaseUserModel
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &databaseUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := databaseUserModel.Username.ValueString()
	projectID := databaseUserModel.ProjectID.ValueString()
	authDatabaseName := databaseUserModel.AuthDatabaseName.ValueString()

	conn := d.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Get(ctx, authDatabaseName, projectID, username)
	if err != nil {
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diag := newTFDatabaseUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
