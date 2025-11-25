package networkpeering_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkNetworkPeering_basicAWS(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0")
	mig.CreateAndRunTest(t, basicAWSTestCase(t))
}
