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

func TestMigNetworkContainerRS_basicAWS(t *testing.T) {
	var (
		projectID = mig.ProjectIDGlobal(t)
		randInt   = acctest.RandIntRange(0, 255)
		cidrBlock = fmt.Sprintf("10.8.%d.0/24", randInt)
		config    = configBasic(projectID, cidrBlock, constant.AWS, "US_EAST_1")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", constant.AWS),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigNetworkContainerRS_basicAzure(t *testing.T) {
	var (
		projectID = mig.ProjectIDGlobal(t)
		randInt   = acctest.RandIntRange(0, 255)
		cidrBlock = fmt.Sprintf("10.8.%d.0/24", randInt)
		config    = configBasic(projectID, cidrBlock, constant.AZURE, "US_EAST_2")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", constant.AZURE),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigNetworkContainerRS_basicGCP(t *testing.T) {
	var (
		projectID    = mig.ProjectIDGlobal(t)
		randInt      = acctest.RandIntRange(0, 255)
		gcpCidrBlock = fmt.Sprintf("10.%d.0.0/18", randInt)
		config       = configBasic(projectID, gcpCidrBlock, constant.GCP, "")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", constant.GCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
