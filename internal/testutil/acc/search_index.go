package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func CheckDestroySearchIndex(state *terraform.State) error {
	if ExistingClusterUsed() {
		return nil
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_search_index" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if _, _, err := ConnV2().AtlasSearchAPI.GetClusterSearchIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"]).Execute(); err == nil {
			return fmt.Errorf("index id (%s) still exists", ids["index_id"])
		}
	}
	return nil
}
