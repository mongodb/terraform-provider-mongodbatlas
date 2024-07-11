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
		clusterInfo               = acc.GetClusterInfo(t, &acc.ClusterRequest{
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{AutoScalingDiskGbEnabled: true},
			},
			Tags: map[string]string{
				"ArchiveTest": "true", "Owner": "test",
			},
		})
		clusterName          = clusterInfo.ClusterName
		projectID            = clusterInfo.ProjectID
		clusterTerraformStr  = clusterInfo.ClusterTerraformStr
		clusterResourceName  = clusterInfo.ClusterResourceName
		deleteExpirationDays = 0
	)
	if mig.IsProviderVersionAtLeast("1.12.2") {
		deleteExpirationDays = 7
	}
	config := configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, deleteExpirationDays)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            clusterTerraformStr,
				Check: resource.ComposeAggregateTestCheckFunc(
					populateWithSampleData(clusterResourceName, projectID, clusterName),
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
