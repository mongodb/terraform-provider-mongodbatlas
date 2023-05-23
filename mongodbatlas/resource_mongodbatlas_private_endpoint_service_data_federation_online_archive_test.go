package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	resourceName = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID   = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
)

func TestAccMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	testCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists("mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := client.DataLakes.GetPrivateLinkEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"])

		if err == nil {
			return fmt.Errorf("Private endpoint service data federation online archive still exists")
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Private endpoint service data federation online archive not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Private endpoint service data federation online archive ID not set")
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := client.DataLakes.GetPrivateLinkEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"])

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveConfig(projectID, endpointID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  type						= "DATA_LAKE"
	  comment					= "Terraform Acceptance Test"
	}
	`, projectID, endpointID)
}
