package federatedsettingsidentityprovider_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	authTypeGroup        = "GROUP"
	authTypeUser         = "USER"
	resourceName         = "mongodbatlas_federated_settings_identity_provider.test"
	dataSourceName       = "data.mongodbatlas_federated_settings_identity_provider.test"
	dataSourcePluralName = "data.mongodbatlas_federated_settings_identity_providers.test"
)

func TestAccFederatedSettingsIdentityProvider_createError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configSAMLBasic("not-used", "not-used", "not-used", "not-used"),
				ExpectError: regexp.MustCompile("create is only supported by OIDC, SAML must be imported"),
			},
		},
	})
}

func TestAccFederatedSettingsIdentityProviderRS_basic(t *testing.T) {
	// SAML IdP can be deleted but not created through the API. If this test is run the resource will be deleted and will have to be created through the UI
	acc.SkipTestForCI(t)
	resource.ParallelTest(t, *basicSAMLTestCase(t))
}

func TestAccFederatedSettingsIdentityProvider_OIDCWorkforce(t *testing.T) {
	resource.ParallelTest(t, *basicOIDCWorkforceTestCase(t))
}

func TestAccFederatedSettingsIdentityProvider_OIDCWorkload(t *testing.T) {
	resource.ParallelTest(t, *basicOIDCWorkloadTestCase(t))
}

func basicSAMLTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID                = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
		ssoURL               = os.Getenv("MONGODB_ATLAS_FEDERATED_SSO_URL")
		issuerURI            = os.Getenv("MONGODB_ATLAS_FEDERATED_ISSUER_URI")
		associatedDomain     = os.Getenv("MONGODB_ATLAS_FEDERATED_SETTINGS_ASSOCIATED_DOMAIN")
		config               = configSAMLBasic(federationSettingsID, ssoURL, issuerURI, associatedDomain)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       resourceName,
				ImportStateIdFunc:  importStateIDFunc(federationSettingsID, idpID),
				ImportState:        true,
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "name", "SAML-test"),
				),
			},
			{
				Config:            config,
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(federationSettingsID, idpID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
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
	checks := []resource.TestCheckFunc{checkExistsManaged(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attrMapCheck)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMapCheck)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description1, audience1),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description2, audience2),
				Check: resource.ComposeTestCheckFunc(
					checkExistsManaged(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description2),
					resource.TestCheckResourceAttr(resourceName, "audience", audience2),
					resource.TestCheckResourceAttr(resourceName, "name", "OIDC-CRUD-test"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "OIDC-CRUD-test"),
				),
			},
			{
				Config:            configOIDCWorkforceBasic(federationSettingsID, associatedDomain, description2, audience2),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFuncManaged(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
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
		checkExistsManaged(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", "OIDC-workload-CRUD"),
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
			{
				Config:            configOIDCWorkloadBasic(federationSettingsID, description2, audience2, authTypeUser),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFuncManaged(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func checkExists(resourceName, idpID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.ConnV2().FederatedAuthenticationApi.GetIdentityProvider(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			idpID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("identity provider (%s) does not exist", idpID)
	}
}

func checkExistsManaged(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		federationSettingsID, idpID, err := readIDsFromState(s, resourceName)
		if err != nil {
			return err
		}
		_, _, err = acc.ConnV2().FederatedAuthenticationApi.GetIdentityProvider(context.Background(),
			federationSettingsID,
			idpID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("identity provider (%s) does not exist", idpID)
	}
}

func readIDsFromState(s *terraform.State, resourceName string) (federationSettingsID, idpID string, err error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", "", fmt.Errorf("no ID is set")
	}
	id := rs.Primary.ID
	if id == "" {
		return "", "", fmt.Errorf("ID is empty")
	}
	federationSettingsID, idpID = federatedsettingsidentityprovider.DecodeIDs(id)
	return federationSettingsID, idpID, nil
}

func importStateIDFunc(federationSettingsID, idpID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		return fmt.Sprintf("%s-%s", federationSettingsID, idpID), nil
	}
}

func importStateIDFuncManaged(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		federationSettingsID, idpID, err := readIDsFromState(s, resourceName)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s-%s", federationSettingsID, idpID), nil
	}
}

func configSAMLBasic(federationSettingsID, ssoURL, issuerURI, associatedDomain string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider" "test" {
        federation_settings_id 		= %[1]q
        name 						= "SAML-test"
        associated_domains     		= [%[4]q]
        sso_debug_enabled 			= true
        status 						= "ACTIVE"
        sso_url 					= %[2]q
        issuer_uri 					= %[3]q
        request_binding 			= "HTTP-POST"
        response_signature_algorithm = "SHA-256"
	  }`, federationSettingsID, ssoURL, issuerURI, associatedDomain)
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

func configOIDCWorkloadBasic(federationSettingsID, description, audience, authorizationType string) string {
	groupsClaimRaw := `"groups"`
	if authorizationType == authTypeUser {
		groupsClaimRaw = `null`
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider" "test" {
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
	  
	  data "mongodbatlas_federated_settings_identity_provider" "test" {
		federation_settings_id = mongodbatlas_federated_settings_identity_provider.test.federation_settings_id
		identity_provider_id   = mongodbatlas_federated_settings_identity_provider.test.idp_id
	  }
	  data "mongodbatlas_federated_settings_identity_providers" "test" {
		federation_settings_id 	= mongodbatlas_federated_settings_identity_provider.test.federation_settings_id
		idp_types 				= [%[4]q]
		protocols 				= [%[7]q]
		depends_on 				= [mongodbatlas_federated_settings_identity_provider.test]
	  }
	  `, federationSettingsID, audience, description, federatedsettingsidentityprovider.WORKLOAD, authorizationType, groupsClaimRaw, federatedsettingsidentityprovider.OIDC)
}
