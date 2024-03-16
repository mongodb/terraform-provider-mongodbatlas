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
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	resourceName = "mongodbatlas_project_api_key.test"
	roleName     = "GROUP_OWNER"
)

func TestAccConfigRSProjectAPIKey_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_basicWithLegacyRootProjectID(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, description, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_changingSingleProject(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID1   = acc.ProjectIDExecution(t)
		projectName2 = acc.RandomProjectName()
		description  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configChangingProject(orgID, projectName2, description, fmt.Sprintf("%q", projectID1)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				Config: configChangingProject(orgID, projectName2, description, "mongodbatlas_project.proj2.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_removingOptionalRootProjectID(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, description, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				Config: configBasic(projectID, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_multiple(t *testing.T) {
	var (
		dataSourceName  = "data.mongodbatlas_project_api_key.test"
		dataSourcesName = "data.mongodbatlas_project_api_keys.test"
		projectID       = acc.ProjectIDExecution(t)
		description     = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configMultiple(projectID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.role_names.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_assignment.0.role_names.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "results.0.project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourcesName, "results.0.project_assignment.0.role_names.0"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_updateDescription(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		description        = acc.RandomName()
		updatedDescription = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: configBasic(projectID, updatedDescription, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_importBasic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, description, roleName, false),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_recreateWhenDeletedExternally(t *testing.T) {
	var (
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID         = acc.ProjectIDExecution(t)
		descriptionPrefix = acc.RandomName()
		description       = descriptionPrefix + "-" + acc.RandomName()
	)

	projectAPIKeyConfig := configBasic(projectID, description, roleName, false)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectAPIKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				PreConfig: func() {
					if err := deleteAPIKeyManually(orgID, descriptionPrefix); err != nil {
						t.Fatalf("failed to manually delete API key resource: %s", err)
					}
				},
				Config:             projectAPIKeyConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // should detect that api key has to be recreated
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_deleteProjectAndAssignment(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID1   = acc.ProjectIDExecution(t)
		projectName2 = acc.RandomProjectName()
		description  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDeletedProjectAndAssignment(orgID, projectID1, projectName2, description, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.1.project_id"),
				),
			},
			{
				Config: configDeletedProjectAndAssignment(orgID, projectID1, projectName2, description, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_invalidRole(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		description = fmt.Sprintf("desc-%s", projectID)
		roleName    = "INVALID_ROLE"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, description, roleName, false),
				ExpectError: regexp.MustCompile("INVALID_ENUM_VALUE"),
			},
		},
	})
}

func deleteAPIKeyManually(orgID, descriptionPrefix string) error {
	list, _, err := acc.Conn().APIKeys.List(context.Background(), orgID, &matlas.ListOptions{})
	if err != nil {
		return err
	}
	for _, key := range list {
		if strings.HasPrefix(key.Desc, descriptionPrefix) {
			if _, err := acc.Conn().APIKeys.Delete(context.Background(), orgID, key.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_api_key" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectAPIKeys, _, err := acc.Conn().ProjectAPIKeys.List(context.Background(), ids["project_id"], nil)
		if err != nil {
			return nil
		}
		for _, val := range projectAPIKeys {
			if val.ID == ids["api_key_id"] {
				return fmt.Errorf("Project API Key (%s) still exists", ids["role_name"])
			}
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

		projectID := rs.Primary.Attributes["project_assignment.0.project_id"]

		return fmt.Sprintf("%s-%s", projectID, rs.Primary.Attributes["api_key_id"]), nil
	}
}

func configBasic(projectID, description, roleNames string, includeRootProjID bool) string {
	var rootProjectID string
	if includeRootProjID {
		rootProjectID = fmt.Sprintf("project_id = %q", projectID)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[2]q
			project_assignment  {
				project_id = %[1]q
				role_names = [%[3]q]
			}
			%[4]s
		}
	`, projectID, description, roleNames, rootProjectID)
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

func configMultiple(projectID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[2]q
			project_assignment  {
				project_id = %[1]q
				role_names = [%[3]q]
			  }
		}
		data "mongodbatlas_project_api_key" "test" {
			project_id      = %[1]q
			api_key_id  = mongodbatlas_project_api_key.test.api_key_id
		}
		
		data "mongodbatlas_project_api_keys" "test" {
			project_id = %[1]q
		}
	`, projectID, description, roleNames)
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
