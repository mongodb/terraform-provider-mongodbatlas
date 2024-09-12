package acc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func ImportStateIDFuncProjectIDClusterName(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return IDWithProjectIDClusterName(rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
	}
}

func IDWithProjectIDClusterName(projectID, clusterName string) (string, error) {
	if err := conversion.ValidateProjectID(projectID); err != nil {
		return "", err
	}
	if err := conversion.ValidateClusterName(clusterName); err != nil {
		return "", err
	}
	return projectID + "-" + clusterName, nil
}
