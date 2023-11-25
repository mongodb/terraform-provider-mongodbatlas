package acc

import (
	"fmt"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DatabaseUserConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "%[4]s"
			password           = "test-acc-password"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}

func DatabaseUserWithX509TypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel, x509Type string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			x509_type          = "%[7]s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel, x509Type)
}

func DatabaseUserWithLabelsConfig(projectName, orgID, roleName, username string, labels []matlas.Label) string {
	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			password           = "test-acc-password"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			%[5]s

		}
	`, projectName, orgID, roleName, username, labelsConf)
}

func DatabaseUserWithRoles(username, password, projectName, orgID string, rolesArr []*matlas.Role) string {
	var roles string

	for _, role := range rolesArr {
		var roleName, databaseName, collection string

		if role.RoleName != "" {
			roleName = fmt.Sprintf(`role_name = %q`, role.RoleName)
		}

		if role.DatabaseName != "" {
			databaseName = fmt.Sprintf(`database_name = %q`, role.DatabaseName)
		}

		if role.CollectionName != "" {
			collection = fmt.Sprintf(`collection_name = %q`, role.CollectionName)
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
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			%s

		}
	`, projectName, orgID, username, password, roles)
}

func DatabaseUserWithAWSIAMTypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			aws_iam_type       = "USER"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}

func DatabaseUserWithScopes(username, password, projectName, orgID, roleName, clusterName string, scopesArr []*matlas.Scope) string {
	var scopes string

	for _, scope := range scopesArr {
		var scopeType string

		if scope.Type != "" {
			scopeType = fmt.Sprintf(`type = %q`, scope.Type)
		}

		scopes += fmt.Sprintf(`
			scopes {
				name = "${mongodbatlas_cluster.my_cluster.name}"
				%s
			}
		`, scopeType)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "${mongodbatlas_project.test.id}"
			name         = "%s"
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			cloud_backup                = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%s"
				database_name = "admin"
			}

			%s

		}
	`, projectName, orgID, clusterName, username, password, roleName, scopes)
}

func DatabaseUserWithLDAPAuthTypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			ldap_auth_type     = "USER"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}
