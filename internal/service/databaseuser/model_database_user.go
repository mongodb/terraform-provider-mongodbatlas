package databaseuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func NewMongoDBDatabaseUser(ctx context.Context, statePasswordValue types.String, dbUserModel *TfDatabaseUserModel) (*admin.CloudDatabaseUser, diag.Diagnostics) {
	var rolesModel []*TfRoleModel
	var labelsModel []*TfLabelModel
	var scopesModel []*TfScopeModel

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

	result := admin.CloudDatabaseUser{
		GroupId:      dbUserModel.ProjectID.ValueString(),
		Username:     dbUserModel.Username.ValueString(),
		X509Type:     dbUserModel.X509Type.ValueStringPointer(),
		AwsIAMType:   dbUserModel.AWSIAMType.ValueStringPointer(),
		OidcAuthType: dbUserModel.OIDCAuthType.ValueStringPointer(),
		LdapAuthType: dbUserModel.LDAPAuthType.ValueStringPointer(),
		DatabaseName: dbUserModel.AuthDatabaseName.ValueString(),
		Roles:        NewMongoDBAtlasRoles(rolesModel),
		Labels:       NewMongoDBAtlasLabels(labelsModel),
		Scopes:       NewMongoDBAtlasScopes(scopesModel),
	}

	if statePasswordValue != dbUserModel.Password {
		// Password value has been modified or no previous state was present. Password is only updated if changed in the terraform configuration CLOUDP-235738
		result.Password = dbUserModel.Password.ValueStringPointer()
	}
	return &result, nil
}

func NewTfDatabaseUserModel(ctx context.Context, model *TfDatabaseUserModel, dbUser *admin.CloudDatabaseUser) (*TfDatabaseUserModel, diag.Diagnostics) {
	rolesSet, diagnostic := types.SetValueFrom(ctx, RoleObjectType, NewTFRolesModel(dbUser.GetRoles()))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	labelsSet, diagnostic := types.SetValueFrom(ctx, LabelObjectType, NewTFLabelsModel(dbUser.GetLabels()))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	scopesSet, diagnostic := types.SetValueFrom(ctx, ScopeObjectType, NewTFScopesModel(dbUser.GetScopes()))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	// ID is encoded to preserve format defined in previous versions.
	encodedID := conversion.EncodeStateID(map[string]string{
		"project_id":         dbUser.GroupId,
		"username":           dbUser.Username,
		"auth_database_name": dbUser.DatabaseName,
	})
	databaseUserModel := &TfDatabaseUserModel{
		ID:               types.StringValue(encodedID),
		ProjectID:        types.StringValue(dbUser.GroupId),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		X509Type:         types.StringValue(dbUser.GetX509Type()),
		OIDCAuthType:     types.StringValue(dbUser.GetOidcAuthType()),
		LDAPAuthType:     types.StringValue(dbUser.GetLdapAuthType()),
		AWSIAMType:       types.StringValue(dbUser.GetAwsIAMType()),
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

func NewTFDatabaseDSUserModel(ctx context.Context, dbUser *admin.CloudDatabaseUser) (*TfDatabaseUserDSModel, diag.Diagnostics) {
	databaseID := fmt.Sprintf("%s-%s-%s", dbUser.GroupId, dbUser.Username, dbUser.DatabaseName)
	databaseUserModel := &TfDatabaseUserDSModel{
		ID:               types.StringValue(databaseID),
		ProjectID:        types.StringValue(dbUser.GroupId),
		AuthDatabaseName: types.StringValue(dbUser.DatabaseName),
		Username:         types.StringValue(dbUser.Username),
		X509Type:         types.StringValue(dbUser.GetX509Type()),
		OIDCAuthType:     types.StringValue(dbUser.GetOidcAuthType()),
		LDAPAuthType:     types.StringValue(dbUser.GetLdapAuthType()),
		AWSIAMType:       types.StringValue(dbUser.GetAwsIAMType()),
		Roles:            NewTFRolesModel(dbUser.GetRoles()),
		Labels:           NewTFLabelsModel(dbUser.GetLabels()),
		Scopes:           NewTFScopesModel(dbUser.GetScopes()),
	}

	return databaseUserModel, nil
}

func NewTFDatabaseUsersModel(ctx context.Context, projectID string, dbUsers []admin.CloudDatabaseUser) (*TfDatabaseUsersDSModel, diag.Diagnostics) {
	results := make([]*TfDatabaseUserDSModel, len(dbUsers))
	for i := range dbUsers {
		dbUserModel, d := NewTFDatabaseDSUserModel(ctx, &dbUsers[i])
		if d.HasError() {
			return nil, d
		}
		results[i] = dbUserModel
	}

	return &TfDatabaseUsersDSModel{
		ProjectID: types.StringValue(projectID),
		Results:   results,
		ID:        types.StringValue(id.UniqueId()),
	}, nil
}

func NewTFScopesModel(scopes []admin.UserScope) []TfScopeModel {
	out := make([]TfScopeModel, len(scopes))
	for i, v := range scopes {
		out[i] = TfScopeModel{
			Name: types.StringValue(v.Name),
			Type: types.StringValue(v.Type),
		}
	}
	return out
}

func NewMongoDBAtlasLabels(labels []*TfLabelModel) *[]admin.ComponentLabel {
	out := make([]admin.ComponentLabel, len(labels))
	for i, v := range labels {
		out[i] = admin.ComponentLabel{
			Key:   v.Key.ValueStringPointer(),
			Value: v.Value.ValueStringPointer(),
		}
	}
	return &out
}

func NewTFLabelsModel(labels []admin.ComponentLabel) []TfLabelModel {
	out := make([]TfLabelModel, len(labels))
	for i, v := range labels {
		out[i] = TfLabelModel{
			Key:   types.StringValue(v.GetKey()),
			Value: types.StringValue(v.GetValue()),
		}
	}
	return out
}

func NewMongoDBAtlasScopes(scopes []*TfScopeModel) *[]admin.UserScope {
	out := make([]admin.UserScope, len(scopes))
	for i, v := range scopes {
		out[i] = admin.UserScope{
			Name: v.Name.ValueString(),
			Type: v.Type.ValueString(),
		}
	}
	return &out
}

func NewTFRolesModel(roles []admin.DatabaseUserRole) []TfRoleModel {
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

func NewMongoDBAtlasRoles(roles []*TfRoleModel) *[]admin.DatabaseUserRole {
	out := make([]admin.DatabaseUserRole, len(roles))
	for i, v := range roles {
		out[i] = admin.DatabaseUserRole{
			RoleName:       v.RoleName.ValueString(),
			DatabaseName:   v.DatabaseName.ValueString(),
			CollectionName: v.CollectionName.ValueStringPointer(),
		}
	}
	return &out
}
