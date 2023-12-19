package projectapikey_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationConfigRSProjectAPIKey_RemovingOptionalRootProjectID(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		description  = fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
		roleName     = "GROUP_OWNER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasProjectAPIKeyConfigBasic(orgID, projectName, description, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttr(resourceName, "project_assignment.#", "1"),
				),
			},
		},
	})
}
