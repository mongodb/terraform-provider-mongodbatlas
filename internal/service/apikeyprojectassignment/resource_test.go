package apikeyprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName    = "mongodbatlas_api_key_project_assignment.test"
	roleName        = "GROUP_OWNER"
	roleNameUpdated = "GROUP_READ_ONLY"
)

func TestAccApiKeyProjectAssignmentRS_basic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName1 = acc.RandomProjectName()
		projectName2 = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{

		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		// CheckDestroy:             checkDestroyApiKeyProjectAssignment,
		Steps: []resource.TestStep{
			{
				Config: apiKeyProjectAssignmentConfig(orgID, roleName, projectName1, projectName2),
				Check:  apiKeyProjectAssignmentAttributeChecks(projectName1, roleName),
			},
			{
				Config: apiKeyProjectAssignmentConfig(orgID, roleNameUpdated, projectName1, projectName2),
				Check:  apiKeyProjectAssignmentAttributeChecks(projectName1, roleNameUpdated),
			},
			{
				Config:            apiKeyProjectAssignmentConfig(orgID, roleNameUpdated, projectName1, projectName2),
				ResourceName:      resourceName,
				ImportStateIdFunc: checkApiKeyProjectAssignmentImportStateIDFunc(resourceName, "project_id", "api_key_id"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkApiKeyProjectAssignmentImportStateIDFunc(resourceName, attrNameProjectID, attrNameApiKeyID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return IDWithProjectIDApiKeyID(rs.Primary.Attributes[attrNameProjectID], rs.Primary.Attributes[attrNameApiKeyID])
	}
}

func IDWithProjectIDApiKeyID(projectID, apiKeyID string) (string, error) {
	if err := conversion.ValidateProjectID(projectID); err != nil {
		return "", err
	}
	return projectID + "-" + apiKeyID, nil
}

func apiKeyProjectAssignmentAttributeChecks(projectNameOrID, roleNames string) resource.TestCheckFunc {
	roles := getRoleNames(roleNames)
	attrsMap := map[string]string{
		"role_names.#": strconv.Itoa(len(roles)),
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttrWith(resourceName, "project_id", acc.IsProjectNameOrID(projectNameOrID)),
	}
	for _, role := range roles {
		checks = append(checks,
			resource.TestCheckTypeSetElemAttr(resourceName, "role_names.*", role))
	}
	return acc.CheckRSAndDS(resourceName, nil, nil, []string{}, attrsMap, checks...)
}

func apiKeyProjectAssignmentConfig(orgID, roleName, projectName1, projectName2 string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
			org_id     = "%[1]s"
			description  = "Test API Key"

			role_names = ["ORG_READ_ONLY"]
		}

		data "mongodbatlas_roles_org_id" "test" {}

		resource "mongodbatlas_project" "test1" {
			name   = "%[3]s"
			org_id = data.mongodbatlas_roles_org_id.test.org_id
        }

		resource "mongodbatlas_project" "test2" {
			name   = "%[4]s"
			org_id = data.mongodbatlas_roles_org_id.test.org_id
        }

		resource "mongodbatlas_api_key_project_assignment" "test1" {
			project_id = mongodbatlas_project.test1.id
			api_key_id = mongodbatlas_api_key.test.id

			role_names  = ["GROUP_OWNER"]
		}
		
		resource "mongodbatlas_api_key_project_assignment" "test2" {
			project_id = mongodbatlas_project.test2.id
			api_key_id = mongodbatlas_api_key.test.id

			role_names  = ["GROUP_OWNER"]
		}

		resource "mongodbatlas_api_key_project_assignment" "test" {
		    depends_on [mongodbatlas_api_key_project_assignment.test1, mongodbatlas_api_key_project_assignment.test2]
			project_id = mongodbatlas_project.test2.id
			api_key_id = mongodbatlas_api_key.test.id

			role_names  = ["%[2]s"]
		}
	`, orgID, roleName, projectName1, projectName2)
}

func getRoleNames(roleNames string) []string {
	var ret []string
	for _, role := range strings.Split(roleNames, ",") {
		ret = append(ret, strings.TrimSpace(role))
	}
	return ret
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		apiKeyID := ids["api_key_id"]
		orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
		if found, _, _ := acc.ConnV2().ProgrammaticAPIKeysApi.GetApiKey(context.Background(), orgID, apiKeyID).Execute(); found == nil {
			return fmt.Errorf("API Key (%s) does not exist", apiKeyID)
		}
		return nil
	}
}
