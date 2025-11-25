package projectipaccesslist_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_project_ip_access_list.test"
	dataSourceName = "data.mongodbatlas_project_ip_access_list.test"
)

func TestAccProjectIPAccesslist_settingIPAddress(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		ipAddress        = acc.RandomIP(179, 154, 226)
		comment          = fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
		updatedIPAddress = acc.RandomIP(179, 154, 228)
		updatedComment   = fmt.Sprintf("TestAcc for ipAddress updated (%s)", updatedIPAddress)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithIPAddress(projectID, ipAddress, comment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(ipAddress, "", "", comment)...),
			},
			{
				Config: configWithIPAddress(projectID, updatedIPAddress, updatedComment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(updatedIPAddress, "", "", updatedComment)...),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectIPAccessList_settingCIDRBlock(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		cidrBlock        = acc.RandomIP(179, 154, 226) + "/32"
		comment          = fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)
		updatedCIDRBlock = acc.RandomIP(179, 154, 228) + "/32"
		updatedComment   = fmt.Sprintf("TestAcc for cidrBlock updated (%s)", updatedCIDRBlock)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithCIDRBlock(projectID, cidrBlock, comment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks("", cidrBlock, "", comment)...),
			},
			{
				Config: configWithCIDRBlock(projectID, updatedCIDRBlock, updatedComment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks("", updatedCIDRBlock, "", updatedComment)...),
			},
		},
	})
}

func TestAccProjectIPAccessList_settingAWSSecurityGroup(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		vpcID            = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock     = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID     = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion        = os.Getenv("AWS_REGION")
		awsSGroup        = os.Getenv("AWS_SECURITY_GROUP_1")
		updatedAWSSgroup = os.Getenv("AWS_SECURITY_GROUP_2")
		providerName     = "AWS"
		comment          = fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)
		updatedComment   = fmt.Sprintf("TestAcc for awsSecurityGroup updated (%s)", updatedAWSSgroup)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPeeringEnvAWS(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks("", "", awsSGroup, comment)...),
			},
			{
				Config: configWithAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, updatedAWSSgroup, updatedComment),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks("", "", updatedAWSSgroup, updatedComment)...),
			},
		},
	})
}

func TestAccProjectIPAccessList_settingMultiple(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		resourceFmt      = "mongodbatlas_project_ip_access_list.test_%d"
		ipWhiteListCount = 20
		accessList       = []map[string]string{}
		checks           = []resource.TestCheckFunc{}
	)

	for i := range ipWhiteListCount {
		entry := make(map[string]string)
		entryName := ""
		ipAddr := ""

		if i%2 == 0 {
			entryName = "cidr_block"
			entry["cidr_block"] = acc.RandomIP(byte(i), 2, 3) + "/32"
			ipAddr = entry["cidr_block"]
		} else {
			entryName = "ip_address"
			entry["ip_address"] = acc.RandomIP(byte(i), 2, 3)
			ipAddr = entry["ip_address"]
		}
		entry["comment"] = fmt.Sprintf("TestAcc for %s (%s)", entryName, ipAddr)

		accessList = append(accessList, entry)
		checks = append(checks, checkExists(fmt.Sprintf(resourceFmt, i)))
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMultiple(projectID, accessList, false),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				Config: configWithMultiple(projectID, accessList, true),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccProjectIPAccessList_importIncorrectId(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		ipAddress = acc.RandomIP(179, 154, 226)
		comment   = fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithIPAddress(projectID, ipAddress, comment),
			},
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: "incorrect_id_without_project_id_and_dash",
				ExpectError:   regexp.MustCompile("import format error"),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ProjectIPAccessListApi.GetAccessListEntry(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err != nil {
			return fmt.Errorf("project ip access list entry (%s) does not exist", ids["entry"])
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_access_list" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ProjectIPAccessListApi.GetAccessListEntry(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err == nil {
			return fmt.Errorf("project ip access list entry (%s) still exists", ids["entry"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["entry"]), nil
	}
}

func commonChecks(ipAddress, cidrBlock, awsSGroup, comment string) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		checkExists(dataSourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "comment", comment),
		resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
	}
	if ipAddress != "" {
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
			resource.TestCheckResourceAttr(dataSourceName, "ip_address", ipAddress))
	}
	if cidrBlock != "" {
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
			resource.TestCheckResourceAttr(dataSourceName, "cidr_block", cidrBlock))
	}
	if awsSGroup != "" {
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
			resource.TestCheckResourceAttr(dataSourceName, "aws_security_group", awsSGroup))
	}
	return checks
}

func configWithIPAddress(projectID, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = %[1]q
			ip_address = %[2]q
			comment    = %[3]q
		}

		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			ip_address = mongodbatlas_project_ip_access_list.test.ip_address
		}
	`, projectID, ipAddress, comment)
}

func configWithCIDRBlock(projectID, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = %[1]q
			cidr_block = %[2]q
			comment    = %[3]q
		}

		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			cidr_block = mongodbatlas_project_ip_access_list.test.cidr_block
		}
	`, projectID, cidrBlock, comment)
}

func configWithAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = %[1]q
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = %[2]q
			region_name			  = %[6]q
		}

		resource "mongodbatlas_network_peering" "test" {
			accepter_region_name	  = "us-east-1"
			project_id    			    = %[1]q
			container_id            = mongodbatlas_network_container.test.container_id
			provider_name           = %[2]q
			vpc_id					        = %[3]q
			aws_account_id	        = %[4]q
			route_table_cidr_block  = %[5]q
		}

		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id         = %[1]q
			aws_security_group = %[7]q
			comment            = %[8]q

			depends_on = ["mongodbatlas_network_peering.test"]
		}

		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = %[1]q
			aws_security_group = mongodbatlas_project_ip_access_list.test.aws_security_group
		}
	`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}

func configWithMultiple(projectID string, accessList []map[string]string, isUpdate bool) string {
	var config strings.Builder
	for i, entry := range accessList {
		comment := entry["comment"]
		if isUpdate {
			comment = entry["comment"] + " update"
		}

		if cidr, ok := entry["cidr_block"]; ok {
			config.WriteString(fmt.Sprintf(`
				resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
					project_id   = %[2]q
					cidr_block = %[3]q
					comment    = %[4]q
				}
			`, i, projectID, cidr, comment))
		} else {
			config.WriteString(fmt.Sprintf(`
				resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
					project_id   = %[2]q
					ip_address = %[3]q
					comment    = %[4]q
				}
			`, i, projectID, entry["ip_address"], comment))
		}
	}
	return config.String()
}
