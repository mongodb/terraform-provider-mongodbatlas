package networkcontainer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourceName = "data.mongodbatlas_network_container.test"
)

func TestAccNetworkContainerDS_basicAWS(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicAWS(projectName, orgID, cidrBlock, providerNameAws, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerDS_basicAzure(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicAzure(projectName, orgID, cidrBlock, providerNameAzure, "US_EAST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerDS_basicGCP(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicGCP(projectName, orgID, gcpCidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),
				),
			},
		},
	})
}

func dataSourceConfigBasicAWS(projectName, orgID, cidrBlock, providerName, region string) string {
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

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}

	`, projectName, orgID, cidrBlock, providerName, region)
}

func dataSourceConfigBasicAzure(projectName, orgID, cidrBlock, providerName, region string) string {
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

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}
	`, projectName, orgID, cidrBlock, providerName, region)
}

func dataSourceConfigBasicGCP(projectName, orgID, cidrBlock, providerName string) string {
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

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}
	`, projectName, orgID, cidrBlock, providerName)
}
