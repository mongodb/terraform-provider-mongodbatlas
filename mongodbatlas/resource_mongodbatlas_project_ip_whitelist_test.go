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

func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for ipAddress updated (%s)", updatedIPAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, updatedIPAddress, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
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

func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for cidrBlock updated (%s)", updatedCIDRBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, updatedCIDRBlock, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
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

func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingAWSSecurityGroup(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "mongodbatlas_project_ip_whitelist.test"
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
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, updatedAWSSgroup, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
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

func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingMultiple(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_whitelist.test_%d"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	whitelist := make([]map[string]string, 0)

	for i := 0; i < 100; i++ {
		entry := make(map[string]string)
		entryName := ""
		ipAddr := ""

		if i%2 == 0 {
			entryName = "cidr_block"
			entry["cidr_block"] = fmt.Sprintf("%d.2.3.%d/32", i, acctest.RandIntRange(0, 255))
			ipAddr = entry["cidr_block"]
		} else {
			entryName = "ip_address"
			entry["ip_address"] = fmt.Sprintf("%d.2.3.%d", i, acctest.RandIntRange(0, 255))
			ipAddr = entry["ip_address"]
		}
		entry["comment"] = fmt.Sprintf("TestAcc for %s (%s)", entryName, ipAddr)

		whitelist = append(whitelist, entry)
	}
	//TODO: make testAccCheckMongoDBAtlasProjectIPWhitelistExists dynamic
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingMultiple(projectID, whitelist, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 0)),
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 1)),
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 2)),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingMultiple(projectID, whitelist, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 0)),
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 1)),
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(fmt.Sprintf(resourceName, 2)),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
	resourceName := "mongodbatlas_project_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectIPWhitelistImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName string) resource.TestCheckFunc {
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

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			return fmt.Errorf("project ip whitelist entry (%s) does not exist", ids["entry"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasProjectIPWhitelistDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_whitelist" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err == nil {
			return fmt.Errorf("project ip whitelist entry (%s) still exists", ids["entry"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectIPWhitelistImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["entry"]), nil
	}
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id = "%s"
			ip_address = "%s"
			comment    = "%s"
		}
	`, projectID, ipAddress, comment)
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id = "%s"
			cidr_block = "%s"
			comment    = "%s"
		}
	`, projectID, cidrBlock, comment)
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
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

		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id         = "%[1]s"
			aws_security_group = "%[7]s"
			comment            = "%[8]s"

			depends_on = ["mongodbatlas_network_peering.test"]
		}
	`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingMultiple(projectID string, whitelist []map[string]string, isUpdate bool) string {
	config := ""

	for i, entry := range whitelist {
		comment := entry["comment"]

		if isUpdate {
			comment = entry["comment"] + " update"
		}

		if cidr, ok := entry["cidr_block"]; ok {
			config += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_whitelist" "test_%d" {
				project_id = "%s"
				cidr_block = "%s"
				comment    = "%s"
			}
		`, i, projectID, cidr, comment)
		} else {
			config += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_whitelist" "test_%d" {
				project_id = "%s"
				ip_address = "%s"
				comment    = "%s"
			}
		`, i, projectID, entry["ip_address"], comment)
		}
	}
	return config
}
