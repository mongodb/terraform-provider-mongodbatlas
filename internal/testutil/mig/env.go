package mig

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func PreCheckBasic(tb testing.TB) {
	acc.PreCheckBasic(tb)
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func PreCheck(tb testing.TB) {
	acc.PreCheck(tb)
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

func PreCheckBasicOwnerID(tb testing.TB) {
	PreCheckBasic(tb)
	if os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PROJECT_OWNER_ID` must be set ")
	}
}

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
