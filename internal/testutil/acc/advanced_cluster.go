package acc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
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
		"retain_backups_enabled",   // This field is TF specific and not returned by Atlas, so Import can't fill it in.
		"mongo_db_major_version",   // Risks plan change of 8 --> 8.0 (always normalized to `major.minor`)
		"state_name",               // Cluster state can change from IDLE to UPDATING and risks making the test flaky
		"delete_on_create_timeout", // This field is TF specific and not returned by Atlas, so Import can't fill it in.
	)

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
		if _, _, err := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute(); err != nil {
			return fmt.Errorf("cluster(%s:%s) does not exist: %w", projectID, clusterName, err)
		}
		return nil
	}
}

func CheckFCVPinningConfig(resourceName, dataSourceName, pluralDataSourceName string, mongoDBMajorVersion int, pinningExpirationDate *string, fcvVersion *int) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"mongo_db_major_version": fmt.Sprintf("%d.0", mongoDBMajorVersion),
	}

	if pinningExpirationDate != nil {
		mapChecks["pinned_fcv.expiration_date"] = *pinningExpirationDate
	} else {
		mapChecks["pinned_fcv.%"] = "0"
	}

	if fcvVersion != nil {
		mapChecks["pinned_fcv.version"] = fmt.Sprintf("%d.0", *fcvVersion)
	}

	additionalCheck := resource.TestCheckResourceAttrWith(resourceName, "mongo_db_version", MatchesExpression(fmt.Sprintf("%d..*", mongoDBMajorVersion)))

	return CheckRSAndDS(resourceName, admin.PtrString(dataSourceName), admin.PtrString(pluralDataSourceName), []string{}, mapChecks, additionalCheck)
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

func ConfigBasicDedicated(projectID, name, zoneName string) string {
	zoneNameLine := ""
	if zoneName != "" {
		zoneNameLine = fmt.Sprintf("zone_name = %q", zoneName)
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "REPLICASET"
		
		replication_specs = [{
			region_configs = [{
				priority        = 7
				provider_name = "AWS"
				region_name     = "US_EAST_1"
				electable_specs = {
					node_count = 3
					instance_size = "M10"
				}
			}]
			%[3]s
		}]
	}
	data "mongodbatlas_advanced_cluster" "test" {
		project_id = mongodbatlas_advanced_cluster.test.project_id
		name 	     = mongodbatlas_advanced_cluster.test.name
		depends_on = [mongodbatlas_advanced_cluster.test]
	}
			
	data "mongodbatlas_advanced_clusters" "test" {
		project_id = mongodbatlas_advanced_cluster.test.project_id
		depends_on = [mongodbatlas_advanced_cluster.test]
	}
	`, projectID, name, zoneNameLine)
}

func JoinQuotedStrings(list []string) string {
	quoted := make([]string, len(list))
	for i, item := range list {
		quoted[i] = fmt.Sprintf("%q", item)
	}
	return strings.Join(quoted, ", ")
}
