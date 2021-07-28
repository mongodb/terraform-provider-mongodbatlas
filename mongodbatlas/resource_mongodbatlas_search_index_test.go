package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasSearchIndex_basic(t *testing.T) {
	var (
		index        matlas.SearchIndex
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = "name_test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasSearchIndex_withMapping(t *testing.T) {
	var (
		index           matlas.SearchIndex
		resourceName    = "mongodbatlas_search_index.test"
		clusterName     = acctest.RandomWithPrefix("test-acc-index")
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name            = "name_test"
		updatedAnalyzer = "lucene.simple"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfigAdvanced(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "analyzer", updatedAnalyzer),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasSearchIndex_importBasic(t *testing.T) {
	var (
		index        matlas.SearchIndex
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = "name_test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
			{
				Config:            testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasSearchIndexExists(resourceName string, index *matlas.SearchIndex) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		indexResponse, _, err := conn.Search.GetIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"])
		if err == nil {
			*index = *indexResponse
			return nil
		}

		return fmt.Errorf("index (%s) does not exist", ids["index_id"])
	}
}

func testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "aws_conf" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10
			replication_factor           = 3
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "US_EAST_2"

		}

		resource "mongodbatlas_search_index" "test" {
			project_id         = mongodbatlas_cluster.aws_conf.project_id
			cluster_name       = mongodbatlas_cluster.aws_conf.name
			analyzer = "lucene.standard"
			collection_name = "collection_test"
			database = "database_test"
			mappings_dynamic = "true"
			name = "name_test"
			search_analyzer = "lucene.standard"
		}

	
	`, projectID, clusterName)
}

func testAccMongoDBAtlasSearchIndexConfigAdvanced(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "aws_conf" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10
			replication_factor           = 3
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "US_EAST_2"

		}

		resource "mongodbatlas_search_index" "test" {
			project_id         = mongodbatlas_cluster.aws_conf.project_id
			cluster_name       = mongodbatlas_cluster.aws_conf.name

			analyzer = "lucene.simple"
			collection_name = "collection_test"
			database = "database_test"
			mappings_dynamic = false
			mappings_fields = <<-EOF
							 {
				  "address": {
					"type": "document",
					"fields": {
					  "city": {
						"type": "string",
						"analyzer": "lucene.simple",
						"ignoreAbove": 255
					  },
					  "state": {
						"type": "string",
						"analyzer": "lucene.english"
					  }
					}
				  },
				  "company": {
					"type": "string",
					"analyzer": "lucene.whitespace",
					"multi": {
					  "mySecondaryAnalyzer": {
						"type": "string",
						"analyzer": "lucene.french"
					  }
					}
				  },
				  "employees": {
					"type": "string",
					"analyzer": "lucene.standard"
				  }
				}
   			EOF
			name = "name_test"
			search_analyzer = "lucene.standard"
			analyzers {
				name = "index_analyzer_test_name"
				char_filters {
					type = "mapping"
					mappings = <<-EOF
					{"\\" : "/"}
					EOF
				}
				tokenizer {
					type = "nGram"
					min_gram = 2
					max_gram = 5
				}

			token_filters {
				type = "length"
				min = 20
				max = 33
			}
		}
	}
	
	`, projectID, clusterName)
}

func testAccCheckMongoDBAtlasSearchIndexDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_search_index" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		searchIndex, _, err := conn.Search.GetIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"])
		if err == nil && searchIndex != nil {
			return fmt.Errorf("index id (%s) still exists", ids["index_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["cluster_name"], ids["index_id"]), nil
	}
}
