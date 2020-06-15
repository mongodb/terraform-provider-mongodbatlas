package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasNetworkContainer_basicAWS(t *testing.T) {
	var container matlas.Container

	randInt := acctest.RandIntRange(0, 255)
	randIntUpdated := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_network_container.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("10.8.%d.0/24", randInt)
	cidrBlockUpdated := fmt.Sprintf("10.8.%d.0/24", randIntUpdated)

	providerName := "AWS"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectID, cidrBlockUpdated, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasNetworkContainer_basicAzure(t *testing.T) {
	var container matlas.Container

	randInt := acctest.RandIntRange(0, 255)
	randIntUpdated := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_network_container.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("10.8.%d.0/24", randInt)
	cidrBlockUpdated := fmt.Sprintf("10.8.%d.0/24", randIntUpdated)

	providerName := "AZURE"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAzure(projectID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAzure(projectID, cidrBlockUpdated, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasNetworkContainer_basicGCP(t *testing.T) {
	var container matlas.Container

	randInt := acctest.RandIntRange(0, 255)
	randIntUpdated := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_network_container.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("10.8.%d.0/24", randInt)
	cidrBlockUpdated := fmt.Sprintf("10.8.%d.0/24", randIntUpdated)

	providerName := "GCP"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigGCP(projectID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigGCP(projectID, cidrBlockUpdated, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasNetworkContainer_importBasic(t *testing.T) {
	randInt := acctest.RandIntRange(0, 255)
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	resourceName := "mongodbatlas_network_container.test"
	cidrBlock := fmt.Sprintf("10.8.%d.0/24", randInt)
	providerName := "AWS"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectID, cidrBlock, providerName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasNetworkContainerImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccCheckMongoDBAtlasNetworkContainerImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]), nil
	}
}

func testAccCheckMongoDBAtlasNetworkContainerExists(resourceName string, container *matlas.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if containerResp, _, err := conn.Containers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]); err == nil {
			*container = *containerResp
			return nil
		}
		return fmt.Errorf("container(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
	}
}

func testAccCheckMongoDBAtlasNetworkContainerAttributes(container *matlas.Container, providerName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if container.ProviderName != providerName {
			return fmt.Errorf("bad provider name: %s", container.ProviderName)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasNetworkContainerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_container" {
			continue
		}

		// Try to find the container
		_, _, err := conn.Containers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])

		if err == nil {
			return fmt.Errorf("container (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasNetworkContainerConfigAWS(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
			region_name			 = "US_EAST_1"
		}
	`, projectID, cidrBlock, providerName)
}

func testAccMongoDBAtlasNetworkContainerConfigAzure(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
			region_name			 = "US_EAST_2"
		}
	`, projectID, cidrBlock, providerName)
}

func testAccMongoDBAtlasNetworkContainerConfigGCP(projectID, cidrBlock, providerName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "%s"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
		}
	`, projectID, cidrBlock, providerName)
}
