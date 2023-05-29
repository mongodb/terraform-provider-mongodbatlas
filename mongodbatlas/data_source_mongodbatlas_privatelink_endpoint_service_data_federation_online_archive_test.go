package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
)

func TestAccDataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	testCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "comment"),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "type"),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "provider_name"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveConfig(projectID, endpointID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= "Terraform Acceptance Test"
	}

	data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
		project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
		endpoint_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.endpoint_id
	}
	`, projectID, endpointID)
}
