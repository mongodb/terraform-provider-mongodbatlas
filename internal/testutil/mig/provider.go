package mig

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func SkipIfVersionBelow(tb testing.TB, minVersion string) {
	tb.Helper()
	if !IsProviderVersionAtLeast(minVersion) {
		tb.Skipf("Skipping because version %s below %s", versionConstraint(), minVersion)
	}
}

func IsProviderVersionAtLeast(minVersion string) bool {
	vProvider, errProvider := version.NewVersion(versionConstraint())
	vMin, errMin := version.NewVersion(minVersion)
	return errProvider == nil && errMin == nil && vProvider.GreaterThanOrEqual(vMin)
}

func ExternalProviders() map[string]resource.ExternalProvider {
	return acc.ExternalProviders(versionConstraint())
}

func ExternalProvidersWithAWS() map[string]resource.ExternalProvider {
	return acc.ExternalProvidersWithAWS(versionConstraint())
}

func checkLastVersion(tb testing.TB) {
	tb.Helper()
	validateConflictingEnvVars()
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func versionConstraint() string {
	validateConflictingEnvVars()
	return os.Getenv("MONGODB_ATLAS_LAST_VERSION")
}

func validateConflictingEnvVars() {
	if os.Getenv("TF_CLI_CONFIG_FILE") != "" {
		log.Fatal("found `TF_CLI_CONFIG_FILE` in env-var when running migration tests, this might override the terraform provider for MONGODB_ATLAS_LAST_VERSION and instead use local latest build!")
	}
}
