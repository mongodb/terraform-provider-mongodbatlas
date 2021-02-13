package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCustomDNSConfigurationAWS_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_dns_configuration_cluster_aws.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(projectID, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCustomDNSConfigurationAWS_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_dns_configuration_cluster_aws.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasCustomDNSConfigurationAWSStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		_, _, err := conn.CustomAWSDNS.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return nil
		}

		return fmt.Errorf("custom dns configuration cluster(%s) does not exist", rs.Primary.ID)
	}
}
func testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_custom_dns_configuration_cluster_aws" {
			continue
		}

		// Try to find the Custom DNS Configuration for Atlas Clusters on AWS
		resp, _, err := conn.CustomAWSDNS.Get(context.Background(), rs.Primary.ID)
		if err != nil && resp != nil && resp.Enabled {
			return fmt.Errorf("custom dns configuration cluster aws (%s) still enabled", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasCustomDNSConfigurationAWSStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id     = "%s"
			enabled       = %t
		}`, projectID, enabled)
}
