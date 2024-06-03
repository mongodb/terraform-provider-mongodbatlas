package federatedsettingsorgconfig_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrg_createError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configBasic("not-used", "not-used", "not-used", "not-used", false, false),
				ExpectError: regexp.MustCompile("this resource must be imported"),
			},
		},
	})
}

func TestAccFederatedSettingsOrg_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // will delete the MONGODB_ATLAS_FEDERATED_ORG_ID on finish, no workaround: https://github.com/hashicorp/terraform-plugin-testing/issues/85

	var (
		resourceName         = "mongodbatlas_federated_settings_org_config.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		idpID                = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
		associatedDomain     = os.Getenv("MONGODB_ATLAS_FEDERATED_SETTINGS_ASSOCIATED_DOMAIN")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettingsIdentityProvider(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:             configBasic(federationSettingsID, orgID, idpID, associatedDomain, false, false),
				ResourceName:       resourceName,
				ImportStateIdFunc:  importStateIDFunc(federationSettingsID, orgID),
				ImportState:        true,
				ImportStateVerify:  false,
				ImportStatePersist: true, // ensure update will be tested in the next step
			},
			{
				Config: configBasic(federationSettingsID, orgID, idpID, associatedDomain, true, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "domain_restriction_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "domain_allow_list.0", associatedDomain),
					resource.TestCheckResourceAttr(resourceName, "data_access_identity_provider_ids.#", "1"),
				),
			},
			{
				Config: configBasic(federationSettingsID, orgID, idpID, associatedDomain, true, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "data_access_identity_provider_ids.#", "0"),
				),
			},
			{
				Config:            configBasic(federationSettingsID, orgID, idpID, associatedDomain, false, false),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(federationSettingsID, orgID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.ConnV2().FederatedAuthenticationApi.GetConnectedOrgConfig(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			rs.Primary.Attributes["org_id"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("connected org  (%s) does not exist", rs.Primary.Attributes["org_id"])
	}
}

func importStateIDFunc(federationSettingsID, orgID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ID := conversion.EncodeStateID(map[string]string{
			"federation_settings_id": federationSettingsID,
			"org_id":                 orgID,
		})

		ids := conversion.DecodeStateID(ID)
		return fmt.Sprintf("%s-%s", ids["federation_settings_id"], ids["org_id"]), nil
	}
}

func configBasic(federationSettingsID, orgID, identityProviderID, associatedDomain string, withWorkload, useWorkload bool) string {
	var workload string
	var workloadValue = "[]"
	if withWorkload {
		workload = fmt.Sprintf(`
		resource "mongodbatlas_federated_settings_identity_provider" "oidc_workload" {
			federation_settings_id 		= %[1]q
			audience 					= "some-aud"
			authorization_type			= "GROUP"
			description 				= "oidc-for-testing-org-update"
			issuer_uri 					= "https://gitlab.com"
			idp_type 					= %[2]q
			name 						= "OIDC-workload-org-update"
			protocol 					= %[3]q
			groups_claim				= "groups"
			user_claim 					= "sub"
		  }
		`, federationSettingsID, federatedsettingsidentityprovider.WORKLOAD, federatedsettingsidentityprovider.OIDC)
	}
	if useWorkload {
		workloadValue = fmt.Sprintf("[%1s]", "mongodbatlas_federated_settings_identity_provider.oidc_workload.idp_id")
	}
	return fmt.Sprintf(`
	%[5]s
	resource "mongodbatlas_federated_settings_org_config" "test" {
		federation_settings_id     			= %[1]q
		org_id                     			= %[2]q
		domain_restriction_enabled 			= false
		domain_allow_list          			= [%[4]q]
		identity_provider_id       			= %[3]q
		data_access_identity_provider_ids 	= %[6]s
	  }`, federationSettingsID, orgID, identityProviderID, associatedDomain, workload, workloadValue)
}
