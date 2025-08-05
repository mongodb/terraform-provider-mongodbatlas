package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func TestAccStreamDSStreamworkspaces_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_workspaces.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
	)

	checks := paginatedAttrChecks(dataSourceName, nil, nil)
	// created workspace is present in results
	checks = append(checks, resource.TestCheckResourceAttrWith(dataSourceName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "results.*", map[string]string{
			"workspace_name": workspaceName,
		}))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamworkspacesDataSourceConfig(projectID, workspaceName, region, cloudProvider),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccStreamDSStreamworkspaces_withPageConfig(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_workspaces.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
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
				Config: streamworkspacesWithPageAttrDataSourceConfig(projectID, workspaceName, region, cloudProvider, pageNumber),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func streamworkspacesDataSourceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspaces" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
		}
	`, acc.StreamInstanceConfig(projectID, workspaceName, region, cloudProvider))
}

func streamworkspacesWithPageAttrDataSourceConfig(projectID, workspaceName, region, cloudProvider string, pageNum int) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspaces" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
			page_num = %d
			items_per_page = 1
		}
	`, acc.StreamInstanceConfig(projectID, workspaceName, region, cloudProvider), pageNum)
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
