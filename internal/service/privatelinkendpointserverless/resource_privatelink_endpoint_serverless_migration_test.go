package privatelinkendpointserverless_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigServerlessPrivateLinkEndpoint_basic(t *testing.T) {
	acc.SkipTestForCI(t) // Serverless Instances now create Flex clusters
	mig.CreateAndRunTest(t, basicTestCase(t))
}
