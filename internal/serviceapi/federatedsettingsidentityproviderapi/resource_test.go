package federatedsettingsidentityproviderapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	oidcProtocol         = "OIDC"
	samlProtocol         = "SAML"
	authTypeUser         = "USER"
	authTypeGroup        = "GROUP"
	resourceName         = "mongodbatlas_federated_settings_identity_provider_api.test"
	dataSourceName       = "data.mongodbatlas_federated_settings_identity_provider_api.test"
	dataSourcePluralName = "data.mongodbatlas_federated_settings_identity_providers_api.test"
)

func TestAccFederatedSettingsIdentityProviderAPI_OIDCWorkforce(t *testing.T) {
	resource.Test(t, *basicOIDCWorkforceTestCase(t))
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
	checks = acc.AddAttrChecks("data.mongodbatlas_federated_settings_identity_provider_api.test", checks, attrMapCheck)

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
					resource.TestCheckResourceAttr("data.mongodbatlas_federated_settings_identity_provider_api.test", "display_name", "OIDC-CRUD-test"),
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

func TestAccFederatedSettingsIdentityProvider_OIDCWorkload(t *testing.T) {
	resource.ParallelTest(t, *basicOIDCWorkloadTestCase(t))
}

func basicOIDCWorkloadTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		audience1            = "audience"
		audience2            = "audience-updated"
		description1         = "tf-acc-test"
		description2         = "tf-acc-test-updated"
		attrMapCheckGroup    = map[string]string{
			"audience":               audience1,
			"authorization_type":     authTypeGroup,
			"description":            description1,
			"federation_settings_id": federationSettingsID,
			"groups_claim":           "groups",
			"issuer_uri":             "https://token.actions.githubusercontent.com",
			"idp_type":               federatedsettingsidentityprovider.WORKLOAD,
			"protocol":               "OIDC",
		}
		attrMapCheckUser = map[string]string{
			"audience":               audience2,
			"authorization_type":     authTypeUser,
			"description":            description2,
			"federation_settings_id": federationSettingsID,
			"issuer_uri":             "https://token.actions.githubusercontent.com",
			"idp_type":               federatedsettingsidentityprovider.WORKLOAD,
			"protocol":               "OIDC",
			"user_claim":             "sub",
		}
	)
	nameChecks := []resource.TestCheckFunc{
		// checkExistsManaged(resourceName),
		resource.TestCheckResourceAttr(resourceName, "display_name", "OIDC-workload-CRUD"),
		resource.TestCheckResourceAttr(dataSourceName, "display_name", "OIDC-workload-CRUD"),
		resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.display_name", "OIDC-workload-CRUD"),
		resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", "1"),
	}
	checks := acc.AddAttrChecks(resourceName, nameChecks, attrMapCheckGroup)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMapCheckGroup)
	checks = acc.AddAttrChecksPrefix(dataSourcePluralName, checks, attrMapCheckGroup, "results.0", "federation_settings_id")

	checks2 := acc.AddAttrChecks(resourceName, nameChecks, attrMapCheckUser)
	checks2 = acc.AddAttrChecks(dataSourceName, checks2, attrMapCheckUser)
	checks2 = acc.AddAttrChecksPrefix(dataSourcePluralName, checks2, attrMapCheckUser, "results.0", "federation_settings_id")

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOIDCWorkloadBasic(federationSettingsID, description1, audience1, authTypeGroup),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				Config: configOIDCWorkloadBasic(federationSettingsID, description2, audience2, authTypeUser),
				Check:  resource.ComposeAggregateTestCheckFunc(checks2...),
			},
			// {
			// 	Config:            configOIDCWorkloadBasic(federationSettingsID, description2, audience2, authTypeUser),
			// 	ResourceName:      resourceName,
			// 	ImportStateIdFunc: importStateIDFuncManaged(resourceName),
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	}
}

func configOIDCWorkloadBasic(federationSettingsID, description, audience, authorizationType string) string {
	groupsClaimRaw := `"groups"`
	if authorizationType == authTypeUser {
		groupsClaimRaw = `null`
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider_api" "test" {
        federation_settings_id 		= %[1]q
		audience 					= %[2]q
		authorization_type			= %[5]q
		description 				= %[3]q
		issuer_uri 					= "https://token.actions.githubusercontent.com"
		idp_type 					= %[4]q
		name 						= "OIDC-workload-CRUD"
		protocol 					= "OIDC"
		groups_claim				= %[6]s
		user_claim 					= "sub"
	  }
	  
	  data "mongodbatlas_federated_settings_identity_provider_api" "test" {
		federation_settings_id = mongodbatlas_federated_settings_identity_provider.test.federation_settings_id
		identity_provider_id   = mongodbatlas_federated_settings_identity_provider_api.test.idp_id
	  }
	  data "mongodbatlas_federated_settings_identity_providers_api" "test" {
		federation_settings_id 	= mongodbatlas_federated_settings_identity_provider_api.test.federation_settings_id
		idp_type 				= [%[4]q]
		protocol 				= [%[7]q]
		depends_on 				= [mongodbatlas_federated_settings_identity_provider_api.test]
	  }
	  `, federationSettingsID, audience, description, federatedsettingsidentityprovider.WORKLOAD, authorizationType, groupsClaimRaw, federatedsettingsidentityprovider.OIDC)
}

func configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description, audience string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider_api" "test" {
        federation_settings_id 		= %[1]q
		associated_domains 			= [%[3]q]
		audience 					= %[2]q
		authorization_type			= "GROUP"
		client_id 					= "clientId"
		description 				= %[4]q
		groups_claim				= "groups"
		issuer_uri 					= "https://token.actions.githubusercontent.com"
		display_name 						= "OIDC-CRUD-test"
		protocol 					= "OIDC"
		requested_scopes 			= ["profiles"]
		user_claim 					= "sub"
		idp_type 					= "WORKFORCE"
	  }
	  
	  data "mongodbatlas_federated_settings_identity_provider_api" "test" {
		federation_settings_id = mongodbatlas_federated_settings_identity_provider_api.test.federation_settings_id
		id                     = mongodbatlas_federated_settings_identity_provider_api.test.id
	  }`, federationSettingsID, audience, associatedDomain, description)
}

func TestAccFederatedSettingsIdentityProvidersDS_basic(t *testing.T) {
	var (
		dataSourceName      = "data.mongodbatlas_federated_settings_identity_providers_api.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr(federatedsettingsidentityprovider.WORKFORCE), []string{oidcProtocol, samlProtocol}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "2"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr(federatedsettingsidentityprovider.WORKFORCE), []string{samlProtocol}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "SAML-test"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr(federatedsettingsidentityprovider.WORKFORCE), []string{oidcProtocol}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "OIDC-test"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr(federatedsettingsidentityprovider.WORKFORCE), []string{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "SAML-test"), // if no protocol is specified, it defaults to SAML
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr(federatedsettingsidentityprovider.WORKLOAD), []string{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "0"),
				),
			},
		},
	})
}

func configPluralDS(federatedSettingsID string, idpType *string, protocols []string) string {
	var protocolString string
	if len(protocols) > 1 {
		protocolString = fmt.Sprintf(`protocol = [%[1]q, %[2]q]`, protocols[0], protocols[1])
	} else if len(protocols) > 0 {
		protocolString = fmt.Sprintf(`protocol = [%[1]q]`, protocols[0])
	}
	var idpTypeString string
	if idpType != nil {
		idpTypeString = fmt.Sprintf(`idp_type = [%[1]q]`, *idpType)
	}

	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_identity_providers_api" "test" {
			federation_settings_id = "%[1]s"
			%[2]s
			%[3]s
		}
`, federatedSettingsID, protocolString, idpTypeString)
}
