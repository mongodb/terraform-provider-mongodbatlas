package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServerlessPrivateLinkEndpoint_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint_serverless.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName = "serverlessplink"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServerlessConfig(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
				),
			},
		},
	})
}

func TestAccServerlessPrivateLinkEndpoint_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint_serverless.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName = "serverlessimport"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServerlessConfig(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
				),
			},
			{
				Config:                  testAccMongoDBAtlasPrivateLinkEndpointServerlessConfig(orgID, projectName, instanceName, false),
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_strings_private_endpoint_srv"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_serverless" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		privateLink, _, err := conn.ServerlessPrivateEndpoints.Get(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"])
		if err == nil && privateLink != nil {
			return fmt.Errorf("endpoint_id (%s) still exists", ids["endpoint_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasPrivateLinkEndpointServerlessConfig(orgID, projectName, instanceName string, ignoreConnectionStrings bool) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id    = mongodbatlas_serverless_instance.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	}

	%s
	`, testAccMongoDBAtlasServerlessInstanceConfig(orgID, projectName, instanceName, ignoreConnectionStrings))
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ServerlessPrivateEndpoints.Get(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("endpoint_id (%s) does not exist", ids["endpoint_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServerlessImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["instance_name"], ids["endpoint_id"]), nil
	}
}
