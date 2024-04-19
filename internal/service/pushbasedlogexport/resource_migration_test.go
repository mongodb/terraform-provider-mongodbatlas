package pushbasedlogexport_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigPushBasedLogExport_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0")
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigPushBasedLogExport_noPrefixPath(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0")
	mig.CreateAndRunTest(t, noPrefixPathTestCase(t))
}
