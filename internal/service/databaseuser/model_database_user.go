package databaseuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func NewMongoDBDatabaseUser(ctx context.Context, dbUserModel *TfDatabaseUserModel) (*admin.CloudDatabaseUser, diag.Diagnostics) {
	var rolesModel []*TfRoleModel
	var labelsModel []*tfLabelModel
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

	return &admin.CloudDatabaseUser{
		GroupId:      dbUserModel.ProjectID.ValueString(),
		Username:     dbUserModel.Username.ValueString(),
		Password:     dbUserModel.Password.ValueStringPointer(),
		X509Type:     dbUserModel.X509Type.ValueStringPointer(),
		AwsIAMType:   dbUserModel.AWSIAMType.ValueStringPointer(),
		OidcAuthType: dbUserModel.OIDCAuthType.ValueStringPointer(),
		LdapAuthType: dbUserModel.LDAPAuthType.ValueStringPointer(),
		DatabaseName: dbUserModel.AuthDatabaseName.ValueString(),
		Roles:        NewMongoDBAtlasRoles(rolesModel),
		Labels:       NewMongoDBAtlasLabels(labelsModel),
		Scopes:       NewMongoDBAtlasScopes(scopesModel),
	}, nil
}

func NewTfDatabaseUserModel(ctx context.Context, model *TfDatabaseUserModel, dbUser *admin.CloudDatabaseUser) (*TfDatabaseUserModel, diag.Diagnostics) {
	rolesSet, diagnostic := types.SetValueFrom(ctx, RoleObjectType, NewTFRolesModel(dbUser.Roles))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	labelsSet, diagnostic := types.SetValueFrom(ctx, LabelObjectType, NewTFLabelsModel(dbUser.Labels))
	if diagnostic.HasError() {
		return nil, diagnostic
	}

	scopesSet, diagnostic := types.SetValueFrom(ctx, ScopeObjectType, NewTFScopesModel(dbUser.Scopes))
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

func NewTFScopesModel(scopes []admin.UserScope) []TfScopeModel {
	if len(scopes) == 0 {
		return nil
	}

	out := make([]TfScopeModel, len(scopes))
	for i, v := range scopes {
		out[i] = TfScopeModel{
			Name: types.StringValue(v.Name),
			Type: types.StringValue(v.Type),
		}
	}

	return out
}

func NewMongoDBAtlasLabels(labels []*tfLabelModel) []admin.ComponentLabel {
	if len(labels) == 0 {
		return []admin.ComponentLabel{}
	}

	out := make([]admin.ComponentLabel, len(labels))
	for i, v := range labels {
		out[i] = admin.ComponentLabel{
			Key:   v.Key.ValueStringPointer(),
			Value: v.Value.ValueStringPointer(),
		}
	}

	return out
}

func NewMongoDBAtlasScopes(scopes []*TfScopeModel) []admin.UserScope {
	if len(scopes) == 0 {
		return []admin.UserScope{}
	}

	out := make([]admin.UserScope, len(scopes))
	for i, v := range scopes {
		out[i] = admin.UserScope{
			Name: v.Name.ValueString(),
			Type: v.Type.ValueString(),
		}
	}

	return out
}
