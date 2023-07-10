package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkDSPrivateLinkEndpointServiceADL_basic(t *testing.T) {
	datasourceName := "data.mongodbatlas_privatelink_endpoint_service_adl.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID := "vpce-jjg5e24qp93513h03"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointADLDataSourceConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(datasourceName, "type", "DATA_LAKE"),
					resource.TestCheckResourceAttr(datasourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(datasourceName, "comment", "private link adl comment"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateLinkEndpointADLDataSourceConfig(projectID, endpointID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint_service_adl" "test" {
			project_id   = "%[1]s"
			endpoint_id  = "%[2]s"
			comment      = "private link adl comment"
			type		 = "DATA_LAKE"
			provider_name	 = "AWS"
		}

		data "mongodbatlas_privatelink_endpoint_service_adl" "test" {
			project_id      = mongodbatlas_privatelink_endpoint_service_adl.test.project_id
			endpoint_id = mongodbatlas_privatelink_endpoint_service_adl.test.endpoint_id
		}
	`, projectID, endpointID)
}
