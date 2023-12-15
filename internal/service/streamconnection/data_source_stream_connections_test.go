package streamconnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func TestAccStreamDSStreamConnections_basic(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-instance")
		dataSourceName = "data.mongodbatlas_stream_connections.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionsDataSourceConfig(kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false)),
				Check:  streamConnectionsAttributeChecks(dataSourceName, nil, nil, 1),
			},
		},
	})
}

func TestAccStreamDSStreamConnections_withPageConfig(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-instance")
		dataSourceName = "data.mongodbatlas_stream_connections.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionsWithPageAttrDataSourceConfig(kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false)),
				Check:  streamConnectionsAttributeChecks(dataSourceName, admin.PtrInt(2), admin.PtrInt(1), 0),
			},
		},
	})
}

func streamConnectionsDataSourceConfig(streamConnectionConfig string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_connections" "test" {
			project_id = mongodbatlas_stream_connection.test.project_id
			instance_name = mongodbatlas_stream_connection.test.instance_name
		}
	`, streamConnectionConfig)
}

func streamConnectionsWithPageAttrDataSourceConfig(streamConnectionConfig string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_connections" "test" {
			project_id = mongodbatlas_stream_connection.test.project_id
			instance_name = mongodbatlas_stream_connection.test.instance_name
			page_num = 2
			items_per_page = 1
		}
	`, streamConnectionConfig)
}

func streamConnectionsAttributeChecks(resourceName string, pageNum, itemsPerPage *int, totalCount int) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
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
