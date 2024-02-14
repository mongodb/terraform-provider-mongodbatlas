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
func SkipInUnitTest(tb testing.TB) {
	tb.Helper()
	if os.Getenv("TF_ACC") == "" {
		tb.Skip()
	}
}
