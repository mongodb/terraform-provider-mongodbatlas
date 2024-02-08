package networkcontainer_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationNetworkContainerRS_basicAWS(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configAWS(projectName, orgID, cidrBlock, providerNameAws, "US_EAST_1"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAws),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStep(configAWS(projectName, orgID, cidrBlock, providerNameAws, "US_EAST_1")),
		},
	})
}

func TestAccMigrationNetworkContainerRS_basicAzure(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configAzure(projectName, orgID, cidrBlock, providerNameAzure),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameAzure),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStep(configAzure(projectName, orgID, cidrBlock, providerNameAzure)),
		},
	})
}

func TestAccMigrationNetworkContainerRS_basicGCP(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configGCP(projectName, orgID, gcpCidrBlock, providerNameGCP),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &container),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerNameGCP),
					resource.TestCheckResourceAttrSet(resourceName, "provisioned"),
				),
			},
			mig.TestStep(configGCP(projectName, orgID, gcpCidrBlock, providerNameGCP)),
		},
	})
}
