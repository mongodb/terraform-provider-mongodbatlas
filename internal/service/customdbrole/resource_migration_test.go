package customdbrole_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigCustomDBRoles_Basic(t *testing.T) {
	mig.CreateAndRunTest(t, basicTestCase(t))
}
