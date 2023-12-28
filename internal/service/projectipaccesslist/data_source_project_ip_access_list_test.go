package projectipaccesslist_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectDSProjectIPAccessList_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccProjectDSProjectIPAccessList_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func TestAccProjectDSProjectIPAccessList_SettingAWSSecurityGroup(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	projectName := acctest.RandomWithPrefix("test-acc-project-ds-aws")
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")
	providerName := "AWS"

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	awsSGroup := os.Getenv("AWS_SECURITY_GROUP_1")
	comment := fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithAWSSecurityGroup(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
		},
	})
}

func testAccDataMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project.test.id
			ip_address = %[3]q
			comment    = %[4]q
		}

		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			ip_address = mongodbatlas_project_ip_access_list.test.ip_address
		}
	`, orgID, projectName, ipAddress, comment)
}

func testAccDataMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project.test.id
			cidr_block = %[3]q
			comment    = %[4]q
		}
		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			cidr_block = mongodbatlas_project_ip_access_list.test.cidr_block
		}
	`, orgID, projectName, cidrBlock, comment)
}
