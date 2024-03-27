package networkpeering_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkNetworkPeering_basicAWS(t *testing.T) {
	mig.CreateAndRunTest(t, basicAWSTestCase(t))
}
