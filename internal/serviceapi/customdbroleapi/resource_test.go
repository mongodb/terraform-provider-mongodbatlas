package customdbroleapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_custom_db_role_api.test"

func TestAccCustomDBRoleAPI_basic(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		groupName     = acc.RandomProjectName()
		roleName      = acc.RandomName()
		databaseName1 = acc.RandomClusterName()
		databaseName2 = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, groupName, roleName, "INSERT", databaseName1),
				Check:  checkBasic(roleName, "INSERT", databaseName1),
			},
			{
				Config: configBasic(orgID, groupName, roleName, "UPDATE", databaseName2),
				Check:  checkBasic(roleName, "UPDATE", databaseName2),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"actions.0.resources.0.cluster"},
				ImportStateVerifyIdentifierAttribute: "role_name",
			},
		},
	})
}

func configBasic(orgID, groupName, roleName, action, db string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
		resource "mongodbatlas_custom_db_role_api" "test" {
			group_id = mongodbatlas_project.test.id
			role_name  = %[3]q
			actions = [
				{
					action = %[4]q
					resources = [
						{
							collection	= ""
							cluster			= false
							db   				= %[5]q
						}
					]
				}
			]
		}
	`, orgID, groupName, roleName, action, db)
}

func checkBasic(roleName, action, db string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"role_name":                        roleName,
		"actions.#":                        "1",
		"actions.0.action":                 action,
		"actions.0.resources.#":            "1",
		"actions.0.resources.0.db":         db,
		"actions.0.resources.0.collection": "",
		"actions.0.resources.0.cluster":    "false",
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		roleName := rs.Primary.Attributes["role_name"]
		if groupID == "" || roleName == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().CustomDatabaseRolesApi.GetCustomDatabaseRole(context.Background(), groupID, roleName).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("custom DB Role (%s/%s) does not exist", groupID, roleName)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_custom_db_role_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		roleName := rs.Primary.Attributes["role_name"]
		if groupID == "" || roleName == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", rs.Type)
		}
		_, _, err := acc.ConnV2().CustomDatabaseRolesApi.GetCustomDatabaseRole(context.Background(), groupID, roleName).Execute()
		if err == nil {
			return fmt.Errorf("custom DB Role (%s) still exists", roleName)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		roleName := rs.Primary.Attributes["role_name"]
		if groupID == "" || roleName == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", groupID, roleName), nil
	}
}
