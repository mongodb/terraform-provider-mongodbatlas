package ldapverify_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigLDAPVerify_basic(t *testing.T) {
	mig.CreateAndRunTest(t, basicTestCase(t))
}
