package organization_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSOrganization_Basic(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-organization-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccConfigRSOrganization_importBasic(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-import-organization-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasOrganizationImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasOrganizationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		organizationOptions := &matlas.OrganizationsListOptions{}
		orgs, _, err := conn.Organizations.List(context.Background(), organizationOptions)
		if err == nil {
			for _, val := range orgs.Results {
				if val.ID == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["role_name"])
				}
			}
			return nil
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasOrganizationDestroy(s *terraform.State) error {
	conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		organizationOptions := &matlas.OrganizationsListOptions{}
		orgs, _, err := conn.Organizations.List(context.Background(), organizationOptions)
		if err == nil {
			for _, val := range orgs.Results {
				if val.ID == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["role_name"])
				}
			}
			return nil
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasOrganizationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.Attributes["org_id"], nil
	}
}

func testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleNames string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = "%s"
		name = "%s"
		description = "%s"
		role_names = ["%s"]
	  }
	`, orgOwnerID, name, description, roleNames)
}
