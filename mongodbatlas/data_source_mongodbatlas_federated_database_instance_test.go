package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceFederatedDatabaseInstance_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName      = "data.mongodbatlas_federated_database_instance.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acctest.RandomWithPrefix("test-acc")
		name              = acctest.RandomWithPrefix("test-acc")
		policyName        = acctest.RandomWithPrefix("test-acc")
		roleName          = acctest.RandomWithPrefix("test-acc")
		testS3Bucket      = os.Getenv("AWS_S3_BUCKET")
		region            = "VIRGINIA_USA"
		federatedInstance = matlas.DataFederationInstance{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "4.66.1",
						Source:            "hashicorp/aws",
					},
				},
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseInstanceDataSourceConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedDatabaseDataSourceInstanceExists(resourceName, &federatedInstance),
					testAccCheckMongoDBAtlasFederatedDabaseInstanceAttributes(&federatedInstance, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasFederatedDatabaseDataSourceInstanceExists(resourceName string, dataFederatedInstance *matlas.DataFederationInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		if dataLakeResp, _, err := conn.DataFederation.Get(context.Background(), ids["project_id"], ids["name"]); err == nil {
			*dataFederatedInstance = *dataLakeResp
			return nil
		}

		return fmt.Errorf("federated database instance (%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasFederatedDabaseInstanceAttributes(dataFederatedInstance *matlas.DataFederationInstance, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dataFederatedInstance.Name: %s , username : %s", dataFederatedInstance.Name, name)
		if dataFederatedInstance.Name != name {
			return fmt.Errorf("bad data federated instance name: %s", dataFederatedInstance.Name)
		}

		return nil
	}
}

func testAccMongoDBAtlasFederatedDatabaseInstanceDataSourceConfig(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string) string {
	stepConfig := testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstStep(name, testS3Bucket)
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
resource "mongodbatlas_project" "test" {
   name   = %[3]q
   org_id = %[4]q
}


resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = mongodbatlas_project.test.id
   provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id = mongodbatlas_project.test.id
   role_id =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   aws {
      iam_assumed_role_arn = aws_iam_role.test_role.arn
   }
}

%s
	`, policyName, roleName, projectName, orgID, stepConfig)
}
func testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstStep(name, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_federated_database_instance" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q
   aws {
     role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
     test_s3_bucket = %[2]q
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
	project_id = mongodbatlas_project.test.id
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
	project_id = mongodbatlas_project.test.id
	provider = "atlas"
	read_preference {
		mode = "secondary"
	}
   }
}

data "mongodbatlas_federated_database_instance" "test" {
	project_id           = mongodbatlas_federated_database_instance.test.project_id
	name = mongodbatlas_federated_database_instance.test.name
}
	`, name, testS3Bucket)
}
