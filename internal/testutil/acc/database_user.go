package acc

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"
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

func ConfigDatabaseUserWithLabels(projectID, username, roleName string, labels []admin.ComponentLabel) string {
	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = %q
				value = %q
			}
		`, label.GetKey(), label.GetValue())
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
		}
	`, projectID, username, roleName, labelsConf)
}

func ConfigDatabaseUserWithRoles(projectID, username, password string, rolesArr []*admin.DatabaseUserRole) string {
	var roles string
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
		roles += fmt.Sprintf(`
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
	`, projectID, username, password, roles)
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
	var scopes string
	for _, scope := range scopesArr {
		scopes += fmt.Sprintf(`
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
	`, projectID, username, password, roleName, scopes)
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

		roles {
			role_name     = %[5]q
			database_name = "admin"
		}
	}
`, projectID, username, authType, databaseName, roleName)
}
