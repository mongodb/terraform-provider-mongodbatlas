package accesslistapikey_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationProjectAccesslistAPIKey_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_access_list_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	description := fmt.Sprintf("test-acc-access_list-api_key-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configWithIPAddress(orgID, description, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithIPAddress(orgID, description, ipAddress),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccMigrationProjectAccesslistAPIKey_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_access_list_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	description := fmt.Sprintf("test-acc-access_list-api_key-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configWithCIDRBlock(orgID, description, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithCIDRBlock(orgID, description, cidrBlock),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
