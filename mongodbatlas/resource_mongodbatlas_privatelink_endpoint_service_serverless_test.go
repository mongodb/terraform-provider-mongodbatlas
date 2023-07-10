package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerlessPrivateLinkEndpointService_basic(t *testing.T) {
	var (
		resourceName            = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceName          = "data.mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceEndpointsName = "data.mongodbatlas_privatelink_endpoints_service_serverless.test"
		orgID                   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName             = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName            = acctest.RandomWithPrefix("test-acc-serverless")
		commentOrigin           = "this is a comment for serverless private link endpoint"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceServerlessConfig(orgID, projectName, instanceName, commentOrigin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttr(datasourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "instance_name"),
				),
			},
		},
	})
}

func TestAccServerlessPrivateLinkEndpointService_importBasic(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName   = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName  = acctest.RandomWithPrefix("test-acc-serverless")
		commentOrigin = "this is a comment for serverless private link endpoint"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceServerlessConfig(orgID, projectName, instanceName, commentOrigin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
				),
			},
			{
				Config:            testAccMongoDBAtlasPrivateLinkEndpointServiceServerlessConfig(orgID, projectName, instanceName, commentOrigin),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_serverless" {
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

func testAccMongoDBAtlasPrivateLinkEndpointServiceServerlessConfig(orgID, projectName, instanceName, comment string) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id   = mongodbatlas_project.test.id
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	}


	resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_serverless.test.project_id
		instance_name = mongodbatlas_privatelink_endpoint_serverless.test.instance_name
		endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
		provider_name = "AWS"
		comment = %[4]q
	}

	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q
		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true

		lifecycle {
			ignore_changes = [connection_strings_private_endpoint_srv]
		}
	}

	data "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		name         = mongodbatlas_serverless_instance.test.name
	}

	data "mongodbatlas_privatelink_endpoints_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
	  }

	data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		endpoint_id = mongodbatlas_privatelink_endpoint_service_serverless.test.endpoint_id
	}

	`, orgID, projectName, instanceName, comment)
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceServerlessImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["instance_name"], ids["endpoint_id"]), nil
	}
}
