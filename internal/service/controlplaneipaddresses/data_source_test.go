package controlplaneipaddresses_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccControlPlaneIpAddressesDS_basic(t *testing.T) {
	dataSourceName := "data.mongodbatlas_control_plane_ip_addresses.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "outbound.aws.us-east-1.0", acc.CIDRBlockExpression()),
				),
			},
		},
	})
}

const configBasic = `
data "mongodbatlas_control_plane_ip_addresses" "test" {
}
`
