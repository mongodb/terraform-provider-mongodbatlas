package advancedclustertpf_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	mockConfig = &unit.MockHTTPDataConfig{AllowMissingRequests: true, SideEffect: shortenRetries, IsDiffMustSubstrings: []string{"/clusters"}}
)

func shortenRetries() error {
	advancedclustertpf.RetryMinTimeout = 100 * time.Millisecond
	advancedclustertpf.RetryDelay = 100 * time.Millisecond
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	return nil
}

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	testCase := BasicTenantTestCase(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	testCase := SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
	testCase := SymmetricShardedOldSchema(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccClusterAdvancedCluster_tenantUpgrade(t *testing.T) {
	testCase := TenantUpgrade(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccAdvancedCluster_replicasetAdvConfigUpdate(t *testing.T) {
	testCase := ReplicasetAdvConfigUpdate(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccAdvancedCluster_shardedBasic(t *testing.T) {
	testCase := shardedBasic(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}
