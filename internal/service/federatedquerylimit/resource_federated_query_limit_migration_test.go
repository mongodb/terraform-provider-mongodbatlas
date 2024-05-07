package federatedquerylimit_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigFederatedDatabaseQueryLimit_basic(t *testing.T) {
	mig.CreateTestAndRunUseExternalProvider(t, basicTestCase(t), acc.ExternalProvidersOnlyAWS(), nil)
}
