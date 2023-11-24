package mongodbatlas_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSProjectAPIKey_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		description  = fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
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
		projectName     = acctest.RandomWithPrefix("test-acc")
		description     = fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
		roleName        = "GROUP_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigMultiple(orgID, projectName, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
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
		projectName        = acctest.RandomWithPrefix("test-acc")
		description        = fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
		updatedDescription = fmt.Sprintf("test-acc-project-api_key-updated-%s", acctest.RandString(5))
		roleName           = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, updatedDescription, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
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
		projectName  = acctest.RandomWithPrefix("test-acc")
		description  = fmt.Sprintf("test-acc-import-project-api_key-%s", acctest.RandString(5))
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName),
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
		projectName       = acctest.RandomWithPrefix("test-acc")
		descriptionPrefix = "test-acc-project-to-delete-api-key"
		description       = fmt.Sprintf("%s-%s", descriptionPrefix, acctest.RandString(5))
		roleName          = "GROUP_OWNER"
	)

	projectAPIKeyConfig := testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectAPIKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
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
		projectName       = acctest.RandomWithPrefix("test-acc")
		secondProjectName = acctest.RandomWithPrefix("test-acc")
		description       = fmt.Sprintf("%s-%s", "test-acc-project", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
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
	conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas
	list, _, err := conn.APIKeys.List(context.Background(), orgID, &matlas.ListOptions{})
	if err != nil {
		return err
	}
	for _, key := range list {
		if strings.HasPrefix(key.Desc, descriptionPrefix) {
			if _, err := conn.APIKeys.Delete(context.Background(), orgID, key.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasProjectAPIKeyDestroy(s *terraform.State) error {
	conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_api_key" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		projectAPIKeys, _, err := conn.ProjectAPIKeys.List(context.Background(), ids["project_id"], nil)
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

func testAccCheckMongoDBAtlasProjectAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["api_key_id"]), nil
	}
}

func testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_api_key" "test" {
			project_id     = mongodbatlas_project.test.id
			description  = %[3]q
			project_assignment  {
				project_id = mongodbatlas_project.test.id
				role_names = [%[4]q]
			}
		}
	`, orgID, projectName, description, roleNames)
}

func testAccMongoDBAtlasProjectAPIKeyConfigMultiple(orgID, projectName, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_api_key" "test" {
			project_id     = mongodbatlas_project.test.id
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
			project_id     = mongodbatlas_project.project1.id
			description  = %[4]q
			project_assignment  {
				project_id = mongodbatlas_project.project1.id
				role_names = ["GROUP_OWNER"]
			}
			%[5]s
		}
	`, orgID, projectName, secondProject, description, secondProjectAssignment)
}
