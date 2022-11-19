package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkDSPrivateLinkEndpointsServiceADL_basic(t *testing.T) {
	datasourceName := "data.mongodbatlas_privatelink_endpoints_service_adl.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID := "vpce-jjg5e24qp93513h03"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointsADLDataSourceConfig(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.endpoint_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.type"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.comment"),
					resource.TestCheckResourceAttr(datasourceName, "total_count", "1"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateLinkEndpointsADLDataSourceConfig(projectID, endpointID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint_service_adl" "test" {
			project_id   = "%[1]s"
			endpoint_id  = "%[2]s"
			comment      = "private link adl comment"
			type		 = "DATA_LAKE"
			provider_name	 = "AWS"
		}

		data "mongodbatlas_privatelink_endpoints_service_adl" "test" {
			project_id      = mongodbatlas_privatelink_endpoint_service_adl.test.project_id
		}
	`, projectID, endpointID)
}
