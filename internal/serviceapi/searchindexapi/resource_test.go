package searchindexapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_search_index_api.test"
	database     = "sample_airbnb"
	collection   = "listingsAndReviews"
)

func TestAccSearchIndexAPI_basic(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, indexName),
				Check:  checkBasic(projectID, clusterName, indexName),
			},
			{
				Config:                               configBasic(projectID, clusterName, indexName),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"delete_on_create_timeout", "definition.%", "definition.mappings.%", "definition.mappings.dynamic"},
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccSearchIndexAPI_withMappingAndAnalyzer(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMappingAndAnalyzer(projectID, clusterName, indexName, true),
				Check:  checkWithMappingAndAnalyzer(projectID, clusterName, indexName, true),
			},
		},
	})
}

func TestAccSearchIndexAPI_withSynonymsUpdatedToEmpty(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			// configWithSynonyms requires a source collection to be setup in the database, creation reachs READY but does have an error.
			// Any follow up steps fail with an error unexpected state 'FAILED', wanted target 'READY, STEADY'.
			{
				Config: configWithSynonyms(projectID, clusterName, indexName, true),
				Check:  checkWithSynonyms(projectID, clusterName, indexName, true),
			},
			// {
			// 	Config: configWithSynonyms(projectID, clusterName, indexName, false),
			// 	Check:  checkWithSynonyms(projectID, clusterName, indexName, false),
			// },
		},
	})
}

func TestAccSearchIndexAPI_updatedToEmptyAnalyzers(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMappingAndAnalyzer(projectID, clusterName, indexName, true),
				Check:  checkWithMappingAndAnalyzer(projectID, clusterName, indexName, true),
			},
			// Currently fails due to Invalid definition: "typeSets" cannot be empty. CLOUDP to allow configuration for sending null in list (and other) types
			// {
			// 	Config: configWithMappingAndAnalyzer(projectID, clusterName, indexName, false),
			// 	Check:  checkWithMappingAndAnalyzer(projectID, clusterName, indexName, false),
			// },
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceFalse(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceBool(projectID, clusterName, indexName, false),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, false),
			},
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceTrue(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceBool(projectID, clusterName, indexName, true),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, true),
			},
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceInclude(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceJSON(projectID, clusterName, indexName, `{"include":["include1","include2"]}`),
				Check:  checkStoredSourceJSON(projectID, clusterName, indexName, `{"include":["include1","include2"]}`),
			},
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceExclude(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceJSON(projectID, clusterName, indexName, `{"exclude":["exclude1","exclude2"]}`),
				Check:  checkStoredSourceJSON(projectID, clusterName, indexName, `{"exclude":["exclude1","exclude2"]}`),
			},
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceUpdateEmptyType(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceBool(projectID, clusterName, indexName, false),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, false),
			},
			{
				Config: configWithStoredSourceBool(projectID, clusterName, indexName, true),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, true),
			},
		},
	})
}

func TestAccSearchIndexAPI_withStoredSourceUpdateSearchType(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithStoredSourceBoolAndType(projectID, clusterName, indexName, "search", false),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, false),
			},
			{
				Config: configWithStoredSourceBoolAndType(projectID, clusterName, indexName, "search", true),
				Check:  checkStoredSourceBool(projectID, clusterName, indexName, true),
			},
		},
	})
}

func TestAccSearchIndexAPI_withVector(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configVector(projectID, clusterName, indexName),
				Check:  checkVector(projectID, clusterName, indexName),
			},
		},
	})
}

func TestAccSearchIndexAPI_withTypeSets_ConfigurableDynamic(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, true)
		indexName              = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithTypeSets(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`, []string{`{"type":"string"}`}),
				Check:  checkTypeSets(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`, 1),
			},
			{
				Config: configWithTypeSets(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`, []string{`{"type":"string"}`, `{"type":"number"}`}),
				Check:  checkTypeSets(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`, 2),
			},
			{
				Config: configWithTypeSetsOmitted(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`),
				Check:  checkTypeSetsOmitted(projectID, clusterName, indexName, `{"typeSet":"ts_acc"}`),
			},
		},
	})
}

func configBasic(projectID, clusterName, indexName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q

			definition = {
				mappings = {
					dynamic = jsonencode(true)
				}
			}
		}
	`, projectID, clusterName, indexName, database, collection)
}

func configWithMappingAndAnalyzer(projectID, clusterName, indexName string, includeAnalyzers bool) string {
	var analyzers string
	if includeAnalyzers {
		analyzers = `
				analyzers = [{
					name = "index_analyzer_test_name"
					char_filters = [
						jsonencode({
							type     = "mapping"
							mappings = {"\\\\"="/"}
						})
					]
					tokenizer = {
						type   = jsonencode("nGram")
						minGram = 2
						maxGram = 5
					}
					token_filters = [
						jsonencode({
							type = "length"
							min  = 20
							max  = 33
						})
					]
				}]`
	}
	return fmt.Sprintf(`
        resource "mongodbatlas_search_index_api" "test" {
            group_id        = %[1]q
            cluster_name    = %[2]q
            name            = %[3]q
            database        = %[4]q
            collection_name = %[5]q

            definition = {
                mappings = {
                    dynamic = jsonencode(false)
                    fields = {
						address = jsonencode({
							type = "document"
							fields = {
								city = {
									type = "string"
									analyzer = "lucene.simple"
									ignoreAbove = 255
								}
								state = {
									type = "string"
									analyzer = "lucene.english"
								}
							}
						})
						company = jsonencode({
							type = "string"
							analyzer = "lucene.whitespace"
							multi = {
								mySecondaryAnalyzer = {
									type = "string"
									analyzer = "lucene.french"
								}
							}
						})
						employees = jsonencode({
							type = "string"
							analyzer = "lucene.standard"
						})
                	}
                }
				analyzer = "lucene.standard"
				%[6]s
            }
        }
    `, projectID, clusterName, indexName, database, collection, analyzers)
}

func checkBasic(projectID, clusterName, indexName string) resource.TestCheckFunc {
	attributes := map[string]string{
		"group_id":        projectID,
		"cluster_name":    clusterName,
		"name":            indexName,
		"database":        database,
		"collection_name": collection,
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "index_id")
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkWithMappingAndAnalyzer(projectID, clusterName, indexName string, expectAnalyzers bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkBasic(projectID, clusterName, indexName),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.mappings.dynamic", "false"),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.mappings.fields.%", "3"),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.analyzer", "lucene.standard"),
	}
	if expectAnalyzers {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "latest_definition.analyzers.#", "1"))
	} else {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "latest_definition.analyzers.#", "0"))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		indexID := rs.Primary.Attributes["index_id"]
		if groupID == "" || clusterName == "" || indexID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().AtlasSearchApi.GetClusterSearchIndex(context.Background(), groupID, clusterName, indexID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("search index(%s/%s/%s) does not exist", groupID, clusterName, indexID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_search_index_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		indexID := rs.Primary.Attributes["index_id"]
		if groupID == "" || clusterName == "" || indexID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().AtlasSearchApi.GetClusterSearchIndex(context.Background(), groupID, clusterName, indexID).Execute()
		if err == nil {
			return fmt.Errorf("search index (%s/%s/%s) still exists", groupID, clusterName, indexID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		indexID := rs.Primary.Attributes["index_id"]
		if groupID == "" || clusterName == "" || indexID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", groupID, clusterName, indexID), nil
	}
}

func configWithSynonyms(projectID, clusterName, indexName string, with bool) string {
	var synonyms string
	if with {
		synonyms = fmt.Sprintf(`
				synonyms = [{
					analyzer = "lucene.simple"
					name     = "synonym_test"
					source = {
						collection = %q
					}
				}]`, collection)
	} else {
		synonyms = "synonyms = []"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q
			definition = {
				search_analyzer = "lucene.standard"
				mappings = {
					dynamic = jsonencode(true)
					fields = {}
				}
				%[6]s
			}
		}
	`, projectID, clusterName, indexName, database, collection, synonyms)
}

func configWithStoredSourceBool(projectID, clusterName, indexName string, val bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q

			definition = {
				mappings = { dynamic = jsonencode(true) }
				stored_source = jsonencode(%[6]t)
			}
		}
	`, projectID, clusterName, indexName, database, collection, val)
}

func configWithStoredSourceBoolAndType(projectID, clusterName, indexName, indexType string, val bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q
			type            = %q

			definition = {
				mappings = { dynamic = jsonencode(true) }
				stored_source = jsonencode(%[7]t)
			}
		}
	`, projectID, clusterName, indexName, database, collection, indexType, val)
}

func configWithStoredSourceJSON(projectID, clusterName, indexName, json string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q

			definition = {
				mappings = { dynamic = jsonencode(true) }
				stored_source = jsonencode(%[6]s)
			}
		}
	`, projectID, clusterName, indexName, database, collection, json)
}

func configVector(projectID, clusterName, indexName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q
			type            = "vectorSearch"

			definition = {
				fields = [
					jsonencode({
						type          = "vector"
						path          = "plot_embedding"
						numDimensions = 1536
						similarity    = "euclidean"
					})
				]
			}
		}
	`, projectID, clusterName, indexName, database, collection)
}

func configWithTypeSets(projectID, clusterName, indexName, dynamicJSON string, types []string) string {
	var typesStr string
	for i, t := range types {
		if i > 0 {
			typesStr += ","
		}
		typesStr += fmt.Sprintf("jsonencode(%s)", t)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q
			type            = "search"

			definition = {
				mappings = {
					dynamic = jsonencode(%[6]s)
				}
				type_sets = [{
					name  = "ts_acc"
					types = [%[7]s]
				}]
			}
		}
	`, projectID, clusterName, indexName, database, collection, dynamicJSON, typesStr)
}

func configWithTypeSetsOmitted(projectID, clusterName, indexName, dynamicJSON string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index_api" "test" {
			group_id        = %[1]q
			cluster_name    = %[2]q
			name            = %[3]q
			database        = %[4]q
			collection_name = %[5]q
			type            = "search"

			definition = {
				mappings = {
					dynamic = jsonencode(%[6]s)
				}
			}
		}
	`, projectID, clusterName, indexName, database, collection, dynamicJSON)
}

// ---------------------------
// Check helpers
// ---------------------------

func checkAttrs(projectID, clusterName, indexName string) resource.TestCheckFunc {
	attributes := map[string]string{
		"group_id":        projectID,
		"cluster_name":    clusterName,
		"name":            indexName,
		"database":        database,
		"collection_name": collection,
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "index_id")
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkWithSynonyms(projectID, clusterName, indexName string, has bool) resource.TestCheckFunc {
	count := "0"
	extra := []resource.TestCheckFunc{}
	if has {
		count = "1"
		extra = []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "latest_definition.synonyms.0.analyzer", "lucene.simple"),
			resource.TestCheckResourceAttr(resourceName, "latest_definition.synonyms.0.name", "synonym_test"),
			resource.TestCheckResourceAttr(resourceName, "latest_definition.synonyms.0.source.collection", collection),
		}
	}
	checks := []resource.TestCheckFunc{
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.synonyms.#", count),
	}
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkStoredSourceBool(projectID, clusterName, indexName string, val bool) resource.TestCheckFunc {
	boolJSON := "false"
	if val {
		boolJSON = "true"
	}
	return resource.ComposeAggregateTestCheckFunc(
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttrWith(resourceName, "latest_definition.stored_source", acc.JSONEquals(boolJSON)),
	)
}

func checkStoredSourceJSON(projectID, clusterName, indexName, json string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttrWith(resourceName, "latest_definition.stored_source", acc.JSONEquals(json)),
	)
}

func checkVector(projectID, clusterName, indexName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttr(resourceName, "type", "vectorSearch"),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.fields.#", "1"),
	)
}

func checkTypeSets(projectID, clusterName, indexName, dynamicJSON string, typeCount int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttrWith(resourceName, "latest_definition.mappings.dynamic", acc.JSONEquals(dynamicJSON)),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.type_sets.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.type_sets.0.name", "ts_acc"),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.type_sets.0.types.#", fmt.Sprintf("%d", typeCount)),
	)
}

func checkTypeSetsOmitted(projectID, clusterName, indexName, dynamicJSON string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkAttrs(projectID, clusterName, indexName),
		resource.TestCheckResourceAttrWith(resourceName, "latest_definition.mappings.dynamic", acc.JSONEquals(dynamicJSON)),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.type_sets.#", "0"),
	)
}
