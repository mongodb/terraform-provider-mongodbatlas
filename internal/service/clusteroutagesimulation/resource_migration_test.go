package clusteroutagesimulation_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigOutageSimulationCluster_SingleRegion_basic(t *testing.T) {
	mig.CreateAndRunTest(t, singleRegionTestCase(t))
}

func TestMigOutageSimulationCluster_MultiRegion_basic(t *testing.T) {
	mig.CreateAndRunTest(t, multiRegionTestCase(t))
}
