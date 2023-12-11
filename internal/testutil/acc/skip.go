package acc

import (
	"os"
	"strings"
	"testing"
)

func SkipTest(tb testing.TB) {
	tb.Helper()
	if strings.EqualFold(os.Getenv("SKIP_TEST"), "true") {
		tb.Skip()
	}
}

// SkipIfTFAccNotDefined is added to acceptance tests were you do not want any preparation code executed if the resulting steps will not run.
// Keep in mind that if TF_ACC is empty, go still runs acceptance tests but terraform-plugin-testing does not execute the resulting steps.
func SkipIfTFAccNotDefined(tb testing.TB) {
	tb.Helper()
	if strings.EqualFold(os.Getenv("TF_ACC"), "") {
		tb.Skip()
	}
}

// SkipTestForCI is added to tests that cannot run as part of a CI
func SkipTestForCI(tb testing.TB) {
	tb.Helper()
	if strings.EqualFold(os.Getenv("CI"), "true") {
		tb.Skip()
	}
}

func SkipTestExtCred(tb testing.TB) {
	tb.Helper()
	if strings.EqualFold(os.Getenv("SKIP_TEST_EXTERNAL_CREDENTIALS"), "true") {
		tb.Skip()
	}
}
