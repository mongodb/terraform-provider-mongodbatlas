package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccClusterDSAdvancedClusters_basic(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = "data.mongodbatlas_advanced_clusters.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasAdvancedClustersConfig(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.termination_protection_enabled"),
				),
			},
		},
	})
}

func TestAccClusterDSAdvancedClusters_advancedConf(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = "data.mongodbatlas_advanced_clusters.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
		processArgs    = &matlas.ProcessArgs{
			FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_1",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasAdvancedClustersConfigAdvancedConf(projectID, name, processArgs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
				),
			},
		},
	})
}

func TestAccClusterDSAdvancedClusters_multicloud(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = "data.mongodbatlas_advanced_clusters.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasAdvancedClustersMultiCloudConfig(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasAdvancedClustersConfig(projectID, name string) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_advanced_clusters" "test" {
  project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, testAccMongoDBAtlasAdvancedClusterConfigTenant(projectID, name))
}

func testAccDataSourceMongoDBAtlasAdvancedClustersMultiCloudConfig(projectID, name string) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_advanced_clusters" "test" {
  project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, name))
}

func testAccDataSourceMongoDBAtlasAdvancedClustersConfigAdvancedConf(projectID, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_advanced_clusters" "test" {
  project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(projectID, name, p))
}
