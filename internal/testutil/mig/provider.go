package mig

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func IsProviderVersionAtLeast(minVersion string) bool {
	vProvider, errProvider := version.NewVersion(versionConstraint())
	vMin, errMin := version.NewVersion(minVersion)
	return errProvider == nil && errMin == nil && vProvider.GreaterThanOrEqual(vMin)
}

func isProviderVersionGreaterThan(minVersion string) bool {
	vProvider, errProvider := version.NewVersion(versionConstraint())
	vMin, errMin := version.NewVersion(minVersion)
	return errProvider == nil && errMin == nil && vProvider.GreaterThan(vMin)
}

func SkipIfProviderVersionGreaterThan(tb testing.TB, maxVersion string) {
	tb.Helper()
	if isProviderVersionGreaterThan(maxVersion) {
		tb.Skipf("skipping test for provider version %s, maxVersion : %s", versionConstraint(), maxVersion)
	}
}

func ExternalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": *providerAtlas(),
	}
}

func ExternalProvidersOnlyAWS() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"aws": *providerAWS(),
	}
}

func ExternalProvidersWithAWS() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": *providerAtlas(),
		"aws":          *providerAWS(),
	}
}

func providerAtlas() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: versionConstraint(),
		Source:            "mongodb/mongodbatlas",
	}
}

func providerAWS() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: "5.1.0",
		Source:            "hashicorp/aws",
	}
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
