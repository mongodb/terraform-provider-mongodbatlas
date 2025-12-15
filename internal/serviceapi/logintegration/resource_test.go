package logintegration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName       = "mongodbatlas_log_integration.test"
	dataSourceName     = "data.mongodbatlas_log_integration.test"
	nonEmptyPrefixPath = "push-log-prefix-v3"
)

func TestAccLogIntegration_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID            = acc.ProjectIDExecution(tb)
		s3BucketName         = acc.RandomS3BucketName()
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketName)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		kmsKey               = os.Getenv("AWS_KMS_KEY_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb); acc.PreCheckAwsEnvBasic(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithDataSource(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(commonChecks(s3BucketName, nonEmptyPrefixPath), dataSourceChecks(s3BucketName, nonEmptyPrefixPath)...)...,
				),
			},
			{
				Config: configBasic(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true, true, kmsKey),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName, nonEmptyPrefixPath)...),
			},
			{
				Config: configBasic(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true, false, ""),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName, nonEmptyPrefixPath)...),
			},
			{
				Config:            configBasic(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true, false, ""),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func dataSourceChecks(s3BucketName, prefixPath string) []resource.TestCheckFunc {
	mapChecks := map[string]string{
		"bucket_name": s3BucketName,
		"prefix_path": prefixPath,
		"type":        "S3_LOG_EXPORT",
		"log_types.#": "1",
	}
	checks := acc.AddAttrChecks(dataSourceName, nil, mapChecks)
	return acc.AddAttrSetChecks(dataSourceName, checks, "project_id", "iam_role_id", "id")
}

func configWithDataSource(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_log_integration" "test" {
			project_id = mongodbatlas_log_integration.test.project_id
			id         = mongodbatlas_log_integration.test.id
		}
	`, configBasic(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath, true, false, ""))
}

func commonChecks(s3BucketName, prefixPath string) []resource.TestCheckFunc {
	mapChecks := map[string]string{
		"bucket_name": s3BucketName,
		"prefix_path": prefixPath,
		"type":        "S3_LOG_EXPORT",
		"log_types.#": "1",
		"log_types.0": "MONGOD_AUDIT",
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	return acc.AddAttrSetChecks(resourceName, checks, "project_id", "iam_role_id", "id")
}

func configBasic(projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath, useKmsKey bool, kmsKey string) string {
	return fmt.Sprintf(`
	 	locals {
				project_id = %[1]q
		 		s3_bucket_name = %[2]q
		 		s3_bucket_policy_name = %[3]q
		 		aws_iam_role_policy_name = %[4]q
		 		aws_iam_role_name = %[5]q
		 	  }

			   %[6]s

			   %[7]s		
	`, projectID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(s3BucketName), logIntegrationConfig(usePrefixPath, useKmsKey, prefixPath, kmsKey))
}

// logIntegrationConfig returns config for mongodbatlas_log_integration resource.
// This method uses the project and S3 bucket created in awsIAMroleAuthAndS3Config()
func logIntegrationConfig(usePrefixPath, useKmsKey bool, prefixPath, kmsKey string) string {
	prefixPathAttr := ""
	if usePrefixPath {
		prefixPathAttr = fmt.Sprintf("prefix_path   = %[1]q", prefixPath)
	}
	kmsKeyAttr := ""
	if useKmsKey {
		kmsKeyAttr = fmt.Sprintf("kms_key   = %[1]q", kmsKey)
	}

	return fmt.Sprintf(`resource "mongodbatlas_log_integration" "test" {
		project_id  = local.project_id
		bucket_name = aws_s3_bucket.log_bucket.bucket
		iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
		type        = "S3_LOG_EXPORT"
		log_types   = ["MONGOD_AUDIT"]
		%[1]s
		%[2]s
	}
	`, prefixPathAttr, kmsKeyAttr)
}

// awsIAMroleAuthAndS3Config returns config for required IAM roles and authorizes them (sets up cloud provider access) with a mongodbatlas_project
// This method also creates two S3 buckets and sets up required access policy for them
func awsIAMroleAuthAndS3Config(bucketName string) string {
	bucketResourceName := "arn:aws:s3:::" + bucketName
	return fmt.Sprintf(`
		// Create IAM role & policy to authorize with Atlas
resource "aws_iam_role_policy" "test_policy" {
  name = local.aws_iam_role_policy_name
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:GetObject",
          "s3:ListBucket",
          "s3:GetObjectVersion"
        ],
        "Resource": "*"
      },
      {
       "Effect": "Allow",
        "Action": "s3:*",
        "Resource": [
          %[1]q
        ]
      }
    ]
  }
  EOF
}



resource "aws_iam_role" "test_role" {
  name = local.aws_iam_role_name
  max_session_duration = 43200

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

// Set up cloud provider access in Atlas for a project using the created IAM role
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = local.project_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id    = local.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

// Create S3 buckets
resource "aws_s3_bucket" "log_bucket" {
  bucket        = local.s3_bucket_name
  force_destroy = true  // required as atlas creates a test folder in the bucket when push-based log export is set up 
}

// Add authorization policy to existing IAM role
resource "aws_iam_role_policy" "s3_bucket_policy" {
  name   = local.s3_bucket_policy_name
  role   = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:ListBucket",
          "s3:PutObject",
          "s3:GetObject",
          "s3:GetBucketLocation"
        ],
        "Resource": [
          "${aws_s3_bucket.log_bucket.arn}",
          "${aws_s3_bucket.log_bucket.arn}/*"
        ]
      }
    ]
  }
  EOF
}
		`, bucketResourceName)
}

func checkDestroy(state *terraform.State) error {
	if projectDestroyedErr := acc.CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_push_based_log_export_api" {
			_, _, err := acc.ConnV2().PushBasedLogExportApi.GetGroupLogIntegration(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]).Execute()
			if err == nil {
				return fmt.Errorf("push-based log export for project_id %s with id %s still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"])
			}
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
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]), nil
	}
}
