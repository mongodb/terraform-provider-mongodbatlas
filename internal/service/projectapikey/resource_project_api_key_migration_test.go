package projectapikey_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigProjectAPIKey_basic(t *testing.T) {
	mig.CreateAndRunTest(t, basicTestCase(t))
}
