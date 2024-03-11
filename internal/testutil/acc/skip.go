package acc

import (
	"os"
	"strings"
	"testing"
)

// SkipTestForCI is added to tests that cannot run as part of CI, e.g. in Github actions.
func SkipTestForCI(tb testing.TB) {
	tb.Helper()
	if strings.EqualFold(os.Getenv("CI"), "true") {
		tb.Skip()
	}
}

// SkipInUnitTest is rarely needed, it is used in acc and mig tests to make sure that they don't run in unit test mode.
// TF skipping test when TF_ACC=1 is not set is implemented inside resource.Test / ParallelTest functions.
// SkipInUnitTest allows to call functions in the test body that must not run in unit test mode, only in acc/mig mode.
// As an example it is used in ProjectIDGlobal so it can be called from the test methdod,
// or in TestAccConfigDSAtlasUser_ByUserID so it can call fetchUser.
func SkipInUnitTest(tb testing.TB) {
	tb.Helper()
	if os.Getenv("TF_ACC") == "" {
		tb.Skip()
	}
}
