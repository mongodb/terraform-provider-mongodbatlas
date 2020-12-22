package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCustomDNSConfigurationAWS_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_custom_dns_configuration_cluster_aws.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSDataSourceConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasCustomDNSConfigurationAWSDataSourceConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id     = "%s"
			enabled       = %t
		}

		data "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id      = mongodbatlas_custom_dns_configuration_cluster_aws.test.id
		}
	`, projectID, enabled)
}
