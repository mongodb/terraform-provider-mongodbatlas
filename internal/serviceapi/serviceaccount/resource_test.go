package serviceaccount_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

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
				Config: configBasic(orgID, name1, description1, []string{"ORG_OWNER"}, 24),
				Check:  checkBasic(true),
			},
			{
				Config: configBasic(orgID, name2, description2, []string{"ORG_READ_ONLY"}, 24),
				Check:  checkBasic(false),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"secret_expires_after_hours"},
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
				Config: configBasic(orgID, name, "Roles Update", []string{"ORG_BILLING_ADMIN", "ORG_READ_ONLY"}, 24),
			},
			{
				// change order
				Config: configBasic(orgID, name, "Roles Update", []string{"ORG_READ_ONLY", "ORG_BILLING_ADMIN"}, 24),
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
				Config: configBasic(orgID, name, "description", []string{"ORG_READ_ONLY"}, 24),
				Check:  checkExists(resourceName),
			},
			{
				Config:      configBasic(orgID, name, "description", []string{"ORG_READ_ONLY"}, 48),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("secret_expires_after_hours cannot be updated"),
			},
			{
				Config:      configBasic("updated-org-id", name, "description", []string{"ORG_READ_ONLY"}, 24),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("org_id cannot be updated"),
			},
		},
	})
}

func configBasic(orgID, name, description string, roles []string, secretExpiresAfterHours int) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	rolesHCL := fmt.Sprintf("[%s]", rolesStr)
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = %[3]q
			roles                      = %[4]s
			secret_expires_after_hours = %[5]d
		}

		data "mongodbatlas_service_account" "test" {
			org_id = %[1]q
			client_id = mongodbatlas_service_account.test.client_id
		}
		data "mongodbatlas_service_accounts" "test" {
			org_id = %[1]q
			depends_on = [mongodbatlas_service_account.test]
		}
	`, orgID, name, description, rolesHCL, secretExpiresAfterHours)
}

func checkBasic(isCreate bool) resource.TestCheckFunc {
	commonAttrsSet := []string{"client_id", "created_at", "secrets.0.secret_id", "secrets.0.created_at", "secrets.0.expires_at"}
	commonAttrsMap := map[string]string{"secrets.#": "1"}

	checks := acc.CheckRSAndDS(resourceName, admin.PtrString(dataSourceName), nil, commonAttrsSet, commonAttrsMap, checkExists(resourceName))

	additionalChecks := []resource.TestCheckFunc{}
	if isCreate {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.secret")
	} else {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.masked_secret_value")
	}

	additionalChecks = acc.AddAttrSetChecks(dataSourceName, additionalChecks, "secrets.0.masked_secret_value")
	additionalChecks = acc.AddAttrSetChecksPrefix(dataSourcePluralName, additionalChecks, []string{"secrets.0.masked_secret_value"}, "results.0")

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
		_, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
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

		_, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
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
