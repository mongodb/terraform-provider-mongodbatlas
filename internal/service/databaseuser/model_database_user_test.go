package databaseuser_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

var (
	projectID        = "projectID"
	authDatabaseName = "AuthDatabaseName"
	username         = "Username"
	password         = "Password"
	x509Type         = "X509Type"
	oidCAuthType     = "OIDCAuthType"
	ldapAuthType     = "LDAPAuthType"
	awsIAMType       = "AWSIAMType"
	roleName         = "roleName"
	collectionName   = "collectionName"
	databaseName     = "databaseName"
	key              = "key"
	value            = "value"
	name             = "name"
	typeVar          = "type"
	tfUserRole       = databaseuser.TfRoleModel{
		RoleName:       types.StringValue(roleName),
		CollectionName: types.StringValue(collectionName),
		DatabaseName:   types.StringValue(databaseName),
	}
	tfLabel = databaseuser.TfLabelModel{
		Key:   types.StringValue(key),
		Value: types.StringValue(value),
	}
	tfScope = databaseuser.TfScopeModel{
		Name: types.StringValue(name),
		Type: types.StringValue(typeVar),
	}
	sdkScope = admin.UserScope{
		Name: name,
		Type: typeVar,
	}
	sdkLabel = admin.ComponentLabel{
		Key:   &key,
		Value: &value,
	}
	sdkRole = admin.DatabaseUserRole{
		CollectionName: &collectionName,
		DatabaseName:   databaseName,
		RoleName:       roleName,
	}
	rolesSet, _       = types.SetValueFrom(context.Background(), databaseuser.RoleObjectType, []databaseuser.TfRoleModel{tfUserRole})
	labelsSet, _      = types.SetValueFrom(context.Background(), databaseuser.LabelObjectType, []databaseuser.TfLabelModel{tfLabel})
	scopesSet, _      = types.SetValueFrom(context.Background(), databaseuser.ScopeObjectType, []databaseuser.TfScopeModel{tfScope})
	wrongRoleSet, _   = types.SetValueFrom(context.Background(), databaseuser.LabelObjectType, []databaseuser.TfRoleModel{tfUserRole})
	wrongLabelSet, _  = types.SetValueFrom(context.Background(), databaseuser.RoleObjectType, []databaseuser.TfLabelModel{tfLabel})
	wrongScopeSet, _  = types.SetValueFrom(context.Background(), databaseuser.RoleObjectType, []databaseuser.TfScopeModel{tfScope})
	cloudDatabaseUser = &admin.CloudDatabaseUser{
		GroupId:      projectID,
		DatabaseName: authDatabaseName,
		Username:     username,
		Description:  conversion.Pointer(""),
		Password:     &password,
		X509Type:     &x509Type,
		OidcAuthType: &oidCAuthType,
		LdapAuthType: &ldapAuthType,
		AwsIAMType:   &awsIAMType,
		Roles:        &[]admin.DatabaseUserRole{sdkRole},
		Labels:       &[]admin.ComponentLabel{sdkLabel},
		Scopes:       &[]admin.UserScope{sdkScope},
	}
	cloudDatabaseUserWithoutPassword = &admin.CloudDatabaseUser{
		GroupId:      projectID,
		DatabaseName: authDatabaseName,
		Description:  conversion.Pointer(""),
		Username:     username,
		X509Type:     &x509Type,
		OidcAuthType: &oidCAuthType,
		LdapAuthType: &ldapAuthType,
		AwsIAMType:   &awsIAMType,
		Roles:        &[]admin.DatabaseUserRole{sdkRole},
		Labels:       &[]admin.ComponentLabel{sdkLabel},
		Scopes:       &[]admin.UserScope{sdkScope},
	}
)

func TestNewMongoDBDatabaseUser(t *testing.T) {
	testCases := []struct {
		tfDatabaseUserModel databaseuser.TfDatabaseUserModel
		passwordStateValue  types.String
		expectedResult      *admin.CloudDatabaseUser
		name                string
		expectedError       bool
	}{
		{
			name:                "CloudDatabaseUser for create",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, labelsSet, scopesSet, types.StringValue(password)),
			passwordStateValue:  types.StringNull(),
			expectedResult:      cloudDatabaseUser,
			expectedError:       false,
		},
		{
			name:                "CloudDatabaseUser with new password in model when password value is modified",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, labelsSet, scopesSet, types.StringValue(password)),
			passwordStateValue:  types.StringValue("oldPassword"),
			expectedResult:      cloudDatabaseUser,
			expectedError:       false,
		},
		{
			name:                "CloudDatabaseUser with no password in model when password value remains the same",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, labelsSet, scopesSet, types.StringValue(password)),
			passwordStateValue:  types.StringValue(password),
			expectedResult:      cloudDatabaseUserWithoutPassword,
			expectedError:       false,
		},
		{
			name:                "CloudDatabaseUser with no password in model when password value is removed",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, labelsSet, scopesSet, types.StringNull()),
			passwordStateValue:  types.StringValue(password),
			expectedResult:      cloudDatabaseUserWithoutPassword,
			expectedError:       false,
		},
		{
			name:                "Roles fail",
			tfDatabaseUserModel: *getDatabaseUserModel(wrongRoleSet, labelsSet, scopesSet, types.StringValue(password)),
			expectedError:       true,
		},
		{
			name:                "Labels fail",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, wrongLabelSet, scopesSet, types.StringValue(password)),
			expectedError:       true,
		},
		{
			name:                "Scopes fail",
			tfDatabaseUserModel: *getDatabaseUserModel(rolesSet, labelsSet, wrongScopeSet, types.StringValue(password)),
			expectedError:       true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, err := databaseuser.NewMongoDBDatabaseUser(t.Context(), tc.passwordStateValue, types.StringValue(""), &testCases[i].tfDatabaseUserModel)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewTfDatabaseUserModel(t *testing.T) {
	testCases := []struct {
		expectedResult  *databaseuser.TfDatabaseUserModel
		currentModel    databaseuser.TfDatabaseUserModel
		sdkDatabaseUser *admin.CloudDatabaseUser
		name            string
		expectedError   bool
	}{
		{
			name:            "Success TfDatabaseUserModel",
			sdkDatabaseUser: cloudDatabaseUser,
			currentModel:    databaseuser.TfDatabaseUserModel{Password: types.StringValue(password), Description: types.StringValue("")},
			expectedResult:  getDatabaseUserModel(rolesSet, labelsSet, scopesSet, types.StringValue(password)),
			expectedError:   false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, err := databaseuser.NewTfDatabaseUserModel(t.Context(), &testCases[i].currentModel, testCases[i].sdkDatabaseUser)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewMongoDBAtlasScopes(t *testing.T) {
	testCases := []struct {
		expectedResult *[]admin.UserScope
		name           string
		currentScopes  []*databaseuser.TfScopeModel
	}{
		{
			name:           "Success TfScopeModel",
			currentScopes:  []*databaseuser.TfScopeModel{&tfScope},
			expectedResult: &[]admin.UserScope{sdkScope},
		},
		{
			name:           "Empty TfScopeModel",
			currentScopes:  []*databaseuser.TfScopeModel{},
			expectedResult: &[]admin.UserScope{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewMongoDBAtlasScopes(tc.currentScopes)

			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewTFScopesModel(t *testing.T) {
	testCases := []struct {
		name           string
		currentScopes  []admin.UserScope
		expectedResult []databaseuser.TfScopeModel
	}{
		{
			name:           "Success TfScopeModel",
			currentScopes:  []admin.UserScope{sdkScope},
			expectedResult: []databaseuser.TfScopeModel{tfScope},
		},
		{
			name:           "Empty TfScopeModel",
			currentScopes:  []admin.UserScope{},
			expectedResult: []databaseuser.TfScopeModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewTFScopesModel(tc.currentScopes)

			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewMongoDBAtlasLabels(t *testing.T) {
	testCases := []struct {
		expectedResult *[]admin.ComponentLabel
		name           string
		currentLabels  []*databaseuser.TfLabelModel
	}{
		{
			name:           "Success TfLabelModel",
			currentLabels:  []*databaseuser.TfLabelModel{&tfLabel},
			expectedResult: &[]admin.ComponentLabel{sdkLabel},
		},
		{
			name:           "Empty TfLabelModel",
			currentLabels:  []*databaseuser.TfLabelModel{},
			expectedResult: &[]admin.ComponentLabel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewMongoDBAtlasLabels(tc.currentLabels)

			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewTFLabelsModel(t *testing.T) {
	testCases := []struct {
		name           string
		currentLabels  []admin.ComponentLabel
		expectedResult []databaseuser.TfLabelModel
	}{
		{
			name:           "Success TfLabelModel",
			currentLabels:  []admin.ComponentLabel{sdkLabel},
			expectedResult: []databaseuser.TfLabelModel{tfLabel},
		},
		{
			name:           "Empty TfLabelModel",
			currentLabels:  []admin.ComponentLabel{},
			expectedResult: []databaseuser.TfLabelModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewTFLabelsModel(tc.currentLabels)

			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewMongoDBAtlasRoles(t *testing.T) {
	testCases := []struct {
		expectedResult *[]admin.DatabaseUserRole
		name           string
		currentRoles   []*databaseuser.TfRoleModel
	}{
		{
			name:           "Success DatabaseUserRole",
			currentRoles:   []*databaseuser.TfRoleModel{&tfUserRole},
			expectedResult: &[]admin.DatabaseUserRole{sdkRole},
		},
		{
			name:           "Empty DatabaseUserRole",
			currentRoles:   []*databaseuser.TfRoleModel{},
			expectedResult: &[]admin.DatabaseUserRole{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewMongoDBAtlasRoles(tc.currentRoles)

			assert.Equal(t, tc.expectedResult, resultModel, "created SDK model did not match expected output")
		})
	}
}

func TestNewTFRolesModel(t *testing.T) {
	testCases := []struct {
		name           string
		currentRoles   []admin.DatabaseUserRole
		expectedResult []databaseuser.TfRoleModel
	}{
		{
			name:           "Success DatabaseUserRole",
			currentRoles:   []admin.DatabaseUserRole{sdkRole},
			expectedResult: []databaseuser.TfRoleModel{tfUserRole},
		},
		{
			name:           "Empty DatabaseUserRole",
			currentRoles:   []admin.DatabaseUserRole{},
			expectedResult: []databaseuser.TfRoleModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := databaseuser.NewTFRolesModel(tc.currentRoles)

			assert.Equal(t, tc.expectedResult, resultModel, "created terraform model did not match expected output")
		})
	}
}

func getDatabaseUserModel(roles, labels, scopes basetypes.SetValue, password types.String) *databaseuser.TfDatabaseUserModel {
	encodedID := conversion.EncodeStateID(map[string]string{
		"project_id":         projectID,
		"username":           username,
		"auth_database_name": authDatabaseName,
	})
	return &databaseuser.TfDatabaseUserModel{
		ID:               types.StringValue(encodedID),
		ProjectID:        types.StringValue(projectID),
		AuthDatabaseName: types.StringValue(authDatabaseName),
		Username:         types.StringValue(username),
		Description:      types.StringValue(""),
		Password:         password,
		X509Type:         types.StringValue(x509Type),
		OIDCAuthType:     types.StringValue(oidCAuthType),
		LDAPAuthType:     types.StringValue(ldapAuthType),
		AWSIAMType:       types.StringValue(awsIAMType),
		Roles:            roles,
		Labels:           labels,
		Scopes:           scopes,
	}
}

func TestSplitDatabaseUserImportID(t *testing.T) {
	tests := map[string]struct {
		importID    string
		projectID   string
		username    string
		authDBName  string
		errorString string
	}{
		"valid input": {
			importID:   "664619d870c247237f4b86a6/my-username-dash/my-db-name",
			projectID:  "664619d870c247237f4b86a6",
			username:   "my-username-dash",
			authDBName: "my-db-name",
		},
		"valid input legacy": {
			importID:   "664619d870c247237f4b86a6-myUsernameCamel-mydbname",
			projectID:  "664619d870c247237f4b86a6",
			username:   "myUsernameCamel",
			authDBName: "mydbname",
		},
		"invalid input projectID": {
			importID:    "part1/part2/part3",
			projectID:   "part1",
			username:    "part2",
			authDBName:  "part3",
			errorString: "project_id must be a 24 character hex string: part1",
		},
		"invalid input with more parts": {
			importID:    "part1/part2/part3/part4",
			projectID:   "",
			username:    "",
			authDBName:  "",
			errorString: databaseuser.ErrorImportFormat,
		},
		"invalid input with less parts": {
			importID:    "part1/part2",
			projectID:   "",
			username:    "",
			authDBName:  "",
			errorString: databaseuser.ErrorImportFormat,
		},
		"empty input": {
			importID:    "",
			projectID:   "",
			username:    "",
			authDBName:  "",
			errorString: databaseuser.ErrorImportFormat,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			part1, part2, part3, err := databaseuser.SplitDatabaseUserImportID(tc.importID)
			if tc.errorString != "" {
				require.EqualError(t, err, tc.errorString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.projectID, part1)
			assert.Equal(t, tc.username, part2)
			assert.Equal(t, tc.authDBName, part3)
		})
	}
}
