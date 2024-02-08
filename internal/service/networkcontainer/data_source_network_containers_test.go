package networkcontainer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkContainerDSPlural_basicAWS(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourcePluralConfigBasicAWS(projectName, orgID, cidrBlock, providerNameAws, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceContainersName, &container),
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
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourcePluralConfigBasicAzure(projectName, orgID, cidrBlock, providerNameAzure, "US_EAST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceContainersName, &container),
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
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: dataSourcePluralConfigBasicGCP(projectName, orgID, gcpCidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceContainersName, &container),
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

func dataSourcePluralConfigBasicAWS(projectName, orgID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block 	 = "%s"
			provider_name		 = "%s"
			region_name			 = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AWS"
		}

	`, projectName, orgID, cidrBlock, providerName, region)
}

func dataSourcePluralConfigBasicAzure(projectName, orgID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
			region			     = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AZURE"
		}
	`, projectName, orgID, cidrBlock, providerName, region)
}

func dataSourcePluralConfigBasicGCP(projectName, orgID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}
		
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "GCP"
		}
	`, projectName, orgID, cidrBlock, providerName)
}
