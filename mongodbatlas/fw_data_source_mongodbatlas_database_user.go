package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/description"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

type DatabaseUserDS struct {
	client *MongoDBClient
}

func NewDatabaseUserDS() datasource.DataSource {
	return &DatabaseUserDS{}
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

var _ datasource.DataSource = &DatabaseUserDS{}
var _ datasource.DataSourceWithConfigure = &DatabaseUserDS{}

func (d *DatabaseUserDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, databaseUserResourceName)
}

func (d *DatabaseUserDS) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := ConfigureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
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
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.ProjectID,
			},
			"auth_database_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.ProjectID,
			},
			"username": schema.StringAttribute{
				Required:            true,
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
			"oidc_auth_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.OIDC,
			},
			"ldap_auth_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.LDAPAuthYype,
			},
			"aws_iam_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: description.AWSIAMType,
			},
			"roles": schema.ListNestedAttribute{
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
			"labels": schema.ListNestedAttribute{
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
			"scopes": schema.ListNestedAttribute{
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
		MarkdownDescription: description.DatabaseUserDS,
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
		LDAPAuthType:     types.StringValue(dbUser.LDAPAuthType),
		AWSIAMType:       types.StringValue(dbUser.AWSIAMType),
		Roles:            newTFRolesModel(dbUser.Roles),
		Labels:           newTFLabelsModel(dbUser.Labels),
		Scopes:           newTFScopesModel(dbUser.Scopes),
	}

	return databaseUserModel, nil
}
