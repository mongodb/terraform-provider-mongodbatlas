package mig

import (
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ExternalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": {
			VersionConstraint: versionConstraint(),
			Source:            "mongodb/mongodbatlas",
		},
	}
}

func ExternalProvidersWithAWS(awsVersion string) map[string]resource.ExternalProvider {
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
