package pushbasedlogexport_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName       = "mongodbatlas_push_based_log_export.test"
	datasourceName     = "data.mongodbatlas_push_based_log_export.test"
	nonEmptyPrefixPath = "push-log-prefix"
	defaultPrefixPath  = ""
)

func TestAccPushBasedLogExport_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID            = acc.ProjectIDExecution(tb)
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
				Config: configBasic(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true, "", nil),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName1, nonEmptyPrefixPath)...),
			},
			{
				Config: configBasicUpdated(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName2, nonEmptyPrefixPath)...),
			},
			{
				Config:                               configBasicUpdated(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportStateVerifyIgnore:              []string{"delete_on_create_timeout"},
			},
		},
	}
}

func TestAccPushBasedLogExport_noPrefixPath(t *testing.T) {
	resource.Test(t, *noPrefixPathTestCase(t))
}

func noPrefixPathTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID            = acc.ProjectIDExecution(tb)
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
				Config: configBasic(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, defaultPrefixPath, false, "", nil),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(s3BucketName1, defaultPrefixPath)...),
			},
		},
	}
}

func TestAccPushBasedLogExport_createFailure(t *testing.T) {
	resource.Test(t, *createFailure(t))
}

func createFailure(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID = acc.ProjectIDExecution(tb)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      pushBasedLogExportInvalidConfig(projectID),
				ExpectError: regexp.MustCompile("CLOUD_PROVIDER_ACCESS_ROLE_NOT_FOUND"),
			},
		},
	}
}

func TestAccPushBasedLogExport_createTimeoutWithDeleteOnCreateTimeout(t *testing.T) {
	resource.Test(t, *createTimeoutWithDeleteOnCreateTimeout(t))
}

func createTimeoutWithDeleteOnCreateTimeout(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID             = acc.ProjectIDExecution(tb)
		s3BucketNamePrefix    = acc.RandomS3BucketName()
		s3BucketName1         = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2         = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName    = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName        = acc.RandomIAMRole()
		awsIAMRolePolicyName  = fmt.Sprintf("%s-policy", awsIAMRoleName)
		createTimeout         = "1s"
		deleteOnCreateTimeout = true
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true, acc.TimeoutConfig(&createTimeout, nil, nil), &deleteOnCreateTimeout),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	}
}

func pushBasedLogExportInvalidConfig(projectID string) string {
	return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
		project_id  = %[1]q
		bucket_name = "invalidBucket"
		iam_role_id = "aaaaaaaaaa99999999999999"
	}
	`, projectID)
}

func commonChecks(s3BucketName, prefixPath string) []resource.TestCheckFunc {
	attributes := map[string]string{
		"bucket_name": s3BucketName,
		"prefix_path": prefixPath,
	}
	checks := addAttrChecks(nil, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "project_id", "iam_role_id")
	return acc.AddAttrSetChecks(datasourceName, checks, "project_id", "iam_role_id")
}

func addAttrChecks(checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	checks = acc.AddAttrChecks(resourceName, checks, mapChecks)
	return acc.AddAttrChecks(datasourceName, checks, mapChecks)
}

func configBasic(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool, timeoutConfig string, deleteOnCreateTimeout *bool) string {
	test := fmt.Sprintf(`
	 	locals {
				project_id = %[1]q
		 		s3_bucket_name_1 = %[2]q
				s3_bucket_name_2 = %[3]q
		 		s3_bucket_policy_name = %[4]q
		 		aws_iam_role_policy_name = %[5]q
		 		aws_iam_role_name = %[6]q
		 	  }

			   %[7]s

			   %[8]s		
	`, projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(s3BucketName1, s3BucketName2), pushBasedLogExportConfig(false, usePrefixPath, prefixPath, timeoutConfig, deleteOnCreateTimeout))
	return test
}

func configBasicUpdated(projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool) string {
	test := fmt.Sprintf(`
	 	locals {
				project_id = %[1]q
		 		s3_bucket_name_1 = %[2]q
				s3_bucket_name_2 = %[3]q
		 		s3_bucket_policy_name = %[4]q
		 		aws_iam_role_policy_name = %[5]q
		 		aws_iam_role_name = %[6]q
		 	  }

			   %[7]s

			   %[8]s
	`, projectID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(s3BucketName1, s3BucketName2), pushBasedLogExportConfig(true, usePrefixPath, prefixPath, "", nil)) // updating the S3 bucket to use for push-based log config
	return test
}

// pushBasedLogExportConfig returns config for mongodbatlas_push_based_log_export resource and data source.
// This method uses the project and S3 bucket created in awsIAMroleAuthAndS3Config()
func pushBasedLogExportConfig(useBucket2, usePrefixPath bool, prefixPath, timeoutConfig string, deleteOnCreateTimeout *bool) string {
	deleteOnCreateTimeoutAttr := ""
	if deleteOnCreateTimeout != nil {
		deleteOnCreateTimeoutAttr = fmt.Sprintf("delete_on_create_timeout = %[1]t", *deleteOnCreateTimeout)
	}
	bucketNameAttr := "bucket_name = aws_s3_bucket.log_bucket_1.bucket"
	if useBucket2 {
		bucketNameAttr = "bucket_name = aws_s3_bucket.log_bucket_2.bucket"
	}
	if usePrefixPath {
		return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
			project_id  = local.project_id
			%[1]s
			iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			prefix_path = %[2]q
			%[4]s
			%[5]s
		}
		
		%[3]s
		`, bucketNameAttr, prefixPath, pushBasedLogExportDataSourceConfig(), deleteOnCreateTimeoutAttr, timeoutConfig)
	}

	return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
		project_id  = local.project_id
		%[1]s
		iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
	}
	
	%[2]s
	`, bucketNameAttr, pushBasedLogExportDataSourceConfig())
}

func pushBasedLogExportDataSourceConfig() string {
	return `data "mongodbatlas_push_based_log_export" "test" {
		project_id    = local.project_id

		depends_on = ["mongodbatlas_push_based_log_export.test"]
	  }`
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
		if rs.Type == "mongodbatlas_push_based_log_export" {
			resp, _, err := acc.ConnV2().PushBasedLogExportApi.GetLogExport(context.Background(), rs.Primary.Attributes["project_id"]).Execute()
			if err == nil && *resp.State != pushbasedlogexport.UnconfiguredState {
				return fmt.Errorf("push-based log export for project_id %s still configured with state %s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["state"])
			}
			if err != nil {
				return fmt.Errorf("push-based log export for project_id %s still configured", rs.Primary.Attributes["project_id"])
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
		return rs.Primary.Attributes["project_id"], nil
	}
}
