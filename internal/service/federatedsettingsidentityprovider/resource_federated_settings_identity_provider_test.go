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

func TestAccFederatedSettingsIdentityProvider_createError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configSAMLBasic("not-used", "not-used", "not-used", "not-used"),
				ExpectError: regexp.MustCompile("this resource must be imported"),
			},
		},
	})
}

func TestAccFederatedSettingsIdentityProviderRS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicSAMLTestCase(t))
}

func TestAccFederatedSettingsIdentityProviderRS_OIDCWorkforce(t *testing.T) {
	resource.ParallelTest(t, *basicOIDCWorkforceTestCase(t))
}

func basicSAMLTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		resourceName         = "mongodbatlas_federated_settings_identity_provider.test"
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
		resourceName         = "mongodbatlas_federated_settings_identity_provider.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		associatedDomain     = os.Getenv("MONGODB_ATLAS_FEDERATED_SETTINGS_ASSOCIATED_DOMAIN")
		audience             = "audience"
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOIDCWorkforceBasic(federationSettingsID, associatedDomain, &audience),
				Check: resource.ComposeTestCheckFunc(
					checkExistsManaged(resourceName),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "name", "OIDC-CRUD-test"),
				),
			},
			{
				Config:            configOIDCWorkforceBasic(federationSettingsID, associatedDomain, &audience),
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

func configOIDCWorkforceBasic(federationSettingsID, associatedDomain string, audience *string) string {
	var audienceString string
	if audience != nil {
		audienceString = fmt.Sprintf(`audience = %[1]q`, *audience)
	}

	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider" "test" {
        federation_settings_id 		= %[1]q
		associated_domains 			= [%[3]q]
		authorization_type			= "GROUP"
		client_id 					= "clientId"
		description 				= "tf-acc-test"
		groups_claim				= "groups"
		issuer_uri 					= "https://token.actions.githubusercontent.com"
		protocol 					= "OIDC"
		requested_scopes 			= ["profiles"]
		user_claim 					= "sub"
        name 						= "OIDC-CRUD-test"

		%[2]s

	  }`, federationSettingsID, audienceString, associatedDomain)
}
