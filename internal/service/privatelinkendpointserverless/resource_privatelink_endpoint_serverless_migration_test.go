package privatelinkendpointserverless_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigServerlessPrivateLinkEndpoint_basic(t *testing.T) {
	acc.SkipTestForCI(t)
	mig.CreateAndRunTest(t, basicTestCase(t))
}
