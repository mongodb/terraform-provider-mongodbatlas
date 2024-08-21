package apikey_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAPIKeys_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_api_key.test"
		dataSourceName = "data.mongodbatlas_api_keys.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description    = acc.RandomName()
		roleName       = "ORG_MEMBER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDSPlural(orgID, description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func configDSPlural(orgID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
		  org_id = "%s"
		  description  = "%s"
		  role_names  = ["%s"]
		}

		data "mongodbatlas_api_keys" "test" {
		  org_id = mongodbatlas_api_key.test.org_id
		}
	`, orgID, description, roleNames)
}
