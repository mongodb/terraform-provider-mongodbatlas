package mig

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ExternalProviders(tb testing.TB) map[string]resource.ExternalProvider {
	checkLastVersion(tb)
	return map[string]resource.ExternalProvider{
		"mongodbatlas": {
			VersionConstraint: versionConstraint(),
			Source:            "mongodb/mongodbatlas",
		},
	}
}

func ExternalProvidersWithAWS(tb testing.TB, awsVersion string) map[string]resource.ExternalProvider {
	checkLastVersion(tb)
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

func versionConstraint() string {
	return os.Getenv("MONGODB_ATLAS_LAST_VERSION")
}
