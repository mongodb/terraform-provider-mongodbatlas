package acc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"

	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func CheckDatabaseUserExists(resourceName string, dbUser *admin.CloudDatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		connV2 := TestMongoDBClient.(*config.MongoDBClient).AtlasV2

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no project_id is set")
		}

		if rs.Primary.Attributes["auth_database_name"] == "" {
			return fmt.Errorf("no auth_database_name is set")
		}

		if rs.Primary.Attributes["username"] == "" {
			return fmt.Errorf("no username is set")
		}

		authDB := rs.Primary.Attributes["auth_database_name"]
		projectID := rs.Primary.Attributes["project_id"]
		username := rs.Primary.Attributes["username"]

		if dbUserResp, _, err := connV2.DatabaseUsersApi.GetDatabaseUser(context.Background(), projectID, authDB, username).Execute(); err == nil {
			*dbUser = *dbUserResp
			return nil
		}

		return fmt.Errorf("database user(%s-%s-%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"], rs.Primary.Attributes["auth_database_name"])
	}
}

func CheckDatabaseUserAttributes(dbUser *admin.CloudDatabaseUser, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dbUser.Username: %s , username : %s", dbUser.Username, username)
		if dbUser.Username != username {
			return fmt.Errorf("bad username: %s", dbUser.Username)
		}

		return nil
	}
}

func CheckDestroyDatabaseUser(s *terraform.State) error {
	connV2 := TestMongoDBClient.(*config.MongoDBClient).AtlasV2

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user" {
			continue
		}

		projectID, username, authDatabaseName, err := databaseuser.SplitDatabaseUserImportID(rs.Primary.ID)
		if err != nil {
			continue
		}
		// Try to find the database user
		_, _, err = connV2.DatabaseUsersApi.GetDatabaseUser(context.Background(), projectID, authDatabaseName, username).Execute()
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", projectID)
		}
	}

	return nil
}

func ConfigDatabaseUserBasic(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
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

func ConfigDatabaseUserWithX509Type(projectName, orgID, roleName, username, keyLabel, valueLabel, x509Type string) string {
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

func ConfigDatabaseUserWithLabels(projectName, orgID, roleName, username string, labels []admin.ComponentLabel) string {
	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = "%s"
				value = "%s"
			}
		`, label.GetKey(), label.GetValue())
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

func ConfigDatabaseUserWithRoles(username, password, projectName, orgID string, rolesArr []*admin.DatabaseUserRole) string {
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

func ConfigDatabaseUserWithAWSIAMType(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
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

func ConfigDatabaseUserWithScopes(username, password, roleName, projectId, clusterName, clusterTerraformStr string, scopesArr []*admin.UserScope) string {
	var scopes string

	for _, scope := range scopesArr {
		var scopeType string

		if scope.Type != "" {
			scopeType = fmt.Sprintf(`type = %q`, scope.Type)
		}

		scopes += fmt.Sprintf(`
			scopes {
				name = "%s"
				%s
			}
		`, clusterName, scopeType)
	}

	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = %s
			auth_database_name = "admin"

			roles {
				role_name     = "%s"
				database_name = "admin"
			}

			%s

		}
	`, username, password, projectId, roleName, scopes)
}

func ConfigDatabaseUserWithLDAPAuthType(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
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
