package searchindex_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	collectionName = "collection_test"
	searchAnalyzer = "lucene.standard"
	resourceName   = "mongodbatlas_search_index.test"
	datasourceName = "data.mongodbatlas_search_index.data_index"
)

func TestAccSearchIndexRS_basic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", ""),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterInfo.ClusterName),
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
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", "search"),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterInfo.ClusterName),
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
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigMapping(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", ""),
					resource.TestCheckResourceAttrSet(resourceName, "mappings_fields"),
					resource.TestCheckResourceAttrSet(resourceName, "analyzers"),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterInfo.ClusterName),
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
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigSynonyms(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "synonyms.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(resourceName, "synonyms.0.source_collection", collectionName),
					resource.TestCheckResourceAttr(resourceName, "type", ""),

					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterInfo.ClusterName),
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
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigBasic(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
				),
			},
			{
				Config:            testAccSearchIndexConfigBasic(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr, false),
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
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		indexName    = acctest.RandomWithPrefix("test-acc-index")
		databaseName = acctest.RandomWithPrefix("test-acc-db")
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchIndexConfigVector(clusterInfo.ProjectIDStr, indexName, databaseName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchIndexExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "vectorSearch"),
					resource.TestCheckResourceAttrSet(resourceName, "fields"),
					resource.TestCheckResourceAttrWith(resourceName, "fields", acc.JSONEquals(fields)),

					resource.TestCheckResourceAttr(datasourceName, "type", "vectorSearch"),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(datasourceName, "database", databaseName),
					resource.TestCheckResourceAttr(datasourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(datasourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(datasourceName, "index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "fields"),
					resource.TestCheckResourceAttrWith(datasourceName, "fields", acc.JSONEquals(fields)),
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
		ids := conversion.DecodeStateID(rs.Primary.ID)

		connV2 := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).AtlasV2
		_, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"]).Execute()
		if err != nil {
			return fmt.Errorf("index (%s) does not exist", ids["index_id"])
		}
		return nil
	}
}

func testAccSearchIndexConfigBasic(projectIDStr, indexName, databaseName, clusterNameStr, clusterTerraformStr string, explicitType bool) string {
	var indexType string
	if explicitType {
		indexType = `type="search"`
	}
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = "true"
			%[7]s
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]s
			project_id       = %[2]s
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectIDStr, indexName, databaseName, collectionName, searchAnalyzer, indexType)
}

func testAccSearchIndexConfigMapping(projectIDStr, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
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
			project_id       = %[2]s
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectIDStr, indexName, databaseName, collectionName, searchAnalyzer)
}

func testAccSearchIndexConfigSynonyms(projectIDStr, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
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
			project_id       = %[2]s
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectIDStr, indexName, databaseName, collectionName, searchAnalyzer)
}

func testAccSearchIndexConfigVector(projectIDStr, indexName, databaseName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
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
			project_id       = %[2]s
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterNameStr, projectIDStr, indexName, databaseName, collectionName)
}

func testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["cluster_name"], ids["index_id"]), nil
	}
}
