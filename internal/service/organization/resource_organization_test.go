package organization_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName   = "mongodbatlas_organization.test"
	pluralDSName   = "data.mongodbatlas_organizations.test"
	datasourceName = "data.mongodbatlas_organization.test"

	defaultSettings = &admin.OrganizationSettings{
		ApiAccessListRequired:   conversion.Pointer(false),
		MultiFactorAuthRequired: conversion.Pointer(false),
		RestrictEmployeeAccess:  conversion.Pointer(false),
		GenAIFeaturesEnabled:    conversion.Pointer(true),
	}
)

func TestAccConfigRSOrganization_Basic(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		orgOwnerID  = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name        = acc.RandomName()
		updatedName = acc.RandomName()
		description = "test Key for Acceptance tests"
		roleName    = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgOwnerID, name, description, roleName, false, nil),
				Check: checkAggr(orgOwnerID, name, description, defaultSettings,
					resource.TestCheckResourceAttr(resourceName, "skip_default_alerts_settings", "true")),
			},
			{
				Config: configBasic(orgOwnerID, updatedName, description, roleName, true, conversion.Pointer(false)),
				Check: checkAggr(orgOwnerID, updatedName, description, defaultSettings,
					resource.TestCheckResourceAttr(resourceName, "skip_default_alerts_settings", "false")),
			},
			{
				Config: configBasic(orgOwnerID, updatedName, description, roleName, true, conversion.Pointer(true)),
				Check: checkAggr(orgOwnerID, updatedName, description, defaultSettings,
					resource.TestCheckResourceAttr(resourceName, "skip_default_alerts_settings", "true")),
			},
		},
	})
}

func TestAccConfigRSOrganization_BasicAccess(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		orgOwnerID            = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name                  = acc.RandomName()
		description           = "test Key for Acceptance tests"
		roleName              = "ORG_BILLING_ADMIN"
		roleNameCorrectAccess = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(orgOwnerID, name, description, roleName, false, nil),
				ExpectError: regexp.MustCompile("API Key must have the ORG_OWNER role"),
			},
			{
				Config: configBasic(orgOwnerID, name, description, roleNameCorrectAccess, true, conversion.Pointer(false)),
				Check: checkAggr(orgOwnerID, name, description, defaultSettings,
					resource.TestCheckResourceAttr(resourceName, "skip_default_alerts_settings", "false")),
			},
		},
	})
}

func TestAccConfigRSOrganization_Settings(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		orgOwnerID  = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name        = acc.RandomName()
		nameUpdated = "org-name-updated"
		description = "test Key for Acceptance tests"
		roleName    = "ORG_OWNER"

		settingsConfig = &admin.OrganizationSettings{
			ApiAccessListRequired:   conversion.Pointer(false),
			MultiFactorAuthRequired: conversion.Pointer(true),
			GenAIFeaturesEnabled:    conversion.Pointer(false),
			SecurityContact:         conversion.StringPtr("test@mongodb.com"),
		}

		settingsConfigUpdated = &admin.OrganizationSettings{
			ApiAccessListRequired:   conversion.Pointer(false),
			MultiFactorAuthRequired: conversion.Pointer(true),
			RestrictEmployeeAccess:  conversion.Pointer(false),
			GenAIFeaturesEnabled:    conversion.Pointer(true),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithSettings(orgOwnerID, name, description, roleName, settingsConfig),
				Check:  checkAggr(orgOwnerID, name, description, settingsConfig),
			},
			{
				Config: configWithSettings(orgOwnerID, name, description, roleName, settingsConfigUpdated),
				Check:  checkAggr(orgOwnerID, name, description, settingsConfigUpdated),
			},
			{
				Config: configBasic(orgOwnerID, nameUpdated, description, roleName, false, nil),
				Check: checkAggr(orgOwnerID, nameUpdated, description, settingsConfigUpdated,
					resource.TestCheckResourceAttr(resourceName, "skip_default_alerts_settings", "true")),
			},
		},
	})
}

func TestAccConfigRSOrganizationCreate_Errors(t *testing.T) {
	var (
		roleName    = "ORG_OWNER"
		unknownUser = "65def6160f722a1507105aaa"
	)
	acc.SkipTestForCI(t) // test will fail in CI since API_KEY_MUST_BE_ASSOCIATED_WITH_PAYING_ORG is returned
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(unknownUser, acc.RandomName(), "should fail since user is not found", roleName, false, nil),
				ExpectError: regexp.MustCompile(`USER_NOT_FOUND`),
			},
		},
	})
}

func TestAccConfigDSOrganization_noAccessShouldFail(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configWithPluralDS("555") + acc.ConfigOrgMemberProvider(),
				ExpectError: regexp.MustCompile("error getting organizations information:"),
			},
		},
	})
}

func TestAccConfigDSOrganization_basic(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPluralDS(orgID),
				Check: checkAggrDS(resource.TestCheckResourceAttr(datasourceName, "gen_ai_features_enabled", "true"),
					resource.TestCheckResourceAttr(pluralDSName, "results.0.gen_ai_features_enabled", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "users.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "users.0.id")),
			},
		},
	})
}

func TestAccConfigDSOrganization_users(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPluralDS(orgID),
				Check: checkAggrDS(
					resource.TestCheckResourceAttrWith(datasourceName, "users.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrSet(datasourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "users.0.roles.0.org_roles.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "users.0.roles.0.project_role_assignments.#"),
					resource.TestCheckResourceAttrWith(datasourceName, "users.0.username", acc.IsUsername()),
					resource.TestCheckResourceAttrWith(datasourceName, "users.0.last_auth", acc.IsTimestamp()),
					resource.TestCheckResourceAttrWith(datasourceName, "users.0.created_at", acc.IsTimestamp()),

					resource.TestCheckResourceAttrWith(pluralDSName, "results.0.users.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.0.users.0.id"),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.0.users.0.roles.0.org_roles.#"),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.0.users.0.roles.0.project_role_assignments.#"),
					resource.TestCheckResourceAttrWith(pluralDSName, "results.0.users.0.username", acc.IsUsername()),
					resource.TestCheckResourceAttrWith(pluralDSName, "results.0.users.0.last_auth", acc.IsTimestamp()),
				),
			},
		},
	})
}

func TestAccConfigDSOrganizations_withPagination(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPagination(1, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(pluralDSName, "results.#"),
				),
			},
		},
	})
}

func TestAccConfigRSOrganization_import(t *testing.T) {
	acc.SkipInUnitTest(t) // needed so OrganizationsApi is not called in unit tests
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	resp, _, _ := acc.ConnV2().OrganizationsApi.GetOrg(t.Context(), orgID).Execute()
	orgName := resp.GetName()
	require.NotEmpty(t, orgName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configImportSet(orgID, orgName), // Use import so a new organization is not created, the resource must exist in a step before import state is verified.
			},
			{
				ResourceName:                         resourceName,
				ImportStateId:                        orgID,
				ImportState:                          true, // Do the import check.
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "org_id",
			},
			{
				// Use removed block so the organization is not deleted.
				// Even if something goes wrong, the organization wouldn't be deleted if it has some projects, it would return ORG_NOT_EMPTY error.
				Config: acc.ConfigRemove(resourceName),
			},
		},
	})
}

func configWithPluralDS(orgID string) string {
	cfg := fmt.Sprintf(`
		
		data "mongodbatlas_organization" "test" {
			org_id = %[1]q
		}
	`, orgID)
	return fmt.Sprintf(`	
		data "mongodbatlas_organizations" "test" {
		}

		%s
	`, cfg)
}

func configWithPagination(pageNum, itemPage int) string {
	return fmt.Sprintf(`
		data "mongodbatlas_organizations" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, pageNum, itemPage)
}

func configImportSet(orgID, orgName string) string {
	return fmt.Sprintf(`
		import {
			id = %[1]q
			to = mongodbatlas_organization.test
		}

		resource "mongodbatlas_organization" "test" {
			name = %[2]q
			lifecycle {
   		 prevent_destroy = true
  		}
		}
	`, orgID, orgName)
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

		ids := conversion.DecodeStateID(rs.Primary.ID)
		conn, err := getTestClientWithNewOrgCreds(rs)
		if err != nil {
			return err
		}

		orgs, _, err := conn.OrganizationsApi.ListOrgs(context.Background()).Execute()
		if err == nil {
			for _, val := range orgs.GetResults() {
				if val.GetId() == ids["org_id"] {
					return nil
				}
			}
			return fmt.Errorf("Organization (%s) doesn't exist", ids["org_id"])
		}

		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		conn, err := getTestClientWithNewOrgCreds(rs)
		if err != nil {
			return err
		}

		orgs, _, err := conn.OrganizationsApi.ListOrgs(context.Background()).Execute()
		if err == nil {
			for _, val := range orgs.GetResults() {
				if val.GetId() == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["org_id"])
				}
			}
			return nil
		}
	}

	return nil
}

func configBasic(orgOwnerID, name, description, roleNames string, useSkipDefaultAlertSettings bool, skipDefaultAlertSettings *bool) string {
	skipDefaultAlertSettingsStr := ""

	if useSkipDefaultAlertSettings && skipDefaultAlertSettings != nil {
		skipDefaultAlertSettingsStr = fmt.Sprintf("skip_default_alerts_settings = %t", *skipDefaultAlertSettings)
	}

	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = %q
		name = %q
		description = %q
		role_names = [%q]
		
		%s
	  }
	`, orgOwnerID, name, description, roleNames, skipDefaultAlertSettingsStr)
}

func configWithSettings(orgOwnerID, name, description, roleNames string, settingsConfig *admin.OrganizationSettings) string {
	settingsStr := getSettingsConfig(settingsConfig)

	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = %q
		name = %q
		description = %q
		role_names = [%q]
		%s
	  }
	`, orgOwnerID, name, description, roleNames, settingsStr)
}

func getSettingsConfig(settings *admin.OrganizationSettings) string {
	var configs []string

	if settings.ApiAccessListRequired != nil {
		configs = append(configs, fmt.Sprintf("api_access_list_required = %t", *settings.ApiAccessListRequired))
	}
	if settings.MultiFactorAuthRequired != nil {
		configs = append(configs, fmt.Sprintf("multi_factor_auth_required = %t", *settings.MultiFactorAuthRequired))
	}
	if settings.RestrictEmployeeAccess != nil {
		configs = append(configs, fmt.Sprintf("restrict_employee_access = %t", *settings.RestrictEmployeeAccess))
	}
	if settings.GenAIFeaturesEnabled != nil {
		configs = append(configs, fmt.Sprintf("gen_ai_features_enabled = %t", *settings.GenAIFeaturesEnabled))
	}
	if settings.SecurityContact != nil {
		configs = append(configs, fmt.Sprintf("security_contact = %q", *settings.SecurityContact))
	}

	return strings.Join(configs, "\n")
}

// getTestClientWithNewOrgCreds creates a new Atlas client with credentials for the newly created organization which
// is required to call relevant API methods for the new organization, for example ListOrganizations requires that the requesting API
// key must have the Organization Member role. So we cannot invoke API methods on the new organization with credentials configured in the provider.
func getTestClientWithNewOrgCreds(rs *terraform.ResourceState) (*admin.APIClient, error) {
	if rs.Primary.Attributes["public_key"] == "" {
		return nil, fmt.Errorf("no public_key is set")
	}

	if rs.Primary.Attributes["private_key"] == "" {
		return nil, fmt.Errorf("no private_key is set")
	}

	cfg := config.Config{
		PublicKey:  rs.Primary.Attributes["public_key"],
		PrivateKey: rs.Primary.Attributes["private_key"],
		BaseURL:    acc.MongoDBClient.Config.BaseURL,
	}
	clients, _ := cfg.NewClient(context.Background())
	return clients.AtlasV2, nil
}

func TestValidateAPIKeyIsOrgOwner(t *testing.T) {
	tests := []struct {
		name    string
		roles   []string
		wantErr bool
	}{
		{
			name:    "Contains OrgOwner",
			roles:   []string{"ORG_MEMBER", "ORG_OWNER", "ORG_READ_ONLY"},
			wantErr: false,
		},
		{
			name:    "Does Not Contain OrgOwner",
			roles:   []string{"ORG_MEMBER", "READ_ONLY"},
			wantErr: true,
		},
		{
			name:    "Empty Roles",
			roles:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := organization.ValidateAPIKeyIsOrgOwner(tt.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAPIKeyIsOrgOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func checkAggr(orgOwnerID, name, description string, settings *admin.OrganizationSettings, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attributes := map[string]string{
		"name":                       name,
		"org_owner_id":               orgOwnerID,
		"description":                description,
		"api_access_list_required":   strconv.FormatBool(settings.GetApiAccessListRequired()),
		"multi_factor_auth_required": strconv.FormatBool(settings.GetMultiFactorAuthRequired()),
		"restrict_employee_access":   strconv.FormatBool(settings.GetRestrictEmployeeAccess()),
		"gen_ai_features_enabled":    strconv.FormatBool(settings.GetGenAIFeaturesEnabled()),
		"security_contact":           settings.GetSecurityContact(),
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "role_names.#")
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkAggrDS(extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc

	singularKeys := []string{
		"name",
		"id",
		"restrict_employee_access",
		"multi_factor_auth_required",
		"api_access_list_required",
		"skip_default_alerts_settings",
	}
	checks = acc.AddAttrSetChecks(datasourceName, checks, singularKeys...)

	pluralKeys := getPluralDSAttrKeys(singularKeys)
	checks = acc.AddAttrSetChecks(pluralDSName, checks, pluralKeys...)

	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func getPluralDSAttrKeys(singularKeys []string) []string {
	pluralKeys := []string{"results.#"}
	for _, key := range singularKeys {
		pluralKeys = append(pluralKeys, "results.0."+key)
	}
	return pluralKeys
}
