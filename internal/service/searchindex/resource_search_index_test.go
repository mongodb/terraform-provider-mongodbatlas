package searchindex_test

import (
	"context"
	"fmt"
	"regexp"
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
				Config: configBasic(projectID, clusterName, indexName, "search", databaseName, ""),
				Check:  checkBasic(projectID, clusterName, indexName, "search", databaseName, ""),
			},
		},
	})
}

func TestAccSearchIndex_withMapping(t *testing.T) {
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
				Config: configWithMapping(projectID, indexName, databaseName, clusterName),
				Check:  checkWithMapping(projectID, indexName, databaseName, clusterName),
			},
		},
	})
}

func TestAccSearchIndex_withSynonyms(t *testing.T) {
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
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, with),
				Check:  checkWithSynonyms(projectID, indexName, databaseName, clusterName, with),
			},
		},
	})
}

func TestAccSearchIndex_updatedToEmptySynonyms(t *testing.T) {
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
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, with),
				Check:  checkWithSynonyms(projectID, indexName, databaseName, clusterName, with),
			},
			{
				Config: configWithSynonyms(projectID, indexName, databaseName, clusterName, without),
				Check:  checkWithSynonyms(projectID, indexName, databaseName, clusterName, without),
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
				Check:  checkAdditionalAnalyzers(projectID, indexName, databaseName, clusterName, true),
			},
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, ""),
				Check:  checkAdditionalAnalyzers(projectID, indexName, databaseName, clusterName, false),
			},
			{
				Config:      configAdditional(projectID, indexName, databaseName, clusterName, incorrectFormatAnalyzersTF),
				ExpectError: regexp.MustCompile("cannot unmarshal search index attribute `analyzers` because it has an incorrect format"),
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
				Check:  checkAdditionalMappingsFields(projectID, indexName, databaseName, clusterName, true),
			},
			{
				Config: configAdditional(projectID, indexName, databaseName, clusterName, ""),
				Check:  checkAdditionalMappingsFields(projectID, indexName, databaseName, clusterName, false),
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
				Config: configBasic(projectID, clusterName, indexName, "", databaseName, ""),
				Check:  checkBasic(projectID, clusterName, indexName, "", databaseName, ""),
			},
			{
				Config:            configBasic(projectID, clusterName, indexName, "", databaseName, ""),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccSearchIndex_withStoredSourceFalse(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCase(t, "false"))
}

func TestAccSearchIndex_withStoredSourceTrue(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCase(t, "true"))
}

func TestAccSearchIndex_withStoredSourceInclude(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCase(t, storedSourceIncludeJSON))
}

func TestAccSearchIndex_withStoredSourceExclude(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCase(t, storedSourceExcludeJSON))
}

func TestAccSearchIndex_withStoredSourceUpdateEmptyType(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCaseUpdate(t, ""))
}

func TestAccSearchIndex_withStoredSourceUpdateSearchType(t *testing.T) {
	resource.ParallelTest(t, *storedSourceTestCaseUpdate(t, "search"))
}

func storedSourceTestCase(tb testing.TB, storedSource string) *resource.TestCase {
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
				Config: configBasic(projectID, clusterName, indexName, "search", databaseName, storedSource),
				Check:  checkBasic(projectID, clusterName, indexName, "search", databaseName, storedSource),
			},
		},
	}
}

func storedSourceTestCaseUpdate(tb testing.TB, searchType string) *resource.TestCase {
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
				Config: configBasic(projectID, clusterName, indexName, searchType, databaseName, "false"),
				Check:  checkBasic(projectID, clusterName, indexName, searchType, databaseName, "false"),
			},
			{
				Config: configBasic(projectID, clusterName, indexName, searchType, databaseName, "true"),
				Check:  checkBasic(projectID, clusterName, indexName, searchType, databaseName, "true"),
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

func checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attributes := map[string]string{
		"project_id":      projectID,
		"cluster_name":    clusterName,
		"name":            indexName,
		"type":            indexType,
		"database":        databaseName,
		"collection_name": collectionName,
	}
	if indexType != "vectorSearch" {
		attributes["mappings_dynamic"] = mappingsDynamic
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrChecks(datasourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "index_id")
	checks = acc.AddAttrSetChecks(datasourceName, checks, "index_id")
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
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

func configBasic(projectID, clusterName, indexName, indexType, databaseName, storedSource string) string {
	var extra string
	if indexType != "" {
		extra += fmt.Sprintf("type=%q\n", indexType)
	}
	if storedSource != "" {
		if storedSource == "true" || storedSource == "false" {
			extra += fmt.Sprintf("stored_source=%q\n", storedSource)
		} else {
			extra += fmt.Sprintf("stored_source= <<-EOF\n%s\nEOF\n", storedSource)
		}
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
			cluster_name     = mongodbatlas_search_index.test.cluster_name
			project_id       = mongodbatlas_search_index.test.project_id
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, extra)
}

func checkBasic(projectID, clusterName, indexName, indexType, databaseName, storedSource string) resource.TestCheckFunc {
	mappingsDynamic := "true"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "stored_source", storedSource),
		resource.TestCheckResourceAttr(datasourceName, "stored_source", storedSource),
	}
	if storedSource != "" && storedSource != "true" && storedSource != "false" {
		checks = []resource.TestCheckFunc{
			resource.TestCheckResourceAttrWith(resourceName, "stored_source", acc.JSONEquals(storedSource)),
			resource.TestCheckResourceAttrWith(datasourceName, "stored_source", acc.JSONEquals(storedSource)),
		}
	}
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic, checks...)
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
			cluster_name     = mongodbatlas_search_index.test.cluster_name
			project_id       = mongodbatlas_search_index.test.project_id
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, analyzersTF, mappingsFieldsTF)
}

func checkWithMapping(projectID, indexName, databaseName, clusterName string) resource.TestCheckFunc {
	indexType := ""
	mappingsDynamic := "false"
	attrNames := []string{"mappings_fields", "analyzers"}
	checks := acc.AddAttrSetChecks(resourceName, nil, attrNames...)
	checks = acc.AddAttrSetChecks(datasourceName, checks, attrNames...)
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic, checks...)
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
			cluster_name     = mongodbatlas_search_index.test.cluster_name
			project_id       = mongodbatlas_search_index.test.project_id
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, synonymsStr)
}

func checkWithSynonyms(projectID, indexName, databaseName, clusterName string, has bool) resource.TestCheckFunc {
	indexType := ""
	mappingsDynamic := "true"
	attrs := map[string]string{"synonyms.#": "0"}
	if has {
		attrs = map[string]string{
			"synonyms.#":                   "1",
			"synonyms.0.analyzer":          "lucene.simple",
			"synonyms.0.name":              "synonym_test",
			"synonyms.0.source_collection": collectionName,
		}
	}
	checks := acc.AddAttrChecks(resourceName, nil, attrs)
	checks = acc.AddAttrChecks(datasourceName, checks, attrs)
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic, checks...)
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

		data "mongodbatlas_search_index" "data_index" {
			cluster_name     = mongodbatlas_search_index.test.cluster_name
			project_id       = mongodbatlas_search_index.test.project_id
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, searchAnalyzer, additional)
}

func checkAdditionalAnalyzers(projectID, indexName, databaseName, clusterName string, has bool) resource.TestCheckFunc {
	indexType := ""
	mappingsDynamic := "true"
	check := resource.TestCheckResourceAttr(resourceName, "analyzers", "")
	if has {
		check = resource.TestCheckResourceAttrWith(resourceName, "analyzers", acc.JSONEquals(analyzersJSON))
	}
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic, check)
}

func checkAdditionalMappingsFields(projectID, indexName, databaseName, clusterName string, has bool) resource.TestCheckFunc {
	indexType := ""
	mappingsDynamic := "true"
	check := resource.TestCheckResourceAttr(resourceName, "mappings_fields", "")
	if has {
		check = resource.TestCheckResourceAttrWith(resourceName, "mappings_fields", acc.JSONEquals(mappingsFieldsJSON))
	}
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic, check)
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
			cluster_name     = mongodbatlas_search_index.test.cluster_name
			project_id       = mongodbatlas_search_index.test.project_id
			index_id 				 = mongodbatlas_search_index.test.index_id
		}
	`, clusterName, projectID, indexName, databaseName, collectionName, fieldsJSON)
}

func checkVector(projectID, indexName, databaseName, clusterName string) resource.TestCheckFunc {
	indexType := "vectorSearch"
	mappingsDynamic := "true"
	return checkAggr(projectID, clusterName, indexName, indexType, databaseName, mappingsDynamic,
		resource.TestCheckResourceAttrWith(resourceName, "fields", acc.JSONEquals(fieldsJSON)),
		resource.TestCheckResourceAttrWith(datasourceName, "fields", acc.JSONEquals(fieldsJSON)))
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

	analyzersTF                = "\nanalyzers = <<-EOF\n" + analyzersJSON + "\nEOF\n"
	incorrectFormatAnalyzersTF = "\nanalyzers = <<-EOF\n" + incorrectFormatAnalyzersJSON + "\nEOF\n"
	mappingsFieldsTF           = "\nmappings_fields = <<-EOF\n" + mappingsFieldsJSON + "\nEOF\n"

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

	incorrectFormatAnalyzersJSON = `
		[
			{
				"wrongField":[
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

	storedSourceIncludeJSON = `
		{ 
			"include": ["include1","include2"]
		}	
	`

	storedSourceExcludeJSON = `
		{
			"exclude": ["exclude1", "exclude2"]
		}	
	`
)
