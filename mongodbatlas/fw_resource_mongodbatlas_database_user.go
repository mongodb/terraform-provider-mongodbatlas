package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/description"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	databaseUserResourceName = "database_user"
)

var _ resource.Resource = &DatabaseUserRS{}
var _ resource.ResourceWithImportState = &DatabaseUserRS{}

type tfDatabaseUserModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	AuthDatabaseName types.String `tfsdk:"auth_database_name"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	X509Type         types.String `tfsdk:"x509_type"`
	LDAPAuthType     types.String `tfsdk:"ldap_auth_type"`
	AWSIAMType       types.String `tfsdk:"aws_iam_type"`
	Roles            types.Set    `tfsdk:"roles"`
	Labels           types.Set    `tfsdk:"labels"`
	Scopes           types.Set    `tfsdk:"scopes"`
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

type DatabaseUserRS struct {
	client *MongoDBClient
}

var RoleObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"roleName":       types.StringType,
	"collectionName": types.StringType,
	"databaseName":   types.StringType,
}}

var LabelObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

var ScopeObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
}}

func (r *DatabaseUserRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: description.ID,
				Description:         description.ID,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: description.ProjectID,
				Description:         description.ProjectID,
			},
			"auth_database_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: description.ProjectID,
				Description:         description.ProjectID,
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: description.Username,
				Description:         description.Username,
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: description.Password,
				Description:         description.Password,
			},
			"x509_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "MANAGED", "CUSTOMER"),
				},
				MarkdownDescription: description.X509Type,
				Description:         description.X509Type,
			},
			"ldap_auth_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "USER", "GROUP"),
				},
				MarkdownDescription: description.LDAPAuthYype,
				Description:         description.LDAPAuthYype,
			},
			"aws_iam_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "USER", "ROLE"),
				},
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
							Optional:            true,
							Computed:            true,
							MarkdownDescription: description.CollectionName,
							Description:         description.CollectionName,
						},
						"database_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: description.DatabaseName,
							Description:         description.DatabaseName,
						},
						"role_name": schema.StringAttribute{
							Optional:            true,
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
							Optional:            true,
							Computed:            true,
							MarkdownDescription: description.Key,
							Description:         description.Key,
						},
						"value": schema.StringAttribute{
							Optional:            true,
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
							Optional:            true,
							Computed:            true,
							MarkdownDescription: description.Name,
							Description:         description.Name,
						},
						"type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: description.Type,
							Description:         description.Type,
						},
					},
				},
			},
		},
		MarkdownDescription: description.DatabaseUserResource,
		Description:         description.DatabaseUserResource,
	}
}

func NewDatabaseUserRS() resource.Resource {
	return &DatabaseUserRS{}
}

func (r *DatabaseUserRS) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, databaseUserResourceName)
}

func (r *DatabaseUserRS) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := ConfigureClientInResource(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.client = client
}

func (r *DatabaseUserRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var databaseUserModel *tfDatabaseUserModel

	diags := req.Plan.Get(ctx, &databaseUserModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, d := newMongoDBDatabaseUser(ctx, databaseUserModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Create(ctx, databaseUserModel.ProjectID.ValueString(), dbUserReq)
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseUserRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var databaseUserModel *tfDatabaseUserModel
	var err error
	resp.Diagnostics.Append(req.State.Get(ctx, &databaseUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := databaseUserModel.Username.ValueString()
	projectID := databaseUserModel.ProjectID.ValueString()
	authDatabaseName := databaseUserModel.AuthDatabaseName.ValueString()

	if databaseUserModel.ID.ValueString() != "" {
		projectID, username, authDatabaseName, err = splitDatabaseUserImportID(databaseUserModel.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error splitting database User info from ID", err.Error())
			return
		}
	}

	conn := r.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Get(ctx, authDatabaseName, projectID, username)
	if err != nil {
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseUserRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var databaseUserModel *tfDatabaseUserModel

	diags := req.Plan.Get(ctx, &databaseUserModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, d := newMongoDBDatabaseUser(ctx, databaseUserModel)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Update(ctx, databaseUserModel.ProjectID.ValueString(), databaseUserModel.Username.ValueString(), dbUserReq)
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseUserRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var dbUserModel *tfDatabaseUserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := deleteProject(ctx, r.client.Atlas, dbUserModel.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorProjectDelete, dbUserModel.ProjectID.ValueString(), err.Error()))
		return
	}
}

func (r *DatabaseUserRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func newMongoDBDatabaseUser(ctx context.Context, dbUserModel *tfDatabaseUserModel) (*matlas.DatabaseUser, diag.Diagnostics) {
	var rolesModel []*tfRoleModel
	var labelsModel []*tfLabelModel
	var scopesModel []*tfScopeModel

	diags := dbUserModel.Roles.ElementsAs(ctx, &rolesModel, false)
	if diags.HasError() {
		return nil, diags
	}

	diags = dbUserModel.Labels.ElementsAs(ctx, &labelsModel, false)
	if diags.HasError() {
		return nil, diags
	}

	diags = dbUserModel.Scopes.ElementsAs(ctx, &scopesModel, false)
	if diags.HasError() {
		return nil, diags
	}

	return &matlas.DatabaseUser{
		GroupID:      dbUserModel.ProjectID.ValueString(),
		Username:     dbUserModel.Username.ValueString(),
		Password:     dbUserModel.Password.ValueString(),
		X509Type:     dbUserModel.X509Type.ValueString(),
		AWSIAMType:   dbUserModel.AWSIAMType.ValueString(),
		LDAPAuthType: dbUserModel.LDAPAuthType.ValueString(),
		DatabaseName: dbUserModel.AuthDatabaseName.ValueString(),
		Roles:        newMongoDBAtlasRoles(rolesModel),
		Labels:       newMongoDBAtlasLabels(labelsModel),
		Scopes:       newMongoDBAtlasScopes(scopesModel),
	}, nil
}

func newTFDatabaseUserModel(ctx context.Context, dbUser *matlas.DatabaseUser) (*tfDatabaseUserModel, diag.Diagnostics) {
	rolesSet, diagnostic := types.SetValueFrom(ctx, RoleObjectType, newTFRolesModel(dbUser.Roles))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	labelsSet, diagnostic := types.SetValueFrom(ctx, LabelObjectType, newTFLabelsModel(dbUser.Labels))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	scopesSet, diagnostic := types.SetValueFrom(ctx, ScopeObjectType, newTFScopesModel(dbUser.Scopes))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	id := fmt.Sprintf("%s-%s-%s", dbUser.GroupID, dbUser.Username, dbUser.DatabaseName)
	databaseUserModel := &tfDatabaseUserModel{
		ID:               types.StringValue(id),
		ProjectID:        types.StringValue(dbUser.GroupID),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		Password:         types.StringValue(dbUser.Password),
		X509Type:         types.StringValue(dbUser.X509Type),
		LDAPAuthType:     types.StringValue(dbUser.LDAPAuthType),
		AWSIAMType:       types.StringValue(dbUser.AWSIAMType),
		Roles:            rolesSet,
		Labels:           labelsSet,
		Scopes:           scopesSet,
	}

	return databaseUserModel, nil
}

func newTFScopesModel(scopes []matlas.Scope) []tfScopeModel {
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

func newTFLabelsModel(labels []matlas.Label) []tfLabelModel {
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

func newTFRolesModel(roles []matlas.Role) []tfRoleModel {
	if len(roles) == 0 {
		return nil
	}

	out := make([]tfRoleModel, len(roles))
	for i, v := range roles {
		out[i] = tfRoleModel{
			RoleName:       types.StringValue(v.RoleName),
			DatabaseName:   types.StringValue(v.DatabaseName),
			CollectionName: types.StringValue(v.CollectionName),
		}
	}

	return out
}

func newMongoDBAtlasRoles(roles []*tfRoleModel) []matlas.Role {
	if len(roles) == 0 {
		return nil
	}

	out := make([]matlas.Role, len(roles))
	for i, v := range roles {
		out[i] = matlas.Role{
			RoleName:       v.RoleName.ValueString(),
			DatabaseName:   v.DatabaseName.ValueString(),
			CollectionName: v.CollectionName.ValueString(),
		}
	}

	return out
}

func newMongoDBAtlasLabels(labels []*tfLabelModel) []matlas.Label {
	if len(labels) == 0 {
		return nil
	}

	out := make([]matlas.Label, len(labels))
	for i, v := range labels {
		out[i] = matlas.Label{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}

	return out
}

func newMongoDBAtlasScopes(scopes []*tfScopeModel) []matlas.Scope {
	if len(scopes) == 0 {
		return nil
	}

	out := make([]matlas.Scope, len(scopes))
	for i, v := range scopes {
		out[i] = matlas.Scope{
			Name: v.Name.ValueString(),
			Type: v.Type.ValueString(),
		}
	}

	return out
}

func splitDatabaseUserImportID(id string) (projectID, username, authDatabaseName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([$a-z]{1,15})$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("import format error: to import a Database User, use the format {project_id}-{username}-{auth_database_name}")
		return
	}

	projectID = parts[1]
	username = parts[2]
	authDatabaseName = parts[3]

	return
}
