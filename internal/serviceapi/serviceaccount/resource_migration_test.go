package serviceaccount_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccServiceAccount_moveFromAPIResource(t *testing.T) {
	var (
		orgID   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name    = acc.RandomName()
		descr   = "moved-from-api-resource acceptance"
		oldName = "mongodbatlas_api_resource.test"
	)
	if orgID == "" {
		t.Skip("MONGODB_ATLAS_ORG_ID not set")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyByClientID(orgID),
		Steps: []resource.TestStep{
			{
				Config: configAPIResourceSA(orgID, name, descr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(oldName, "output.clientId"),
				),
			},
			{
				// Same SA, now managed through the typed resource via moved block.
				// If MoveState is broken, the framework would destroy+create — and
				// the test would fail because identity wouldn't be preserved.
				Config: configTypedSAMoved(orgID, name, descr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", descr),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
				),
			},
		},
	})
}

func configAPIResourceSA(orgID, name, descr string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_api_resource" "test" {
  path         = "/api/atlas/v2/orgs/%[1]s/serviceAccounts"
  id_attribute = ["clientId"]

  body = {
    name                    = %[2]q
    description             = %[3]q
    roles                   = ["ORG_MEMBER"]
    secretExpiresAfterHours = 24
  }
}
`, orgID, name, descr)
}

func configTypedSAMoved(orgID, name, descr string) string {
	return fmt.Sprintf(`
moved {
  from = mongodbatlas_api_resource.test
  to   = mongodbatlas_service_account.test
}

resource "mongodbatlas_service_account" "test" {
  org_id      = %[1]q
  name        = %[2]q
  description = %[3]q
  roles       = ["ORG_MEMBER"]
}
`, orgID, name, descr)
}

func checkDestroyByClientID(orgID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			var clientID string
			switch rs.Type {
			case "mongodbatlas_service_account":
				clientID = rs.Primary.Attributes["client_id"]
			case "mongodbatlas_api_resource":
				clientID = rs.Primary.Attributes["output.clientId"]
			default:
				continue
			}
			if clientID == "" {
				continue
			}
			_, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
			if err == nil {
				return fmt.Errorf("service account %s/%s still exists", orgID, clientID)
			}
		}
		return nil
	}
}
