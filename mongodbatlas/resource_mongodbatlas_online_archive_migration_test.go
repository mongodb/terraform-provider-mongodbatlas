package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/testutils"
)

func TestAccMigrationBackupRSOnlineArchiveWithNoChangeBetweenVersions(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acctest.RandomWithPrefix("test-acc")
		name                      = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		lastVersionConstraint     = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccMigrationPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, name, 1, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "partition_fields.0.field_name", "last_review"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, name, 1, 7),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						testutils.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationBackupRSOnlineArchiveWithProcessRegionNoChangeBetweenVersions(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acctest.RandomWithPrefix("test-acc")
		name                      = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		lastVersionConstraint     = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
		cloudProvider             = "AWS"
		region                    = "SA_EAST_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccMigrationPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						testutils.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccBackupRSOnlineArchiveConfigWithProcessRegion(orgID, projectName, name, cloudProvider, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.region", region),
				),
			},
		},
	})
}
