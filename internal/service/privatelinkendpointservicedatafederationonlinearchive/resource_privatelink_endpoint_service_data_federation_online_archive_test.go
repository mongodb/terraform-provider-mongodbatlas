package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	projectID                                                         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID                                                        = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive),
					resource.TestCheckResourceAttr(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "comment"),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "type"),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "provider_name"),
				),
			},
			{
				ResourceName:      resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive,
				ImportStateIdFunc: testAccCheckMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveFunc(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchiveFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["endpoint_id"]), nil
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().DataLakes.GetPrivateLinkEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"])
		if err == nil {
			return fmt.Errorf("Private endpoint service data federation online archive still exists")
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private endpoint service data federation online archive not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Private endpoint service data federation online archive ID not set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().DataLakes.GetPrivateLinkEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"])
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
	  comment					= "Terraform Acceptance Test"
	}
	`, projectID, endpointID)
}
