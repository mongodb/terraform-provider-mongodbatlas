package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		config     = resourceConfigBasicAWS(projectID, endpointID, comment)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckPrivateEndpoint(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: checkResourceOnlyAggr(projectID, endpointID,
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_optionalAttrsOmitted(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		config     = resourceConfigBasicAWSNoOptional(projectID, endpointID)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckPrivateEndpoint(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkResourceOnlyAggr(projectID, endpointID),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
