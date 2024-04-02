package acc

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const AwsProviderVersion = "5.1.0"

func ExternalProviders(versionAtlasProvider string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": *providerAtlas(versionAtlasProvider),
	}
}

func ExternalProvidersWithAWS(versionAtlasProvider string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": *providerAtlas(versionAtlasProvider),
		"aws":          *providerAWS(),
	}
}

func ExternalProvidersOnlyAWS() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"aws": *providerAWS(),
	}
}

func providerAtlas(versionAtlasProvider string) *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: versionAtlasProvider,
		Source:            "mongodb/mongodbatlas",
	}
}

func providerAWS() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: AwsProviderVersion,
		Source:            "hashicorp/aws",
	}
}
