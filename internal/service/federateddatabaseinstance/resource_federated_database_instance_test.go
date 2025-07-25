package federateddatabaseinstance_test

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_federated_database_instance.test"
	dataSourceName = "data.mongodbatlas_federated_database_instance.test"
)

func TestAccFederatedDatabaseInstance_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		name        = acc.RandomName()
	)

	valueChecks := map[string]string{
		"name":                                                      name,
		"data_process_region.0.cloud_provider":                      "AWS",
		"data_process_region.0.region":                              "OREGON_USA",
		"storage_stores.0.read_preference.0.tag_sets.#":             "2",
		"storage_stores.0.read_preference.0.tag_sets.0.tags.#":      "2",
		"storage_databases.0.collections.0.data_sources.0.database": "sample_airbnb",
	}
	setChecks := []string{"project_id", "storage_stores.0.read_preference.0.tag_sets.#"}
	firstStepChecks := acc.AddAttrChecks(resourceName, nil, valueChecks)
	firstStepChecks = acc.AddAttrSetChecks(resourceName, firstStepChecks, setChecks...)
	firstStepChecks = acc.AddAttrChecks(dataSourceName, firstStepChecks, valueChecks)
	firstStepChecks = acc.AddAttrSetChecks(dataSourceName, firstStepChecks, setChecks...)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				Config: configFirstSteps(name, projectName, orgID),
				Check:  resource.ComposeAggregateTestCheckFunc(firstStepChecks...),
			},
			{
				Config: configFirstStepsUpdate(name, projectName, orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
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
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"storage_stores.0.allow_insecure", "storage_stores.0.include_tags", "storage_stores.0.read_preference.0.max_staleness_seconds",
					"storage_stores.1.allow_insecure", "storage_stores.1.include_tags", "storage_stores.1.read_preference.0.max_staleness_seconds"},
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_s3bucket(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		name         = acc.RandomName()
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
				Config:                   configWithS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:             resourceName,
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				ImportStateIdFunc:        importStateIDFuncS3Bucket(resourceName, testS3Bucket),
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_azureCloudProviderConfig(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		name               = acc.RandomName()
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		tenantID           = os.Getenv("AZURE_TENANT_ID")
	)

	extraChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "cloud_provider_config.0.aws.#", "0"),
		resource.TestCheckResourceAttrSet(resourceName, "cloud_provider_config.0.azure.0.role_id"),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCloudProviderAccessAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				Config: configAzureCloudProvider(projectID, name, atlasAzureAppID, servicePrincipalID, tenantID),
				Check: checkAttrs(
					projectID,
					name,
					map[string]string{
						"cloud_provider_config.0.azure.0.atlas_app_id":         atlasAzureAppID,
						"cloud_provider_config.0.azure.0.service_principal_id": servicePrincipalID,
						"cloud_provider_config.0.azure.0.tenant_id":            tenantID,
					},
					extraChecks...,
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFederatedDatabaseInstance_atlasCluster(t *testing.T) {
	var (
		specs = []acc.ReplicationSpecRequest{
			{Region: "EU_WEST_2"},
		}
		clusterRequest = acc.ClusterRequest{
			ReplicationSpecs: specs,
		}
		name            = acc.RandomName()
		clusterInfo     = acc.GetClusterInfo(t, &clusterRequest)
		projectID       = clusterInfo.ProjectID
		clusterRequest2 = acc.ClusterRequest{
			ProjectID:        projectID,
			ReplicationSpecs: specs,
			ResourceSuffix:   "cluster2",
		}
		cluster2Info        = acc.GetClusterInfo(t, &clusterRequest2)
		dependencyTerraform = fmt.Sprintf("%s\n%s", clusterInfo.TerraformStr, cluster2Info.TerraformStr)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithCluster(dependencyTerraform, projectID, clusterInfo.ResourceName, cluster2Info.ResourceName, name),
				Check: resource.ComposeAggregateTestCheckFunc(
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

func configWithCluster(terraformStr, projectID, cluster1ResourceName, cluster2ResourceName, name string) string {
	return fmt.Sprintf(`
	  %[1]s

	  resource "mongodbatlas_federated_database_instance" "test" {
			project_id = %[2]q
			name       = %[5]q
			storage_databases {
				name = "VirtualDatabase0"
				collections {
				name = "VirtualCollection0"
				data_sources {
					collection = "listingsAndReviews"
					database   = "sample_airbnb"
					store_name = %[3]s.name
				}
				data_sources {

					collection = "listingsAndReviews"
					database   = "sample_airbnb"
					store_name = %[4]s.name
				}
			}
		}
	  
		storage_stores {
		  name         = %[3]s.name
		  cluster_name = %[3]s.name
		  project_id   = %[2]q
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
		  name         = %[4]s.name
		  cluster_name = %[4]s.name
		  project_id   = %[2]q
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
	`, terraformStr, projectID, cluster1ResourceName, cluster2ResourceName, name)
}

func importStateIDFuncS3Bucket(resourceName, s3Bucket string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["name"], s3Bucket), nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["name"]), nil
	}
}

func configWithS3Bucket(policyName, roleName, projectName, orgID, name, testS3Bucket string) string {
	stepConfig := configFirstStepS3Bucket(name, testS3Bucket)
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

func configFirstStepS3Bucket(name, testS3Bucket string) string {
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

func configAzureCloudProvider(projectID, name, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	azureCloudProviderAccess := acc.ConfigSetupAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID)

	return azureCloudProviderAccess + fmt.Sprintf(`

resource "mongodbatlas_federated_database_instance" "test" {
  project_id = %[1]q
  name       = %[2]q
  
  cloud_provider_config {
    azure {
      role_id = mongodbatlas_cloud_provider_access_setup.test.role_id
    }
  }

  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "VirtualCollection0"
      data_sources {
        collection = "listingsAndReviews"
        database   = "sample_airbnb"
        store_name = "ClusterTest"
      }
    }
  }

  storage_stores {
    name         = "ClusterTest"
    cluster_name = "ClusterTest"
    project_id   = %[1]q
    provider     = "atlas"
    read_preference {
      mode = "secondary"
    }
  }
}

data "mongodbatlas_federated_database_instance" "test" {
  project_id = mongodbatlas_federated_database_instance.test.project_id
  name       = mongodbatlas_federated_database_instance.test.name
}
`, projectID, name)
}

func checkAttrs(projectID, name string, extraAttrs map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attrsMap := map[string]string{
		"project_id": projectID,
		"name":       name,

		"cloud_provider_config.#":         "1",
		"cloud_provider_config.0.azure.#": "1",
	}

	maps.Copy(attrsMap, extraAttrs)
	check := acc.CheckRSAndDS(resourceName, conversion.Pointer(dataSourceName), nil, nil, attrsMap, extra...)
	checks := slices.Concat(extra, []resource.TestCheckFunc{check})
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configFirstSteps(federatedInstanceName, projectName, orgID string) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[3]q
	}

resource "mongodbatlas_federated_database_instance" "test" {
    project_id         = mongodbatlas_project.test.id
    name = %[1]q

	data_process_region {
		cloud_provider = "AWS"
    	region         = "OREGON_USA"
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

data "mongodbatlas_federated_database_instance" "test" {
	project_id           = mongodbatlas_federated_database_instance.test.project_id
	name = mongodbatlas_federated_database_instance.test.name
}
	`, federatedInstanceName, projectName, orgID)
}

func configFirstStepsUpdate(federatedInstanceName, projectName, orgID string) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[3]q
	}

resource "mongodbatlas_federated_database_instance" "test" {
   project_id         = mongodbatlas_project.test.id
   name = %[1]q

   data_process_region {
	   cloud_provider = "AWS"
	   region         = "OREGON_USA"
   }

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
