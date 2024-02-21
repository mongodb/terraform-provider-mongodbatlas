package streaminstance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func TestAccStreamDSStreamInstances_basic(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-name")
		dataSourceName = "data.mongodbatlas_stream_instances.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstancesDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider),
				Check:  streamInstancesAttributeChecks(dataSourceName, nil, nil, 1),
			},
		},
	})
}

func TestAccStreamDSStreamInstances_withPageConfig(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-name")
		dataSourceName = "data.mongodbatlas_stream_instances.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstancesWithPageAttrDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider),
				Check:  streamInstancesAttributeChecks(dataSourceName, admin.PtrInt(2), admin.PtrInt(1), 0),
			},
		},
	})
}

func streamInstancesDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
		}
	`, acc.StreamInstanceConfig(orgID, projectName, instanceName, region, cloudProvider))
}

func streamInstancesWithPageAttrDataSourceConfig(orgID, projectName, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			page_num = 2
			items_per_page = 1
		}
	`, acc.StreamInstanceConfig(orgID, projectName, instanceName, region, cloudProvider))
}

func streamInstancesAttributeChecks(resourceName string, pageNum, itemsPerPage *int, totalCount int) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "total_count"),
		resource.TestCheckResourceAttr(resourceName, "results.#", fmt.Sprint(totalCount)),
	}
	if pageNum != nil {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "page_num", fmt.Sprint(*pageNum)))
	}
	if itemsPerPage != nil {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "items_per_page", fmt.Sprint(*itemsPerPage)))
	}
	return resource.ComposeTestCheckFunc(resourceChecks...)
}
