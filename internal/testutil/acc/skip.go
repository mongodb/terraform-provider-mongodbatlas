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
