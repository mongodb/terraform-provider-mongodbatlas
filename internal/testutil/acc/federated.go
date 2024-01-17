package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func CheckDestroyFederatedDatabaseInstance(s *terraform.State) error {
	connV2 := TestAccProviderSdkV2.Meta().(*config.MongoDBClient).AtlasV2

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_federated_database_instance" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := connV2.DataFederationApi.GetFederatedDatabase(context.Background(), ids["project_id"], ids["name"]).Execute()
		if err == nil {
			return fmt.Errorf("federated database instance (%s) still exists", ids["project_id"])
		}
	}

	return nil
}
