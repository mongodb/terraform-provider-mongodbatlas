package cloudbackupcollectionrestorejob_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName            = "mongodbatlas_cloud_backup_collection_restore_job.test"
	dataSourceName          = "data.mongodbatlas_cloud_backup_collection_restore_job.test"
	pluralDataSourceName    = "data.mongodbatlas_cloud_backup_collection_restore_jobs.test"
	collectionsDSName       = "data.mongodbatlas_cloud_backup_collection_restore_job_collections.test"
	writeStrategyCreate     = "CREATE_NEW"
	indexStrategyAll        = "ALL"
	stateSuccessful         = "SUCCESSFUL"
	stateCanceled           = "CANCELED"
	createTimeout1m         = "1m"
	databaseSuffixPIT       = "_pit"
	collectionSuffixPIT     = "_restored"
	sourceProjectIDRef      = "data.mongodbatlas_advanced_cluster.source.project_id"
	sourceClusterNameRef    = "data.mongodbatlas_advanced_cluster.source.name"
	destProjectIDRef        = "mongodbatlas_project.dest.id"
	destClusterNameRef      = "mongodbatlas_advanced_cluster.dest.name"
	destClusterProjectIDRef = "mongodbatlas_advanced_cluster.dest.project_id"
)

func TestAccCloudBackupCollectionRestoreJob_snapshotSameClusterDatabaseRename(t *testing.T) {
	fixture := acc.CloudBackupCollectionRestoreFixture(t)
	renamedDB := acc.RandomName()
	cfg := restoreJobConfig{
		prefixHCL:         sourceClusterHCL(t, fixture),
		snapshotID:        fixture.SnapshotID,
		targetProjectID:   sourceProjectIDRef,
		targetClusterName: sourceClusterNameRef,
		databaseSource:    acc.CollectionRestoreSeedDatabase,
		databaseTarget:    renamedDB,
		collectionSource:  acc.CollectionRestoreSeedRestaurantsNS,
		withDataSources:   true,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: cfg.HCL(),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", stateSuccessful),
					resource.TestCheckResourceAttr(dataSourceName, "state", stateSuccessful),
					resource.TestCheckResourceAttrSet(resourceName, "job_id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(pluralDataSourceName, "target_cluster_name", knownvalue.StringExact(fixture.ClusterName), map[string]knownvalue.Check{
						"state": knownvalue.StringExact(stateSuccessful),
					}),
					acc.PluralResultCheck(collectionsDSName, "source_namespace", knownvalue.StringExact(acc.CollectionRestoreSeedDatabase+"."+acc.CollectionRestoreSeedCollectionName), map[string]knownvalue.Check{
						"effective_target_namespace": knownvalue.StringExact(renamedDB + "." + acc.CollectionRestoreSeedCollectionName),
						"state":                      knownvalue.StringExact(stateSuccessful),
					}),
					acc.PluralResultCheck(collectionsDSName, "source_namespace", knownvalue.StringExact(acc.CollectionRestoreSeedRestaurantsNS), map[string]knownvalue.Check{
						"effective_target_namespace": knownvalue.StringRegexp(regexp.MustCompile(`^sample_restaurants\.restaurants`)),
						"state":                      knownvalue.StringExact(stateSuccessful),
					}),
				},
			},
		},
	})
}

func TestAccCloudBackupCollectionRestoreJob_pitDestClusterSuffixes(t *testing.T) {
	fixture := acc.CloudBackupCollectionRestoreFixture(t)
	destProject := destProjectHCL()
	dest := destClusterInfo(t, destProjectIDRef)
	cfg := restoreJobConfig{
		prefixHCL:         sourceAndDestHCL(t, fixture, &dest, destProject),
		targetProjectID:   destProjectIDRef,
		targetClusterName: destClusterNameRef,
		pointInTime:       fixture.AfterSnapshotUTCSeconds,
		databaseSource:    acc.CollectionRestoreSeedDatabase,
		databaseSuffix:    databaseSuffixPIT,
		collectionSuffix:  collectionSuffixPIT,
		withDataSources:   true,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &dest, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: cfg.HCL(),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", stateSuccessful),
					resource.TestCheckResourceAttr(dataSourceName, "state", stateSuccessful),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(collectionsDSName, "source_namespace", knownvalue.StringExact(acc.CollectionRestoreSeedDatabase+"."+acc.CollectionRestoreSeedCollectionName), map[string]knownvalue.Check{
						"effective_target_namespace": knownvalue.StringExact(acc.CollectionRestoreSeedDatabase + databaseSuffixPIT + "." + acc.CollectionRestoreSeedCollectionName + collectionSuffixPIT),
						"state":                      knownvalue.StringExact(stateSuccessful),
					}),
				},
			},
		},
	})
}

func TestAccCloudBackupCollectionRestoreJob_pitMissingCollection(t *testing.T) {
	fixture := acc.CloudBackupCollectionRestoreFixture(t)
	dest := destClusterInfo(t, fixture.ProjectID)
	cfg := restoreJobConfig{
		prefixHCL:         sourceAndDestHCL(t, fixture, &dest, ""),
		targetProjectID:   destClusterProjectIDRef,
		targetClusterName: destClusterNameRef,
		pointInTime:       fixture.AfterSnapshotUTCSeconds,
		collectionSource:  acc.CollectionRestoreMissingCollectionNS,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &dest, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      cfg.HCL(),
				ExpectError: regexp.MustCompile(`collection was not found`),
			},
		},
	})
}

func TestAccCloudBackupCollectionRestoreJob_createTimeoutPlanCreate(t *testing.T) {
	fixture := acc.CloudBackupCollectionRestoreFixture(t)
	dest := destClusterInfo(t, fixture.ProjectID)
	cfg := restoreJobConfig{
		prefixHCL:         sourceAndDestHCL(t, fixture, &dest, ""),
		snapshotID:        fixture.SnapshotID,
		targetProjectID:   destClusterProjectIDRef,
		targetClusterName: destClusterNameRef,
		databaseSource:    acc.CollectionRestoreSeedDatabase,
		createTimeout:     createTimeout1m,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &dest, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyWaitJob(fixture.ProjectID, fixture.ClusterName, dest.Name, stateCanceled),
		Steps: []resource.TestStep{
			{
				Config:      cfg.HCL(),
				ExpectError: regexp.MustCompile(`timeout while waiting for state to become`),
			},
			{
				Config:             cfg.HCL(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

type restoreJobConfig struct {
	prefixHCL         string
	snapshotID        string
	targetProjectID   string
	targetClusterName string
	databaseSource    string
	databaseTarget    string
	collectionSource  string
	databaseSuffix    string
	collectionSuffix  string
	createTimeout     string
	pointInTime       int64
	withDataSources   bool
}

func (c *restoreJobConfig) HCL() string {
	optional := ""
	if c.snapshotID != "" {
		optional += fmt.Sprintf("\n\tsnapshot_id = %q", c.snapshotID)
	}
	if c.pointInTime != 0 {
		optional += fmt.Sprintf("\n\tpoint_in_time_utc_seconds = %d", c.pointInTime)
	}
	if c.databaseSuffix != "" {
		optional += fmt.Sprintf("\n\tdatabase_suffix = %q", c.databaseSuffix)
	}
	if c.collectionSuffix != "" {
		optional += fmt.Sprintf("\n\tcollection_suffix = %q", c.collectionSuffix)
	}
	if c.createTimeout != "" {
		optional += fmt.Sprintf("\n\ttimeouts = { create = %q }", c.createTimeout)
	}

	databases := ""
	if c.databaseSource != "" {
		target := ""
		if c.databaseTarget != "" {
			target = fmt.Sprintf("\n\t\ttarget_namespace = %q", c.databaseTarget)
		}
		databases = fmt.Sprintf(`
	databases = [{
		source_namespace = %q%s
	}]`, c.databaseSource, target)
	}

	collections := ""
	if c.collectionSource != "" {
		collections = fmt.Sprintf(`
	collections = [{
		source_namespace = %q
	}]`, c.collectionSource)
	}

	ds := ""
	if c.withDataSources {
		ds = `
	data "mongodbatlas_cloud_backup_collection_restore_job" "test" {
		project_id   = mongodbatlas_cloud_backup_collection_restore_job.test.project_id
		cluster_name = mongodbatlas_cloud_backup_collection_restore_job.test.cluster_name
		job_id       = mongodbatlas_cloud_backup_collection_restore_job.test.job_id
	}

	data "mongodbatlas_cloud_backup_collection_restore_jobs" "test" {
		project_id   = mongodbatlas_cloud_backup_collection_restore_job.test.project_id
		cluster_name = mongodbatlas_cloud_backup_collection_restore_job.test.cluster_name
		depends_on   = [mongodbatlas_cloud_backup_collection_restore_job.test]
	}

	data "mongodbatlas_cloud_backup_collection_restore_job_collections" "test" {
		project_id   = mongodbatlas_cloud_backup_collection_restore_job.test.project_id
		cluster_name = mongodbatlas_cloud_backup_collection_restore_job.test.cluster_name
		job_id       = mongodbatlas_cloud_backup_collection_restore_job.test.job_id
	}`
	}

	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_cloud_backup_collection_restore_job" "test" {
		project_id          = %[2]s
		cluster_name        = %[3]s
		target_project_id   = %[4]s
		target_cluster_name = %[5]s
		write_strategy      = %[6]q
		index_strategy      = %[7]q
		%[8]s
		%[9]s
		%[10]s
	}
	%[11]s
	`, c.prefixHCL, sourceProjectIDRef, sourceClusterNameRef, c.targetProjectID, c.targetClusterName, writeStrategyCreate, indexStrategyAll, optional, databases, collections, ds)
}

func sourceClusterHCL(t *testing.T, fixture *acc.CollectionRestoreFixture) string {
	t.Helper()
	hcl, _, _, err := acc.ClusterDatasourceHcl(&acc.ClusterRequest{
		ProjectID:      fixture.ProjectID,
		ClusterName:    fixture.ClusterName,
		ResourceSuffix: "source",
	})
	require.NoError(t, err)
	return hcl
}

func destClusterInfo(t *testing.T, projectID string) acc.ClusterInfo {
	t.Helper()
	return acc.GetClusterInfo(t, &acc.ClusterRequest{
		ProjectID:      projectID,
		ResourceSuffix: "dest",
		ReplicationSpecs: []acc.ReplicationSpecRequest{
			{Region: acc.NextCollectionRestoreDestRegion()},
		},
	})
}

func destProjectHCL() string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "dest" {
  org_id = %q
  name   = %q
}
`, os.Getenv("MONGODB_ATLAS_ORG_ID"), acc.RandomProjectName())
}

func sourceAndDestHCL(t *testing.T, fixture *acc.CollectionRestoreFixture, dest *acc.ClusterInfo, destProject string) string {
	t.Helper()
	hcl := sourceClusterHCL(t, fixture)
	if destProject != "" {
		hcl += destProject
	}
	if dest != nil {
		hcl += dest.TerraformStr
	}
	return hcl
}

func checkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		projectID := rs.Primary.Attributes["project_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		jobID := rs.Primary.Attributes["job_id"]
		if projectID == "" || clusterName == "" || jobID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", name)
		}
		_, err := acc.GetCollectionRestoreJob(context.Background(), projectID, clusterName, jobID)
		return err
	}
}

func checkDestroyWaitJob(sourceProjectID, sourceClusterName, destClusterName, wantState string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx := context.Background()
		jobID, err := acc.LatestCollectionRestoreJobID(ctx, sourceProjectID, sourceClusterName, destClusterName)
		if err != nil {
			return err
		}
		return acc.WaitCollectionRestoreJobState(ctx, sourceProjectID, sourceClusterName, jobID, wantState)
	}
}
