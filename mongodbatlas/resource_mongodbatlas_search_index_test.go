package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSearchIndexRS_basic(t *testing.T) {
	var (
		resourceName                                     = "mongodbatlas_search_index.test"
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(t, projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index-name")
		datasourceIndexesName                            = "data.mongodbatlas_search_indexes.data_index"
		datasourceName                                   = "data.mongodbatlas_search_indexes.data_index"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfig(projectID, indexName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
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

func TestAccSearchIndexRS_withMapping(t *testing.T) {
	var (
		resourceName                                     = "mongodbatlas_search_index.test"
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(t, projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index-name")
		updatedAnalyzer                                  = "lucene.simple"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigAdvanced(projectID, indexName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "analyzer", updatedAnalyzer),
				),
			},
		},
	})
}

func TestAccSearchIndexRS_withSynonyms(t *testing.T) {
	var (
		resourceName                                     = "mongodbatlas_search_index.test"
		datasourceName                                   = "data.mongodbatlas_search_indexes.data_index"
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(t, projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index-name")
		updatedAnalyzer                                  = "lucene.standard"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigSynonyms(projectID, indexName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
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

func TestAccSearchIndexRS_importBasic(t *testing.T) {
	var (
		resourceName                                     = "mongodbatlas_search_index.test"
		projectID                                        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName, clusterNameStr, clusterTerraformStr = getClusterInfo(t, projectID)
		indexName                                        = acctest.RandomWithPrefix("test-acc-index-name")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckSearchIndex(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfig(projectID, indexName, clusterNameStr, clusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
			{
				Config:            testAccSearchIndexConfig(projectID, indexName, clusterNameStr, clusterTerraformStr),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccSearchIndexConfig(projectID, indexName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			database         = "database_test"
			collection_name  = "collection_test"
			mappings_dynamic = "true"
			search_analyzer  = "lucene.standard"
		}

		data "mongodbatlas_search_indexes" "data_index" {
			cluster_name 		= %[1]s
			project_id  		= mongodbatlas_search_index.test.project_id
			database   			= "database_test"
			collection_name = "collection_test"			
		}
	`, clusterNameStr, projectID, indexName)
}

func testAccSearchIndexConfigAdvanced(projectID, indexName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
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
	`, clusterNameStr, projectID, indexName)
}

func testAccSearchIndexConfigSynonyms(projectID, indexName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			name             = %[3]q
			analyzer         = "lucene.standard"
			collection_name  = "collection_test"
			database         = "database_test"
			mappings_dynamic = "true"
			search_analyzer  = "lucene.standard"
			synonyms {
				analyzer          = "lucene.simple"
				name              = "synonym_test"
				source_collection = "collection_test"
			}
		}

		data "mongodbatlas_search_indexes" "data_index" {
			cluster_name   	= %[1]s
			project_id    	= mongodbatlas_search_index.test.project_id
			database   			= "database_test"
			collection_name = "collection_test"
		}

		data "mongodbatlas_search_index" "test_two" {
			cluster_name	= %[1]s
			project_id		= mongodbatlas_search_index.test.project_id
			index_id 			= mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectID, indexName)
}

func testAccCheckMongoDBAtlasSearchIndexDestroy(state *terraform.State) error {
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

func getClusterInfo(t *testing.T, projectID string) (clusterName, clusterNameStr, clusterTerraformStr string) {
	// Allows faster test execution in local, don't use in CI
	clusterName = os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	if clusterName != "" {
		clusterNameStr = fmt.Sprintf("%q", clusterName)
		t.Logf("DONT DO THIS IN CI ONLY IN LOCAL, using exisiting cluster name: %s", clusterName)
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
