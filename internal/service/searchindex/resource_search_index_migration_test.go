package searchindex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigSearchIndex_basic(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
		config                 = configBasic(projectID, indexName, databaseName, clusterName, false)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "search_analyzer", searchAnalyzer),
					resource.TestCheckResourceAttr(resourceName, "type", ""),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigSearchIndex_withVector(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		indexName              = acc.RandomName()
		databaseName           = acc.RandomName()
		config                 = configVector(projectID, indexName, databaseName, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroySearchIndex,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", indexName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "collection_name", collectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "vectorSearch"),
					resource.TestCheckResourceAttrSet(resourceName, "fields"),
					resource.TestCheckResourceAttrWith(resourceName, "fields", acc.JSONEquals(fieldsJSON)),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
