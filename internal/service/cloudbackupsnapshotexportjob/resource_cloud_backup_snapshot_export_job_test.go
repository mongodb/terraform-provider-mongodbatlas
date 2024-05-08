package cloudbackupsnapshotexportjob_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName         = "mongodbatlas_cloud_backup_snapshot_export_job.test"
	dataSourceName       = "data.mongodbatlas_cloud_backup_snapshot_export_job.test"
	dataSourcePluralName = "data.mongodbatlas_cloud_backup_snapshot_export_jobs.test"
)

func TestAccBackupSnapshotExportJob_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs AWS IAM role and S3 bucket

	var (
		clusterInfo = acc.GetClusterInfo(tb, &acc.ClusterRequest{CloudBackup: true})
		bucketName  = os.Getenv("AWS_S3_BUCKET")
		iamRoleID   = os.Getenv("IAM_ROLE_ID")
		projectID   = acc.ProjectIDExecution(tb)
		clusterName = clusterInfo.ClusterName
		attrsSet    = []string{
			"id",
			"export_job_id",
			"project_id",
			"cluster_name",
			"snapshot_id",
			"export_bucket_id",
		}
		attrsMapWithProject = map[string]string{
			"project_id": projectID,
		}
		attrsPluralDS = map[string]string{
			"project_id":                    projectID,
			"results.0.custom_data.0.key":   "exported by",
			"results.0.custom_data.0.value": "tf-acc-test",
		}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attrsMapWithProject)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrsMapWithProject)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, attrsPluralDS)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(tb); acc.PreCheckS3Bucket(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, bucketName, iamRoleID, clusterName),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID, clusterName, exportJobID, err := readRequired(rs, resourceName)
		if err != nil {
			return err
		}
		_, _, err = acc.Conn().CloudProviderSnapshotExportJobs.Get(context.Background(), projectID, clusterName, exportJobID)
		if err == nil {
			return nil
		}
		return fmt.Errorf("snapshot export job (%s) does not exist", exportJobID)
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID, clusterName, exportJobID, err := readRequired(rs, resourceName)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s--%s--%s", projectID, clusterName, exportJobID), err
	}
}

func readRequired(rs *terraform.ResourceState, resourceName string) (projectID, clusterName, exportJobID string, err error) {
	projectID, ok := rs.Primary.Attributes["project_id"]
	if !ok {
		err = fmt.Errorf("project_id not defined in resource: %s", resourceName)
	}
	clusterName, ok = rs.Primary.Attributes["cluster_name"]
	if !ok {
		err = fmt.Errorf("cluster_name not defined in resource: %s", resourceName)
	}
	exportJobID, ok = rs.Primary.Attributes["export_job_id"]
	if !ok {
		err = fmt.Errorf("export_job_id not defined in resource: %s", resourceName)
	}
	return projectID, clusterName, exportJobID, err
}

func configBasic(projectID, bucketName, iamRoleID, clusterName string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = %[1]q
  cluster_name      = %[4]q
  description       = "tf-acc-test"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id     = %[1]q
  iam_role_id    = "%[3]s"
  bucket_name    = "%[2]s"
  cloud_provider = "AWS"
}

resource "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id   		= %[1]q
  cluster_name 		= %[4]q
  snapshot_id		= mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  export_bucket_id 	= mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id

  custom_data {
    key   = "exported by"
    value = "tf-acc-test"
  }
}

data "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
    project_id 		= %[1]q
    cluster_name 	= %[4]q
    export_job_id 	= mongodbatlas_cloud_backup_snapshot_export_job.test.export_job_id
}
  
data "mongodbatlas_cloud_backup_snapshot_export_jobs" "test" {
    depends_on 	= [mongodbatlas_cloud_backup_snapshot_export_job.test] 
    project_id   	= %[1]q
    cluster_name 	= %[4]q
}

`, projectID, bucketName, iamRoleID, clusterName)
}
