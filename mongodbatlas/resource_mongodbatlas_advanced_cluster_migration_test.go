package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc/todoacc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccMigrationAdvancedClusterRS_singleAWSProvider(t *testing.T) {
	var (
		cluster               matlas.AdvancedCluster
		resourceName          = "mongodbatlas_advanced_cluster.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		rName                 = acctest.RandomWithPrefix("test-acc")
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: todoacc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationAdvancedClusterRS_multiCloud(t *testing.T) {
	var (
		cluster               matlas.AdvancedCluster
		resourceName          = "mongodbatlas_advanced_cluster.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		rName                 = acctest.RandomWithPrefix("test-acc")
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: todoacc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
