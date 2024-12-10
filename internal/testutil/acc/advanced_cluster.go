package acc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
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
	ClusterLabelsMap1 = map[string]string{
		"key":   "label key 1",
		"value": "label value 1",
	}

	ClusterLabelsMap2 = map[string]string{
		"key":   "label key 2",
		"value": "label value 2",
	}

	ClusterLabelsMap3 = map[string]string{
		"key":   "label key 3",
		"value": "label value 3",
	}
)

func TestStepImportCluster(resourceName string, ignoredFields ...string) resource.TestStep {
	return resource.TestStep{
		ResourceName:                         resourceName,
		ImportStateIdFunc:                    ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
		ImportState:                          true,
		ImportStateVerify:                    true,
		ImportStateVerifyIdentifierAttribute: "name",
		ImportStateVerifyIgnore:              ignoredFields,
	}
}

func CheckDestroyCluster(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" && rs.Type != "mongodbatlas_advanced_cluster" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		clusterName := rs.Primary.Attributes["name"]
		if projectID == "" || clusterName == "" {
			return fmt.Errorf("projectID or clusterName is empty: %s, %s", projectID, clusterName)
		}
		resp, _, _ := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute()
		if resp.GetId() != "" {
			return fmt.Errorf("cluster (%s:%s) still exists", clusterName, rs.Primary.ID)
		}
	}
	return nil
}

func CheckExistsCluster(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		clusterName := rs.Primary.Attributes["name"]
		if projectID == "" || clusterName == "" {
			return fmt.Errorf("projectID or clusterName is empty: %s, %s", projectID, clusterName)
		}
		if err := CheckExistsClusterHandlingRetry(projectID, clusterName); err != nil {
			return fmt.Errorf("cluster(%s:%s) does not exist: %w", projectID, clusterName, err)
		}
		return nil
	}
}

func CheckExistsClusterHandlingRetry(projectID, clusterName string) error {
	return retry.RetryContext(context.Background(), 3*time.Minute, func() *retry.RetryError {
		if _, _, err := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute(); err != nil {
			if admin.IsErrorCode(err, "SERVICE_UNAVAILABLE") {
				// retrying get operation because for migration test it can be the first time new API is called for a cluster so API responds with temporary error as it transition to enabling ISS FF
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
}
