package flexrestorejob_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourceName       = "data.mongodbatlas_flex_restore_job.test"
	dataSourcePluralName = "data.mongodbatlas_flex_restore_job.test"
)

func TestAccFlexRestoreJob_basic(t *testing.T) {
	tc := basicTestCase(t)
	resource.Test(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = "2c64aaf6f7ec54c6e8b18c9c"
		clusterName = acc.RandomName()
		snapshotID  = "2c64aaf6f7ec54c6e8b18c9c"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, snapshotID),
				Check:  checksFlexSnapshot(),
			},
		},
	}
}

func configBasic(projectID, clusterName, restoreJobID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_flex_restore_job" "test" {
			project_id = %[1]q
			name = %[2]q
			restore_job_id = %[3]q
		}

		data "mongodbatlas_flex_restore_job" "test" {
			project_id = %[1]q
			name =  %[2]q
		}`, projectID, clusterName, restoreJobID)
}

func checksFlexSnapshot() resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	attrSet := []string{
		"delivery_type",
		"expiration_date",
		"restore_finished_date",
		"restore_scheduled_date",
		"snapshot_finished_date",
		"snapshot_id",
		"snapshot_url",
		"status",
		"target_deployment_item_name",
		"target_project_id",
	}
	pluralMap := []string{
		"project_id",
		"name",
		"results.#",
	}
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, pluralMap...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrSet...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
