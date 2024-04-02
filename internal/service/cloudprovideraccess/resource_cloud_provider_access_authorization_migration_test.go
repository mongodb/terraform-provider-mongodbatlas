package cloudprovideraccess_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudProviderAccessAuthorizationAWS_basic(t *testing.T) {
	mig.CreateTestAndRunUseExternalProvider(t, basicAuthorizationTestCase(t), mig.ExternalProvidersWithAWS(), acc.ExternalProvidersOnlyAWS())
}
