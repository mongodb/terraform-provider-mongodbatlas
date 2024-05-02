package privatelinkendpointservice_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkNetworkPeering_basicAWS(t *testing.T) {
	// can only be one privatelinkendpointservice per project
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicAWSTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
