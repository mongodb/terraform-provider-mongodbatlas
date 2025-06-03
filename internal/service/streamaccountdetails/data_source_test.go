package streamaccountdetails_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamAccountDetailsDS_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_account_details.test_details"
		clusterRequest = acc.ClusterRequest{
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{Region: "US_EAST_1"}, // A cluster in a region supported by Streams tenants is required
			},
		}
		clusterInfo = acc.GetClusterInfo(t, &clusterRequest)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: StreamAccountDetailsConfig(clusterInfo.TerraformStr, clusterInfo.ProjectID, "aws", "US_EAST_1", clusterInfo.ResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "project_id", clusterInfo.ProjectID),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_account_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
					resource.TestCheckNoResourceAttr(dataSourceName, "azure_subscription_id"),
					resource.TestCheckNoResourceAttr(dataSourceName, "virtual_network_name"),
				),
			},
		},
	})
}

func StreamAccountDetailsConfig(clusterInfoStr, projectID, cloudProvider, regionName, clusterInfoNameRef string) string {
	return fmt.Sprintf(`
		%[1]s

		data "mongodbatlas_stream_account_details" "test_details" {
			project_id 		= %[2]q
			cloud_provider	= %[3]q
			region_name 	= %[4]q
			depends_on = [
				%[5]q
			]
		}
`, clusterInfoStr, projectID, cloudProvider, regionName, clusterInfoNameRef)
}
