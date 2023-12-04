package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func CheckDestroyNetworkPeering(s *terraform.State) error {
	conn := TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_network_peering" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := conn.Peers.Get(context.Background(), ids["project_id"], ids["peer_id"])

		if err == nil {
			return fmt.Errorf("peer (%s) still exists", ids["peer_id"])
		}
	}

	return nil
}
