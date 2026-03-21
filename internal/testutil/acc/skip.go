package acc

import (
	"os"
	"strconv"
	"testing"
)

// SkipTestForCI is added to tests that cannot run as part of CI, e.g. in Github actions.
func SkipTestForCI(tb testing.TB) {
	tb.Helper()
	if InCI() {
		tb.Skip()
	}
}

func InCI() bool {
	val, _ := strconv.ParseBool(os.Getenv("CI"))
	return val
}

// SkipInUnitTest allows skipping a test entirely when TF_ACC=1 is not defined.
// This can be useful for acceptance tests that define logic prior to resource.Test/resource.ParallelTest functions as this code would always be run.
func SkipInUnitTest(tb testing.TB) {
	tb.Helper()
	if InUnitTest() {
		tb.Skip()
	}
}

func InUnitTest() bool {
	return os.Getenv("TF_ACC") == ""
}

func HasPAKCreds() bool {
	return os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") != "" || os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") != ""
}

func HasSACreds() bool {
	return os.Getenv("MONGODB_ATLAS_CLIENT_ID") != "" && os.Getenv("MONGODB_ATLAS_CLIENT_SECRET") != ""
}

func HasAccessToken() bool {
	return os.Getenv("MONGODB_ATLAS_ACCESS_TOKEN") != ""
}

func SkipInSA(tb testing.TB, description string) {
	tb.Helper()
	if HasSACreds() {
		tb.Skip(description)
	}
}

func SkipInPAK(tb testing.TB, description string) {
	tb.Helper()
	if HasPAKCreds() {
		tb.Skip(description)
	}
}

func SkipInAccessToken(tb testing.TB, description string) {
	tb.Helper()
	if HasAccessToken() {
		tb.Skip(description)
	}
}

// SkipIfNotSA skips the test if Service Account credentials are not configured.
func SkipIfNotSA(tb testing.TB) {
	tb.Helper()
	if !HasSACreds() {
		tb.Skip("MONGODB_ATLAS_CLIENT_ID and MONGODB_ATLAS_CLIENT_SECRET required")
	}
}
