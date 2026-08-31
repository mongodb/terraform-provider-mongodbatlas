package cloudbackupcollectionrestorejobcollection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const collectionsDSName = "data.mongodbatlas_cloud_backup_collection_restore_job_collections.test"

func TestAccCloudBackupCollectionRestoreJobCollections_envRead(t *testing.T) {
	acc.SkipTestForCI(t) // needs a finished restore job via env vars
	var (
		projectID   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName = os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
		jobID       = os.Getenv("MONGODB_ATLAS_COLLECTION_RESTORE_JOB_ID")
		sourceNS    = os.Getenv("MONGODB_ATLAS_SOURCE_NAMESPACE")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
			if projectID == "" || clusterName == "" || jobID == "" {
				t.Fatal("`MONGODB_ATLAS_PROJECT_ID`, `MONGODB_ATLAS_CLUSTER_NAME`, and `MONGODB_ATLAS_COLLECTION_RESTORE_JOB_ID` must be set for acceptance testing")
			}
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configCollectionsEnvRead(projectID, clusterName, jobID, sourceNS),
				Check:  collectionsEnvReadChecks(sourceNS),
			},
		},
	})
}

func collectionsEnvReadChecks(sourceNS string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrWith(collectionsDSName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(collectionsDSName, "results.0.source_namespace"),
		resource.TestCheckResourceAttrSet(collectionsDSName, "results.0.state"),
	}
	if sourceNS != "" {
		checks = append(checks,
			resource.TestCheckResourceAttr(collectionsDSName, "source_namespace", sourceNS),
			resource.TestCheckResourceAttr(collectionsDSName, "results.0.source_namespace", sourceNS),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configCollectionsEnvRead(projectID, clusterName, jobID, sourceNS string) string {
	filter := ""
	if sourceNS != "" {
		filter = fmt.Sprintf("\n  source_namespace = %q", sourceNS)
	}
	return fmt.Sprintf(`
data "mongodbatlas_cloud_backup_collection_restore_job_collections" "test" {
  project_id   = %q
  cluster_name = %q
  job_id       = %q%s
}
`, projectID, clusterName, jobID, filter)
}
