package federateddatabaseinstance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedDatabaseInstanceDSPlural_basic(t *testing.T) {
	var (
		resourceName = "data.mongodbatlas_federated_database_instances.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		firstName    = acc.RandomName()
		secondName   = acc.RandomName()
		policyName   = acc.RandomName()
		roleName     = acc.RandomIAMRole()
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t); acc.PreCheckS3Bucket(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configDSPlural(policyName, roleName, projectName, orgID, firstName, secondName, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func configDSPlural(policyName, roleName, projectName, orgID, firstName, secondName, testS3Bucket string) string {
	stepConfig := configDSPluralFirstStep(firstName, secondName, testS3Bucket)
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
func configDSPluralFirstStep(firstName, secondName, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_federated_database_instance" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q
   cloud_provider_config {
	aws {
		role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
		test_s3_bucket = %[3]q
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
					store_name = %[3]q
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

   storage_stores {
	bucket = %[3]q
	delimiter = "/"
	name = %[3]q
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

resource "mongodbatlas_federated_database_instance" "test2" {
	project_id         = mongodbatlas_project.test.id
	name = %[2]q
 
	storage_databases {
	 name = "VirtualDatabase0"
	 collections {
			 name = "VirtualCollection0"
			 data_sources {
					 collection = "listingsAndReviews"
					 database = "sample_airbnb"
					 store_name =  "ClusterTest2"
			 }
	 }
	}
 
	storage_stores {
	 name = "ClusterTest2"
	 cluster_name = "ClusterTest2"
	 project_id = mongodbatlas_project.test.id
	 provider = "atlas"
	 read_preference {
		 mode = "secondary"
	 }
	}
 
	storage_stores {
	 name = "dataStore0"
	 cluster_name = "ClusterTest2"
	 project_id = mongodbatlas_project.test.id
	 provider = "atlas"
	 read_preference {
		 mode = "secondary"
	 }
	}
 }

data "mongodbatlas_federated_database_instances" "test" {
	project_id           = mongodbatlas_federated_database_instance.test.project_id
}
	`, firstName, secondName, testS3Bucket)
}
