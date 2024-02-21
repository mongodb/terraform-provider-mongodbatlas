package accesslistapikey_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAccesslistAPIKeys_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_access_list_api_key.test"
		dataSourceName = "data.mongodbatlas_access_list_api_keys.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description    = acc.RandomName()
		ipAddress      = acc.RandomIP(179, 154, 226)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDSPlural(orgID, description, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func configDSPlural(orgID, description, ipAddress string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_access_list_api_keys" "test" {
		org_id     = %[1]q
		api_key_id = mongodbatlas_access_list_api_key.test.api_key_id
	  }
	  
	  resource "mongodbatlas_api_key" "test" {
		org_id = %[1]q
		description = %[2]q
		role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	  }
	  
	  resource "mongodbatlas_access_list_api_key" "test" {
		org_id     = %[1]q
		ip_address = %[3]q
	    api_key_id = mongodbatlas_api_key.test.api_key_id
	  }
	`, orgID, description, ipAddress)
}
