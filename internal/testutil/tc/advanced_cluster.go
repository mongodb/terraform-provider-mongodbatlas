package tc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_advanced_cluster.test"
	dataSourceName       = "data.mongodbatlas_advanced_cluster.test"
	dataSourcePluralName = "data.mongodbatlas_advanced_clusters.test"
)

func SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T, projectName, clusterName string) *resource.TestCase {
	t.Helper()
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, clusterName, 50),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(50),
			},
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, clusterName, 55),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(55),
			},
		},
	}
}

func configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, name string, diskSizeGB int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "SHARDED"

			replication_specs = [{
				num_shards = 2

				region_configs = [{
				electable_specs = {
					instance_size = "M10"
					node_count    = 3
					disk_size_gb  = %[4]d
				}
				analytics_specs = {
					instance_size = "M10"
					node_count    = 0
					disk_size_gb  = %[4]d
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				},
				]
			}]
		}
	`, orgID, projectName, name, diskSizeGB)
}

func checkShardedOldSchemaDiskSizeGBElectableLevel(diskSizeGB int) resource.TestCheckFunc {
	return checkAggr(
		[]string{},
		map[string]string{
			"replication_specs.0.num_shards": "2",
			"disk_size_gb":                   fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.analytics_specs.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
		})
}

func checkAggr(attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	// checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks := []resource.TestCheckFunc{}
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	// checks = acc.AddAttrChecks(dataSourceName, checks, attrsMap)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	// checks = acc.AddAttrSetChecks(dataSourceName, checks, attrsSet...)
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// func checkExists(resourceName string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[resourceName]
// 		if !ok {
// 			return fmt.Errorf("not found: %s", resourceName)
// 		}
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("no ID is set")
// 		}
// 		ids := conversion.DecodeStateID(rs.Primary.ID)
// 		err := acc.CheckClusterExistsHandlingRetry(ids["project_id"], ids["cluster_name"])
// 		if err == nil {
// 			return nil
// 		}
// 		return fmt.Errorf("cluster(%s:%s) does not exist: %w", rs.Primary.Attributes["project_id"], rs.Primary.ID, err)
// 	}
// }
