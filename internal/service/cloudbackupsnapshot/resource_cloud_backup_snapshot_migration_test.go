package cloudbackupsnapshot_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationBackupRSCloudBackupSnapshot_basic(t *testing.T) {
	var (
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acctest.RandomWithPrefix("test-acc")
		clusterName     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description     = "My description in my cluster"
		retentionInDays = "4"
		config          = configBasic(orgID, projectName, clusterName, description, retentionInDays)
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
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
			mig.TestStep(config),
		},
	})
}
