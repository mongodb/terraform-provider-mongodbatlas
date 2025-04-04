package cluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCluster_basicAWS_simple(t *testing.T) {
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCluster_partial_advancedConf(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.24.0") // version where tls_cipher_config_mode was introduced
	mig.CreateAndRunTest(t, partialAdvancedConfTestCase(t))
}

// func TestMigDefaultWriteReadAdvancedConf_advancedConf(t *testing.T) {
// 	mig.SkipIfVersionBelow(t, "1.24.0") // version where tls_cipher_config_mode was introduced
// 	mig.CreateAndRunTest(t, basicDefaultWriteReadAdvancedConfTestCase(t))
// }
