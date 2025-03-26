package resourcepolicy_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigResourcePolicy_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.22.0") // this feature was introduced in provider version 1.21.0, plural data source schema was changed in 1.22.0

	var description *string
	if mig.IsProviderVersionAtLeast("1.32.0") {
		description = descriptionPtr
	}

	mig.CreateAndRunTestNonParallel(t, basicTestCase(t, description))
}
