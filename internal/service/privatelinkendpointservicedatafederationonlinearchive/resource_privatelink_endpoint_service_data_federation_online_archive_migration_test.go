package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	acc.SkipTestForCI(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            resourceConfigBasic(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive),
					resource.TestCheckResourceAttr(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "comment"),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "type"),
					resource.TestCheckResourceAttrSet(resourceNamePrivatelinkEdnpointServiceDataFederationOnlineArchive, "provider_name"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   resourceConfigBasic(projectID, endpointID),
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
