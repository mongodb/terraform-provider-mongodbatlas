package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFederatedDatabaseQueryLimits_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "data.mongodbatlas_federated_query_limits.test"
		// orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		// projectName  = acctest.RandomWithPrefix("test-acc-project")
		// tenantName   = acctest.RandomWithPrefix("test-acc-tenant-name")
		// limitName    = "bytesProcessed.monthly"

		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc-project")
		tenantName  = acctest.RandomWithPrefix("test-acc-tenant")
		// limitName   = "bytesProcessed.monthly"

		// name              = acctest.RandomWithPrefix("test-acc")
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
						VersionConstraint: "4.66.1",
						Source:            "hashicorp/aws",
					},
				},
				ProviderFactories: testAccProviderFactories,
				Config:            testAccMongoDBAtlasFederatedDatabaseQueryLimitsDataSourceConfig2(policyName, roleName, projectName, orgID, tenantName, testS3Bucket, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_name"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasFederatedDatabaseQueryLimitsDataSourceConfig2(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string) string {
	stepConfig := testAccMongoDBAtlasFederatedDatabaseQueryLimitsConfigDataSourceFirstStep(name, testS3Bucket)
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

func testAccMongoDBAtlasFederatedDatabaseQueryLimitsConfigDataSourceFirstStep(name, testS3Bucket string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_federated_database_instance" "db_instance" {
   project_id         = mongodbatlas_project.test_project.id
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
	region = "US_EAST_1"
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

  data "mongodbatlas_federated_query_limits" "test" {
	project_id = mongodbatlas_project.test_project.id
	tenant_name = mongodbatlas_federated_database_instance.db_instance.name
  }

	`, name, testS3Bucket)
}

// func testAccMongoDBAtlasFederatedDatabaseQueryLimitsDataSourceConfig(projectName, orgID, tenantName, limitName string) string {
// 	return fmt.Sprintf(`
// 	resource "mongodbatlas_project" "project-tf" {
// 		name   = %[1]q
// 		org_id = %[2]q
// 	 }

// 	 resource "mongodbatlas_cluster" "cluster-1" {
// 		project_id = mongodbatlas_project.project-tf.id
// 		provider_name               = "AWS"
// 		name                        = "tfCluster0"
// 		backing_provider_name       = "AWS"
// 		provider_region_name        = "US_EAST_1"
// 		provider_instance_size_name = "M10"
// 	  }

// 	  resource "mongodbatlas_cluster" "cluster-2" {
// 		project_id = mongodbatlas_project.project-tf.id
// 		provider_name               = "AWS"
// 		name                        = "tfCluster1"
// 		backing_provider_name       = "AWS"
// 		provider_region_name        = "US_EAST_1"
// 		provider_instance_size_name = "M10"
// 	  }

// 	  resource "mongodbatlas_federated_database_instance" "db-instance" {
// 		project_id = mongodbatlas_project.project-tf.id
// 		name       = %[3]q
// 		aws {
// 		  role_id        = ""
// 		  test_s3_bucket = ""
// 		}
// 		storage_databases {
// 		  name = "VirtualDatabase0"
// 		  collections {
// 			name = "VirtualCollection0"
// 			data_sources {
// 			  collection = "listingsAndReviews"
// 			  database   = "sample_airbnb"
// 			  store_name = mongodbatlas_cluster.cluster-1.name
// 			}
// 			data_sources {
// 			  collection = "listingsAndReviews"
// 			  database   = "sample_airbnb"
// 			  store_name = mongodbatlas_cluster.cluster-2.name
// 			}
// 		  }
// 		}

// 		storage_stores {
// 		  name         = mongodbatlas_cluster.cluster-1.name
// 		  cluster_name = mongodbatlas_cluster.cluster-1.name
// 		  project_id   = mongodbatlas_project.project-tf.id
// 		  provider     = "atlas"
// 		  read_preference {
// 			mode = "secondary"
// 		  }
// 		}

// 		storage_stores {
// 		  name         = mongodbatlas_cluster.cluster-2.name
// 		  cluster_name = mongodbatlas_cluster.cluster-2.name
// 		  project_id   = mongodbatlas_project.project-tf.id
// 		  provider     = "atlas"
// 		  read_preference {
// 			mode = "secondary"
// 		  }
// 		}
// 	  }

// 	  resource "mongodbatlas_federated_query_limit" "query-limit" {
// 		project_id = mongodbatlas_project.project-tf.id
// 		tenant_name = mongodbatlas_federated_database_instance.db-instance.name
// 		limit_name = %[4]q
// 		overrun_policy = "BLOCK"
// 		value = 5147483648
// 	  }

// 	  data "mongodbatlas_federated_query_limits" "test" {
// 		project_id = mongodbatlas_project.project-tf.id
// 		tenant_name = mongodbatlas_federated_database_instance.db-instance.name
// 	  }
// 	`, projectName, orgID, tenantName, limitName)
// }
