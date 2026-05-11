package apiresource_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_api_resource.test"
	dataSourceName = "data.mongodbatlas_api_resource.test"
)

// TestAccAPIResource_serviceAccount_basic mirrors TestAccServiceAccount_basic
// using the generic mongodbatlas_api_resource. The typed resource owns the
// endpoint contract; this test asserts the generic resource produces an
// equivalent lifecycle against the same Atlas API.
func TestAccAPIResource_serviceAccount_basic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name1        = acc.RandomName()
		name2        = fmt.Sprintf("%s-updated", name1)
		description1 = "Acceptance Test SA generic"
		description2 = "Updated Description"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyServiceAccount(orgID),
		Steps: []resource.TestStep{
			{
				Config: configServiceAccount(orgID, name1, description1, []string{"ORG_READ_ONLY"}, 24),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsServiceAccount(resourceName, orgID),
					resource.TestCheckResourceAttrSet(resourceName, "output.clientId"),
					resource.TestCheckResourceAttrSet(resourceName, "output.createdAt"),
					resource.TestCheckResourceAttrSet(resourceName, "output.secrets.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "output.secrets.0.secret"),
					resource.TestCheckResourceAttr(resourceName, "body.name", name1),
					resource.TestCheckResourceAttr(resourceName, "body.description", description1),
					resource.TestCheckResourceAttrSet(dataSourceName, "output.clientId"),
				),
			},
			{
				Config: configServiceAccount(orgID, name2, description2, []string{"ORG_READ_ONLY"}, 24),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsServiceAccount(resourceName, orgID),
					resource.TestCheckResourceAttr(resourceName, "body.name", name2),
					resource.TestCheckResourceAttr(resourceName, "body.description", description2),
					// After refresh the API stops returning the secret; masked value is populated.
					resource.TestCheckResourceAttrSet(resourceName, "output.secrets.0.maskedSecretValue"),
				),
			},
		},
	})
}

// TestAccAPIResource_serviceAccount_rolesOrdering verifies the reshape engine
// + semantic-equality plan modifier do not raise spurious diffs when roles
// are reordered.
func TestAccAPIResource_serviceAccount_rolesOrdering(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyServiceAccount(orgID),
		Steps: []resource.TestStep{
			{
				Config: configServiceAccount(orgID, name, "Roles ordering", []string{"ORG_BILLING_ADMIN", "ORG_READ_ONLY"}, 24),
			},
			{
				Config: configServiceAccount(orgID, name, "Roles ordering", []string{"ORG_READ_ONLY", "ORG_BILLING_ADMIN"}, 24),
			},
		},
	})
}

func configServiceAccount(orgID, name, description string, roles []string, secretExpiresAfterHours int) string {
	rolesHCL := `["` + strings.Join(roles, `", "`) + `"]`
	return fmt.Sprintf(`
		resource "mongodbatlas_api_resource" "test" {
			path                  = "/api/atlas/v2/orgs/%[1]s/serviceAccounts"
			id_attribute          = ["clientId"]
			create_only_body_keys = ["secretExpiresAfterHours"]

			body = {
				name                    = %[2]q
				description             = %[3]q
				roles                   = %[4]s
				secretExpiresAfterHours = %[5]d
			}
		}

		data "mongodbatlas_api_resource" "test" {
			path = mongodbatlas_api_resource.test.id
		}
	`, orgID, name, description, rolesHCL, secretExpiresAfterHours)
}

func checkExistsServiceAccount(rsName, orgID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("not found: %s", rsName)
		}
		clientID := rs.Primary.Attributes["output.clientId"]
		if clientID == "" {
			return fmt.Errorf("checkExists: output.clientId not set for %s", rsName)
		}
		if _, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute(); err != nil {
			return fmt.Errorf("service account (%s/%s) does not exist: %s", orgID, clientID, err)
		}
		return nil
	}
}

func checkDestroyServiceAccount(orgID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_api_resource" {
				continue
			}
			clientID := rs.Primary.Attributes["output.clientId"]
			if clientID == "" {
				continue
			}
			if _, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute(); err == nil {
				return fmt.Errorf("service account (%s/%s) still exists", orgID, clientID)
			}
		}
		return nil
	}
}
