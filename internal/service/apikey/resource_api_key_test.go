package apikey_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSAPIKey_basic(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_api_key.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description       = acc.RandomName()
		descriptionUpdate = acc.RandomName()
		roleName          = "ORG_MEMBER"
		roleNameUpdated   = "ORG_BILLING_ADMIN"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: configBasic(orgID, descriptionUpdate, roleNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
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
		_, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.GetOrgApiKey(context.Background(), ids["org_id"], ids["api_key_id"]).Execute()
		if err != nil {
			return fmt.Errorf("API Key (%s) does not exist", ids["api_key_id"])
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_api_key" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.GetOrgApiKey(context.Background(), ids["org_id"], ids["api_key_id"]).Execute()
		if err == nil {
			return fmt.Errorf("API Key (%s) still exists", ids["role_name"])
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
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["api_key_id"]), nil
	}
}

func configBasic(orgID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
			org_id     = "%s"
			description  = "%s"

			role_names  = ["%s"]
		}
	`, orgID, description, roleNames)
}
