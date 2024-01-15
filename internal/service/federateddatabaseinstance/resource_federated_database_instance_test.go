package federateddatabaseinstance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedDatabaseInstance_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps(name, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage_stores.0.read_preference.0.tag_sets.#"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.0.tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_databases.0.collections.0.data_sources.0.database", "sample_airbnb"),
				),
			},
			{
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstStepsUpdate(name, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage_stores.0.read_preference.0.tag_sets.#"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_databases.0.collections.0.data_sources.0.database_regex", ".sample_airbnb"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"storage_stores.0.allow_insecure", "storage_stores.0.include_tags", "storage_stores.0.read_preference.0.max_staleness_seconds",
					"storage_stores.1.allow_insecure", "storage_stores.1.include_tags", "storage_stores.1.read_preference.0.max_staleness_seconds"},
			},
			{
				Config: testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstStepsUpdate(name, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_S3bucket(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		region       = "VIRGINIA_USA"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasFederatedDatabaseInstanceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:             resourceName,
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				ImportStateIdFunc:        testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFuncS3Bucket(resourceName, testS3Bucket),
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_atlasCluster(t *testing.T) {
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasFederatedDatabaseInstanceAtlasProviderConfig(projectName, orgID, name),
				Check: resource.ComposeTestCheckFunc(
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

func testAccMongoDBAtlasFederatedDatabaseInstanceAtlasProviderConfig(projectName, orgID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "project-tf" {
		name   = %[1]q
		org_id = %[2]q
	  }
	  
	  resource "mongodbatlas_cluster" "cluster-1" {
		project_id = mongodbatlas_project.project-tf.id
		provider_name               = "AWS"
		name                        = "tfCluster0"
		backing_provider_name       = "AWS"
		provider_region_name        = "EU_WEST_2"
		provider_instance_size_name = "M10"
	  }
	  
	  
	  resource "mongodbatlas_cluster" "cluster-2" {
		project_id = mongodbatlas_project.project-tf.id
		provider_name               = "AWS"
		name                        = "tfCluster1"
		backing_provider_name       = "AWS"
		provider_region_name        = "EU_WEST_2"
		provider_instance_size_name = "M10"
	  }

	  resource "mongodbatlas_federated_database_instance" "test" {
		project_id = mongodbatlas_project.project-tf.id
		name       = %[3]q
		storage_databases {
		  name = "VirtualDatabase0"
		  collections {
			name = "VirtualCollection0"
			data_sources {
			  collection = "listingsAndReviews"
			  database   = "sample_airbnb"
			  store_name = mongodbatlas_cluster.cluster-1.name
			}
			data_sources {

			  collection = "listingsAndReviews"
			  database   = "sample_airbnb"
			  store_name = mongodbatlas_cluster.cluster-2.name
			}
		  }
		}
	  
		storage_stores {
		  name         = mongodbatlas_cluster.cluster-1.name
		  cluster_name = mongodbatlas_cluster.cluster-1.name
		  project_id   = mongodbatlas_project.project-tf.id
		  provider     = "atlas"
		  read_preference {
			mode = "secondary"
			tag_sets {
				tags {
					name = "environment"
					value = "development"
				}
				tags {
					name = "application"
					value = "app"
				}
			}
			tag_sets {
				tags {
					name = "environment1"
					value = "development1"
				}
				tags {
					name = "application1"
					value = "app-1"
				}
			}
		  }
		}
	  
		storage_stores {
		  name         = mongodbatlas_cluster.cluster-2.name
		  cluster_name = mongodbatlas_cluster.cluster-2.name
		  project_id   = mongodbatlas_project.project-tf.id
		  provider     = "atlas"
		  read_preference {
			mode = "secondary"
			tag_sets {
				tags {
					name = "environment"
					value = "development"
				}
				tags {
					name = "application"
					value = "app"
				}
			}
			tag_sets {
				tags {
					name = "environment1"
					value = "development1"
				}
				tags {
					name = "application1"
					value = "app-1"
				}
			}
		  }
		}
	  }
	`, projectName, orgID, name)
}

func testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFuncS3Bucket(resourceName, s3Bucket string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["name"], s3Bucket), nil
	}
}

func testAccCheckMongoDBAtlasFederatedDatabaseInstanceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["name"]), nil
	}
}

func testAccMongoDBAtlasFederatedDatabaseInstanceConfigS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket, dataLakeRegion string) string {
	stepConfig := testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps3Bucket(name, testS3Bucket)
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
func testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps3Bucket(name, testS3Bucket string) string {
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
	}
   }
}
	`, name, testS3Bucket)
}

func testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstSteps(federatedInstanceName, projectName, orgID string) string {
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
				name = "environment"
				value = "development"
			}
			tags {
				name = "application"
				value = "app"
			}
		}
		tag_sets {
			tags {
				name = "environment1"
				value = "development1"
			}
			tags {
				name = "application1"
				value = "app-1"
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
				name = "environment"
				value = "development"
			}
			tags {
				name = "application"
				value = "app"
			}
		}
		tag_sets {
			tags {
				name = "environment1"
				value = "development1"
			}
			tags {
				name = "application1"
				value = "app-1"
			}
		}
	}
   }
}
	`, federatedInstanceName, projectName, orgID)
}

func testAccMongoDBAtlasFederatedDatabaseInstanceConfigFirstStepsUpdate(federatedInstanceName, projectName, orgID string) string {
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
					database_regex = ".sample_airbnb"
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
				name = "environment"
				value = "development"
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
				name = "environment"
				value = "development"
			}
		}
		tag_sets {
			tags {
				name = "environment"
				value = "development"
			}
		}
	}
   }
}
	`, federatedInstanceName, projectName, orgID)
}
