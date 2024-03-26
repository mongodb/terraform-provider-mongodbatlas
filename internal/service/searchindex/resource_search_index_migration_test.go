package searchindex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigSearchIndex_basic(t *testing.T) {
	resource.ParallelTest(t, mig.CreateTest(t, basicTestCase(t)))
}

func TestMigSearchIndex_withVector(t *testing.T) {
	resource.ParallelTest(t, mig.CreateTest(t, basicTestCaseVector(t)))
}
