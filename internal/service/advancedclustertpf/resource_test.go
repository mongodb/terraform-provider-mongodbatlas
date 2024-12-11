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

func TestAccMockableAdvancedCluster_basicTenant(t *testing.T) {
	testCase := basicTenantTestCase(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccMockableAdvancedCluster_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	testCase := symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccMockableAdvancedCluster_symmetricShardedOldSchema(t *testing.T) {
	testCase := symmetricShardedOldSchema(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccMockableAdvancedCluster_tenantUpgrade(t *testing.T) {
	testCase := tenantUpgrade(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccMockableAdvancedCluster_replicasetAdvConfigUpdate(t *testing.T) {
	testCase := replicasetAdvConfigUpdate(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}

func TestAccMockableAdvancedCluster_shardedBasic(t *testing.T) {
	testCase := shardedBasic(t)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, testCase)
}
