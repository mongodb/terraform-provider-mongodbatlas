package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasNetworkPeerings_basic(t *testing.T) {

	var peer matlas.Peer

	resourceName := "mongodbatlas_network_peering.test"
	dataSourceName := "data.mongodbatlas_network_peerings.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasNetworkPeeringsConfig(projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkPeeringExists(resourceName, &peer),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "aws_account_id", awsAccountID),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.aws_account_id"),
				),
			},
		},
	})

}

func testAccDSMongoDBAtlasNetworkPeeringsConfig(projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_network_container" "test" {
	project_id   		= "%[1]s"
	atlas_cidr_block    = "192.168.208.0/21"
	provider_name		= "AWS"
	region_name			= "%[5]s"
}

resource "mongodbatlas_network_peering" "test" {
	accepter_region_name	= "us-east-1"
	project_id    			= "%[1]s"
	container_id            = mongodbatlas_network_container.test.id
	provider_name           = "AWS"
	route_table_cidr_block  = "%[4]s"
	vpc_id					= "%[2]s"
	aws_account_id			= "%[3]s"
}

data "mongodbatlas_network_peerings" "test" {
	project_id = mongodbatlas_network_peering.test.project_id
}
`, projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion)
}
