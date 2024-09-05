package cluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCluster_basicAWS_simple(t *testing.T) {
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCluster_partial_advancedConf(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.19.0") // version where change_stream_options_pre_and_post_images_expire_after_seconds was introduced
	mig.CreateAndRunTest(t, partialAdvancedConfTestCase(t))
}
