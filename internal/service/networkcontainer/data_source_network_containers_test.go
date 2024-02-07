package networkcontainer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

var (
	container                matlas.Container
	randInt                  = acctest.RandIntRange(0, 255)
	resourceName             = "mongodbatlas_network_container.test"
	dataSourceContainersName = "data.mongodbatlas_network_containers.test"
	cidrBlock                = fmt.Sprintf("10.8.%d.0/24", randInt)
	gcpCidrBlock             = fmt.Sprintf("10.%d.0.0/18", randInt)
	providerNameAws          = "AWS"
	providerNameAzure        = "AZURE"
	providerNameGCP          = "GCP"
	projectID                = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
)

func TestAccNetworkContainerDSPlural_basicAWS(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourcePluralConfigBasicAWS(projectID, cidrBlock, providerNameAws, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerDSPlural_basicAzure(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourcePluralConfigBasicAzure(projectID, cidrBlock, providerNameAzure, "US_EAST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerDSPlural_basicGCP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicGCP(projectID, gcpCidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provisioned"),
				),
			},
		},
	})
}

func dataSourcePluralConfigBasicAWS(projectID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block 	 = "%s"
			provider_name		 = "%s"
			region_name			 = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AWS"
		}

	`, projectID, cidrBlock, providerName, region)
}

func dataSourcePluralConfigBasicAzure(projectID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
			region			     = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AZURE"
		}
	`, projectID, cidrBlock, providerName, region)
}

func dataSourceConfigBasicGCP(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "GCP"
		}
	`, projectID, cidrBlock, providerName)
}
