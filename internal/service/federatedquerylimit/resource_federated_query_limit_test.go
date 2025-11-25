package federatedquerylimit_test

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
	resourceName         = "mongodbatlas_federated_query_limit.test"
	dataSourceName       = "data.mongodbatlas_federated_query_limit.test"
	dataSourcePluralName = "data.mongodbatlas_federated_query_limits.test"
)

func TestAccFederatedDatabaseQueryLimit_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs an S3 bucket

	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		limitName    = "bytesProcessed.monthly"
		projectName  = acc.RandomProjectName()
		tenantName   = acc.RandomName()
		policyName   = acc.RandomName()
		roleName     = acc.RandomIAMRole()
	)

	return &resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(tb); acc.PreCheckS3Bucket(tb) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configBasic(policyName, roleName, projectName, orgID, tenantName, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttr(resourceName, "limit_name", limitName),

					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_name"),
					resource.TestCheckResourceAttr(dataSourceName, "limit_name", limitName),

					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "tenant_name"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
				),
			},
			{
				ResourceName:             resourceName,
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				ImportStateIdFunc:        importStateIDFunc(resourceName),
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	}
}

func configBasic(policyName, roleName, projectName, orgID, name, testS3Bucket string) string {
	stepConfig := configFirstStep(name, testS3Bucket)
	bucketResourceName := "arn:aws:s3:::" + testS3Bucket
	return fmt.Sprintf(`
resource "aws_iam_role_policy" "test_policy" {
  name = %[1]q
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
  name = %[2]q

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
resource "mongodbatlas_project" "test_project" {
   name   = %[3]q
   org_id = %[4]q
}


resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = mongodbatlas_project.test_project.id
   provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id = mongodbatlas_project.test_project.id
   role_id =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   aws {
      iam_assumed_role_arn = aws_iam_role.test_role.arn
   }
}

%[5]s
    `, policyName, roleName, projectName, orgID, stepConfig, bucketResourceName)
}

func configFirstStep(name, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_federated_database_instance" "db_instance" {
   project_id         = mongodbatlas_project.test_project.id
   name = %[1]q
   cloud_provider_config {
		aws {
			role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
			test_s3_bucket = %[2]q
		}
	}

   storage_databases {
	name = "VirtualDatabase0"
	collections {
			name = "VirtualCollection0"
			data_sources {
					collection = "listingsAndReviews"
					database = "sample_airbnb"
					store_name =  "ClusterTest"
			}
			data_sources {
					store_name = %[2]q
					path = "/{fileName string}.yaml"
			}
	}
   }

   storage_stores {
	name = "ClusterTest"
	cluster_name = "ClusterTest"
	project_id = mongodbatlas_project.test_project.id
	provider = "atlas"
	read_preference {
		mode = "secondary"
	}
   }

   storage_stores {
	bucket = %[2]q
	delimiter = "/"
	name = %[2]q
	prefix = "templates/"
	provider = "s3"
	region = "EU_WEST_1"
   }

   storage_stores {
	name = "dataStore0"
	cluster_name = "ClusterTest"
	project_id = mongodbatlas_project.test_project.id
	provider = "atlas"
	read_preference {
		mode = "secondary"
	}
   }
}

resource "mongodbatlas_federated_query_limit" "test" {
	project_id = mongodbatlas_project.test_project.id
	tenant_name = mongodbatlas_federated_database_instance.db_instance.name
	limit_name = "bytesProcessed.monthly"
	overrun_policy = "BLOCK"
	value = 5147483648
}

data "mongodbatlas_federated_query_limit" "test" {
	project_id = mongodbatlas_project.test_project.id
	tenant_name = mongodbatlas_federated_database_instance.db_instance.name
	limit_name = "bytesProcessed.monthly"
}

data "mongodbatlas_federated_query_limits" "test" {
	project_id = mongodbatlas_project.test_project.id
	tenant_name = mongodbatlas_federated_database_instance.db_instance.name
}
    `, name, testS3Bucket)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["tenant_name"], ids["limit_name"]), nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_federated_query_limit" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataFederationApi.GetDataFederationLimit(context.Background(), ids["project_id"], ids["tenant_name"], ids["limit_name"]).Execute()
		if err == nil {
			return fmt.Errorf("federated database query limit (%s) for project (%s) and tenant (%s)still exists", ids["project_id"], ids["tenant_name"], ids["limit_name"])
		}
	}
	return nil
}
