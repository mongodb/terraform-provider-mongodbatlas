package cloudbackupsnapshotexportbucket_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName         = "mongodbatlas_cloud_backup_snapshot_export_bucket.test"
	dataSourceName       = "data.mongodbatlas_cloud_backup_snapshot_export_bucket.test"
	dataSourcePluralName = "data.mongodbatlas_cloud_backup_snapshot_export_buckets.test"
)

func TestAccBackupSnapshotExportBucket_basicAWS(t *testing.T) {
	resource.Test(t, *basicAWSTestCase(t))
}

func TestAccBackupSnapshotExportBucket_basicAzure(t *testing.T) {
	resource.Test(t, *basicAzureTestCase(t))
}

func basicAWSTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID    = acc.ProjectIDExecution(tb)
		bucketName   = os.Getenv("AWS_S3_BUCKET")
		policyName   = acc.RandomName()
		roleName     = acc.RandomIAMRole()
		attrMapCheck = map[string]string{
			"project_id":     projectID,
			"bucket_name":    bucketName,
			"cloud_provider": "AWS",
		}
		pluralAttrMapCheck = map[string]string{
			"project_id":               projectID,
			"results.#":                "1",
			"results.0.bucket_name":    bucketName,
			"results.0.cloud_provider": "AWS",
		}
		attrsSet = []string{
			"iam_role_id",
		}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attrMapCheck)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMapCheck)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, pluralAttrMapCheck)
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, []string{"results.0.iam_role_id"}...)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb); acc.PreCheckS3Bucket(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWSBasic(projectID, bucketName, policyName, roleName),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
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

func basicAzureTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	var (
		projectID          = acc.ProjectIDExecution(t)
		tenantID           = os.Getenv("AZURE_TENANT_ID")
		bucketName         = os.Getenv("AZURE_BLOB_STORAGE_CONTAINER_NAME")
		serviceURL         = os.Getenv("AZURE_SERVICE_URL")
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		attrMapCheck       = map[string]string{
			"project_id":     projectID,
			"bucket_name":    bucketName,
			"service_url":    serviceURL,
			"tenant_id":      tenantID,
			"cloud_provider": "AZURE",
		}
		pluralAttrMapCheck = map[string]string{
			"project_id":               projectID,
			"results.#":                "1",
			"results.0.bucket_name":    bucketName,
			"results.0.service_url":    serviceURL,
			"results.0.cloud_provider": "AZURE",
			"results.0.tenant_id":      tenantID,
		}
		attrsSet = []string{
			"role_id",
		}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attrMapCheck)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrMapCheck)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrsSet...)
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, pluralAttrMapCheck)
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, []string{"results.0.role_id"}...)

	return &resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
			acc.PreCheckCloudProviderAccessAzure(t)
			acc.PreCheckAzureExportBucket(t)
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAzureBasic(projectID, atlasAzureAppID, servicePrincipalID, tenantID, bucketName, serviceURL),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
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
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().CloudBackupsApi.GetExportBucket(context.Background(), ids["project_id"], ids["id"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("snapshot export bucket (%s) does not exist", ids["id"])
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_export_bucket" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		snapshotExportBucket, _, err := acc.ConnV2().CloudBackupsApi.GetExportBucket(context.Background(), ids["project_id"], ids["id"]).Execute()
		if err == nil && snapshotExportBucket != nil {
			return fmt.Errorf("snapshot export bucket (%s) still exists", ids["id"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["id"]), nil
	}
}

func configAWSBasic(projectID, bucketName, policyName, roleName string) string {
	return fmt.Sprintf(`
    resource "aws_iam_role_policy" "test_policy" {
        name = %[3]q
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
        name = %[4]q

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


        resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
            project_id     = %[1]q
            iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
            bucket_name    = %[2]q
            cloud_provider = "AWS"
        }

        data "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
            project_id   =  mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
            export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
        }

        data "mongodbatlas_cloud_backup_snapshot_export_buckets" "test" {
            project_id   =  mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
        }
    `, projectID, bucketName, policyName, roleName)
}

func configAzureBasic(projectID, atlasAzureAppID, servicePrincipalID, tenantID, bucketName, serviceURL string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
			project_id    = %[1]q
			provider_name = "AZURE"
			azure_config {
				atlas_azure_app_id = %[2]q
				service_principal_id = %[3]q
				tenant_id = %[4]q
			}
      	}

		resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
			
			azure {
				atlas_azure_app_id = %[2]q
				service_principal_id = %[3]q
				tenant_id = %[4]q
			}
		}


        resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
            project_id     = %[1]q
            bucket_name    = %[5]q
            cloud_provider = "AZURE"
			service_url	   = %[6]q
			role_id		   = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
        }

        data "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
            project_id   =  mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
            export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
        }

        data "mongodbatlas_cloud_backup_snapshot_export_buckets" "test" {
            project_id   =  mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
        }
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID, bucketName, serviceURL)
}
