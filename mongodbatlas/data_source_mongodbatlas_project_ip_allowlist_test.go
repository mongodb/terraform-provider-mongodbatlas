package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataMongoDBAtlasProjectIPAllowlist_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_allowlist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccDataMongoDBAtlasProjectIPAllowlist_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_allowlist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingCIDRBlock(projectID, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccDataMongoDBAtlasProjectIPAllowlist_SettingAWSSecurityGroup(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_allowlist.test"
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")
	providerName := "AWS"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	awsSGroup := os.Getenv("AWS_SECURITY_GROUP_ID")
	comment := fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_allowlist" "test" {
			project_id = "%s"
			ip_address = "%s"
			comment    = "%s"
		}

		data "mongodbatlas_project_ip_allowlist" "test" {
			project_id = mongodbatlas_project_ip_allowlist.test.project_id
			ip_address = mongodbatlas_project_ip_allowlist.test.ip_address
		}
	`, projectID, ipAddress, comment)
}

func testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingCIDRBlock(projectID, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_allowlist" "test" {
			project_id = "%s"
			cidr_block = "%s"
			comment    = "%s"
		}
		data "mongodbatlas_project_ip_allowlist" "test" {
			project_id = mongodbatlas_project_ip_allowlist.test.project_id
			cidr_block = mongodbatlas_project_ip_allowlist.test.cidr_block
		}
	`, projectID, cidrBlock, comment)
}

func testAccDataMongoDBAtlasProjectIPAllowlistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = "%[1]s"
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = "%[2]s"
			region_name			  = "%[6]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			accepter_region_name    = lower(replace("%[6]s", "_", "-"))
			project_id    		    = "%[1]s"
			container_id            = mongodbatlas_network_container.test.container_id
			provider_name           = "%[2]s"
			route_table_cidr_block  = "%[5]s"
			vpc_id					        = "%[3]s"
			aws_account_id	        = "%[4]s"
		}

		resource "mongodbatlas_project_ip_allowlist" "test" {
			project_id         = "%[1]s"
			aws_security_group = "%[7]s"
			comment            = "%[8]s"

			depends_on = ["mongodbatlas_network_peering.test"]
		}

		data "mongodbatlas_project_ip_allowlist" "test" {
			project_id = mongodbatlas_project_ip_allowlist.test.project_id
			aws_security_group = mongodbatlas_project_ip_allowlist.test.aws_security_group
		}
	`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}
