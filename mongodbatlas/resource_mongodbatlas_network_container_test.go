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
	var (
		container        matlas.Container
		randInt          = acctest.RandIntRange(0, 255)
		randIntUpdated   = acctest.RandIntRange(0, 255)
		resourceName     = "mongodbatlas_network_container.test"
		cidrBlock        = fmt.Sprintf("10.8.%d.0/24", randInt)
		cidrBlockUpdated = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
		providerName     = "AWS"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acctest.RandomWithPrefix("test-acc")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectName, orgID, cidrBlock, providerName, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectName, orgID, cidrBlockUpdated, providerName, "US_WEST_2"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAzure(projectName, orgID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAzure(projectName, orgID, cidrBlockUpdated, providerName),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigGCP(projectName, orgID, cidrBlock, providerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkContainerExists(resourceName, &container),
					testAccCheckMongoDBAtlasNetworkContainerAttributes(&container, providerName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigGCP(projectName, orgID, cidrBlockUpdated, providerName),
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
	var (
		randInt      = acctest.RandIntRange(0, 255)
		resourceName = "mongodbatlas_network_container.test"
		cidrBlock    = fmt.Sprintf("10.8.%d.0/24", randInt)
		providerName = "AWS"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkContainerConfigAWS(projectName, orgID, cidrBlock, providerName, "US_EAST_1"),
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
			return "", fmt.Errorf("not found: %s", resourceName)
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

func testAccMongoDBAtlasNetworkContainerConfigAWS(projectName, orgID, cidrBlock, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id   		 = "${mongodbatlas_project.test.id}"
			atlas_cidr_block = "%s"
			provider_name		 = "%s"
			region_name			 = "%s"
		}
	`, projectName, orgID, cidrBlock, providerName, region)
}

func testAccMongoDBAtlasNetworkContainerConfigAzure(projectName, orgID, cidrBlock, providerName string) string {
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

func testAccMongoDBAtlasNetworkContainerConfigGCP(projectName, orgID, cidrBlock, providerName string) string {
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
