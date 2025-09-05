package cloudbackupsnapshotexportjob_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
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

	var (
		clusterInfo = acc.GetClusterInfo(tb, &acc.ClusterRequest{CloudBackup: true})
		bucketName  = acc.RandomS3BucketName()
		roleName    = acc.RandomIAMRole()
		policyName  = acc.RandomName()
		projectID   = clusterInfo.ProjectID
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
		PreCheck:                 func() { acc.PreCheckBasicSleep(tb, &clusterInfo, "", ""); mig.PreCheckOldPreviewEnv(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, bucketName, roleName, policyName, clusterInfo.TerraformNameRef, clusterInfo.TerraformStr),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"state"}, // state can change from Queued to InProgress
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
		_, _, err = acc.ConnV2().CloudBackupsApi.GetBackupExport(context.Background(), projectID, clusterName, exportJobID).Execute()
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

func configBasic(projectID, bucketName, roleName, policyName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
resource "aws_iam_role_policy" "test_policy" {
    name = %[4]q
    role = aws_iam_role.test_role.id
    policy = <<-EOF
    {
        "Version": "2012-10-17",
        "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:GetBucketLocation",
            "Resource": "arn:aws:s3:::%[2]s"
        },
        {
            "Effect": "Allow",
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::%[2]s/*"
        }]
    }
    EOF
}

resource "aws_iam_role" "test_role" {
    name = %[3]q
    assume_role_policy = <<EOF
	{
	  "Version": "2012-10-17",
	  "Statement": [
	    {
	      "Effect": "Allow",
	      "Principal": {
	        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
	      },
	      "Action": "sts:AssumeRole",
	      "Condition": {
	        "StringEquals": {
	          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
	        }
	      }
	    }
	  ]
	}
	EOF
}

resource "aws_s3_bucket" "backup" {
	bucket          = %[2]q
	force_destroy   = true
}

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
	project_id    = %[1]q
	provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
	project_id = %[1]q
	role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	aws {
	  iam_assumed_role_arn = aws_iam_role.test_role.arn
	}
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
	project_id        = %[1]q
	cluster_name      = %[5]s
	description       = "tf-acc-test"
	retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
	project_id     = %[1]q
	iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
	bucket_name    = aws_s3_bucket.backup.bucket
	cloud_provider = "AWS"
}

resource "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
	project_id   		= %[1]q
	cluster_name 		= %[5]s
	snapshot_id		= mongodbatlas_cloud_backup_snapshot.test.snapshot_id
	export_bucket_id 	= mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
	custom_data {
		key   = "exported by"
		value = "tf-acc-test"
	}
}

data "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
    project_id 		= %[1]q
    cluster_name 	= %[5]s
    export_job_id 	= mongodbatlas_cloud_backup_snapshot_export_job.test.export_job_id
}
  
data "mongodbatlas_cloud_backup_snapshot_export_jobs" "test" {
    depends_on 	= [mongodbatlas_cloud_backup_snapshot_export_job.test] 
    project_id   	= %[1]q
    cluster_name 	= %[5]s
}

`, projectID, bucketName, roleName, policyName, clusterNameStr)
}
