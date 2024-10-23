package organization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrganizations_basic(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPluralDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.restrict_employee_access"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.multi_factor_auth_required"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.api_access_list_required"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.gen_ai_features_enabled", "true"),
				),
			},
		},
	})
}

func TestAccConfigDSOrganizations_withPagination(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPagination(2, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
				),
			},
		},
	})
}

func configWithPluralDS() string {
	return `	
		data "mongodbatlas_organizations" "test" {
		}
	`
}

func configWithPagination(pageNum, itemPage int) string {
	return fmt.Sprintf(`
		data "mongodbatlas_organizations" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, pageNum, itemPage)
}
