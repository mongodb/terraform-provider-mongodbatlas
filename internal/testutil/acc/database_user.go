package acc

import (
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312027/admin"
)

func ConfigDatabaseUserBasic(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = %[4]q
				value = %[5]q
			}
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithX509Type(projectID, username, x509Type, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			x509_type          = %[3]q
			auth_database_name = "$external"

			roles {
				role_name     = %[4]q
				database_name = "admin"
			}

			labels {
				key   = %[5]q
				value = %[6]q
			}
		}
	`, projectID, username, x509Type, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithLabels(projectID, username, roleName, description string, labels []admin.ComponentLabel) string {
	var labelsConf strings.Builder
	for _, label := range labels {
		fmt.Fprintf(&labelsConf, `
			labels {
				key   = %q
				value = %q
			}
		`, label.GetKey(), label.GetValue())
	}
	var descriptionConf string
	if description != "" {
		descriptionConf = fmt.Sprintf(`	description = %q`, description)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			%[4]s
			%[5]s
		}
	`, projectID, username, roleName, labelsConf.String(), descriptionConf)
}

func ConfigDatabaseUserWithRoles(projectID, username, password string, rolesArr []*admin.DatabaseUserRole) string {
	var roles strings.Builder
	for _, role := range rolesArr {
		var roleName, databaseName, collection string
		if role.RoleName != "" {
			roleName = fmt.Sprintf(`role_name = %q`, role.RoleName)
		}
		if role.DatabaseName != "" {
			databaseName = fmt.Sprintf(`database_name = %q`, role.DatabaseName)
		}
		if role.GetCollectionName() != "" {
			collection = fmt.Sprintf(`collection_name = %q`, role.GetCollectionName())
		}
		fmt.Fprintf(&roles, `
			roles {
				%s
				%s
				%s
			}
		`, roleName, databaseName, collection)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = %[3]q
			auth_database_name = "admin"

			%[4]s
		}
	`, projectID, username, password, roles.String())
}

func ConfigDatabaseUserWithAWSIAMType(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			aws_iam_type       = "USER"
			auth_database_name = "$external"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = %[4]q
				value = %[5]q
			}
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithScopes(projectID, username, password, roleName string, scopesArr []*admin.UserScope) string {
	var scopes strings.Builder
	for _, scope := range scopesArr {
		fmt.Fprintf(&scopes, `
			scopes {
				name = %q
				type = %q
			}
		`, scope.GetName(), scope.GetType())
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = %[3]q
			auth_database_name = "admin"

			roles {
				role_name     = %[4]q
				database_name = "admin"
			}

			%[5]s
		}
	`, projectID, username, password, roleName, scopes.String())
}

func ConfigDatabaseUserWithLDAPAuthType(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			ldap_auth_type     = "USER"
			auth_database_name = "$external"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = %[4]q
				value = %[5]q
			}
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDataBaseUserWithOIDCAuthType(projectID, username, authType, databaseName, roleName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_database_user" "test" {
		project_id         = %[1]q
		username           = %[2]q
		oidc_auth_type     = %[3]q
		auth_database_name = %[4]q

		roles = [{ role_name = %[5]q, database_name = "admin" }]
	}
`, projectID, username, authType, databaseName, roleName)
}

func ConfigDatabaseUserPasswordWo(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id          = %[1]q
			username            = %[2]q
			password_wo         = "test-acc-password-wo"
			password_wo_version = 1
			auth_database_name  = "admin"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`, projectID, username)
}

func ConfigDatabaseUserPasswordWoVersionUpdate(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id          = %[1]q
			username            = %[2]q
			password_wo         = "test-acc-password-wo-updated"
			password_wo_version = 2
			auth_database_name  = "admin"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`, projectID, username)
}

// TPF-syntax variants of the configs above: nested collections written as assignments
// (`roles = [{ ... }]`) instead of blocks. The autogenerated implementation models nested
// collections as nested attributes, which reject block syntax, so acceptance tests for the
// autogen resource and the step that runs it in migration tests use these variants
// (customers rewriting their configs from blocks to assignments is the accepted breaking
// change; see the spike doc).

func ConfigDatabaseUserBasicTPF(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			auth_database_name = "admin"

			roles  = [{ role_name = %[3]q, database_name = "admin" }]
			labels = [{ key = %[4]q, value = %[5]q }]
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithX509TypeTPF(projectID, username, x509Type, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			x509_type          = %[3]q
			auth_database_name = "$external"

			roles  = [{ role_name = %[4]q, database_name = "admin" }]
			labels = [{ key = %[5]q, value = %[6]q }]
		}
	`, projectID, username, x509Type, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithLabelsTPF(projectID, username, roleName, description string, labels []admin.ComponentLabel) string {
	var labelsConf []string
	for _, label := range labels {
		labelsConf = append(labelsConf, fmt.Sprintf(`{ key = %q, value = %q }`, label.GetKey(), label.GetValue()))
	}
	var descriptionConf string
	if description != "" {
		descriptionConf = fmt.Sprintf(`	description = %q`, description)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			auth_database_name = "admin"

			roles  = [{ role_name = %[3]q, database_name = "admin" }]
			labels = [%[4]s]
			%[5]s
		}
	`, projectID, username, roleName, strings.Join(labelsConf, ", "), descriptionConf)
}

func ConfigDatabaseUserWithRolesTPF(projectID, username, password string, rolesArr []*admin.DatabaseUserRole) string {
	var roles []string
	for _, role := range rolesArr {
		var parts []string
		if role.RoleName != "" {
			parts = append(parts, fmt.Sprintf(`role_name = %q`, role.RoleName))
		}
		if role.DatabaseName != "" {
			parts = append(parts, fmt.Sprintf(`database_name = %q`, role.DatabaseName))
		}
		if role.GetCollectionName() != "" {
			parts = append(parts, fmt.Sprintf(`collection_name = %q`, role.GetCollectionName()))
		}
		roles = append(roles, fmt.Sprintf(`{ %s }`, strings.Join(parts, ", ")))
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = %[3]q
			auth_database_name = "admin"

			roles = [%[4]s]
		}
	`, projectID, username, password, strings.Join(roles, ", "))
}

func ConfigDatabaseUserWithAWSIAMTypeTPF(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			aws_iam_type       = "USER"
			auth_database_name = "$external"

			roles  = [{ role_name = %[3]q, database_name = "admin" }]
			labels = [{ key = %[4]q, value = %[5]q }]
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserWithScopesTPF(projectID, username, password, roleName string, scopesArr []*admin.UserScope) string {
	var scopes []string
	for _, scope := range scopesArr {
		scopes = append(scopes, fmt.Sprintf(`{ name = %q, type = %q }`, scope.GetName(), scope.GetType()))
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = %[3]q
			auth_database_name = "admin"

			roles  = [{ role_name = %[4]q, database_name = "admin" }]
			scopes = [%[5]s]
		}
	`, projectID, username, password, roleName, strings.Join(scopes, ", "))
}

func ConfigDatabaseUserWithLDAPAuthTypeTPF(projectID, username, roleName, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			ldap_auth_type     = "USER"
			auth_database_name = "$external"

			roles  = [{ role_name = %[3]q, database_name = "admin" }]
			labels = [{ key = %[4]q, value = %[5]q }]
		}
	`, projectID, username, roleName, keyLabel, valueLabel)
}

func ConfigDatabaseUserPasswordWoTPF(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id          = %[1]q
			username            = %[2]q
			password_wo         = "test-acc-password-wo"
			password_wo_version = 1
			auth_database_name  = "admin"

			roles = [{ role_name = "atlasAdmin", database_name = "admin" }]
		}
	`, projectID, username)
}

func ConfigDatabaseUserPasswordWoVersionUpdateTPF(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id          = %[1]q
			username            = %[2]q
			password_wo         = "test-acc-password-wo-updated"
			password_wo_version = 2
			auth_database_name  = "admin"

			roles = [{ role_name = "atlasAdmin", database_name = "admin" }]
		}
	`, projectID, username)
}
