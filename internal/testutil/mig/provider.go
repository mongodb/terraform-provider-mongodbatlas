package mig

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ExternalProviders(tb testing.TB) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": {
			VersionConstraint: versionConstraint(),
			Source:            "mongodb/mongodbatlas",
		},
	}
}

func ExternalProvidersWithAWS(tb testing.TB, awsVersion string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": {
			VersionConstraint: versionConstraint(),
			Source:            "mongodb/mongodbatlas",
		},
		"aws": {
			VersionConstraint: awsVersion,
			Source:            "hashicorp/aws",
		},
	}
}

func IsProviderVersionAtLeast(minVersion string) bool {
	vProvider, errProvider := version.NewVersion(versionConstraint())
	vMin, errMin := version.NewVersion(minVersion)
	return errProvider == nil && errMin == nil && vProvider.GreaterThanOrEqual(vMin)
}

func checkLastVersion(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func versionConstraint() string {
	return os.Getenv("MONGODB_ATLAS_LAST_VERSION")
}
