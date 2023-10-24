package mongodbatlas

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	autogen "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/autogen/databaseuser/output"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	databaseUserResourceName = "database_user"
)

var _ resource.ResourceWithConfigure = &DatabaseUserRS{}
var _ resource.ResourceWithImportState = &DatabaseUserRS{}

type DatabaseUserRS struct {
	RSCommon
}

func NewDatabaseUserRS() resource.Resource {
	return &DatabaseUserRS{
		RSCommon: RSCommon{
			resourceName: databaseUserResourceName,
		},
	}
}

func (r *DatabaseUserRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = autogen.DatabaseUserResourceSchema(ctx)
}

func (r *DatabaseUserRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var databaseUserPlan *autogen.DatabaseUserModel

	diags := req.Plan.Get(ctx, &databaseUserPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, d := newMongoDBDatabaseUser(ctx, databaseUserPlan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Create(ctx, databaseUserPlan.ProjectId.ValueString(), dbUserReq)
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, databaseUserPlan, dbUser)
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
	var databaseUserState *autogen.DatabaseUserModel
	var err error
	resp.Diagnostics.Append(req.State.Get(ctx, &databaseUserState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := databaseUserState.Username.ValueString()
	projectID := databaseUserState.ProjectId.ValueString()
	authDatabaseName := databaseUserState.AuthDatabaseName.ValueString()

	// Use the ID only with the IMPORT operation
	if databaseUserState.Id.ValueString() != "" && (username == "" || projectID == "" || authDatabaseName == "") {
		projectID, username, authDatabaseName, err = splitDatabaseUserImportID(databaseUserState.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error splitting database User info from ID", err.Error())
			return
		}
	}

	conn := r.client.Atlas
	dbUser, httpResponse, err := conn.DatabaseUsers.Get(ctx, authDatabaseName, projectID, username)
	if err != nil {
		// case 404
		// deleted in the backend case
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddError("resource not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("error getting database user information", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, databaseUserState, dbUser)
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
	var databaseUserPlan *autogen.DatabaseUserModel

	diags := req.Plan.Get(ctx, &databaseUserPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbUserReq, d := newMongoDBDatabaseUser(ctx, databaseUserPlan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	dbUser, _, err := conn.DatabaseUsers.Update(ctx, databaseUserPlan.ProjectId.ValueString(), databaseUserPlan.Username.ValueString(), dbUserReq)
	if err != nil {
		resp.Diagnostics.AddError("error during database user creation", err.Error())
		return
	}

	dbUserModel, diagnostic := newTFDatabaseUserModel(ctx, databaseUserPlan, dbUser)
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
	var databaseUserState *autogen.DatabaseUserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &databaseUserState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	_, err := conn.DatabaseUsers.Delete(ctx, databaseUserState.AuthDatabaseName.ValueString(), databaseUserState.ProjectId.ValueString(), databaseUserState.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error when destroying the database user resource", err.Error())
		return
	}
}

func (r *DatabaseUserRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func newMongoDBDatabaseUser(ctx context.Context, dbUserModel *autogen.DatabaseUserModel) (*matlas.DatabaseUser, diag.Diagnostics) {
	var rolesModel []autogen.RolesValue
	var labelsModel []autogen.LabelsValue
	var scopesModel []autogen.ScopesValue

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
		GroupID:      dbUserModel.ProjectId.ValueString(),
		Username:     dbUserModel.Username.ValueString(),
		Password:     dbUserModel.Password.ValueString(),
		X509Type:     dbUserModel.X509Type.ValueString(),
		AWSIAMType:   dbUserModel.AwsIamType.ValueString(),
		OIDCAuthType: dbUserModel.OidcAuthType.ValueString(),
		LDAPAuthType: dbUserModel.LdapAuthType.ValueString(),
		DatabaseName: dbUserModel.AuthDatabaseName.ValueString(),
		Roles:        newMongoDBAtlasRoles(rolesModel),
		Labels:       newMongoDBAtlasLabels(labelsModel),
		Scopes:       newMongoDBAtlasScopes(scopesModel),
	}, nil
}

func newTFDatabaseUserModel(ctx context.Context, model *autogen.DatabaseUserModel, dbUser *matlas.DatabaseUser) (*autogen.DatabaseUserModel, diag.Diagnostics) {
	rolesSet, diagnostic := types.SetValueFrom(ctx, autogen.RolesValue{}.Type(ctx), newTFRolesModel(ctx, dbUser.Roles))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	labelsSet, diagnostic := types.SetValueFrom(ctx, autogen.LabelsValue{}.Type(ctx), newTFLabelsModel(ctx, dbUser.Labels))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	scopesSet, diagnostic := types.SetValueFrom(ctx, autogen.ScopesValue{}.Type(ctx), newTFScopesModel(dbUser.Scopes))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	// ID is encoded to preserve format defined in previous versions.
	encodedID := encodeStateID(map[string]string{
		"project_id":         dbUser.GroupID,
		"username":           dbUser.Username,
		"auth_database_name": dbUser.DatabaseName,
	})
	databaseUserModel := &autogen.DatabaseUserModel{
		Id:               types.StringValue(encodedID),
		ProjectId:        types.StringValue(dbUser.GroupID),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		X509Type:         types.StringValue(dbUser.X509Type),
		OidcAuthType:     types.StringValue(dbUser.OIDCAuthType),
		LdapAuthType:     types.StringValue(dbUser.LDAPAuthType),
		AwsIamType:       types.StringValue(dbUser.AWSIAMType),
		Roles:            rolesSet,
		Labels:           labelsSet,
		Scopes:           scopesSet,
	}

	if model != nil && model.Password.ValueString() != "" {
		// The Password is not retuned from the endpoint so we use the one provided in the model
		databaseUserModel.Password = model.Password
	}

	return databaseUserModel, nil
}

func newTFScopesModel(scopes []matlas.Scope) []autogen.ScopesValue {
	if len(scopes) == 0 {
		return nil
	}

	out := make([]autogen.ScopesValue, len(scopes))
	for i, v := range scopes {
		out[i] = autogen.ScopesValue{
			Name:  types.StringValue(v.Name),
			Typee: types.StringValue(v.Type),
		}
	}

	return out
}

func newTFLabelsModel(ctx context.Context, labels []matlas.Label) []basetypes.ObjectValuable {
	if len(labels) == 0 {
		return nil
	}

	out := make([]basetypes.ObjectValuable, len(labels))
	for i, v := range labels {
		value := autogen.LabelsValue{
			Key:   types.StringValue(v.Key),
			Value: types.StringValue(v.Value),
		}

		objVal, _ := value.ToObjectValue(ctx)
		objValuable, _ := autogen.LabelsType{}.ValueFromObject(ctx, objVal)
		out[i] = objValuable
	}

	return out
}

func newTFRolesModel(ctx context.Context, roles []matlas.Role) []basetypes.ObjectValuable {
	if len(roles) == 0 {
		return nil
	}

	out := make([]basetypes.ObjectValuable, len(roles))
	for i, v := range roles {
		value := autogen.RolesValue{
			RoleName:     types.StringValue(v.RoleName),
			DatabaseName: types.StringValue(v.DatabaseName),
		}

		if v.CollectionName != "" {
			value.CollectionName = types.StringValue(v.CollectionName)
		}

		objVal, _ := value.ToObjectValue(ctx)
		objValuable, _ := autogen.RolesType{}.ValueFromObject(ctx, objVal)
		out[i] = objValuable
	}

	return out
}

func newMongoDBAtlasRoles(roles []autogen.RolesValue) []matlas.Role {
	if len(roles) == 0 {
		return []matlas.Role{}
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

func newMongoDBAtlasLabels(labels []autogen.LabelsValue) []matlas.Label {
	if len(labels) == 0 {
		return []matlas.Label{}
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

func newMongoDBAtlasScopes(scopes []autogen.ScopesValue) []matlas.Scope {
	if len(scopes) == 0 {
		return []matlas.Scope{}
	}

	out := make([]matlas.Scope, len(scopes))
	for i, v := range scopes {
		out[i] = matlas.Scope{
			Name: v.Name.ValueString(),
			Type: v.Typee.ValueString(),
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
