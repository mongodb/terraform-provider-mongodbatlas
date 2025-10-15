package streamconnection_test

import (
	_ "embed"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA
	mig.CreateAndRunTest(t, testCaseKafkaPlaintext(t, "-mig"))
}

func TestMigStreamRSStreamConnection_cluster(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA
	mig.CreateAndRunTest(t, testCaseCluster(t, "-mig"))
}
