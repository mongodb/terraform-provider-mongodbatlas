package clusteroutagesimulation_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigOutageSimulationCluster_SingleRegion_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // version where advanced_cluster TPF was GA
	mig.CreateAndRunTest(t, singleRegionTestCase(t))
}

func TestMigOutageSimulationCluster_MultiRegion_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // version where advanced_cluster TPF was GA
	mig.CreateAndRunTest(t, multiRegionTestCase(t))
}
