package logintegration_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/atlas-sdk-go/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_log_integration.test"
	dataSourceName       = "data.mongodbatlas_log_integration.test"
	pluralDataSourceName = "data.mongodbatlas_log_integrations.test"
	prefixPath           = "prefix-path"
	datasourcesConfig    = `
		data "mongodbatlas_log_integration" "test" {
			project_id     = mongodbatlas_log_integration.test.project_id
			integration_id = mongodbatlas_log_integration.test.integration_id
		}

		data "mongodbatlas_log_integrations" "test" {
			project_id = mongodbatlas_log_integration.test.project_id
			depends_on = [mongodbatlas_log_integration.test]
		}
	`
)

var (
	logTypesMongoD = []string{"MONGOD"}
	logTypesMongoS = []string{"MONGOS"}
	logTypesAll    = []string{"MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"}
)

type s3Config struct {
	kmsKey            *string
	bucketName        string
	bucketPolicyName  string
	iamRoleName       string
	iamRolePolicyName string
	prefixPath        string
}

func TestAccLogIntegration_basicS3(t *testing.T) {
	var (
		projectID            = acc.ProjectIDExecution(t)
		s3BucketName         = acc.RandomS3BucketName()
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketName)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		kmsKey               = os.Getenv("AWS_KMS_KEY_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckAwsEnvBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicS3(projectID, logTypesMongoD, &s3Config{nil, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath}, true),
				Check:  checkBasicS3(logTypesMongoD, s3BucketName, prefixPath, true),
			},
			{
				Config: configBasicS3(projectID, logTypesAll, &s3Config{&kmsKey, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath}, false),
				Check:  checkBasicS3(logTypesAll, s3BucketName, prefixPath, false),
			},
			{
				Config: configBasicS3(projectID, logTypesMongoS, &s3Config{nil, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath}, false),
				Check:  checkBasicS3(logTypesMongoS, s3BucketName, prefixPath, false),
			},
			{
				Config:                               configBasicS3(projectID, logTypesMongoS, &s3Config{&kmsKey, s3BucketName, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, prefixPath}, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "integration_id",
			},
		},
	})
}

func checkBasicS3(logTypes []string, bucketName, prefixPath string, withDS bool) resource.TestCheckFunc {
	setChecks := []string{"iam_role_id", "integration_id"}
	mapChecks := map[string]string{
		"bucket_name": bucketName,
		"prefix_path": prefixPath,
		"type":        "S3_LOG_EXPORT",
		"log_types.#": strconv.Itoa(len(logTypes)),
		"log_types.0": logTypes[0],
	}
	checks := []resource.TestCheckFunc{}
	var dsName *string
	if withDS {
		dsName = admin.PtrString(dataSourceName)
		checks = append(checks, resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)))
	}
	checks = append(checks, acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configBasicS3(projectID string, logTypes []string, config *s3Config, withDS bool) string {
	logTypesStr := fmt.Sprintf("[%s]", `"`+strings.Join(logTypes, `", "`)+`"`)
	kmsKeyHCL := ""
	if config.kmsKey != nil {
		kmsKeyHCL = fmt.Sprintf("kms_key = %q", *config.kmsKey)
	}
	dsConfig := ""
	if withDS {
		dsConfig = datasourcesConfig
	}
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_log_integration" "test" {
			project_id  = %[2]q
			bucket_name = aws_s3_bucket.log_bucket.bucket
			iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			type        = "S3_LOG_EXPORT"
			log_types   = %[3]s
			prefix_path = %[4]q
			%[5]s
		}

		%[6]s
	`, awsIAMRoleAuthAndS3Config(projectID, config), projectID, logTypesStr, config.prefixPath, kmsKeyHCL, dsConfig)
}

func awsIAMRoleAuthAndS3Config(projectID string, config *s3Config) string {
	return fmt.Sprintf(`
		// Create IAM role & policy to authorize with Atlas
		resource "aws_iam_role_policy" "test_policy" {
		    name = %[4]q
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
								"arn:aws:s3:::%[2]s"
							]
						}
					]
				}
				EOF
		}

		resource "aws_iam_role" "test_role" {
		    name = %[1]q
		    max_session_duration = 43200

		    assume_role_policy = <<-EOF
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

		// Create S3 buckets
		resource "aws_s3_bucket" "log_bucket" {
		    bucket        = %[2]q
		    force_destroy = true  // required as atlas creates a test folder in the bucket when push-based log export is set up 

			lifecycle {
				ignore_changes = [tags, tags_all]
			}
		}

		// Add authorization policy to existing IAM role
		resource "aws_iam_role_policy" "s3_bucket_policy" {
		    name   = %[5]q
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
	`, projectID, config.bucketName, config.iamRoleName, config.iamRolePolicyName, config.bucketName)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().PushBasedLogExportApi.GetGroupLogIntegration(context.Background(), projectID, integrationID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("log integration for project_id %s with id %s does not exist", projectID, integrationID)
	}
}

func checkDestroy(state *terraform.State) error {
	for name, rs := range state.RootModule().Resources {
		if name != resourceName {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().PushBasedLogExportApi.GetGroupLogIntegration(context.Background(), projectID, integrationID).Execute()
		if err == nil {
			return fmt.Errorf("log integration for project_id %s with id %s still exists", projectID, integrationID)
		}
		return nil
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if projectID == "" || integrationID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, integrationID), nil
	}
}
