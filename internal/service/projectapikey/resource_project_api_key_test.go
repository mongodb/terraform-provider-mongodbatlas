package projectapikey_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName          = "mongodbatlas_project_api_key.test"
	dataSourceName        = "data.mongodbatlas_project_api_key.test"
	pluralDataSourcesName = "data.mongodbatlas_project_api_keys.test"
	roleName              = "GROUP_OWNER"
	updatedRoleName       = "GROUP_READ_ONLY"
)

func TestAccProjectAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID   = acc.ProjectIDExecution(tb)
		description = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config: configBasic(description, projectID, roleName),
				Check:  checkAggr(description, projectID, roleName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	}
}

func TestAccProjectAPIKey_changingSingleProject(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID1   = acc.ProjectIDExecution(t)
		projectName2 = acc.RandomProjectName()
		description  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID1),
		Steps: []resource.TestStep{
			{
				Config: configChangingProject(orgID, projectName2, description, fmt.Sprintf("%q", projectID1)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				Config: configChangingProject(orgID, projectName2, description, "mongodbatlas_project.proj2.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccProjectAPIKey_updateDescription(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		description        = acc.RandomName()
		updatedDescription = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config: configBasic(description, projectID, roleName),
				Check:  checkAggr(description, projectID, roleName),
			},
			{
				Config: configBasic(updatedDescription, projectID, roleName),
				Check:  checkAggr(updatedDescription, projectID, roleName),
			},
		},
	})
}

func TestAccProjectAPIKey_updateRole(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config: configBasic(description, projectID, roleName),
				Check:  checkAggr(description, projectID, roleName),
			},
			{
				Config: configBasic(description, projectID, updatedRoleName),
				Check:  checkAggr(description, projectID, updatedRoleName),
			},
		},
	})
}

func TestAccProjectAPIKey_duplicateProject(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config:      configTwoAssignments(description, projectID, roleName, projectID, updatedRoleName),
				ExpectError: regexp.MustCompile("duplicated projectID in assignments: " + projectID),
			},
		},
	})
}

func TestAccProjectAPIKey_recreateWhenDeletedExternally(t *testing.T) {
	var (
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID         = acc.ProjectIDExecution(t)
		descriptionPrefix = acc.RandomName()
		description       = descriptionPrefix + "-" + acc.RandomName()
		config            = configBasic(description, projectID, roleName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checkAggr(description, projectID, roleName),
			},
			{
				PreConfig: func() {
					if err := deleteAPIKeyManually(orgID, descriptionPrefix); err != nil {
						t.Fatalf("failed to manually delete API key resource: %s", err)
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // should detect that api key has to be recreated
			},
		},
	})
}

func TestAccProjectAPIKey_deleteProjectAndAssignment(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID1   = acc.ProjectIDExecution(t)
		projectName2 = acc.RandomProjectName()
		description  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID1),
		Steps: []resource.TestStep{
			{
				Config: configDeletedProjectAndAssignment(orgID, projectID1, projectName2, description, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.1.project_id"),
				),
			},
			{
				Config: configDeletedProjectAndAssignment(orgID, projectID1, projectName2, description, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
				),
			},
		},
	})
}

func TestAccProjectAPIKey_invalidRole(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = fmt.Sprintf("desc-%s", projectID)
		roleName    = "INVALID_ROLE"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config:      configBasic(description, projectID, roleName),
				ExpectError: regexp.MustCompile("INVALID_ENUM_VALUE"),
			},
		},
	})
}

func deleteAPIKeyManually(orgID, descriptionPrefix string) error {
	list, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListApiKeys(context.Background(), orgID).Execute()
	if err != nil {
		return err
	}
	for _, key := range list.GetResults() {
		if strings.HasPrefix(key.GetDesc(), descriptionPrefix) {
			if _, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.DeleteApiKey(context.Background(), orgID, key.GetId()).Execute(); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkDestroy(projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_project_api_key" {
				continue
			}
			projectAPIKeys, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListProjectApiKeys(context.Background(), projectID).Execute()
			if err != nil {
				return nil
			}
			ids := conversion.DecodeStateID(rs.Primary.ID)
			for _, val := range projectAPIKeys.GetResults() {
				if val.GetId() == ids["api_key_id"] {
					return fmt.Errorf("Project API Key (%s) still exists", ids["api_key_id"])
				}
			}
		}
		return nil
	}
}

func checkExists(resourceName, projectID string) resource.TestCheckFunc {
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
		list, _, _ := acc.ConnV2().ProgrammaticAPIKeysApi.ListProjectApiKeys(context.Background(), projectID).Execute()
		for _, val := range list.GetResults() {
			if val.GetId() == apiKeyID {
				return nil
			}
		}
		return fmt.Errorf("API Key (%s) does not exist", apiKeyID)
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["api_key_id"], nil
	}
}

func configBasic(description, projectID, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[1]q
			project_assignment  {
				project_id = %[2]q
				role_names = [%[3]q]
			}
		}
		%[4]s
	`, description, projectID, roleNames, configDataSources(projectID))
}

func configTwoAssignments(description, projectID1, roleNames1, projectID2, roleName2 string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[1]q
			project_assignment  {
				project_id = %[2]q
				role_names = [%[3]q]
			}
			project_assignment  {
				project_id = %[4]q
				role_names = [%[5]q]
			}
		}
	`, description, projectID1, roleNames1, projectID2, roleName2)
}

func configChangingProject(orgID, projectName2, description, assignedProject string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "proj2" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_project_api_key" "test" {
			description  = %[3]q
			project_assignment  {
				project_id = %[4]s
				role_names = ["GROUP_OWNER"]
			}
		}
	`, orgID, projectName2, description, assignedProject)
}

func configDeletedProjectAndAssignment(orgID, projectID1, projectName2, description string, includeSecondProject bool) string {
	var secondProject, secondProjectAssignment string
	if includeSecondProject {
		secondProject = fmt.Sprintf(`
		resource "mongodbatlas_project" "project2" {
			org_id = %[1]q
			name   = %[2]q
		}
		`, orgID, projectName2)
		secondProjectAssignment = `
			project_assignment  {
				project_id = mongodbatlas_project.project2.id
				role_names = ["GROUP_OWNER"]
			}
		`
	}
	return fmt.Sprintf(`
		 %[3]s
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[2]q
			project_assignment  {
				project_id = %[1]q
				role_names = ["GROUP_OWNER"]
			}
			%[4]s
		}
	`, projectID1, description, secondProject, secondProjectAssignment)
}

func configDataSources(projectID string) string {
	return fmt.Sprintf(`
			data "mongodbatlas_project_api_key" "test" {
				project_id      = %[1]q
				api_key_id  = mongodbatlas_project_api_key.test.api_key_id
			}
			
			data "mongodbatlas_project_api_keys" "test" {
				project_id      = %[1]q
				depends_on = [mongodbatlas_project_api_key.test]
			}
		`, projectID)
}

func checkAggr(description, projectID, roleName string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attributes := map[string]string{
		"description":                       description,
		"project_assignment.#":              "1",
		"project_assignment.0.project_id":   projectID,
		"project_assignment.0.role_names.#": "1",
		"project_assignment.0.role_names.0": roleName,
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName, projectID),
		resource.TestCheckResourceAttrSet(pluralDataSourcesName, "results.0.project_assignment.0.project_id"),
		resource.TestCheckResourceAttrSet(pluralDataSourcesName, "results.0.project_assignment.0.role_names.0"),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrChecks(dataSourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "public_key", "private_key")
	checks = acc.AddAttrSetChecks(dataSourceName, checks, "public_key", "private_key")
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
