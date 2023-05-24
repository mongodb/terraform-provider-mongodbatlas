package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives.test"
)

func TestAccDataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchives_basic(t *testing.T) {
	testCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchivesConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchiveExists(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives, "results.#"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasPrivateEndpointServiceDataFederationOnlineArchivesConfig(projectID, endpointID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  type						= "DATA_LAKE"
	  comment					= "Terraform Acceptance Test"
	}

	data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "test" {
		project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
	}
	`, projectID, endpointID)
}
