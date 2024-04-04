package federatedsettingsidentityprovider_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigFederatedSettingsIdentityProviderRS_basic(t *testing.T) {
	acc.SkipTestForCI(t) // this resource can only be imported
	mig.CreateAndRunTest(t, basicTestCase(t))
}
