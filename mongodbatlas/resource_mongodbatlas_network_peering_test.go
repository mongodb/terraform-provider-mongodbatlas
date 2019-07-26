package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasNetworkPeering_basic(t *testing.T) {

	var peer matlas.Peer

	resourceName := "mongodbatlas_network_peering.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	vpcID := os.Getenv("AWS_VPC_ID")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnv(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkPeeringConfig(projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkPeeringExists(resourceName, &peer),
					testAccCheckMongoDBAtlasNetworkPeeringAttributes(&peer, vpcCIDRBlock),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "aws_account_id", awsAccountID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasNetworkPeeringImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accepter_region_name"},
			},
		},
	})

}

func testAccCheckMongoDBAtlasNetworkPeeringImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["peer_id"]), nil
	}
}

func testAccCheckMongoDBAtlasNetworkPeeringExists(resourceName string, peer *matlas.Peer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["peer_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if peerResp, _, err := conn.Peers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["peer_id"]); err == nil {
			*peer = *peerResp
			return nil
		}

		return fmt.Errorf("peer(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["peer_id"])
	}
}

func testAccCheckMongoDBAtlasNetworkPeeringAttributes(peer *matlas.Peer, vpcCIDRBlock string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if peer.RouteTableCIDRBlock != vpcCIDRBlock {
			return fmt.Errorf("bad vpcCIDRBlock: %s", peer.RouteTableCIDRBlock)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasNetworkPeeringDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_network_peering" {
			continue
		}

		// Try to find the peer
		_, _, err := conn.Peers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["peer_id"])

		if err == nil {
			return fmt.Errorf("peer (%s) still exists", rs.Primary.Attributes["peer_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasNetworkPeeringConfig(projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = "%[1]s"
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = "AWS"
			region_name			  = "%[5]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			accepter_region_name	  = "us-east-1"	
			project_id    			    = "%[1]s"
			container_id            = mongodbatlas_network_container.test.container_id
			provider_name           = "AWS"
			route_table_cidr_block  = "%[4]s"
			vpc_id					= "%[2]s"
			aws_account_id	= "%[3]s"
		}
	`, projectID, vpcID, awsAccountID, vpcCIDRBlock, awsRegion)
}
