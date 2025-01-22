package onlinearchive_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupRSOnlineArchiveWithNoChangeBetweenVersions(t *testing.T) {
	var (
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		clusterInfo               = acc.GetClusterInfo(t, clusterRequest())
		clusterName               = clusterInfo.Name
		projectID                 = clusterInfo.ProjectID
		clusterTerraformStr       = clusterInfo.TerraformStr
		clusterResourceName       = clusterInfo.ResourceName
		deleteExpirationDays      = 0
	)
	if mig.IsProviderVersionAtLeast("1.12.2") {
		deleteExpirationDays = 7
	}
	config := configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, deleteExpirationDays)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            clusterTerraformStr,
				Check:             acc.PopulateWithSampleDataTestCheck(projectID, clusterName),
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
