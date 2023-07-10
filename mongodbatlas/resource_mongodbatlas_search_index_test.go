package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccClusterRSSearchIndex_basic(t *testing.T) {
	var (
		index                 matlas.SearchIndex
		resourceName          = "mongodbatlas_search_index.test"
		clusterName           = acctest.RandomWithPrefix("test-acc-index")
		projectID             = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name                  = "name_test"
		datasourceIndexesName = "data.mongodbatlas_search_indexes.data_index"
		datasourceName        = "data.mongodbatlas_search_indexes.data_index"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
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
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "collection_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "database"),
					resource.TestCheckResourceAttrSet(datasourceName, "search_analyzer"),
					resource.TestCheckResourceAttrSet(datasourceName, "synonyms.#"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.source_collection", "collection_test"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "cluster_name"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "database"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "results.0.index_id"),
					resource.TestCheckResourceAttrSet(datasourceIndexesName, "results.0.name"),
				),
			},
		},
	})
}

func TestAccClusterRSSearchIndex_withMapping(t *testing.T) {
	var (
		index           matlas.SearchIndex
		resourceName    = "mongodbatlas_search_index.test"
		clusterName     = acctest.RandomWithPrefix("test-acc-index")
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name            = "name_test"
		updatedAnalyzer = "lucene.simple"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
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
			},
		},
	})
}

func TestAccClusterRSSearchIndex_withSynonyms(t *testing.T) {
	var (
		index           matlas.SearchIndex
		resourceName    = "mongodbatlas_search_index.test"
		datasourceName  = "data.mongodbatlas_search_indexes.data_index"
		clusterName     = acctest.RandomWithPrefix("test-acc-index")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acctest.RandomWithPrefix("test-acc")
		name            = "name_test"
		updatedAnalyzer = "lucene.standard"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfigSynonyms(orgID, projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "analyzer", updatedAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "synonyms.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.source_collection", "collection_test"),
					resource.TestCheckResourceAttrSet(datasourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "database"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.analyzer"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.source_collection", "collection_test"),
				),
			},
		},
	})
}

func TestAccClusterRSSearchIndex_importBasic(t *testing.T) {
	var (
		index        matlas.SearchIndex
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = "name_test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
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
		
			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false
		
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
		}

		resource "mongodbatlas_search_index" "test" {
			project_id       = mongodbatlas_cluster.aws_conf.project_id
			cluster_name     = mongodbatlas_cluster.aws_conf.name
			collection_name  = "collection_test"
			database         = "database_test"
			mappings_dynamic = "true"
			name             = "name_test"
			search_analyzer  = "lucene.standard"
		}

		data "mongodbatlas_search_indexes" "data_index" {
			cluster_name           = mongodbatlas_search_index.test.cluster_name
			project_id         = mongodbatlas_search_index.test.project_id
			database   = "database_test"
			collection_name = "collection_test"
			page_num = 1
			items_per_page = 100
			
		}
	`, projectID, clusterName)
}

func testAccMongoDBAtlasSearchIndexConfigAdvanced(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "aws_conf" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10

			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}

			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"

		}

		resource "mongodbatlas_search_index" "test" {
			project_id   = mongodbatlas_cluster.aws_conf.project_id
			cluster_name = mongodbatlas_cluster.aws_conf.name

			analyzer         = "lucene.simple"
			collection_name  = "collection_test"
			database         = "database_test"
			mappings_dynamic = false
			mappings_fields  = <<-EOF
			{
				"address":{
					"type":"document",
					"fields":{
						"city":{
								"type":"string",
								"analyzer":"lucene.simple",
								"ignoreAbove":255
						},
						"state":{
								"type":"string",
								"analyzer":"lucene.english"
						}
					}
				},
				"company":{
					"type":"string",
					"analyzer":"lucene.whitespace",
					"multi":{
						"mySecondaryAnalyzer":{
							"type":"string",
							"analyzer":"lucene.french"
						}
					}
				},
				"employees":{
					"type":"string",
					"analyzer":"lucene.standard"
				}
			}
			EOF
			name             = "name_test"
			search_analyzer  = "lucene.standard"
			analyzers        = <<-EOF
			[
				{
					"name":"index_analyzer_test_name",
					"charFilters":[
						 {
								"type":"mapping",
								"mappings":{
									 "\\":"/"
								}
						 }
					],
					"tokenizer":[
						 {
								"type":"nGram",
								"minGram":2,
								"maxGram":5s
						 }
					],
					"tokenFilters":[
						 {
								"type":"length",
								"min":20,
								"max":33
						 }
					]
				}
			]
			EOF
		}
	`, projectID, clusterName)
}

func testAccMongoDBAtlasSearchIndexConfigSynonyms(orgID, projectName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "test_cluster" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q
			disk_size_gb = 10
		
			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
		
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false
		
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
		
		}
		
		resource "mongodbatlas_search_index" "test" {
			project_id       = mongodbatlas_cluster.test_cluster.project_id
			cluster_name     = mongodbatlas_cluster.test_cluster.name
			analyzer         = "lucene.standard"
			collection_name  = "collection_test"
			database         = "database_test"
			mappings_dynamic = "true"
			name             = "name_test"
			search_analyzer  = "lucene.standard"
			synonyms {
				analyzer          = "lucene.simple"
				name              = "synonym_test"
				source_collection = "collection_test"
			}
		}

		data "mongodbatlas_search_indexes" "data_index" {
			cluster_name           = mongodbatlas_search_index.test.cluster_name
			project_id         = mongodbatlas_search_index.test.project_id
			database   = "database_test"
			collection_name = "collection_test"
			page_num = 1
			items_per_page = 100
		}

		data "mongodbatlas_search_index" "test_two" {
			cluster_name        = mongodbatlas_search_index.test.cluster_name
			project_id          = mongodbatlas_search_index.test.project_id
			index_id 			= mongodbatlas_search_index.test.index_id
		}
	`, orgID, projectName, clusterName)
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
