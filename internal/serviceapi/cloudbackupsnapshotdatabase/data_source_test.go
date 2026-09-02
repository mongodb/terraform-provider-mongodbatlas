package cloudbackupsnapshotdatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	snapshotDatabasesDS  = "data.mongodbatlas_cloud_backup_snapshot_databases.test"
	snapshotCollections1 = "data.mongodbatlas_cloud_backup_snapshot_database_collections.sample_mflix"
	snapshotCollections2 = "data.mongodbatlas_cloud_backup_snapshot_database_collections.sample_restaurants"
	sampleRestaurantsDB  = "sample_restaurants"
	sampleRestaurantsCol = "restaurants"
)

func TestAccCloudBackupSnapshotDatabase_collectionsInSnapshot(t *testing.T) {
	var (
		fixture = acc.CloudBackupCollectionRestoreFixture(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSnapshotDiscovery(fixture.ProjectID, fixture.ClusterName, fixture.SnapshotID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(snapshotDatabasesDS, "results.#", acc.IntGreatThan(0)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(snapshotDatabasesDS, "name", knownvalue.StringExact(acc.CollectionRestoreSeedDatabase), map[string]knownvalue.Check{}),
					acc.PluralResultCheck(snapshotDatabasesDS, "name", knownvalue.StringExact(sampleRestaurantsDB), map[string]knownvalue.Check{}),
					acc.PluralResultCheck(snapshotCollections1, "name", knownvalue.StringExact(acc.CollectionRestoreSeedCollectionName), map[string]knownvalue.Check{}),
					acc.PluralResultCheck(snapshotCollections2, "name", knownvalue.StringExact(sampleRestaurantsCol), map[string]knownvalue.Check{}),
				},
			},
		},
	})
}

func configSnapshotDiscovery(projectID, clusterName, snapshotID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_backup_snapshot_databases" "test" {
			project_id   = %[1]q
			cluster_name = %[2]q
			snapshot_id  = %[3]q
		}

		data "mongodbatlas_cloud_backup_snapshot_database_collections" "sample_mflix" {
			project_id    = %[1]q
			cluster_name  = %[2]q
			snapshot_id   = %[3]q
			database_name = %[4]q
		}

		data "mongodbatlas_cloud_backup_snapshot_database_collections" "sample_restaurants" {
			project_id    = %[1]q
			cluster_name  = %[2]q
			snapshot_id   = %[3]q
			database_name = %[5]q
		}
	`, projectID, clusterName, snapshotID, acc.CollectionRestoreSeedDatabase, sampleRestaurantsDB)
}
