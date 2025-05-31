package streamaccountdetails_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamAccountDetailsDS_basic(t *testing.T) {
	var (
		projectID, _   = acc.ClusterNameExecutionWithRegion(t, constant.UsEast1, false)
		dataSourceName = "data.mongodbatlas_stream_account_details.test_details"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: StreamAccountDetailsConfig(projectID, "aws", "US_EAST_1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_account_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
					resource.TestCheckNoResourceAttr(dataSourceName, "azure_subscription_id"),
					resource.TestCheckNoResourceAttr(dataSourceName, "virtual_network_name"),
				),
			},
		},
	})
}

func StreamAccountDetailsConfig(projectID, cloudProvider, regionName string) string {
	return fmt.Sprintf(`

		data "mongodbatlas_stream_account_details" "test_details" {
			project_id 		= %[1]q
			cloud_provider	= %[2]q
			region_name 	= %[3]q
		}
`, projectID, cloudProvider, regionName)
}
