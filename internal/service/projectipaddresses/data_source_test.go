package projectipaddresses_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectIPAddressesDS_basic(t *testing.T) {
	var (
		projectID      = acc.ProjectIDExecution(t)
		dataSourceName = "data.mongodbatlas_project_ip_addresses.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: ProjectIPAddressesConfig(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "outbound.aws.us-east-1.0", acc.CIDRBlockExpression()),
				),
			},
		},
	})
}

func ProjectIPAddressesConfig(projectID string) string {
	return fmt.Sprintf(`

	data "mongodbatlas_project_ip_addresses" "test" {
		project_id = %[1]q
	}
`, projectID)
}
