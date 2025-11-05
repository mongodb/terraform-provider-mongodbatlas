package projectapikey_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName    = "mongodbatlas_project_api_key.test"
	dataSourceName  = "data.mongodbatlas_project_api_key.test"
	roleName        = "GROUP_OWNER"
	updatedRoleName = "GROUP_READ_ONLY"
)

func TestAccProjectAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config: configBasic(description, projectID, roleName),
				Check:  check(description, projectID, roleName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName, projectID),
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
				Config: configChangingProject(orgID, projectID1, projectName2, description, roleName, true),
				Check:  check(description, projectID1, roleName),
			},
			{
				Config: configChangingProject(orgID, projectID1, projectName2, description, roleName, false),
				Check:  check(description, projectName2, roleName),
			},
			{
				Config: configChangingProject(orgID, projectID1, projectName2, description, roleName+","+updatedRoleName, false),
				Check:  check(description, projectName2, roleName+","+updatedRoleName),
			},
			{
				Config: configChangingProject(orgID, projectID1, projectName2, description, roleName+","+updatedRoleName, true),
				Check:  check(description, projectID1, roleName+","+updatedRoleName),
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
				Check:  check(description, projectID, roleName),
			},
			{
				Config: configBasic(updatedDescription, projectID, roleName),
				Check:  check(updatedDescription, projectID, roleName),
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
				Check:  check(description, projectID, roleName),
			},
			{
				Config: configBasic(description, projectID, updatedRoleName),
				Check:  check(description, projectID, updatedRoleName),
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
				Config:      configDuplicatedProject(description, projectID, roleName, updatedRoleName),
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
				Check:  check(description, projectID, roleName),
			},
			{
				PreConfig: func() {
					if err := deleteAPIKeyManually(orgID, descriptionPrefix); err != nil {
						t.Fatalf("failed to manually delete API key resource: %s", err)
					}
					// Wait longer and verify deletion to ensure API consistency.
					if err := waitForAPIKeyDeletion(orgID, descriptionPrefix, 30*time.Second); err != nil {
						t.Fatalf("failed to verify API key deletion: %s", err)
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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(projectID),
		Steps: []resource.TestStep{
			{
				Config:      configBasic(description, projectID, "INVALID_ROLE"),
				ExpectError: regexp.MustCompile("INVALID_ENUM_VALUE"),
			},
		},
	})
}

func deleteAPIKeyManually(orgID, descriptionPrefix string) error {
	list, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListOrgApiKeys(context.Background(), orgID).Execute()
	if err != nil {
		return err
	}
	for _, key := range list.GetResults() {
		if strings.HasPrefix(key.GetDesc(), descriptionPrefix) {
			if _, err := acc.ConnV2().ProgrammaticAPIKeysApi.DeleteOrgApiKey(context.Background(), orgID, key.GetId()).Execute(); err != nil {
				return err
			}
		}
	}
	return nil
}

func waitForAPIKeyDeletion(orgID, descriptionPrefix string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		list, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListOrgApiKeys(context.Background(), orgID).Execute()
		if err != nil {
			return fmt.Errorf("error listing API keys: %w", err)
		}
		found := false
		for _, key := range list.GetResults() {
			if strings.HasPrefix(key.GetDesc(), descriptionPrefix) {
				found = true
				break
			}
		}
		if !found {
			return nil // API key successfully deleted and confirmed.
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("timeout waiting for API key deletion after %v", timeout)
}

func checkDestroy(projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_project_api_key" {
				continue
			}
			projectAPIKeys, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.ListGroupApiKeys(context.Background(), projectID).Execute()
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
		if found, _, _ := acc.ConnV2().ProgrammaticAPIKeysApi.GetOrgApiKey(context.Background(), orgID, apiKeyID).Execute(); found == nil {
			return fmt.Errorf("API Key (%s) does not exist", apiKeyID)
		}
		return nil
	}
}

func importStateIDFunc(resourceName, projectID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", projectID, rs.Primary.Attributes["api_key_id"]), nil
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
	`, description, projectID, roleNames, configDataSources(fmt.Sprintf("%q", projectID)))
}

func configDuplicatedProject(description, projectID, roleNames1, roleName2 string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[1]q
			project_assignment  {
				project_id = %[2]q
				role_names = [%[3]q]
			}
			project_assignment  {
				project_id = %[2]q
				role_names = [%[4]q]
			}
		}
	`, description, projectID, roleNames1, roleName2)
}

func configChangingProject(orgID, projectID1, projectName2, description, roleNames string, useProject1 bool) string {
	projectIDStr := "mongodbatlas_project.proj2.id"
	if useProject1 {
		projectIDStr = fmt.Sprintf("%q", projectID1)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "proj2" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_project_api_key" "test" {
			description  = %[3]q
			project_assignment  {
				project_id = %[4]s
				role_names = %[5]s
			}
			depends_on = [mongodbatlas_project.proj2]
		} 
		%[6]s
	`, orgID, projectName2, description, projectIDStr, getRoleNamesStr(roleNames), configDataSources(projectIDStr))
}

func getRoleNames(roleNames string) []string {
	var ret []string
	for _, role := range strings.Split(roleNames, ",") {
		ret = append(ret, strings.TrimSpace(role))
	}
	return ret
}

func getRoleNamesStr(roleNames string) string {
	var quoted []string
	for _, role := range strings.Split(roleNames, ",") {
		quoted = append(quoted, fmt.Sprintf("%q", strings.TrimSpace(role)))
	}
	return fmt.Sprintf("[%s]", strings.Join(quoted, ", "))
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

func configDataSources(projectIDStr string) string {
	return fmt.Sprintf(`
			data "mongodbatlas_project_api_key" "test" {
				project_id      = %[1]s
				api_key_id  = mongodbatlas_project_api_key.test.api_key_id
			}
			
			data "mongodbatlas_project_api_keys" "test" {
				project_id      = %[1]s
				depends_on = [mongodbatlas_project_api_key.test]
			}
		`, projectIDStr)
}

func check(description, projectNameOrID, roleNames string) resource.TestCheckFunc {
	roles := getRoleNames(roleNames)
	attrsMap := map[string]string{
		"description":                       description,
		"project_assignment.#":              "1",
		"project_assignment.0.role_names.#": strconv.Itoa(len(roles)),
	}
	attrs := []string{"public_key", "private_key"}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttrWith(resourceName, "project_assignment.0.project_id", acc.IsProjectNameOrID(projectNameOrID)),
		resource.TestCheckResourceAttrWith(dataSourceName, "project_assignment.0.project_id", acc.IsProjectNameOrID(projectNameOrID)),
	}
	for _, role := range roles {
		checks = append(checks,
			resource.TestCheckTypeSetElemAttr(resourceName, "project_assignment.0.role_names.*", role),
			resource.TestCheckTypeSetElemAttr(dataSourceName, "project_assignment.0.role_names.*", role))
	}
	return acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil, attrs, attrsMap, checks...)
}
