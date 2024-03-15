package projectipaccesslist_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectRSProjectIPAccesslist_SettingIPAddress(t *testing.T) {
	var (
		resourceName     = "mongodbatlas_project_ip_access_list.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acc.RandomProjectName()
		ipAddress        = acc.RandomIP(179, 154, 226)
		comment          = fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
		updatedIPAddress = acc.RandomIP(179, 154, 228)
		updatedComment   = fmt.Sprintf("TestAcc for ipAddress updated (%s)", updatedIPAddress)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, updatedIPAddress, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
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
	var (
		resourceName     = "mongodbatlas_project_ip_access_list.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acc.RandomProjectName()
		cidrBlock        = acc.RandomIP(179, 154, 226) + "/32"
		comment          = fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)
		updatedCIDRBlock = acc.RandomIP(179, 154, 228) + "/32"
		updatedComment   = fmt.Sprintf("TestAcc for cidrBlock updated (%s)", updatedCIDRBlock)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithCIDRBlock(orgID, projectName, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: acc.ConfigProjectIPAccessListWithCIDRBlock(orgID, projectName, updatedCIDRBlock, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
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
	var (
		resourceName     = "mongodbatlas_project_ip_access_list.test"
		vpcID            = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock     = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID     = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion        = os.Getenv("AWS_REGION")
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		awsSGroup        = os.Getenv("AWS_SECURITY_GROUP_1")
		updatedAWSSgroup = os.Getenv("AWS_SECURITY_GROUP_2")
		providerName     = "AWS"
		projectName      = acc.RandomProjectName()
		comment          = fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)
		updatedComment   = fmt.Sprintf("TestAcc for awsSecurityGroup updated (%s)", updatedAWSSgroup)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithAWSSecurityGroup(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: acc.ConfigProjectIPAccessListWithAWSSecurityGroup(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, updatedAWSSgroup, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "aws_security_group", updatedAWSSgroup),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_SettingMultiple(t *testing.T) {
	var (
		resourceName     = "mongodbatlas_project_ip_access_list.test_%d"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acc.RandomProjectName()
		ipWhiteListCount = 20
		accessList       = make([]map[string]string, 0)
	)

	for i := 0; i < ipWhiteListCount; i++ {
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
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithMultiple(projectName, orgID, accessList, false),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 0)),
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 1)),
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 2)),
				),
			},
			{
				Config: acc.ConfigProjectIPAccessListWithMultiple(projectName, orgID, accessList, true),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 0)),
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 1)),
					acc.CheckProjectIPAccessListExists(fmt.Sprintf(resourceName, 2)),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_importBasic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		ipAddress    = acc.RandomIP(179, 154, 226)
		comment      = fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
		resourceName = "mongodbatlas_project_ip_access_list.test"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.ImportStateProjecIPAccessListtIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessList_importIncorrectId(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		ipAddress    = acc.RandomIP(179, 154, 226)
		comment      = fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
		resourceName = "mongodbatlas_project_ip_access_list.test"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment),
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
