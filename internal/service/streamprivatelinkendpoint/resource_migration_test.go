package streamprivatelinkendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamPrivatelinkEndpointConfluent_basic(t *testing.T) {
	acc.SkipTestForCI(t)                // needs confluent cloud resources
	mig.SkipIfVersionBelow(t, "1.25.0") // when resource 1st released
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicConfluentTestCase(t), mig.ExternalProvidersWithConfluent(), nil)
}

func TestMigStreamPrivatelinkEndpointMsk_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs an AWS MSK cluster
	mig.SkipIfVersionBelow(t, "1.29.0")
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicMskTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
