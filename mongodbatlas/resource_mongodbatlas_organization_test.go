package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasOrganization_basic(t *testing.T) {
	var organization matlas.Organization

	resourceName := "mongodbatlas_organization.test"
	orgName := fmt.Sprintf("testacc-organization-%s", acctest.RandString(5))
	orgNameUpdated := fmt.Sprintf("testacc-organization-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfig(orgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName, &organization),
					testAccCheckMongoDBAtlasOrganizationAttributes(&organization, orgName),
					resource.TestCheckResourceAttr(resourceName, "name", orgName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			{
				Config: testAccMongoDBAtlasOrganizationConfig(orgNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName, &organization),
					testAccCheckMongoDBAtlasOrganizationAttributes(&organization, orgNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "name", orgNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasOrganization_importBasic(t *testing.T) {

	orgName := fmt.Sprintf("test-acc-%s", acctest.RandString(5))
	resourceName := "mongodbatlas_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfig(orgName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasOrganizationExists(resourceName string, organization *matlas.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] orgID: %s", rs.Primary.Attributes["id"])

		if projectResp, _, err := conn.Organizations.GetOneOrganization(context.Background(), rs.Primary.Attributes["name"]); err == nil {
			*organization = *projectResp
			return nil
		}

		return fmt.Errorf("organization (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasOrganizationAttributes(organization *matlas.Organization, orgName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if organization.Name != orgName {
			return fmt.Errorf("bad organization name: %s", organization.Name)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasOrganizationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		_, err := conn.Organizations.Delete(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("organization (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccMongoDBAtlasOrganizationConfig(orgName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_organization" "test" {
			name = "%s"
		}
	`, orgName)
}
