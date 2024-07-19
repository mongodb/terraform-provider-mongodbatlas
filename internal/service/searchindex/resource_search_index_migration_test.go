package searchindex_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigSearchIndex_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.17.4")
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigSearchIndex_withVector(t *testing.T) {
	mig.CreateAndRunTest(t, basicVectorTestCase(t))
}
