package acc

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const AwsProviderVersion = "5.1.0"
const azurermProviderVersion = "4.60.0"
const azapiProviderVersion = "1.15.0"
const confluentProviderVersion = "2.12.0"
const googleProviderVersion = "7.0.0"

// Cloud Tag Policy: resources created by acceptance tests must carry
// mongodb-owner and mongodb-env tags
const (
	tagOwnerEnvVar = "MONGODB_TAG_OWNER"
	tagEnvEnvVar   = "MONGODB_TAG_ENV"
)

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

func ExternalProvidersWithConfluent(versionAtlasProvider string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"mongodbatlas": *providerAtlas(versionAtlasProvider),
		"confluent":    *providerConfluent(),
	}
}

func ExternalProvidersOnlyAWS() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"aws": *providerAWS(),
	}
}

func ExternalProvidersOnlyAzapi() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azapi": *providerAzapi(),
	}
}

func ExternalProvidersOnlyAzurerm() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azurerm": *providerAzurerm(),
	}
}

func ExternalProvidersOnlyConfluent() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"confluent": *providerConfluent(),
	}
}

func ExternalProvidersOnlyGoogle() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"google": *providerGoogle(),
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

func providerAzapi() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: azapiProviderVersion,
		Source:            "Azure/azapi",
	}
}

func providerAzurerm() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: azurermProviderVersion,
		Source:            "hashicorp/azurerm",
	}
}

func providerConfluent() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: confluentProviderVersion,
		Source:            "confluentinc/confluent",
	}
}

func providerGoogle() *resource.ExternalProvider {
	return &resource.ExternalProvider{
		VersionConstraint: googleProviderVersion,
		Source:            "hashicorp/google",
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

func ConfigOrgMemberProvider() string {
	return configProvider(os.Getenv("MONGODB_ATLAS_PUBLIC_KEY_READ_ONLY"), os.Getenv("MONGODB_ATLAS_PRIVATE_KEY_READ_ONLY"), os.Getenv("MONGODB_ATLAS_BASE_URL"))
}

// configAzapiProvider creates a new azure/azapi provider with credentials explicit in config.
// This will authorize the provider for a client
func ConfigAzapiProvider(subscriptionID, clientID, clientSecret, tenantID string) string {
	return fmt.Sprintf(`
		provider "azapi" {
			subscription_id = %[1]q
			client_id       = %[2]q
			client_secret   = %[3]q
			tenant_id       = %[4]q
		}
	`, subscriptionID, clientID, clientSecret, tenantID)
}

func ConfigAzurermProvider(subscriptionID, clientID, clientSecret, tenantID string) string {
	return fmt.Sprintf(`
		provider "azurerm" {
			features {}
			subscription_id = %[1]q
			client_id       = %[2]q
			client_secret   = %[3]q
			tenant_id       = %[4]q
		}
	`, subscriptionID, clientID, clientSecret, tenantID)
}

// ConfigTags returns the Cloud Tag Policy tags as a resource-level tags attribute.
// azurerm does not support provider-level default_tags, so Azure resources set tags individually: https://github.com/hashicorp/terraform-provider-azurerm/issues/13776.
func ConfigTags() string {
	return fmt.Sprintf(`tags = {
		"mongodb-owner" = %[1]q
		"mongodb-env"   = %[2]q
	}`, tagValue(tagOwnerEnvVar), tagValue(tagEnvEnvVar))
}

func ConfigDefaultTags() string {
	return fmt.Sprintf(`
		default_tags {
			tags = {
				"mongodb-owner" = %[1]q
				"mongodb-env"   = %[2]q
			}
		}
	`, tagValue(tagOwnerEnvVar), tagValue(tagEnvEnvVar))
}

func ConfigAWSProviderWithTags() string {
	return fmt.Sprintf(`
		provider "aws" {
			%[1]s
		}
	`, ConfigDefaultTags())
}

// tagValue reads the tag value from envVar. Without TF_ACC, configs are never applied,
// so a placeholder value is returned instead of failing.
func tagValue(envVar string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	if os.Getenv("TF_ACC") == "" {
		return "unused"
	}
	panic(fmt.Sprintf("%s must be set", envVar))
}

func ConfigConfluentProvider() string {
	return `provider "confluent" {}`
}

func ConfigGoogleProvider(projectID string) string {
	return fmt.Sprintf(`
		provider "google" {
			project = %[1]q
		}
	`, projectID)
}
