package cloudbackupsnapshot_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationBackupRSCloudBackupSnapshot_basic(t *testing.T) {
	var (
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo     = acc.GetClusterInfo(orgID, true)
		description     = "My description in my cluster"
		retentionInDays = "4"
		config          = configBasic(&clusterInfo, description, retentionInDays)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
			mig.TestStep(config),
		},
	})
}
