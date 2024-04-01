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
		searchIndex, _, err := Conn().Search.GetIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"])
		if err == nil && searchIndex != nil && searchIndex.Status != "IN_PROGRESS" { // index can be in progess for some seconds after delete is called
			return fmt.Errorf("index id (%s) still exists", ids["index_id"])
		}
	}
	return nil
}
