package pushbasedlogexportapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName       = "mongodbatlas_push_based_log_export_api.test"
	nonEmptyPrefixPath = "push-log-prefix"
	defaultPrefixPath  = ""
)

func TestAccPushBasedLogExportAPI_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		groupID              = acc.ProjectIDExecution(tb)
		s3BucketNamePrefix   = acc.RandomS3BucketName()
		s3BucketName1        = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2        = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb); acc.PreCheckAwsEnvBasic(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName1, nonEmptyPrefixPath)...),
			},
			{
				Config: configBasicUpdated(groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName2, nonEmptyPrefixPath)...),
			},
			{
				Config:                               configBasicUpdated(groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "group_id",
				ImportStateVerifyIgnore:              []string{"delete_on_create_timeout"},
			},
		},
	}
}

func commonChecks(s3BucketName, prefixPath string) []resource.TestCheckFunc {
	mapChecks := map[string]string{
		"bucket_name": s3BucketName,
		"prefix_path": prefixPath,
		"state":       "ACTIVE",
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	return acc.AddAttrSetChecks(resourceName, checks, "group_id", "iam_role_id")
}

func configBasic(groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool) string {
	return fmt.Sprintf(`
	 	locals {
				group_id = %[1]q
		 		s3_bucket_name_1 = %[2]q
				s3_bucket_name_2 = %[3]q
		 		s3_bucket_policy_name = %[4]q
		 		aws_iam_role_policy_name = %[5]q
		 		aws_iam_role_name = %[6]q
		 	  }

			   %[7]s

			   %[8]s		
	`, groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(s3BucketName1, s3BucketName2), pushBasedLogExportConfig(false, usePrefixPath, prefixPath))
}

func configBasicUpdated(groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool) string {
	test := fmt.Sprintf(`
	 	locals {
				group_id = %[1]q
		 		s3_bucket_name_1 = %[2]q
				s3_bucket_name_2 = %[3]q
		 		s3_bucket_policy_name = %[4]q
		 		aws_iam_role_policy_name = %[5]q
		 		aws_iam_role_name = %[6]q
		 	  }

			   %[7]s

			   %[8]s
	`, groupID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(s3BucketName1, s3BucketName2), pushBasedLogExportConfig(true, usePrefixPath, prefixPath)) // updating the S3 bucket to use for push-based log config
	return test
}

// pushBasedLogExportConfig returns config for mongodbatlas_push_based_log_expor_api resource and data source.
// This method uses the project and S3 bucket created in awsIAMroleAuthAndS3Config()
func pushBasedLogExportConfig(useBucket2, usePrefixPath bool, prefixPath string) string {
	bucketNameAttr := "bucket_name = aws_s3_bucket.log_bucket_1.bucket"
	if useBucket2 {
		bucketNameAttr = "bucket_name = aws_s3_bucket.log_bucket_2.bucket"
	}
	if usePrefixPath {
		return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export_api" "test" {
			group_id  = local.group_id
			%[1]s
			iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			prefix_path = %[2]q
		}
		`, bucketNameAttr, prefixPath)
	}

	return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export_api" "test" {
		group_id  = local.group_id
		%[1]s
		iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
	}
	`, bucketNameAttr)
}

// awsIAMroleAuthAndS3Config returns config for required IAM roles and authorizes them (sets up cloud provider access) with a mongodbatlas_project
// This method also creates two S3 buckets and sets up required access policy for them
func awsIAMroleAuthAndS3Config(firstBucketName, secondBucketName string) string {
	firstBucketResourceName := "arn:aws:s3:::" + firstBucketName
	secondBucketResourceName := "arn:aws:s3:::" + secondBucketName
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
          %[1]q,
          %[2]q
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
# resource "mongodbatlas_project" "project-tf" {
	#   name     = local.project_name
	#   org_id = local.org_id
	# }

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = local.group_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id    = local.group_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

// Create S3 buckets
resource "aws_s3_bucket" "log_bucket_1" {
  bucket        = local.s3_bucket_name_1
  force_destroy = true  // required as atlas creates a test folder in the bucket when push-based log export is set up 
}

resource "aws_s3_bucket" "log_bucket_2" {
  bucket        = local.s3_bucket_name_2
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
          "${aws_s3_bucket.log_bucket_1.arn}",
          "${aws_s3_bucket.log_bucket_1.arn}/*",
          "${aws_s3_bucket.log_bucket_2.arn}",
          "${aws_s3_bucket.log_bucket_2.arn}/*"
        ]
      }
    ]
  }
  EOF
}
		`, firstBucketResourceName, secondBucketResourceName)
}

func checkDestroy(state *terraform.State) error {
	if projectDestroyedErr := acc.CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_push_based_log_export_api" {
			resp, _, err := acc.ConnV2().PushBasedLogExportApi.GetLogExport(context.Background(), rs.Primary.Attributes["group_id"]).Execute()
			if err == nil && *resp.State != pushbasedlogexport.UnconfiguredState {
				return fmt.Errorf("push-based log export for group_id %s still configured with state %s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["state"])
			}
			if err != nil {
				return fmt.Errorf("push-based log export for group_id %s still configured", rs.Primary.Attributes["group_id"])
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
		return rs.Primary.Attributes["group_id"], nil
	}
}
