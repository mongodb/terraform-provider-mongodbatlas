package networkcontainer_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationNetworkContainerRS_basicAWS(t *testing.T) {
	var (
		projectName = acc.RandomProjectName()
		configAWS   = configAWS(projectName, orgID, cidrBlock, providerNameAws, "US_EAST_1")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configAWS,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(configAWS),
		},
	})
}

func TestAccMigrationNetworkContainerRS_basicAzure(t *testing.T) {
	var (
		projectName = acc.RandomProjectName()
		configAzure = configAzure(projectName, orgID, cidrBlock, providerNameAzure)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configAzure,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(configAzure),
		},
	})
}

func TestAccMigrationNetworkContainerRS_basicGCP(t *testing.T) {
	var (
		projectName = acc.RandomProjectName()
		configGCP   = configGCP(projectName, orgID, gcpCidrBlock, providerNameGCP)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configGCP,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStepCheckEmptyPlan(configGCP),
		},
	})
}
