package streamprocessor_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamProcessor_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.19.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCaseMigration(t))
}
