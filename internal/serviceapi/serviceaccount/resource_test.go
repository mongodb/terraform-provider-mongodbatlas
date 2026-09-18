package serviceaccount_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_service_account.test"
const dataSourceName = "data.mongodbatlas_service_account.test"
const dataSourcePluralName = "data.mongodbatlas_service_accounts.test"

func TestAccServiceAccount_basic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name1        = acc.RandomName()
		name2        = fmt.Sprintf("%s-updated", name1)
		description1 = "Acceptance Test SA"
		description2 = "Updated Description"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name1, description1, []string{"ORG_OWNER"}, new(24), false),
				Check:  checkBasic(true, false),
			},
			{
				Config: configBasic(orgID, name2, description2, []string{"ORG_READ_ONLY"}, new(24), false),
				Check:  checkBasic(false, false),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Neither attribute is populated by the API on read: secret_expires_after_hours is create-only and
				// without_initial_secret is request-only. Import state cannot reproduce them from config.
				ImportStateVerifyIgnore: []string{"secret_expires_after_hours", "without_initial_secret"},
			},
		},
	})
}

func TestAccServiceAccount_rolesOrdering(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, "Roles Update", []string{"ORG_BILLING_ADMIN", "ORG_READ_ONLY"}, new(24), false),
			},
			{
				// change order
				Config: configBasic(orgID, name, "Roles Update", []string{"ORG_READ_ONLY", "ORG_BILLING_ADMIN"}, new(24), false),
			},
		},
	})
}

func TestAccServiceAccount_createOnlyAttributes(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, "description", []string{"ORG_READ_ONLY"}, new(24), false),
				Check:  checkExists(resourceName),
			},
			{
				// secret_expires_after_hours is create-only and stored from config, so a change is rejected.
				Config:      configBasic(orgID, name, "description", []string{"ORG_READ_ONLY"}, new(48), false),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("secret_expires_after_hours cannot be updated"),
			},
			{
				Config:      configBasic("updated-org-id", name, "description", []string{"ORG_READ_ONLY"}, new(24), false),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("org_id cannot be updated"),
			},
			{
				// without_initial_secret is optional-only, so a config that omits it stores null. Setting it on
				// update is not rejected (validateCreateOnly skips the check when the state value is null). The
				// resulting plan is non-empty: the attribute changes from null to true and the API ignores the
				// PATCH field, so no server-side change happens.
				Config:             configBasic(orgID, name, "description", []string{"ORG_READ_ONLY"}, new(24), true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccServiceAccount_withoutInitialSecret(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, "Without initial secret", []string{"ORG_READ_ONLY"}, nil, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "secrets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "without_initial_secret", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "client_id"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"without_initial_secret"},
			},
		},
	})
}

func TestAccServiceAccount_pluralDSIncludeSystemManaged(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				// The MCP config creates a system-managed Service Account (egress_client_id),
				// Using an output to verify that it's returned in the plural DS.
				Config: fmt.Sprintf(`
					resource "mongodbatlas_mcp_config" "test" {
						org_id          = %[1]q
						mcp_config_name = %[2]q
						roles           = ["ORG_READ_ONLY"]
					}

					data "mongodbatlas_service_accounts" "test" {
						org_id                = %[1]q
						include_system_managed = true
						depends_on            = [mongodbatlas_mcp_config.test]
					}

					output "includes_egress_service_account" {
						value = contains([for sa in data.mongodbatlas_service_accounts.test.results : sa.client_id], mongodbatlas_mcp_config.test.egress_client_id)
					}
				`, orgID, name),
				Check: resource.TestCheckOutput("includes_egress_service_account", "true"),
			},
		},
	})
}

func configBasic(orgID, name, description string, roles []string, secretExpiresAfterHours *int, withoutInitialSecret bool) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	rolesHCL := fmt.Sprintf("[%s]", rolesStr)
	secretExpiresLine := ""
	if secretExpiresAfterHours != nil {
		secretExpiresLine = fmt.Sprintf("\n\t\t\tsecret_expires_after_hours = %d", *secretExpiresAfterHours)
	}
	withoutInitialSecretLine := ""
	if withoutInitialSecret {
		withoutInitialSecretLine = "\n\t\t\twithout_initial_secret = true"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = %[3]q
			roles                      = %[4]s%[5]s%[6]s
		}

		data "mongodbatlas_service_account" "test" {
			org_id = %[1]q
			client_id = mongodbatlas_service_account.test.client_id
		}
		data "mongodbatlas_service_accounts" "test" {
			org_id = %[1]q
			depends_on = [mongodbatlas_service_account.test]
		}
	`, orgID, name, description, rolesHCL, secretExpiresLine, withoutInitialSecretLine)
}

func checkBasic(isCreate, withoutInitialSecret bool) resource.TestCheckFunc {
	commonAttrsSet := []string{"client_id", "created_at", "secrets.0.secret_id", "secrets.0.created_at", "secrets.0.expires_at"}
	commonAttrsMap := map[string]string{"secrets.#": "1"}

	checks := acc.CheckRSAndDS(resourceName, new(dataSourceName), nil, commonAttrsSet, commonAttrsMap, checkExists(resourceName))

	additionalChecks := []resource.TestCheckFunc{}
	if isCreate {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.secret")
	} else {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.masked_secret_value")
	}

	additionalChecks = acc.AddAttrSetChecks(dataSourceName, additionalChecks, "secrets.0.masked_secret_value")
	additionalChecks = acc.AddAttrSetChecksPrefix(dataSourcePluralName, additionalChecks, []string{"secrets.0.masked_secret_value"}, "results.0")
	if withoutInitialSecret {
		additionalChecks = append(additionalChecks, resource.TestCheckResourceAttr(resourceName, "without_initial_secret", strconv.FormatBool(withoutInitialSecret)))
	} else {
		additionalChecks = append(additionalChecks, resource.TestCheckNoResourceAttr(resourceName, "without_initial_secret"))
	}

	return resource.ComposeAggregateTestCheckFunc(checks, resource.ComposeAggregateTestCheckFunc(additionalChecks...))
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if orgID == "" || clientID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().ServiceAccountsAPI.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("service account (%s/%s) does not exist", orgID, clientID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_service_account" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if orgID == "" || clientID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		_, _, err := acc.ConnV2().ServiceAccountsAPI.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
		if err == nil {
			return fmt.Errorf("service account (%s/%s) still exists", orgID, clientID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if orgID == "" || clientID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", orgID, clientID), nil
	}
}
