package acc

import (
	"fmt"
	"os"

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

// configProvider creates a new provider with credentials explicit in config.
//
// This can be used when you want credentials different from the default env-vars.
func configProvider(publicKey, privateKey, baseURL string) string {
	return fmt.Sprintf(`
provider %[1]q {
	public_key = %[2]q
	private_key = %[3]q
	base_url = %[4]q
}
`, ProviderNameMongoDBAtlas, publicKey, privateKey, baseURL)
}

// ConfigGovProvider creates provider using MONGODB_ATLAS_GOV_* env vars.
//
// Remember to use PreCheckGovBasic when using this.
func ConfigGovProvider() string {
	return configProvider(os.Getenv("MONGODB_ATLAS_GOV_PUBLIC_KEY"), os.Getenv("MONGODB_ATLAS_GOV_PRIVATE_KEY"), os.Getenv("MONGODB_ATLAS_GOV_BASE_URL"))
}

func ConfigRPProvider() string {
	return configProvider(os.Getenv("MONGODB_ATLAS_RP_PUBLIC_KEY"), os.Getenv("MONGODB_ATLAS_RP_PRIVATE_KEY"), os.Getenv("MONGODB_ATLAS_BASE_URL"))
}
