package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFederatedDatabaseQueryLimit_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_federated_query_limit.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-project")
		tenantName   = acctest.RandomWithPrefix("test-acc-tenant")
		limitName    = "bytesProcessed.monthly"
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		region       = "VIRGINIA_USA"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseQueryLimitDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseQueryLimitConfig(policyName, roleName, projectName, orgID, tenantName, testS3Bucket, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttr(resourceName, "limit_name", limitName),
				),
			},
			{
				ResourceName:      resourceName,
				ProviderFactories: testAccProviderFactories,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedDatabaseQueryLimitImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMongoDBAtlasFederatedDatabaseQueryLimitConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string) string {
	stepConfig := testAccMongoDBAtlasFederatedDatabaseQueryLimitConfigFirstStep(name, testS3Bucket)
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
		"Action": "*",
		"Resource": "*"
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

%s
	`, policyName, roleName, projectName, orgID, stepConfig)
}

func testAccMongoDBAtlasFederatedDatabaseQueryLimitConfigFirstStep(name, testS3Bucket string) string {
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
	`, name, testS3Bucket)
}

func testAccCheckMongoDBAtlasFederatedDatabaseQueryLimitImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["tenant_name"], ids["limit_name"]), nil
	}
}

func testAccCheckMongoDBAtlasFederatedDatabaseQueryLimitDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_federated_query_limit" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.DataFederation.GetQueryLimit(context.Background(), ids["project_id"], ids["tenant_name"], ids["limit_name"])
		if err == nil {
			return fmt.Errorf("federated database query limit (%s) for project (%s) and tenant (%s)still exists", ids["project_id"], ids["tenant_name"], ids["limit_name"])
		}
	}

	return nil
}
