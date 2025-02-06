package flexsnapshot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourceName       = "data.mongodbatlas_flex_snapshot.test"
	dataSourcePluralName = "data.mongodbatlas_flex_snapshots.test"
)

func TestAccFlexSnapshot_basic(t *testing.T) {
	tc := basicTestCase(t)
	resource.Test(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.RandomName()
		clusterName = acc.RandomName()
		snapshotID  = acc.RandomName()
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

func configBasic(projectID, clusterName, snapshotID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_flex_snapshot" "test" {
		project_id = %[1]q
		name = %[2]q
		snapshot_id = %[3]q
		}

		data "mongodbatlas_flex_snapshots" "test" {
		project_id = %[1]q
		name =  %[2]q
		}`, projectID, clusterName, snapshotID)
}

func checksFlexSnapshot() resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{} //TODO: check exists?
	attrSet := []string{
		"backup_settings.enabled",
		"cluster_type",
		"connection_strings.standard",
		"create_date",
		"id",
		"mongo_db_version",
		"state_name",
		"version_release_system",
		"provider_settings.provider_name",
	}
	pluralMap := []string{
		"project_id",
		"name",
		"results",
	}
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, pluralMap...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrSet...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
