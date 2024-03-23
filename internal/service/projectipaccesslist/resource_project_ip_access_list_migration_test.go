package projectipaccesslist_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigProjectIPAccessList_settingIPAddress(t *testing.T) {
	var (
		projectID = mig.ProjectIDGlobal(t)
		ipAddress = acc.RandomIP(179, 154, 226)
		comment   = fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
		config    = configWithIPAddress(projectID, ipAddress, comment)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeTestCheckFunc(commonChecks(ipAddress, "", "", comment)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigProjectIPAccessList_settingCIDRBlock(t *testing.T) {
	var (
		projectID = mig.ProjectIDGlobal(t)
		cidrBlock = acc.RandomIP(179, 154, 226) + "/32"
		comment   = fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)
		config    = configWithCIDRBlock(projectID, cidrBlock, comment)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeTestCheckFunc(commonChecks("", cidrBlock, "", comment)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigProjectIPAccessList_settingAWSSecurityGroup(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t) // No mig.ProjectIDGlobal because network container
		awsSGroup    = os.Getenv("AWS_SECURITY_GROUP_1")
		vpcID        = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion    = os.Getenv("AWS_REGION")
		providerName = "AWS"
		comment      = fmt.Sprintf("TestAcc for awsSecurityGroup (%s)", awsSGroup)
		config       = configWithAWSSecurityGroup(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckPeeringEnvAWS(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeTestCheckFunc(commonChecks("", "", awsSGroup, comment)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
