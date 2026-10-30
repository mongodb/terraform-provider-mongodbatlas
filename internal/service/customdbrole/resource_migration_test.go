package customdbrole_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigCustomDBRoles_Basic(t *testing.T) {
	skipSupersededMigTest(t)
	mig.CreateAndRunTest(t, basicTestCase(t))
}

// skipSupersededMigTest skips the manual-era migration tests while mongodbatlas_custom_db_role is served by
// the autogen implementation (spike pilot): their configs use the assignment syntax the released provider
// rejects. Superseded by the two-config TestMigCustomDBRole_* tests in internal/serviceapi/customdbrole,
// restored or removed at the full migration.
func skipSupersededMigTest(t *testing.T) {
	t.Helper()
	t.Skip("superseded by TestMigCustomDBRole_* in internal/serviceapi/customdbrole (autogen lift-and-shift pilot)")
}
