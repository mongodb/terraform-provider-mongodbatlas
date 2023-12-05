package customdnsconfigurationclusteraws_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSCustomDNSConfigurationAWS_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_custom_dns_configuration_cluster_aws.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSDataSourceConfig(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasCustomDNSConfigurationAWSDataSourceConfig(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id     = mongodbatlas_project.test.id
			enabled       = %[3]t
		}

		data "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id      = mongodbatlas_custom_dns_configuration_cluster_aws.test.id
		}
	`, orgID, projectName, enabled)
}
