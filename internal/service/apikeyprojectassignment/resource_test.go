package apikeyprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName    = "mongodbatlas_api_key_project_assignment.test"
	pluralDSName    = "data.mongodbatlas_api_key_project_assignments.plural"
	singularDSName  = "data.mongodbatlas_api_key_project_assignment.singular"
	roleName        = "GROUP_OWNER"
	roleNameUpdated = "GROUP_READ_ONLY"
)

func TestAccApiKeyProjectAssignmentRS_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{

		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: apiKeyProjectAssignmentConfig(orgID, roleName, projectName),
				Check:  apiKeyProjectAssignmentAttributeChecks(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("roles"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(roleName),
						}),
					),
				},
			},
			{
				Config: apiKeyProjectAssignmentConfig(orgID, roleNameUpdated, projectName),
				Check:  apiKeyProjectAssignmentAttributeChecks(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("roles"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(roleNameUpdated),
						}),
					),
				},
			},
			{
				Config:                               apiKeyProjectAssignmentConfig(orgID, roleNameUpdated, projectName),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName, "project_id", "api_key_id"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
}

func importStateIDFunc(resourceName, attrNameProjectID, attrNameAPIKeyID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes[attrNameProjectID], rs.Primary.Attributes[attrNameAPIKeyID]), nil
	}
}

func apiKeyProjectAssignmentAttributeChecks() resource.TestCheckFunc {
	attrsMap := map[string]string{
		"roles.#": "1",
	}
	attrsSet := []string{"project_id", "api_key_id"}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	return acc.CheckRSAndDS(resourceName, conversion.Pointer(singularDSName), conversion.Pointer(pluralDSName), attrsSet, attrsMap, checks...)
}

func apiKeyProjectAssignmentConfig(orgID, roleName, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
			org_id     = %[1]q
			description  = "Test API Key"
			role_names = ["ORG_READ_ONLY"]
		}

		resource "mongodbatlas_project" "test" {
			name   = %[3]q
			org_id = %[1]q
        }

		resource "mongodbatlas_api_key_project_assignment" "test" {
			project_id = mongodbatlas_project.test.id
			api_key_id = mongodbatlas_api_key.test.api_key_id
			roles      = [%[2]q]
		}

		data "mongodbatlas_api_key_project_assignments" "plural" {
			project_id = mongodbatlas_project.test.id
			depends_on = [mongodbatlas_api_key_project_assignment.test]
		}

		data "mongodbatlas_api_key_project_assignment" "singular" {
			project_id = mongodbatlas_project.test.id
			api_key_id = mongodbatlas_api_key.test.api_key_id
			depends_on = [mongodbatlas_api_key_project_assignment.test]
		}
	`, orgID, roleName, projectName)
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
		projectID := rs.Primary.Attributes["project_id"]
		apiKeyID := rs.Primary.Attributes["api_key_id"]
		apiKeys, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListProjectApiKeys(context.Background(), projectID).Execute()
		if err != nil {
			return fmt.Errorf("error fetching API Keys: %w", err)
		}
		for _, apiKey := range apiKeys.GetResults() {
			if apiKey.GetId() == apiKeyID {
				return nil
			}
		}
		return fmt.Errorf("API Key (%s) does not exist in project (%s)", apiKeyID, projectID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_api_key_project_assignment" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		apiKeyID := rs.Primary.Attributes["api_key_id"]
		apiKeys, apiResp, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListProjectApiKeys(context.Background(), projectID).Execute()
		if err != nil {
			if validate.StatusNotFound(apiResp) {
				return nil
			}
			return fmt.Errorf("error fetching Project API Keys: %w", err)
		}
		for _, apiKey := range apiKeys.GetResults() {
			if apiKey.GetId() == apiKeyID {
				return fmt.Errorf("Project API Key (%s) still exists in project (%s)", apiKeyID, projectID)
			}
		}
	}
	return nil
}
