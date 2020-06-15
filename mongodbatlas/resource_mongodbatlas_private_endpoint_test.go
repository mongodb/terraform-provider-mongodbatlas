package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasPrivateEndpoint_basic(t *testing.T) {
	resourceName := "mongodbatlas_private_endpoint.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointConfigBasic(projectID, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasPrivateEndpoint_import(t *testing.T) {
	resourceName := "mongodbatlas_private_endpoint.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointConfigBasic(projectID, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasPrivateEndpointImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"provider_name", "region"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateEndpointImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["private_link_id"]), nil
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := decodeStateID(rs.Primary.ID)

		if _, _, err := conn.PrivateEndpoints.Get(context.Background(), ids["project_id"], ids["private_link_id"]); err == nil {
			return nil
		}

		return fmt.Errorf("MongoDB Private Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["private_link_id"], rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_private_endpoint" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.PrivateEndpoints.Get(context.Background(), ids["project_id"], ids["private_link_id"])
		if err == nil {
			return fmt.Errorf("MongoDB Private Endpoint(%s) still exists", ids["private_link_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasPrivateEndpointConfigBasic(projectID, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}
	`, projectID, providerName, region)
}
