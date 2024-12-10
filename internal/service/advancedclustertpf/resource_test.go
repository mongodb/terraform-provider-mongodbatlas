package advancedclustertpf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/tc"
)

const (
	resourceName = "mongodbatlas_advanced_cluster.test"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	testCase := tc.BasicTenantTestCase(t)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	testCase := tc.SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
	testCase := tc.SymmetricShardedOldSchema(t)
	resource.ParallelTest(t, *testCase)
}

func TestAccClusterAdvancedCluster_tenantUpgrade(t *testing.T) {
	testCase := tc.TenantUpgrade(t)
	resource.ParallelTest(t, *testCase)
}
