package searchindex_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccSearchIndex_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccSearchIndex_withSearchType(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, indexName, databaseName, clusterName, true),
				Check:  checkBasic(projectID, indexName, databaseName, clusterName, true),
			},
		},
	})
}

func TestAccSearchIndex_withMapping(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
		indexType              = ""
		mappingsDynamic        = "false"
	)
	checks := commonChecks(indexName, indexType, mappingsDynamic, databaseName, clusterName, false)
	checks = addAttrSetChecks(checks, "mappings_fields", "analyzers")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configWithMapping(projectID, indexName, databaseName, clusterName),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccSearchIndex_withSynonyms(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
		indexType              = ""
		mappingsDynamic        = "true"
		mapChecks              = map[string]string{
			"synonyms.#":                   "1",
			"synonyms.0.analyzer":          "lucene.simple",
			"synonyms.0.name":              "synonym_test",
			"synonyms.0.source_collection": collectionName,
		}
	)
	checks := commonChecks(indexName, indexType, mappingsDynamic, databaseName, clusterName, false)
	checks = addAttrChecks(checks, mapChecks)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, with),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccSearchIndex_updatedToEmptySynonyms(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
		indexType              = ""
		mappingsDynamic        = "true"
	)
	checks := commonChecks(indexName, indexType, mappingsDynamic, databaseName, clusterName, false)
	checks1 := addAttrChecks(checks, map[string]string{
		"synonyms.#":                   "1",
		"synonyms.0.analyzer":          "lucene.simple",
		"synonyms.0.name":              "synonym_test",
		"synonyms.0.source_collection": collectionName,
	})
	checks2 := addAttrChecks(checks, map[string]string{"synonyms.#": "0"})
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, with),
				Check:  resource.ComposeAggregateTestCheckFunc(checks1...),
			},
			{
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, without),
				Check:  resource.ComposeAggregateTestCheckFunc(checks2...),
			},
		},
	})
}

func TestAccSearchIndex_updatedToEmptyAnalyzers(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, analyzersTF),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrWith(resourceName, "analyzers", acc.JSONEquals(analyzersJSON)),
				),
			},
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "analyzers", ""),
				),
			},
		},
	})
}

func TestAccSearchIndex_updatedToEmptyMappingsFields(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, mappingsFieldsTF),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrWith(resourceName, "mappings_fields", acc.JSONEquals(mappingsFieldsJSON)),
				),
			},
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "mappings_fields", ""),
				),
			},
		},
	})
}

func TestAccSearchIndex_withVector(t *testing.T) {
	resource.ParallelTest(t, *basicVectorTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID, clusterName = acc.ClusterNameExecution(tb)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, indexName, databaseName, clusterName, false),
				Check:  checkBasic(projectID, indexName, databaseName, clusterName, false),
			},
			{
				Config:            configBasic(projectID, indexName, databaseName, clusterName, false),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func basicVectorTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID, clusterName = acc.ClusterNameExecution(tb)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config: configVector(projectID, indexName, databaseName, clusterName),
				Check:  checkVector(projectID, indexName, databaseName, clusterName),
			},
		},
	}
}

func commonChecks(indexName, indexType, mappingsDynamic, databaseName, clusterName string, explicitType bool) []resource.TestCheckFunc {
	attributes := map[string]string{
		"name":             indexName,
		"cluster_name":     clusterName,
		"database":         databaseName,
		"collection_name":  collectionName,
		"mappings_dynamic": mappingsDynamic,
	}
	indexTypeEffective := ""
	if explicitType {
		indexTypeEffective = indexType
	}
	checks := addAttrChecks(nil, attributes)
	checks = append(checks,
		resource.TestCheckResourceAttr(resourceName, "type", indexTypeEffective),
		resource.TestCheckResourceAttr(datasourceName, "type", indexTypeEffective))
	checks = acc.AddAttrSetChecks(resourceName, checks, "project_id", "index_id")
	return acc.AddAttrSetChecks(datasourceName, checks, "project_id", "index_id")
}

func addAttrChecks(checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	checks = acc.AddAttrChecks(resourceName, checks, mapChecks)
	return acc.AddAttrChecks(datasourceName, checks, mapChecks)
}

func addAttrSetChecks(checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	checks = acc.AddAttrSetChecks(resourceName, checks, attrNames...)
	return acc.AddAttrSetChecks(datasourceName, checks, attrNames...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().AtlasSearchApi.GetAtlasSearchIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"]).Execute()
		if err != nil {
			return fmt.Errorf("index (%s) does not exist", ids["index_id"])
		}
		return nil
	}
}

func configBasic(projectID, indexName, databaseName, clusterName string, explicitType bool) string {
	var indexType string
	if explicitType {
		indexType = `type="search"`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = "true"
			%[7]s
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, indexType)
}

func checkBasic(projectID, indexName, databaseName, clusterName string, explicitType bool) resource.TestCheckFunc {
	indexType := "search"
	mappingsDynamic := "true"
	checks := commonChecks(indexName, indexType, mappingsDynamic, databaseName, clusterName, explicitType)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configWithMapping(projectID, indexName, databaseName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = false
			%[7]s
			%[8]s
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, analyzersTF, mappingsFieldsTF)
}

func configWithSynonyms(projectID, indexName, databaseName, clusterName string, has bool) string {
	var synonymsStr string
	if has {
		synonymsStr = fmt.Sprintf(`
			synonyms {
				analyzer          = "lucene.simple"
				name              = "synonym_test"
				source_collection = %q
			}
		`, collectionName)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = true
			%[7]s
		}

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, synonymsStr)
}

func configAdditional(projectID, indexName, databaseName, clusterName, additional string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
			search_analyzer  = %[6]q
			mappings_dynamic = true
			%[7]s
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, additional)
}

func configVector(projectID, indexName, databaseName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			name             = %[3]q
			database         = %[4]q
			collection_name  = %[5]q
		
			type = "vectorSearch"
			
			fields = <<-EOF
	    %[6]s
			EOF
		}
	
		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = %[1]q
			project_id       = %[2]q
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, fieldsJSON)
}

func checkVector(projectID, indexName, databaseName, clusterName string) resource.TestCheckFunc {
	indexType := "vectorSearch"
	attributes := map[string]string{
		"project_id":      projectID,
		"name":            indexName,
		"cluster_name":    clusterName,
		"database":        databaseName,
		"collection_name": collectionName,
		"type":            indexType,
	}
	checks := addAttrChecks(nil, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "index_id")
	checks = acc.AddAttrSetChecks(datasourceName, checks, "index_id")
	checks = append(checks,
		resource.TestCheckResourceAttrWith(resourceName, "fields", acc.JSONEquals(fieldsJSON)),
		resource.TestCheckResourceAttrWith(datasourceName, "fields", acc.JSONEquals(fieldsJSON)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["cluster_name"], ids["index_id"]), nil
	}
}

const (
	collectionName = "collection_test"
	searchAnalyzer = "lucene.standard"
	resourceName   = "mongodbatlas_search_index.test"
	datasourceName = "data.mongodbatlas_search_index.data_index"
	with           = true
	without        = false

	analyzersTF      = "\nanalyzers = <<-EOF\n" + analyzersJSON + "\nEOF\n"
	mappingsFieldsTF = "\nmappings_fields = <<-EOF\n" + mappingsFieldsJSON + "\nEOF\n"

	analyzersJSON = `
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
`

	mappingsFieldsJSON = `
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
	`

	fieldsJSON = `
		[{
			"type": "vector",
			"path": "plot_embedding",
			"numDimensions": 1536,
			"similarity": "euclidean"
		}]	
	`
)
