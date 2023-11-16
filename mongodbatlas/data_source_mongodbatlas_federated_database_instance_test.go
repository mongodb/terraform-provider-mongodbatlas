package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20231115001/admin"
)

func TestAccDataSourceFederatedDatabaseInstance_basic(t *testing.T) {
	var (
		resourceName      = "data.mongodbatlas_federated_database_instance.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acctest.RandomWithPrefix("test-acc")
		name              = acctest.RandomWithPrefix("test-acc")
		federatedInstance = admin.DataLakeTenant{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstSteps(name, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedDatabaseDataSourceInstanceExists(resourceName, &federatedInstance),
					testAccCheckMongoDBAtlasFederatedDabaseInstanceAttributes(&federatedInstance, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage_stores.0.read_preference.0.tag_sets.#"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.0.tags.#", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceFederatedDatabaseInstance_S3Bucket(t *testing.T) {
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
		federatedInstance = admin.DataLakeTenant{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasFederatedDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasFederatedDatabaseInstanceDataSourceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, region),
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

func testAccCheckMongoDBAtlasFederatedDatabaseDataSourceInstanceExists(resourceName string, dataFederatedInstance *admin.DataLakeTenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		connV2 := testAccProviderSdkV2.Meta().(*MongoDBClient).AtlasV2

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		if dataLakeResp, _, err := connV2.DataFederationApi.GetFederatedDatabase(context.Background(), ids["project_id"], ids["name"]).Execute(); err == nil {
			*dataFederatedInstance = *dataLakeResp
			return nil
		}

		return fmt.Errorf("federated database instance (%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasFederatedDabaseInstanceAttributes(dataFederatedInstance *admin.DataLakeTenant, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dataFederatedInstance.Name: %s , username : %s", dataFederatedInstance.GetName(), name)
		if dataFederatedInstance.GetName() != name {
			return fmt.Errorf("bad data federated instance name: %s", dataFederatedInstance.GetName())
		}
		return nil
	}
}

func testAccMongoDBAtlasFederatedDatabaseInstanceDataSourceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string) string {
	stepConfig := testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstStepS3Bucket(name, testS3Bucket)
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
func testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstStepS3Bucket(name, testS3Bucket string) string {
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

func testAccMongoDBAtlasFederatedDatabaseInstanceConfigDataSourceFirstSteps(federatedInstanceName, projectName, orgID string) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[3]q
	}

resource "mongodbatlas_federated_database_instance" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q

   storage_databases {
	name = "VirtualDatabase0"
	collections {
			name = "VirtualCollection0"
			data_sources {
					collection = "listingsAndReviews"
					database = "sample_airbnb"
					store_name =  "ClusterTest"
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
	}
   }
}

data "mongodbatlas_federated_database_instance" "test" {
	project_id           = mongodbatlas_federated_database_instance.test.project_id
	name = mongodbatlas_federated_database_instance.test.name
}
	`, federatedInstanceName, projectName, orgID)
}
