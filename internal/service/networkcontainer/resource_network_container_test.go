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

func TestAccNetworkRSNetworkContainer_basicAWS(t *testing.T) {
	var (
		container                matlas.Container
		randInt                  = acctest.RandIntRange(0, 255)
		randIntUpdated           = acctest.RandIntRange(0, 255)
		resourceName             = "mongodbatlas_network_container.test"
		dataSourceName           = "data.mongodbatlas_network_container.test"
		dataSourceContainersName = "data.mongodbatlas_network_containers.test"
		cidrBlock                = fmt.Sprintf("10.8.%d.0/24", randInt)
		cidrBlockUpdated         = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
		providerName             = "AWS"
		orgID                    = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName              = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWS(projectName, orgID, cidrBlock, providerName, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.atlas_cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceContainersName, "results.0.provisioned"),
				),
			},
			{
				Config: configAWS(projectName, orgID, cidrBlockUpdated, providerName, "US_WEST_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkRSNetworkContainer_basicAzure(t *testing.T) {
	var (
		container        matlas.Container
		randInt          = acctest.RandIntRange(0, 255)
		randIntUpdated   = acctest.RandIntRange(0, 255)
		resourceName     = "mongodbatlas_network_container.test"
		cidrBlock        = fmt.Sprintf("192.168.%d.0/24", randInt)
		cidrBlockUpdated = fmt.Sprintf("192.168.%d.0/24", randIntUpdated)
		providerName     = "AZURE"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAzure(projectName, orgID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: configAzure(projectName, orgID, cidrBlockUpdated, providerName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkRSNetworkContainer_basicGCP(t *testing.T) {
	var (
		container        matlas.Container
		randInt          = acctest.RandIntRange(0, 255)
		randIntUpdated   = acctest.RandIntRange(0, 255)
		resourceName     = "mongodbatlas_network_container.test"
		cidrBlock        = fmt.Sprintf("10.%d.0.0/18", randInt)
		cidrBlockUpdated = fmt.Sprintf("10.%d.0.0/18", randIntUpdated)
		providerName     = "GCP"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGCP(projectName, orgID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: configGCP(projectName, orgID, cidrBlockUpdated, providerName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccNetworkRSNetworkContainer_WithRegionsGCP(t *testing.T) {
	var (
		container                matlas.Container
		randInt                  = acctest.RandIntRange(0, 255)
		resourceName             = "mongodbatlas_network_container.test"
		dataSourceName           = "data.mongodbatlas_network_container.test"
		dataSourceContainersName = "data.mongodbatlas_network_containers.test"
		cidrBlock                = fmt.Sprintf("10.%d.0.0/21", randInt)
		providerName             = "GCP"
		orgID                    = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName              = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGCPWithRegions(projectName, orgID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
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

func TestAccNetworkRSNetworkContainer_importBasic(t *testing.T) {
	var (
		randInt      = acctest.RandIntRange(0, 255)
		resourceName = "mongodbatlas_network_container.test"
		cidrBlock    = fmt.Sprintf("10.8.%d.0/24", randInt)
		providerName = "AWS"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWS(projectName, orgID, cidrBlock, providerName, "US_EAST_1"),
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
		if rs.Type != "mongodbatlas_container" {
			continue
		}

		_, _, err := acc.Conn().Containers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])

		if err == nil {
			return fmt.Errorf("container (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
		}
	}
	return nil
}

func configAWS(projectName, orgID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
			region_name			 = "%s"
		}

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = "${mongodbatlas_network_container.test.project_id}"
			provider_name = "AWS"
		}

	`, projectName, orgID, cidrBlock, providerName, region)
}

func configAzure(projectName, orgID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "${mongodbatlas_project.test.id}"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
			region			     = "US_EAST_2"
		}
	`, projectName, orgID, cidrBlock, providerName)
}

func configGCP(projectName, orgID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "${mongodbatlas_project.test.id}"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
		}
	`, projectName, orgID, cidrBlock, providerName)
}

func configGCPWithRegions(projectName, orgID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = mongodbatlas_project.test.id
			atlas_cidr_block = "%s"
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
	`, projectName, orgID, cidrBlock, providerName)
}
