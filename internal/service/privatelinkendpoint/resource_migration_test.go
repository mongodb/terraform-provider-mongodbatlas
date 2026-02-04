package privatelinkendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigPrivateLinkEndpoint_basicAWS(t *testing.T) {
	mig.CreateAndRunTest(t, basicAWSTestCase(t, "us-west-2"))
}
