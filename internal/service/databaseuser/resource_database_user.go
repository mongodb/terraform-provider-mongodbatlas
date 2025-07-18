package databaseuser

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	databaseUserResourceName = "database_user"
	ErrorImportFormat        = "import format error: to import a Database User, use the format {project_id}-{username}-{auth_database_name} OR {project_id}/{username}/{auth_database_name}"
)

var _ resource.ResourceWithConfigure = &databaseUserRS{}
var _ resource.ResourceWithImportState = &databaseUserRS{}

type databaseUserRS struct {
	config.RSCommon
}

func Resource() resource.Resource {
	return &databaseUserRS{
		RSCommon: config.RSCommon{
			ResourceName: databaseUserResourceName,
		},
	}
}

type TfDatabaseUserModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	AuthDatabaseName types.String `tfsdk:"auth_database_name"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	X509Type         types.String `tfsdk:"x509_type"`
	OIDCAuthType     types.String `tfsdk:"oidc_auth_type"`
	LDAPAuthType     types.String `tfsdk:"ldap_auth_type"`
	AWSIAMType       types.String `tfsdk:"aws_iam_type"`
	Description      types.String `tfsdk:"description"`
	Roles            types.Set    `tfsdk:"roles"`
	Labels           types.Set    `tfsdk:"labels"`
	Scopes           types.Set    `tfsdk:"scopes"`
}

type TfRoleModel struct {
	RoleName       types.String `tfsdk:"role_name"`
	CollectionName types.String `tfsdk:"collection_name"`
	DatabaseName   types.String `tfsdk:"database_name"`
}

type TfLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type TfScopeModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

var RoleObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"role_name":       types.StringType,
	"collection_name": types.StringType,
	"database_name":   types.StringType,
}}

var LabelObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

var ScopeObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
}}

func (r *databaseUserRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth_database_name": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRelative().AtParent().AtName("x509_type"),
						path.MatchRelative().AtParent().AtName("ldap_auth_type"),
						path.MatchRelative().AtParent().AtName("aws_iam_type"),
					}...),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"x509_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "MANAGED", "CUSTOMER"),
				},
			},
			"oidc_auth_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
			},
			"ldap_auth_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "USER", "GROUP"),
				},
			},
			"aws_iam_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "USER", "ROLE"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"roles": schema.SetNestedBlock{
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"collection_name": schema.StringAttribute{
							Optional: true,
						},
						"database_name": schema.StringAttribute{
							Required: true,
						},
						"role_name": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"labels": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"value": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"scopes": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *databaseUserRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *TfDatabaseUserModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, localDiags := NewMongoDBDatabaseUser(ctx, types.StringNull(), types.StringNull(), plan)
	resp.Diagnostics.Append(localDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	dbUser, _, err := connV2.DatabaseUsersApi.CreateDatabaseUser(ctx, plan.ProjectID.ValueString(), dbUserReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := NewTfDatabaseUserModel(ctx, plan, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning("If the password value will be managed externally it is advised to remove the attribute", "More details can be found in resource documentation under the 'password' attribute")

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseUserRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *TfDatabaseUserModel
	var err error
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := state.Username.ValueString()
	projectID := state.ProjectID.ValueString()
	authDatabaseName := state.AuthDatabaseName.ValueString()
	connV2 := r.Client.AtlasV2
	dbUser, httpResponse, err := connV2.DatabaseUsersApi.GetDatabaseUser(ctx, projectID, authDatabaseName, username).Execute()
	if err != nil {
		// case 404
		// deleted in the backend case
		if validate.StatusNotFound(httpResponse) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := NewTfDatabaseUserModel(ctx, state, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseUserRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *TfDatabaseUserModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, localDiags := NewMongoDBDatabaseUser(ctx, state.Password, state.Description, plan)
	resp.Diagnostics.Append(localDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	dbUser, _, err := connV2.DatabaseUsersApi.UpdateDatabaseUser(ctx,
		plan.ProjectID.ValueString(),
		plan.AuthDatabaseName.ValueString(),
		plan.Username.ValueString(), dbUserReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := NewTfDatabaseUserModel(ctx, plan, dbUser)
	resp.Diagnostics.Append(diagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dbUserModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseUserRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TfDatabaseUserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	_, err := connV2.DatabaseUsersApi.DeleteDatabaseUser(
		ctx,
		state.ProjectID.ValueString(),
		state.AuthDatabaseName.ValueString(),
		state.Username.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error when destroying the database user resource", err.Error())
		return
	}
}

func (r *databaseUserRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, username, authDatabaseName, err := SplitDatabaseUserImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting database User info from ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), username)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("auth_database_name"), authDatabaseName)...)
}

func SplitDatabaseUserImportID(id string) (projectID, username, authDatabaseName string, err error) {
	ok, splitParts := conversion.ImportSplit(id, 3)
	if ok {
		projectID, username, authDatabaseName = splitParts[0], splitParts[1], splitParts[2]
		err = conversion.ValidateProjectID(projectID)
		return
	}
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)-([$a-z]{1,15})$`)
	parts := re.FindStringSubmatch(id)
	if len(parts) != 4 {
		err = errors.New(ErrorImportFormat)
		return
	}
	projectID = parts[1]
	username = parts[2]
	authDatabaseName = parts[3]
	return
}
