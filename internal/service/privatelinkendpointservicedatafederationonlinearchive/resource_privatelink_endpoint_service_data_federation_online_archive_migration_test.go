package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	var (
		projectID  = acc.ProjectIDGlobal(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		config     = resourceConfigBasic(projectID, endpointID, comment)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckPrivateEndpoint(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
