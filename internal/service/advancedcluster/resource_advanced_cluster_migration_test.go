package advancedcluster_test

import (
	"fmt"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationAdvancedClusterRS_singleAWSProvider(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            testAccMongoDBAtlasAdvancedClusterConfigSingleProviderSDKv2(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAdvancedClusterConfigSingleProviderTestFuncs(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				// testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
				// testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
				// resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				// resource.TestCheckResourceAttr(resourceName, "name", rName),
				// resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
				// resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
				// resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationAdvancedClusterRS_singleAWSProviderUpdate(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            testAccMongoDBAtlasAdvancedClusterConfigSingleProviderSDKv2(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAdvancedClusterConfigSingleProviderTestFuncs(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				// testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
				// testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
				// resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				// resource.TestCheckResourceAttr(resourceName, "name", rName),
				// resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
				// resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
				// resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccClusterAdvancedClusterTestCheckFuncsMulticloud(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationAdvancedClusterRS_multiCloud(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		rName        = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSDKv2(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),

					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.zone_name", advancedcluster.DefaultZoneName),

					resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
					resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.priority", "7"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),

					resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "NORTH_AMERICA_NORTHEAST_1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func testAccMongoDBAtlasAdvancedClusterConfigSingleProviderSDKv2(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = "true"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }
}
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

	`, orgID, projectName, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSDKv2(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = false

  bi_connector_config {
	enabled = true
  }

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      provider_name = "GCP"
      priority      = 6
      region_name   = "NORTH_AMERICA_NORTHEAST_1"
    }
  }
}

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}
