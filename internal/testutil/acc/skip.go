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

// nonAtlasCLICredentialEnvVars are the environment variables that configure non-AtlasCLI
// credentials, including the aliases read by config.NewEnvVars and AWS Secrets Manager.
var nonAtlasCLICredentialEnvVars = []string{
	"MONGODB_ATLAS_PUBLIC_KEY", "MONGODB_ATLAS_PRIVATE_KEY",
	"MONGODB_ATLAS_PUBLIC_API_KEY", "MONGODB_ATLAS_PRIVATE_API_KEY",
	"MCLI_PUBLIC_API_KEY", "MCLI_PRIVATE_API_KEY",
	"MONGODB_ATLAS_CLIENT_ID", "MONGODB_ATLAS_CLIENT_SECRET",
	"MONGODB_ATLAS_ACCESS_TOKEN",
	"ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN",
}

// SkipIfNonAtlasCLICredentials skips the test when any supported non-AtlasCLI credential
// source is configured (PAK, Service Account, access token, or AWS Secrets Manager), so that
// the AtlasCLI fallback is the only authentication path exercised.
func SkipIfNonAtlasCLICredentials(tb testing.TB, description string) {
	tb.Helper()
	for _, envVar := range nonAtlasCLICredentialEnvVars {
		if os.Getenv(envVar) != "" {
			tb.Skipf("%s: %s is set", description, envVar)
		}
	}
}

// SkipIfNotSA skips the test if Service Account credentials are not configured.
func SkipIfNotSA(tb testing.TB) {
	tb.Helper()
	if !HasSACreds() {
		tb.Skip("MONGODB_ATLAS_CLIENT_ID and MONGODB_ATLAS_CLIENT_SECRET required")
	}
}
