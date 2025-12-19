package orgserviceaccountsecretapi_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccOrgServiceAccountSecretAPI_dataSourceErrors(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configDataSourceInvalidClientID(orgID),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
			{
				Config:      configDataSourceInvalidSecretID(orgID, name),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
		},
	})
}

func configDataSourceInvalidClientID(orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_org_service_account_secret_api" "test" {
			org_id    = %[1]q
			client_id = "000000000000000000000000"
			id        = "000000000000000000000000"
		}
	`, orgID)
}

func configDataSourceInvalidSecretID(orgID, saName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_service_account_api" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for data source errors"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_org_service_account_secret_api" "test" {
			org_id    = mongodbatlas_org_service_account_api.test.org_id
			client_id = mongodbatlas_org_service_account_api.test.client_id
			id        = "000000000000000000000001"
		}
	`, orgID, saName)
}
