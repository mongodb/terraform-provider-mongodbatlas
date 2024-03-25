package searchindex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigSearchIndexRS_basic(t *testing.T) {
	testCase := mig.ConvertToMigrationTest(t, basicTestCase(t))
	resource.ParallelTest(t, testCase)
}

func TestMigSearchIndexRS_withVector(t *testing.T) {
	testCase := mig.ConvertToMigrationTest(t, basicTestCaseVector(t))
	resource.ParallelTest(t, testCase)
}
