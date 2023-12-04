package cloudprovideraccess_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSCloudProviderAccessAuthorizationAWS_basic(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		policyName      = acctest.RandomWithPrefix("tf-acc")
		roleName        = acctest.RandomWithPrefix("tf-acc")
		roleNameUpdated = acctest.RandomWithPrefix("tf-acc")
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.PreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"aws": {
				VersionConstraint: "5.1.0",
				Source:            "hashicorp/aws",
			},
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		// same as regular cloud provider access resource
		CheckDestroy: testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessAuthorizationConfig(projectID, policyName, roleName),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderAccessAuthorizationConfig(projectID, policyName, roleNameUpdated),
			},
		},
	},
	)
}

func TestAccConfigRSCloudProviderAccessAuthorizationAzure_basic(t *testing.T) {
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acctest.RandomWithPrefix("tf-acc")
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		tenantID           = os.Getenv("AZURE_TENANT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCloudProviderAccessAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasProviderAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderAccessAuthorizationAzure(orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID),
			},
		},
	},
	)
}

func testAccMongoDBAtlasCloudProviderAccessAuthorizationConfig(projectID, roleName, policyName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role_policy" "test_policy" {
  name = %[2]q
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
	`, projectID, policyName, roleName)
}

func testAccMongoDBAtlasCloudProviderAccessAuthorizationAzure(orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_project.test.id
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[3]q
			service_principal_id = %[4]q
			tenant_id = %[5]q
		}
	 }

   resource "mongodbatlas_cloud_provider_access_authorization" "test" {
		project_id = mongodbatlas_project.test.id
    role_id = mongodbatlas_cloud_provider_access_setup.test.role_id
		azure {
			atlas_azure_app_id = %[3]q
			service_principal_id = %[4]q
			tenant_id = %[5]q
		}
	 }
	`, orgID, projectName, atlasAzureAppID, servicePrincipalID, tenantID)
}
