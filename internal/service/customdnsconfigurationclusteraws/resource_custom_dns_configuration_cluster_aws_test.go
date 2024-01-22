package customdnsconfigurationclusteraws_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSCustomDNSConfigurationAWS_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_dns_configuration_cluster_aws.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(orgID, projectName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDNSConfigurationAWSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccConfigRSCustomDNSConfigurationAWS_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_dns_configuration_cluster_aws.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(orgID, projectName, true),
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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.Conn().CustomAWSDNS.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return nil
		}
		return fmt.Errorf("custom dns configuration cluster(%s) does not exist", rs.Primary.ID)
	}
}
func testAccCheckMongoDBAtlasCustomDNSConfigurationAWSDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_custom_dns_configuration_cluster_aws" {
			continue
		}

		// Try to find the Custom DNS Configuration for Atlas Clusters on AWS
		resp, _, err := acc.Conn().CustomAWSDNS.Get(context.Background(), rs.Primary.ID)
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

func testAccMongoDBAtlasCustomDNSConfigurationAWSConfig(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id     = mongodbatlas_project.test.id
			enabled       = %[3]t
		}`, orgID, projectName, enabled)
}
