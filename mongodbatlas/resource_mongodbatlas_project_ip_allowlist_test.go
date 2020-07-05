package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProjectIPAllowlist_SettingIPAddress(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_allowlist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for ipAddress updated (%s)", updatedIPAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPAllowlistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, updatedIPAddress, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})

}
func TestAccResourceMongoDBAtlasProjectIPAllowlist_SettingCIDRBlock(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_allowlist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for cidrBlock updated (%s)", updatedCIDRBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPAllowlistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingCIDRBlock(projectID, cidrBlock, comment),
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
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingCIDRBlock(projectID, updatedCIDRBlock, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasProjectIPAllowlist_SettingAWSSecurityGroup(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_allowlist.test"
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")
	providerName := "AWS"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	awsSGroup := "sg-0026348ec11780bd1"
	comment := fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)

	updatedAWSSgroup := "sg-0026348ec11780bd2"
	updatedComment := fmt.Sprintf("TestAcc for awsSecurityGroup updated (%s)", updatedAWSSgroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPAllowlistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
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
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, updatedAWSSgroup, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_security_group", updatedAWSSgroup),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasProjectIPAllowlist_SettingMultiple(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_allowlist.test_%d"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	var entry, comment, entryName string

	// Creating 100 allowlist entriies at the same time
	for i := 0; i < 100; i++ {
		entry = fmt.Sprintf("%d.2.3.%d", i, acctest.RandIntRange(0, 255))
		comment = fmt.Sprintf("TestAcc for %s (%s)", entryName, entry)
		if i%2 == 0 {
			entry = fmt.Sprintf("%d.2.3.%d/32", i, acctest.RandIntRange(0, 255))
			comment = fmt.Sprintf("TestAcc for %s (%s)", entryName, entry)
		}

		t.Run(comment, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasProjectIPAllowlistDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingMultiple(projectID, entry, comment, i),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasProjectIPAllowlistExists(fmt.Sprintf(resourceName, i)),
						),
					},
					{
						Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingMultiple(projectID, entry, comment+" updated", i),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasProjectIPAllowlistExists(fmt.Sprintf(resourceName, i)),
						),
					},
				},
			})
		})
	}
}

func TestAccResourceMongoDBAtlasProjectIPAllowlist_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
	resourceName := "mongodbatlas_project_ip_allowlist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPAllowlistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, ipAddress, comment),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectIPAllowlistImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPAllowlistExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"]) // TODO: Language Inclusivity
		if err != nil {
			return fmt.Errorf("project ip allowlist entry (%s) does not exist", ids["entry"])
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasProjectIPAllowlistDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_allowlist" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"]) // TODO: Language Inclusivity
		if err == nil {
			return fmt.Errorf("project ip allowlist entry (%s) still exists", ids["entry"])
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasProjectIPAllowlistImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["entry"]), nil
	}
}

func testAccMongoDBAtlasProjectIPAllowlistConfigSettingIPAddress(projectID, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_allowlist" "test" {
			project_id = "%s"
			ip_address = "%s"
			comment    = "%s"
		}
	`, projectID, ipAddress, comment)
}

func testAccMongoDBAtlasProjectIPAllowlistConfigSettingCIDRBlock(projectID, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_allowlist" "test" {
			project_id = "%s"
			cidr_block = "%s"
			comment    = "%s"
		}
	`, projectID, cidrBlock, comment)
}

func testAccMongoDBAtlasProjectIPAllowlistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = "%[1]s"
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = "%[2]s"
			region_name			  = "%[6]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			accepter_region_name	  = "us-east-1"	
			project_id    			    = "%[1]s"
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
	`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}

func testAccMongoDBAtlasProjectIPAllowlistConfigSettingMultiple(projectID, entry, comment string, i int) string {
	if i%2 == 0 {
		return fmt.Sprintf(`
			resource "mongodbatlas_project_ip_allowlist" "test_%d" {
				project_id = "%s"
				cidr_block = "%s"
				comment    = "%s"
			}
		`, i, projectID, entry, comment)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_allowlist" "test_%d" {
			project_id = "%s"
			ip_address = "%s"
			comment    = "%s"
		}
	`, i, projectID, entry, comment)
}
