package mig

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

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
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func versionConstraint() string {
	return os.Getenv("MONGODB_ATLAS_LAST_VERSION")
}
