package onlinearchive_test

import (
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupRSOnlineArchiveWithNoChangeBetweenVersions(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acc.RandomProjectName()
		clusterName               = acc.RandomClusterName()
		deleteExpirationDays      = 0
	)
	if mig.IsProviderVersionAtLeast("1.12.2") {
		deleteExpirationDays = 7
	}
	config := configWithDailySchedule(orgID, projectName, clusterName, 1, deleteExpirationDays)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configFirstStep(orgID, projectName, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "partition_fields.0.field_name", "last_review"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
