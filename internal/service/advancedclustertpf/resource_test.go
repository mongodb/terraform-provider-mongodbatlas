package advancedclustertpf_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/tc"
)

const (
	resourceName = "mongodbatlas_advanced_cluster.test"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)
	testCase := tc.BasicTenantTestCase(t, projectID, clusterName, clusterNameUpdated)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)
	testCase := tc.SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t, orgID, projectName, clusterName)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)
	testCase := tc.SymmetricShardedOldSchema(t, orgID, projectName, clusterName)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedCluster_tenantUpgrade(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	testCase := tc.TenantUpgrade(t, projectID, clusterName)
	resource.ParallelTest(t, *testCase)
}
