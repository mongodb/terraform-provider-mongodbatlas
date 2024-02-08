package networkcontainer_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccNetworkContainerRS_basicAWS(t *testing.T) {
	var (
		randIntUpdated   = acctest.RandIntRange(0, 255)
		cidrBlockUpdated = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWS(projectID, cidrBlock, providerNameAws, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: configAWS(projectID, cidrBlockUpdated, providerNameAws, "US_WEST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerRS_basicAzure(t *testing.T) {
	var (
		randIntUpdated   = acctest.RandIntRange(0, 255)
		cidrBlockUpdated = fmt.Sprintf("192.168.%d.0/24", randIntUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAzure(projectID, cidrBlock, providerNameAzure, "US_EAST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: configAzure(projectID, cidrBlockUpdated, providerNameAzure, "US_EAST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerRS_basicGCP(t *testing.T) {
	var (
		randIntUpdated   = acctest.RandIntRange(0, 255)
		cidrBlockUpdated = fmt.Sprintf("10.%d.0.0/18", randIntUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGCP(projectID, gcpCidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: configGCP(projectID, cidrBlockUpdated, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkContainerRS_WithRegionsGCP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGCPWithRegions(projectID, cidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
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

func TestAccNetworkContainerRS_importBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWS(projectID, cidrBlock, providerNameAws, "US_EAST_1"),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]), nil
	}
}

func checkExists(resourceName string, container *matlas.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])
		if containerResp, _, err := acc.Conn().Containers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]); err == nil {
			*container = *containerResp
			return nil
		}
		return fmt.Errorf("container(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_network_container" {
			continue
		}

		_, _, err := acc.Conn().Containers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])

		if err == nil {
			return fmt.Errorf("container (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
		}
	}
	return nil
}

func configAWS(projectID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
			region_name			 = "%s"
		}
	`, projectID, cidrBlock, providerName, region)
}

func configAzure(projectID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
			region			     = "%s"
		}
	`, projectID, cidrBlock, providerName, region)
}

func configGCP(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
		}
	`, projectID, cidrBlock, providerName)
}

func configGCPWithRegions(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block     = "%s"
			provider_name		 = "%s"
			regions = ["US_EAST_4", "US_WEST_3"]
		}

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = mongodbatlas_network_container.test.project_id
			provider_name = "GCP"
		}
	`, projectID, cidrBlock, providerName)
}
