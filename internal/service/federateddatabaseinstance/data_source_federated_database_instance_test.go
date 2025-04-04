package federateddatabaseinstance_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func TestAccFederatedDatabaseInstanceDS_s3Bucket(t *testing.T) {
	var (
		resourceName      = "data.mongodbatlas_federated_database_instance.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acc.RandomProjectName()
		name              = acc.RandomName()
		policyName        = acc.RandomName()
		roleName          = acc.RandomIAMRole()
		testS3Bucket      = os.Getenv("AWS_S3_BUCKET")
		federatedInstance = admin.DataLakeTenant{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t); acc.PreCheckS3Bucket(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configDSWithS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName, &federatedInstance),
					checkAttributes(&federatedInstance, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func checkExists(resourceName string, dataFederatedInstance *admin.DataLakeTenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if dataLakeResp, _, err := acc.ConnV2().DataFederationApi.GetFederatedDatabase(context.Background(), ids["project_id"], ids["name"]).Execute(); err == nil {
			*dataFederatedInstance = *dataLakeResp
			return nil
		}
		return fmt.Errorf("federated database instance (%s) does not exist", ids["project_id"])
	}
}

func checkAttributes(dataFederatedInstance *admin.DataLakeTenant, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dataFederatedInstance.Name: %s , username : %s", dataFederatedInstance.GetName(), name)
		if dataFederatedInstance.GetName() != name {
			return fmt.Errorf("bad data federated instance name: %s", dataFederatedInstance.GetName())
		}
		return nil
	}
}

func configDSWithS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket string) string {
	stepConfig := configDSFirstStepS3Bucket(name, testS3Bucket)
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

%[5]s
	`, policyName, roleName, projectName, orgID, stepConfig, bucketResourceName)
}
func configDSFirstStepS3Bucket(name, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_federated_database_instance" "test" {
   project_id         = mongodbatlas_project.test.id
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
		tag_sets {
			tags {
				name = "environment0"
				value = "development0"
			}
			tags {
				name = "application0"
				value = "app0"
			}
		}
		tag_sets {
			tags {
				name = "environment1"
				value = "development1"
			}
			tags {
				name = "application1"
				value = "app1"
			}
		}
	}
   }
}

data "mongodbatlas_federated_database_instance" "test" {
	project_id           = mongodbatlas_federated_database_instance.test.project_id
	name = mongodbatlas_federated_database_instance.test.name

	cloud_provider_config {
		aws {
			test_s3_bucket = %[2]q
		  }
	}
}
	`, name, testS3Bucket)
}
