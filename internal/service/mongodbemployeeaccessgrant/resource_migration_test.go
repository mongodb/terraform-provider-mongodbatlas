package mongodbemployeeaccessgrant_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigMongoDBEmployeeAccessGrant_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.20.0")                  // this feature was introduced in provider version 1.20.0
	mig.CreateAndRunTestNonParallel(t, basicTestCase(t)) // does not run in parallel to reuse same execution project and cluster
}
