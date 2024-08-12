package acc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

var (
	ClusterTagsMap1 = map[string]string{
		"key":   "key 1",
		"value": "value 1",
	}

	ClusterTagsMap2 = map[string]string{
		"key":   "key 2",
		"value": "value 2",
	}

	ClusterTagsMap3 = map[string]string{
		"key":   "key 3",
		"value": "value 3",
	}
)

func CheckDestroyCluster(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" && rs.Type != "mongodbatlas_advanced_cluster" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		resp, _, _ := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute()
		if resp.GetId() != "" {
			return fmt.Errorf("cluster (%s:%s) still exists", clusterName, rs.Primary.ID)
		}
	}
	return nil
}

func ImportStateClusterIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]), nil
	}
}

func CheckClusterExistsHandlingRetry(projectID, clusterName string) error {
	return retry.RetryContext(context.Background(), 3*time.Minute, func() *retry.RetryError {
		_, _, err := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute()
		if apiError, ok := admin.AsError(err); ok {
			if apiError.GetErrorCode() == "SERVICE_UNAVAILABLE" {
				// retrying get operation because for migration test it can be the first time new API is called for a cluster so API responds with temporary error as it transition to enabling ISS FF
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
}
