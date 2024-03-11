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

// SkipInUnitTest allows skipping a test entirely when TF_ACC=1 is not defined.
// This can be useful for acceptance tests that define logic prior to resource.Test/resource.ParallelTest functions as this code would always be run.
func SkipInUnitTest(tb testing.TB) {
	tb.Helper()
	if os.Getenv("TF_ACC") == "" {
		tb.Skip()
	}
}
