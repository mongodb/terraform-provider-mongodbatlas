package streamaccountdetails_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamAccountDetailsDS_basic(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, false)
		dataSourceName         = "data.mongodbatlas_stream_account_details.test_details"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: StreamAccountDetailsConfig(projectID, "aws", "US_EAST_1", clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_account_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_subscription_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_name"),
				),
			},
		},
	})
}

func StreamAccountDetailsConfig(projectID, cloudProvider, regionName, clusterName string) string {
	instanceName := acc.RandomName()
	streamInstanceConfig := acc.StreamInstanceConfig(projectID, instanceName, "VIRGINIA_USA", "AWS")

	return fmt.Sprintf(`
	%[1]s

	resource "mongodbatlas_stream_connection" "test_connection" {
		    project_id 			= resource.mongodbatlas_stream_instance.test.project_id
			instance_name 		= resource.mongodbatlas_stream_instance.test.instance_name
		 	connection_name 	= %[2]q
			type            	= "Cluster"
			cluster_name    	= %[3]q
			db_role_to_execute 	= {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
		
		resource "mongodbatlas_stream_processor" "test_processor" {
			project_id     	= resource.mongodbatlas_stream_instance.test.project_id
			instance_name  	= resource.mongodbatlas_stream_instance.test.instance_name
			processor_name 	= "testProcessor"
			pipeline       	= jsonencode([
				{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.test_connection.connection_name } },
				{ "$merge"	= {
					"into" = {
						"connectionName" = resource.mongodbatlas_stream_connection.test_connection.connection_name,
						"db" = "randomDb",
						"coll" = "randomColl"
					}
				} }
			])
			state          	= "STARTED"
		}

	data "mongodbatlas_stream_account_details" "test_details" {
  		project_id 		= resource.mongodbatlas_stream_instance.test.project_id
  		cloud_provider	= %[4]q
  		region_name 	= %[5]q
	}
`, streamInstanceConfig, acc.RandomName(), clusterName, cloudProvider, regionName)
}
