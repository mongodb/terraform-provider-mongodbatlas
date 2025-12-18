package federatedsettingsidentityproviderapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_federated_settings_identity_provider_api.test"
const dataSourceName = "data.mongodbatlas_federated_settings_identity_provider_api.test"
const dataSourcePluralName = "data.mongodbatlas_federated_settings_identity_providers_api.test"

func TestAccFederatedSettingsIdentityProviderAPI_OIDCWorkforce(t *testing.T) {
	resource.ParallelTest(t, *basicOIDCWorkforceTestCase(t))
}

func basicOIDCWorkforceTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		associatedDomain     = os.Getenv("MONGODB_ATLAS_FEDERATED_SETTINGS_ASSOCIATED_DOMAIN")
		audience1            = "audience-workforce"
		audience2            = "audience-workforce-updated"
		description1         = "tf-acc-test"
		description2         = "tf-acc-test-updated"
		attrMapCheck         = map[string]string{
			"associated_domains.0":   associatedDomain,
			"audience":               audience1,
			"authorization_type":     "GROUP",
			"client_id":              "clientId",
			"description":            description1,
			"federation_settings_id": federationSettingsID,
			"groups_claim":           "groups",
			"issuer_uri":             "https://token.actions.githubusercontent.com",
			"protocol":               "OIDC",
			"requested_scopes.0":     "profiles",
			"user_claim":             "sub",
			"idp_type":               federatedsettingsidentityprovider.WORKFORCE,
		}
	)
	checks := []resource.TestCheckFunc{}
	checks = acc.AddAttrChecks(resourceName, checks, attrMapCheck)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMapCheck)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description1, audience1),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				Config: configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description2, audience2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// checkExistsManaged(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description2),
					resource.TestCheckResourceAttr(resourceName, "audience", audience2),
					resource.TestCheckResourceAttr(resourceName, "name", "OIDC-CRUD-test"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "OIDC-CRUD-test"),
				),
			},
			// {
			// 	Config:            configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description2, audience2),
			// 	ResourceName:      resourceName,
			// 	ImportStateIdFunc: importStateIDFuncManaged(resourceName),
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	}
}

func configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description, audience string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider" "test" {
        federation_settings_id 		= %[1]q
		associated_domains 			= [%[3]q]
		audience 					= %[2]q
		authorization_type			= "GROUP"
		client_id 					= "clientId"
		description 				= %[4]q
		groups_claim				= "groups"
		issuer_uri 					= "https://token.actions.githubusercontent.com"
		name 						= "OIDC-CRUD-test"
		protocol 					= "OIDC"
		requested_scopes 			= ["profiles"]
		user_claim 					= "sub"
		idp_type 					= "WORKFORCE"
	  }
	  
	  data "mongodbatlas_federated_settings_identity_provider" "test" {
		federation_settings_id = mongodbatlas_federated_settings_identity_provider.test.federation_settings_id
		identity_provider_id   = mongodbatlas_federated_settings_identity_provider.test.idp_id
	  }`, federationSettingsID, audience, associatedDomain, description)
}
