package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasNetworkPeering_basic(t *testing.T) {
	var peer matlas.Peer

	resourceName := "mongodbatlas_network_peering.test"
	projectID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.
	containerID := "5d081429c56c980dc2b810d4"
	vpcID := "vpc-id"
	awsAccountID := "awdAccount"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasNetworkPeeringConfig(projectID, containerID, vpcID, awsAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasNetworkPeeringExists(resourceName, &peer),
					testAccCheckMongoDBAtlasNetworkPeeringAttributes(&peer, containerID),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "container_id", containerID),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "aws_account_id", awsAccountID),
				),
			},
		},
	})

}

func testAccCheckMongoDBAtlasNetworkPeeringExists(resourceName string, peer *matlas.Peer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if peerResp, _, err := conn.Peers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.ID); err == nil {
			*peer = *peerResp
			return nil
		}

		return fmt.Errorf("peer(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasNetworkPeeringAttributes(peer *matlas.Peer, cID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if peer.ContainerID != cID {
			return fmt.Errorf("bad container ID: %s", peer.ContainerID)
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
		_, _, err := conn.Peers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("peer (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasNetworkPeeringConfig(projectID, containerID, vpcID, awsAccountID string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_network_peering" "test" {
	accepter_region_name	= "us-west-1"	
	project_id    			= "%s"
	container_id            = "%s"
	provider_name           = "AWS"
	route_table_cidr_block  = "192.168.0.0/24"
	vpc_id					= "%s"
	aws_account_id			= "%s"
}
`, projectID, containerID, vpcID, awsAccountID)
}
