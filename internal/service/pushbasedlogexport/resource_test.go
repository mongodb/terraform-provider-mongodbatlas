package pushbasedlogexport_test

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
	resourceName       = "mongodbatlas_push_based_log_export.test"
	datasourceName     = "data.mongodbatlas_push_based_log_export.test"
	nonEmptyPrefixPath = "push-log-prefix"
	defaultPrefixPath  = ""
)

func TestAccPushBasedLogExport_basic(t *testing.T) {
	var (
		orgID                = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName          = acc.RandomProjectName()
		s3BucketNamePrefix   = fmt.Sprintf("tf-%s", acc.RandomName())
		s3BucketName1        = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2        = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				Check:  resource.ComposeTestCheckFunc(commonChecks(s3BucketName1, nonEmptyPrefixPath)...),
			},
			{
				Config: configBasicUpdated(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				Check:  resource.ComposeTestCheckFunc(commonChecks(s3BucketName2, nonEmptyPrefixPath)...),
			},
			{
				Config:                               configBasicUpdated(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
}

func TestAccPushBasedLogExport_noPrefixPath(t *testing.T) {
	var (
		orgID                = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName          = acc.RandomProjectName()
		s3BucketNamePrefix   = fmt.Sprintf("tf-%s", acc.RandomName())
		s3BucketName1        = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2        = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, defaultPrefixPath, false),
				Check:  resource.ComposeTestCheckFunc(commonChecks(s3BucketName1, defaultPrefixPath)...),
			},
		},
	})
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

// func addAttrSetChecks(checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
// 	checks = acc.AddAttrSetChecks(resourceName, checks, attrNames...)
// 	return acc.AddAttrSetChecks(datasourceName, checks, attrNames...)
// }

func configBasic(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool) string {
	test := fmt.Sprintf(`
	 	locals {
		 		project_name = %[1]q
		 		org_id = %[2]q
		 		s3_bucket_name_1 = %[3]q
				s3_bucket_name_2 = %[4]q
		 		s3_bucket_policy_name = %[5]q
		 		aws_iam_role_policy_name = %[6]q
		 		aws_iam_role_name = %[7]q
		 	  }

			   %[8]s

			   %[9]s
	`, projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(), pushBasedLogExportConfig(false, usePrefixPath, prefixPath))
	return test
}

func configBasicUpdated(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath string, usePrefixPath bool) string {
	test := fmt.Sprintf(`
	 	locals {
		 		project_name = %[1]q
		 		org_id = %[2]q
		 		s3_bucket_name_1 = %[3]q
				s3_bucket_name_2 = %[4]q
		 		s3_bucket_policy_name = %[5]q
		 		aws_iam_role_policy_name = %[6]q
		 		aws_iam_role_name = %[7]q
		 	  }

			   %[8]s

			   %[9]s
	`, projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName,
		awsIAMroleAuthAndS3Config(), pushBasedLogExportConfig(true, usePrefixPath, prefixPath)) // updating the S3 bucket to use for push-based log config
	return test
}

// func localVarsConfig(projectName, orgID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName string) string {
// 	return fmt.Sprintf(`
// 	locals {
// 		project_name = %[1]q
// 		org_id = %[2]q
// 		s3_bucket_name = %[3]q
// 		s3_bucket_policy_name = %[4]q
// 		aws_iam_role_policy_name = %[5]q
// 		aws_iam_role_name = %[6]q
// 	  }
// `, projectName, orgID, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName)
// }

// pushBasedLogExportConfig returns config for mongodbatlas_push_based_log_export resource and data source.
// This method uses the project and S3 bucket created in awsIAMroleAuthAndS3Config()
func pushBasedLogExportConfig(useBucket2, usePrefixPath bool, prefixPath string) string {
	bucketNameAttr := "bucket_name = aws_s3_bucket.log_bucket_1.bucket"
	if useBucket2 {
		bucketNameAttr = "bucket_name = aws_s3_bucket.log_bucket_2.bucket"
	}
	if usePrefixPath {
		return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
			project_id  = mongodbatlas_project.project-tf.id
			%[1]s
			iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			prefix_path = %[2]q
		}
		
		%[3]s
		`, bucketNameAttr, prefixPath, pushBasedLogExportDataSourceConfig())
	}

	return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
		project_id  = mongodbatlas_project.project-tf.id
		%[1]s
		iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
	}
	
	%[2]s
	`, bucketNameAttr, pushBasedLogExportDataSourceConfig())
}

// func pushBasedLogExportConfigUpdatedBucket(usePrefixPath bool, prefixPath string) string {
// 	if usePrefixPath {
// 		return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
// 			project_id  = mongodbatlas_project.project-tf.id
// 			bucket_name = aws_s3_bucket.log_bucket_2.bucket
// 			iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
// 			prefix_path = %[1]q
// 		}

// 		%[2]s
// 		`, prefixPath, pushBasedLogExportDataSourceConfig())
// 	}

// 	return fmt.Sprintf(`resource "mongodbatlas_push_based_log_export" "test" {
// 		project_id  = mongodbatlas_project.project-tf.id
// 		bucket_name = aws_s3_bucket.log_bucket_2.bucket
// 		iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
// 	}

// 	%[1]s
// 	`, pushBasedLogExportDataSourceConfig())
// }

func pushBasedLogExportDataSourceConfig() string {
	return `data "mongodbatlas_push_based_log_export" "test" {
		project_id    = mongodbatlas_project.project-tf.id

		depends_on = ["mongodbatlas_push_based_log_export.test"]
	  }`
}

// awsIAMroleAuthAndS3Config returns config for required IAM roles and authorizes them (sets up cloud provider access) with a mongodbatlas_project
// This method also creates two S3 buckets and sets up required access policy for them
func awsIAMroleAuthAndS3Config() string {
	return `
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
		"Action": "*",
		"Resource": "*"
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
resource "mongodbatlas_project" "project-tf" {
  name     = local.project_name
  org_id = local.org_id
}

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project-tf.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id    = mongodbatlas_project.project-tf.id
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
		`
}

// func awsS3Config() string {
// 	return `
// 	resource "aws_s3_bucket" "log_bucket" {
// 		bucket = local.s3_bucket_name
// 		force_destroy = true  # required as atlas creates a test folder in the bucket when mongodbatlas_push_based_log_export is set up
// 	  }

// 	  resource "aws_iam_role_policy" "s3_bucket_policy" {
// 		name = local.s3_bucket_policy_name
// 		role = aws_iam_role.test_role.id

// 		policy = <<-EOF
// 		{
// 		  "Version": "2012-10-17",
// 		  "Statement": [
// 			  {
// 				  "Effect": "Allow",
// 				  "Action": [
// 					  "s3:ListBucket",
// 					  "s3:PutObject",
// 					  "s3:GetObject",
// 					  "s3:GetBucketLocation"
// 				  ],
// 				  "Resource": [
// 					  "arn:aws:s3:::maastha-test",
// 					  "arn:aws:s3:::maastha-test/*"
// 				  ]
// 			  }
// 		  ]
// 	  }
// 		EOF
// 	  }
// 	`
// }

// func atlasProjectCloudProviderSetupConfig() string {
// 	return `
// 	resource "mongodbatlas_project" "project-tf" {
// 		name     = local.project_name
// 		org_id = local.org_id
// 	}

// 	resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
// 		project_id    = mongodbatlas_project.project-tf.id
// 		provider_name = "AWS"
// 	}

// 	resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
// 		project_id    = mongodbatlas_project.project-tf.id
// 		role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

// 		aws {
// 		  iam_assumed_role_arn = aws_iam_role.test_role.arn
// 		}
// 	}`
// }

// func awsIAMRoleConfig() string {
// 	return `
// resource "aws_iam_role_policy" "test_policy" {
//   name = local.aws_iam_role_policy_name
//   role = aws_iam_role.test_role.id

//   policy = <<-EOF
//   {
//     "Version": "2012-10-17",
//     "Statement": [
//       {
//         "Effect": "Allow",
// 		"Action": "*",
// 		"Resource": "*"
//       }
//     ]
//   }
//   EOF
// }

// resource "aws_iam_role" "test_role" {
//   name = local.aws_iam_role_name
//   max_session_duration = 43200

//   assume_role_policy = <<EOF
// {
//   "Version": "2012-10-17",
//   "Statement": [
//     {
//       "Effect": "Allow",
//       "Principal": {
//         "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
//       },
//       "Action": "sts:AssumeRole",
//       "Condition": {
//         "StringEquals": {
//           "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
//         }
//       }
//     }
//   ]
// }
// EOF
// }`
// }

func checkDestroy(state *terraform.State) error {
	if projectDestroyedErr := acc.CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_push_based_log_export" {
			_, _, err := acc.ConnV2().PushBasedLogExportApi.GetPushBasedLogConfiguration(context.Background(), rs.Primary.Attributes["project_id"]).Execute()
			if err == nil {
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
