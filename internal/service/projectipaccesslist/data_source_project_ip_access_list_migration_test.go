package projectipaccesslist_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccMigrationProjectDSProjectIPAccessList_SettingIPAddress(t *testing.T) {
	dataSourceName := "data.mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.PreCheckBasicMigration(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccDataMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "comment"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccDataMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectDSProjectIPAccessList_SettingCIDRBlock(t *testing.T) {
	dataSourceName := "data.mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.PreCheckBasicMigration(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccDataMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceName, "comment"),
					resource.TestCheckResourceAttr(dataSourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccDataMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectDSProjectIPAccessList_SettingAWSSecurityGroup(t *testing.T) {
	acc.SkipTestExtCred(t)
	dataSourceName := "data.mongodbatlas_project_ip_access_list.test"
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")
	providerName := "AWS"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	awsSGroup := os.Getenv("AWS_SECURITY_GROUP_ID")
	comment := fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccDataMongoDBAtlasProjectIPAccessListConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectIPAccessListExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_security_group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "comment"),

					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "aws_security_group", awsSGroup),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccDataMongoDBAtlasProjectIPAccessListConfigSettingAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
