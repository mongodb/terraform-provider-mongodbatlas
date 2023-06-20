package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectRSProjectIPAccesslist_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for ipAddress updated (%s)", updatedIPAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, updatedIPAddress, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for cidrBlock updated (%s)", updatedCIDRBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, updatedCIDRBlock, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_SettingAWSSecurityGroup(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "mongodbatlas_project_ip_access_list.test"
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, updatedAWSSgroup, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName),
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

func TestAccProjectRSProjectIPAccessList_SettingMultiple(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test_%d"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	const ipWhiteListCount = 20
	accessList := make([]map[string]string, 0)

	for i := 0; i < ipWhiteListCount; i++ {
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

		accessList = append(accessList, entry)
	}

	//TODO: make testAccCheckMongoDBAtlasProjectIPAccessListExists dynamic
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingMultiple(projectName, orgID, accessList, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 0)),
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 1)),
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 2)),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingMultiple(projectName, orgID, accessList, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 0)),
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 1)),
					testAccCheckMongoDBAtlasProjectIPAccessListExists(fmt.Sprintf(resourceName, 2)),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
	resourceName := "mongodbatlas_project_ip_access_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectIPAccessListImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPAccessListExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPAccessList.Get(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			return fmt.Errorf("project ip access list entry (%s) does not exist", ids["entry"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasProjectIPAccessListDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_access_list" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPAccessList.Get(context.Background(), ids["project_id"], ids["entry"])
		if err == nil {
			return fmt.Errorf("project ip access list entry (%s) still exists", ids["entry"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectIPAccessListImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["entry"]), nil
	}
}

func testAccMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment string) string {
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
	`, orgID, projectName, ipAddress, comment)
}

func testAccMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment string) string {
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
	`, orgID, projectName, cidrBlock, comment)
}

func testAccMongoDBAtlasProjectIPAccessListConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
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

		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id         = "%[1]s"
			aws_security_group = "%[7]s"
			comment            = "%[8]s"

			depends_on = ["mongodbatlas_network_peering.test"]
		}
	`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}

func testAccMongoDBAtlasProjectIPAccessListConfigSettingMultiple(projectName, orgID string, accessList []map[string]string, isUpdate bool) string {
	config := fmt.Sprintf(`
			resource "mongodbatlas_project" "test" {
				name   = %[1]q
				org_id = %[2]q
			}`, projectName, orgID)

	for i, entry := range accessList {
		comment := entry["comment"]

		if isUpdate {
			comment = entry["comment"] + " update"
		}

		if cidr, ok := entry["cidr_block"]; ok {
			config += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
				project_id   = mongodbatlas_project.test.id
				cidr_block = %[2]q
				comment    = %[3]q
			}
		`, i, cidr, comment)
		} else {
			config += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
				project_id   = mongodbatlas_project.test.id
				ip_address = %[2]q
				comment    = %[3]q
			}
		`, i, entry["ip_address"], comment)
		}
	}
	return config
}
