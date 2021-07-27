package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasNetworkContainers_basic(t *testing.T) {
	var container matlas.Container

	randInt := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_network_container.test"
	cidrBlock := fmt.Sprintf("10.8.%d.0/24", randInt)
	dataSourceName := "data.mongodbatlas_network_containers.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	providerName := "AWS"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainersDSConfig(projectName, orgID, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.provisioned"),
				),
			},
		},
	})
}
func TestAccDataSourceMongoDBAtlasNetworkContainers_WithGCPRegions(t *testing.T) {
	var container matlas.Container

	randInt := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_network_container.test"
	cidrBlock := fmt.Sprintf("10.%d.0.0/21", randInt)
	dataSourceName := "data.mongodbatlas_network_containers.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	providerName := "GCP"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainersDSWithGCPRegionsConfig(projectName, orgID, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.provisioned"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasNetworkContainersDSConfig(projectName, orgID, cidrBlock string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "${mongodbatlas_project.test.id}"
			atlas_cidr_block = "%s"
			provider_name		 = "AWS"
			region_name			 = "US_EAST_1"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AWS"
		}
	`, projectName, orgID, cidrBlock)
}

func testAccMongoDBAtlasNetworkContainersDSWithGCPRegionsConfig(projectName, orgID, cidrBlock string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block = "%s"
			provider_name		 = "GCP"
			regions = ["US_EAST_4", "US_WEST_3"]
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = mongodbatlas_network_container.test.project_id
			provider_name = "GCP"
		}
	`, projectName, orgID, cidrBlock)
}
