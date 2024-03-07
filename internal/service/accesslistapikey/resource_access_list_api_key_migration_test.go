package accesslistapikey_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationProjectAccesslistAPIKey_SettingIPAddress(t *testing.T) {
	var (
		resourceName = "mongodbatlas_access_list_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		ipAddress    = acc.RandomIP(179, 154, 226)
		description  = acc.RandomName()
	)

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
	var (
		resourceName = "mongodbatlas_access_list_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		cidrBlock    = acc.RandomIP(179, 154, 226) + "/32"
		description  = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"), // testing provider version before this resource was migrated to Atlas SDK
				Config:            configWithCIDRBlock(orgID, description, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configWithCIDRBlock(orgID, description, cidrBlock),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
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

func TestAccMigrationProjectAccesslistAPIKey_SettingCIDRBlock_WideCIDR_SDKMigration(t *testing.T) {
	var (
		resourceName = "mongodbatlas_access_list_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		cidrBlock    = "100.10.0.0/16"
		description  = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"), // testing provider version before this resource was migrated to Atlas SDK
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
