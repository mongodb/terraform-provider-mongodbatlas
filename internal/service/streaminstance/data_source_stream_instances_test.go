package streaminstance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

func TestAccStreamDSStreamInstances_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_instances.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
	)

	checks := paginatedAttrChecks(dataSourceName, nil, nil)
	// created instance is present in results
	checks = append(checks, resource.TestCheckResourceAttrWith(dataSourceName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "results.*", map[string]string{
			"instance_name": instanceName,
		}))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstancesDataSourceConfig(projectID, instanceName, region, cloudProvider),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccStreamDSStreamInstances_withPageConfig(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_instances.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
		pageNumber     = 1000 // high page number so no results are returned
	)

	checks := paginatedAttrChecks(dataSourceName, admin.PtrInt(pageNumber), admin.PtrInt(1))
	checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, "results.#", "0")) // expecting no results

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstancesWithPageAttrDataSourceConfig(projectID, instanceName, region, cloudProvider, pageNumber),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func streamInstancesDataSourceConfig(projectID, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
		}
	`, acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider))
}

func streamInstancesWithPageAttrDataSourceConfig(projectID, instanceName, region, cloudProvider string, pageNum int) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			page_num = %d
			items_per_page = 1
		}
	`, acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider), pageNum)
}

func paginatedAttrChecks(resourceName string, pageNum, itemsPerPage *int) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "total_count"),
	}
	if pageNum != nil {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "page_num", fmt.Sprint(*pageNum)))
	}
	if itemsPerPage != nil {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "items_per_page", fmt.Sprint(*itemsPerPage)))
	}
	return checks
}
