package cloudbackupsnapshot_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupRSCloudBackupSnapshot_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.29.0") // version when advanced cluster TPF was introduced
	var (
		clusterInfo     = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		description     = "My description in my cluster"
		retentionInDays = "4"
		config          = configBasic(&clusterInfo, description, retentionInDays)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasicSleep(t); mig.PreCheckOldPreviewEnv(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "type", "replicaSet"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "snapshot_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "replica_set_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigBackupRSCloudBackupSnapshot_sharded(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.29.0") // version when advanced cluster TPF was introduced
	var (
		projectID       = acc.ProjectIDExecution(t)
		clusterName     = acc.RandomClusterName()
		description     = "My description in my cluster"
		retentionInDays = "4"
		config          = configSharded(projectID, clusterName, description, retentionInDays)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasicSleep(t); mig.PreCheckOldPreviewEnv(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "type", "shardedCluster"),
					resource.TestCheckResourceAttrWith(resourceName, "members.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrWith(resourceName, "snapshot_ids.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
