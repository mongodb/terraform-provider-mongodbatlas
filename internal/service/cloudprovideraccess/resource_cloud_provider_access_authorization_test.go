package cloudprovideraccess_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudprovideraccess"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccCloudProviderAccessAuthorizationAWS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicAuthorizationTestCase(t))
}

func TestAccCloudProviderAccessAuthorizationAzure_basic(t *testing.T) {
	var (
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		tenantID           = os.Getenv("AZURE_TENANT_ID")
		projectID          = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCloudProviderAccessAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAuthorizationAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID),
			},
		},
	},
	)
}

func TestAccCloudProviderAccessAuthorizationGCP_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_provider_access_authorization.auth_role"
		projectID    = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAuthorizationGCP(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "gcp.0.service_account_for_atlas"),
					resource.TestCheckResourceAttrSet(resourceName, "gcp.0.status"),
				),
			},
		},
	})
}

func configAuthorizationGCP(projectID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "gcp_setup"{
		project_id    = %[1]q
		provider_name = "GCP"
	}

	resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
		project_id = mongodbatlas_cloud_provider_access_setup.gcp_setup.project_id
		role_id    = mongodbatlas_cloud_provider_access_setup.gcp_setup.role_id
	}
	`, projectID)
}

func basicAuthorizationTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		resourceName                  = "mongodbatlas_cloud_provider_access_authorization.auth_role"
		projectID                     = acc.ProjectIDExecution(tb)
		policyName                    = acc.RandomName()
		roleName                      = acc.RandomIAMRole()
		roleNameUpdated               = acc.RandomIAMRole()
		federatedDatabaseInstanceName = acc.RandomName()
		testS3Bucket                  = os.Getenv("AWS_S3_BUCKET")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAuthorizationAWS(projectID, policyName, roleName, federatedDatabaseInstanceName, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.0.iam_assumed_role_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "feature_usages.#"),
				),
			},
			{
				Config: configAuthorizationAWS(projectID, policyName, roleNameUpdated, federatedDatabaseInstanceName, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "aws.0.iam_assumed_role_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "feature_usages.#"),
				),
			},
		},
	}
}

func configAuthorizationAWS(projectID, policyName, roleName, federatedDatabaseInstanceName, testS3Bucket string) string {
	bucketResourceName := "arn:aws:s3:::" + testS3Bucket
	return fmt.Sprintf(`

resource "mongodbatlas_federated_database_instance" "test" {
	project_id         = %[1]q
	name = %[4]q

	cloud_provider_config {
	    aws {
			role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			test_s3_bucket = %[5]q
		}
	}
}

resource "aws_iam_role_policy" "test_policy" {
  name = %[2]q
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
			"Resource": %[6]q
		}
    ]
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
        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_assumed_role_external_id}"
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
	`, projectID, policyName, roleName, federatedDatabaseInstanceName, testS3Bucket, bucketResourceName)
}

func configAuthorizationAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = %[1]q
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[2]q
			service_principal_id = %[3]q
			tenant_id = %[4]q
		}
	 }

   resource "mongodbatlas_cloud_provider_access_authorization" "test" {
		project_id = %[1]q
        role_id = mongodbatlas_cloud_provider_access_setup.test.role_id
		azure {
			atlas_azure_app_id = %[2]q
			service_principal_id = %[3]q
			tenant_id = %[4]q
		}
	 }
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID)
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_access" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)

		id := ids["id"]
		role, _, err := acc.ConnV2().CloudProviderAccessApi.GetCloudProviderAccessRole(context.Background(), ids["project_id"], id).Execute()
		if err != nil {
			return fmt.Errorf(cloudprovideraccess.ErrorGetRead, err)
		}
		if role.GetId() == id || role.GetRoleId() == id {
			return nil
		}
	}
	return nil
}
