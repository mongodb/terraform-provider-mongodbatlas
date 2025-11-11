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

func TestAccSearchIndexAPI_withMappingsFields(t *testing.T) {
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
				Config: configWithMappingsFields(projectID, clusterName, indexName, true),
				Check:  checkWithMappingsFields(projectID, clusterName, indexName, true),
			},
			// TODO: revise update behavior as part of CLOUDP-352324
			// {
			// 	Config: configWithMappingsFields(projectID, clusterName, indexName, false),
			// 	Check:  checkWithMappingsFields(projectID, clusterName, indexName, false),
			// },
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

func configWithMappingsFields(projectID, clusterName, indexName string, with bool) string {
	var fields string
	if with {
		fields = `
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
                }`
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
                    dynamic = jsonencode(true)
                    %[6]s
                }
            }
        }
    `, projectID, clusterName, indexName, database, collection, fields)
}

func checkBasic(projectID, clusterName, indexName string) resource.TestCheckFunc {
	attributes := map[string]string{
		"group_id":                           projectID,
		"cluster_name":                       clusterName,
		"name":                               indexName,
		"database":                           database,
		"collection_name":                    collection,
		"latest_definition.mappings.dynamic": "true",
	}
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "index_id")
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkWithMappingsFields(projectID, clusterName, indexName string, has bool) resource.TestCheckFunc {
	count := "0"
	if has {
		count = "3"
	}
	return resource.ComposeAggregateTestCheckFunc(
		checkBasic(projectID, clusterName, indexName),
		resource.TestCheckResourceAttr(resourceName, "latest_definition.mappings.fields.%", count),
	)
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
