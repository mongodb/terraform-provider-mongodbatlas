package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/testutils"
)

const (
	collectionName = "collection_test"
	searchAnalyzer = "lucene.standard"
	resourceName   = "mongodbatlas_search_index.test"
	datasourceName = "data.mongodbatlas_search_index.data_index"
)

func TestAccSearchIndexRS_basic(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", ""),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "mappings_dynamic", "true"),
					resource.TestCheckResourceAttr(datasourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "type", ""),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
				),
			},
		},
	})
}

func TestAccSearchIndexRS_withSearchType(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", "search"),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "mappings_dynamic", "true"),
					resource.TestCheckResourceAttr(datasourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "type", "search"),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
				),
			},
		},
	})
}

func TestAccSearchIndexRS_withMapping(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigMapping(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", ""),
					resource.TestCheckResourceAttrSet(resourceName, "mappings_fields"),
					resource.TestCheckResourceAttrSet(resourceName, "analyzers"),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(datasourceName, "mappings_dynamic", "false"),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "type", ""),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "mappings_fields"),
					resource.TestCheckResourceAttrSet(datasourceName, "analyzers"),
				),
			},
		},
	})
}

func TestAccSearchIndexRS_withSynonyms(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigSynonyms(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "synonyms.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.source_collection", collectionName),
					resource.TestCheckResourceAttr(resourceName, "type", ""),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(datasourceName, "mappings_dynamic", "true"),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "type", ""),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(datasourceName, "synonyms.0.source_collection", collectionName),
				),
			},
		},
	})
}

func TestAccSearchIndexRS_importBasic(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
			{
				Config:            testAccSearchIndexConfigBasic(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr, false),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSearchIndexRS_withVector(t *testing.T) {
	var (
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index")
		databaseName                                     = acctest.RandomWithPrefix("test-acc-db")
	)
	fields := []map[string]interface{}{
		{
			"type":          "vector",
			"path":          "plot_embedding",
			"numDimensions": float64(1536),
			"similarity":    "euclidean",
		},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigVector(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "vectorSearch"),
					resource.TestCheckResourceAttrSet(resourceName, "fields"),
					resource.TestCheckResourceAttrWith(resourceName, "fields", testutils.JSONEquals(fields)),

					resource.TestCheckResourceAttr(datasourceName, "type", "vectorSearch"),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "fields"),
					resource.TestCheckResourceAttrWith(datasourceName, "fields", testutils.JSONEquals(fields)),
				),
			},
		},
	})
}

func testAccCheckSearchIndexExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := decodeStateID(rs.Primary.ID)

		connV2 := testAccProviderSdkV2.Meta().(*MongoDBClient).AtlasV2
		_, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"]).Execute()
		if err != nil {
			return fmt.Errorf("index (%s) does not exist", ids["index_id"])
		}
		return nil
	}
}

func testAccSearchIndexConfigBasic(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr string, explicitType bool) string {
	var indexType string
	if explicitType {
		indexType = `type="search"`
	}
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = "true"
			%[7]s
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectID, indexName, databaseName, collectionName, searchAnalyzer, indexType)
}

func testAccSearchIndexConfigMapping(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
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
					"analyzer":%[6]q
				}
			}
			EOF
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
					"tokenizer": {
								"type":"nGram",
								"minGram": 2,
								"maxGram": 5
					},
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

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectID, indexName, databaseName, collectionName, searchAnalyzer)
}

func testAccSearchIndexConfigSynonyms(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = true
			synonyms {
				analyzer          = "lucene.simple"
				name              = "synonym_test"
				source_collection = %[5]q
			}
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectID, indexName, databaseName, collectionName, searchAnalyzer)
}

func testAccSearchIndexConfigVector(projectID, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
		
			type = "vectorSearch"
			
			fields = <<-EOF
				[{
					"type": "vector",
					"path": "plot_embedding",
					"numDimensions": 1536,
					"similarity": "euclidean"
				}]
				EOF
		}
	
		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectID, indexName, databaseName, collectionName)
}

func testAccCheckMongoDBAtlasSearchIndexDestroy(state *terraform.State) error {
	if os.Getenv("MONGODB_ATLAS_CLUSTER_NAME") != "" {
		return nil
	}
	conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas

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

func getClusterInfo(projectID string) (clusterName, clusterNameStr, clusterTerraformStr string) {
	// Allows faster test execution in local, don't use in CI
	clusterName = os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	if clusterName != "" {
		clusterNameStr = fmt.Sprintf("%q", clusterName)
	} else {
		clusterName = acctest.RandomWithPrefix("test-acc-index")
		clusterNameStr = "mongodbatlas_cluster.test_cluster.name"
		clusterTerraformStr = fmt.Sprintf(`
			resource "mongodbatlas_cluster" "test_cluster" {
				project_id   									= %[1]q
				name         									= %[2]q
				disk_size_gb 									= 10
				backup_enabled               	= false
				auto_scaling_disk_gb_enabled	= false
				provider_name               	= "AWS"
				provider_instance_size_name 	= "M10"
			
				cluster_type = "REPLICASET"
				replication_specs {
					num_shards = 1
					regions_config {
						region_name     = "US_WEST_2"
						electable_nodes = 3
						priority        = 7
						read_only_nodes = 0
					}
				}
			}
		`, projectID, clusterName)
	}
	return clusterName, clusterNameStr, clusterTerraformStr
}
