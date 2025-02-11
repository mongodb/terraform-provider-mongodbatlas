package acc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
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

	ClusterLabelsMapIgnored = map[string]string{
		"key":   advancedclustertpf.LegacyIgnoredLabelKey,
		"value": "value",
	}
)

func TestStepImportCluster(resourceName string, ignorePrefixFields ...string) resource.TestStep {
	ignorePrefixFields = append(ignorePrefixFields,
		"retain_backups_enabled", // This field is TF specific and not returned by Atlas, so Import can't fill it in.
		"mongo_db_major_version", // Risks plan change of 8 --> 8.0 (always normalized to `major.minor`)
		"state_name",             // Cluster state can change from IDLE to UPDATING and risks making the test flaky
	)

	// auto_scaling & specs (electable_specs, read_only_specs, etc.) are only set in state in SDKv2 if present in the definition.
	// However, as import doesn't have a previous state to compare with, import will always fill them.
	// This will make these fields differ in the state, although the plan change won't be shown to the user as they're computed values.
	if !config.AdvancedClusterV2Schema() {
		ignorePrefixFields = append(ignorePrefixFields, "replication_specs", "id") // TenantUpgrade changes the ID and can make the test flaky
	}
	return resource.TestStep{
		ResourceName:                         resourceName,
		ImportStateIdFunc:                    ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
		ImportState:                          true,
		ImportStateVerify:                    true,
		ImportStateVerifyIdentifierAttribute: "name",
		ImportStateVerifyIgnore:              ignorePrefixFields,
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

func CheckFCVPinningConfig(isAcc bool, resourceName, dataSourceName, pluralDataSourceName string, mongoDBMajorVersion int, pinningExpirationDate *string, fcvVersion *int) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"mongo_db_major_version": fmt.Sprintf("%d.0", mongoDBMajorVersion),
	}

	if pinningExpirationDate != nil {
		mapChecks["pinned_fcv.0.expiration_date"] = *pinningExpirationDate
	} else {
		mapChecks["pinned_fcv.#"] = "0"
	}

	if fcvVersion != nil {
		mapChecks["pinned_fcv.0.version"] = fmt.Sprintf("%d.0", *fcvVersion)
	}

	additionalCheck := resource.TestCheckResourceAttrWith(resourceName, "mongo_db_version", MatchesExpression(fmt.Sprintf("%d..*", mongoDBMajorVersion)))

	return CheckRSAndDSSchemaV2(isAcc, resourceName, admin.PtrString(dataSourceName), admin.PtrString(pluralDataSourceName), []string{}, mapChecks, additionalCheck)
}

func CheckIndependentShardScalingMode(resourceName, clusterName, expectedMode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		projectID := rs.Primary.Attributes["project_id"]
		issMode, _, err := GetIndependentShardScalingMode(context.Background(), projectID, clusterName)
		if err != nil {
			return fmt.Errorf("error getting independent shard scaling mode: %w", err)
		}
		if *issMode != expectedMode {
			return fmt.Errorf("expected independent shard scaling mode to be %s, got %s", expectedMode, *issMode)
		}
		return nil
	}
}

// PopulateWithSampleDataTestCheck is a wrapper around PopulateWithSampleData to be used as a resource.TestCheckFunc
func PopulateWithSampleDataTestCheck(projectID, clusterName string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		return PopulateWithSampleData(projectID, clusterName)
	}
}

// PopulateWithSampleData adds Sample Data to the cluster, otherwise resources like online archive or indexes won't work
func PopulateWithSampleData(projectID, clusterName string) error {
	ctx := context.Background()
	jobLoad, _, err := ConnV2().ClustersApi.LoadSampleDataset(context.Background(), projectID, clusterName).Execute()
	if err != nil || jobLoad == nil {
		return fmt.Errorf("cluster(%s:%s) loading sample data set error: %s", projectID, clusterName, err)
	}
	jobID := jobLoad.GetId()
	stateConf := retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyWorkingState},
		Target:     []string{retrystrategy.RetryStrategyCompletedState},
		Timeout:    15 * time.Minute,
		MinTimeout: 1 * time.Minute,
		Delay:      1 * time.Minute,
		Refresh: func() (result any, state string, err error) {
			job, _, err := ConnV2().ClustersApi.GetSampleDatasetLoadStatus(ctx, projectID, jobID).Execute()
			state = job.GetState()
			return job, state, err
		},
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}
