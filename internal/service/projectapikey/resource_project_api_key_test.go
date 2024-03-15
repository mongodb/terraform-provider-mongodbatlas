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

func TestAccConfigRSProjectAPIKey_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		description  = acc.RandomName()
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_BasicWithLegacyRootProjectID(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		description  = acc.RandomName()
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, true),
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

func TestAccConfigRSProjectAPIKey_ChangingSingleProject(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName1 = acc.RandomProjectName()
		projectName2 = acc.RandomProjectName()
		description  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyChangingProject(orgID, projectName1, projectName2, description, "mongodbatlas_project.proj1.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectAPIKeyChangingProject(orgID, projectName1, projectName2, description, "mongodbatlas_project.proj2.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_RemovingOptionalRootProjectID(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		description  = acc.RandomName()
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_Multiple(t *testing.T) {
	var (
		resourceName    = "mongodbatlas_project_api_key.test"
		dataSourceName  = "data.mongodbatlas_project_api_key.test"
		dataSourcesName = "data.mongodbatlas_project_api_keys.test"
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acc.RandomProjectName()
		description     = acc.RandomName()
		roleName        = "GROUP_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigMultiple(orgID, projectName, description, roleName),
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

func TestAccConfigRSProjectAPIKey_UpdateDescription(t *testing.T) {
	var (
		resourceName       = "mongodbatlas_project_api_key.test"
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName()
		description        = acc.RandomName()
		updatedDescription = acc.RandomName()
		roleName           = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, updatedDescription, roleName, false),
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
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		description  = acc.RandomName()
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectAPIKeyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_RecreateWhenDeletedExternally(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_project_api_key.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acc.RandomProjectName()
		descriptionPrefix = "test-acc-delete-api-key-"
		description       = descriptionPrefix + acc.RandomName()
		roleName          = "GROUP_OWNER"
	)

	projectAPIKeyConfig := testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
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

func TestAccConfigRSProjectAPIKey_DeleteProjectAndAssignment(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_project_api_key.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acc.RandomProjectName()
		secondProjectName = acc.RandomProjectName()
		description       = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigDeletedProjectAndAssignment(orgID, projectName, secondProjectName, description, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.1.project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigDeletedProjectAndAssignment(orgID, projectName, secondProjectName, description, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_assignment.0.project_id"),
				),
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

func testAccCheckMongoDBAtlasProjectAPIKeyDestroy(s *terraform.State) error {
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

func TestAccConfigRSProjectAPIKey_Invalid_Role(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		description = projectName
		roleName    = "INVALID_ROLE"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
				ExpectError: regexp.MustCompile("INVALID_ENUM_VALUE"),
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_assignment.0.project_id"]

		return fmt.Sprintf("%s-%s", projectID, rs.Primary.Attributes["api_key_id"]), nil
	}
}

func testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleNames string, includeRootProjID bool) string {
	var rootProjectID string
	if includeRootProjID {
		rootProjectID = "project_id = mongodbatlas_project.test.id"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
		resource "mongodbatlas_project_api_key" "test" {
			%[3]s
			description  = %[4]q
			project_assignment  {
				project_id = mongodbatlas_project.test.id
				role_names = [%[5]q]
			}
		}
	`, orgID, projectName, rootProjectID, description, roleNames)
}

func testAccMongoDBAtlasProjectAPIKeyChangingProject(orgID, projectName1, projectName2, description, assignedProject string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "proj1" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_project" "proj2" {
			org_id = %[1]q
			name   = %[3]q
		}

		resource "mongodbatlas_project_api_key" "test" {
			description  = %[4]q
			project_assignment  {
				project_id = %[5]s
				role_names = ["GROUP_OWNER"]
			}
		}
	`, orgID, projectName1, projectName2, description, assignedProject)
}

func testAccMongoDBAtlasProjectAPIKeyConfigMultiple(orgID, projectName, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[3]q
			project_assignment  {
				project_id = mongodbatlas_project.test.id
				role_names = [%[4]q]
			  }
		}
		data "mongodbatlas_project_api_key" "test" {
			project_id      = mongodbatlas_project.test.id
			api_key_id  = mongodbatlas_project_api_key.test.api_key_id
		}
		
		data "mongodbatlas_project_api_keys" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, orgID, projectName, description, roleNames)
}

func testAccMongoDBAtlasProjectAPIKeyConfigDeletedProjectAndAssignment(orgID, projectName, secondProjectName, description string, includeSecondProject bool) string {
	var secondProject string
	if includeSecondProject {
		secondProject = fmt.Sprintf(`
		resource "mongodbatlas_project" "project2" {
			org_id = %[1]q
			name   = %[2]q
		}`, orgID, secondProjectName)
	}
	var secondProjectAssignment string
	if includeSecondProject {
		secondProjectAssignment = `
		project_assignment  {
			project_id = mongodbatlas_project.project2.id
			role_names = ["GROUP_OWNER"]
		}
		`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "project1" {
			org_id = %[1]q
			name   = %[2]q
		}
		 %[3]s
		resource "mongodbatlas_project_api_key" "test" {
			description  = %[4]q
			project_assignment  {
				project_id = mongodbatlas_project.project1.id
				role_names = ["GROUP_OWNER"]
			}
			%[5]s
		}
	`, orgID, projectName, secondProject, description, secondProjectAssignment)
}
