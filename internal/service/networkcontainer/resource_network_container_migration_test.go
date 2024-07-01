package networkcontainer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkContainer_basicAWS(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		randInt   = acctest.RandIntRange(0, 255)
		cidrBlock = fmt.Sprintf("10.8.%d.0/24", randInt)
		config    = configBasic(projectID, cidrBlock, constant.AWS, "US_EAST_1")
	)

	// Serial so it doesn't conflict with TestAccNetworkContainer_basicAWS
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigNetworkContainer_basicAzure(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		randInt   = acctest.RandIntRange(0, 255)
		cidrBlock = fmt.Sprintf("10.8.%d.0/24", randInt)
		config    = configBasic(projectID, cidrBlock, constant.AZURE, "US_EAST_2")
	)

	// Serial so it doesn't conflict with TestAccNetworkContainer_basicAzure
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AZURE)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigNetworkContainer_basicGCP(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		randInt      = acctest.RandIntRange(0, 255)
		gcpCidrBlock = fmt.Sprintf("10.%d.0.0/18", randInt)
		config       = configBasic(projectID, gcpCidrBlock, constant.GCP, "")
	)

	// Serial so it doesn't conflict with TestAccNetworkContainer_basicGCP
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeAggregateTestCheckFunc(commonChecks(constant.GCP)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
