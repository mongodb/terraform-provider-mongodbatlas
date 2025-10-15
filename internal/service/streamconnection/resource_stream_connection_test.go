package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	dataSourceConfig = `
data "mongodbatlas_stream_connection" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		connection_name = mongodbatlas_stream_connection.test.connection_name
}
`

	dataSourcePluralConfig = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
}
`
	dataSourcePluralConfigWithPage = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		page_num = 2 # no specific reason for 2, just to test pagination
		items_per_page = 1
	}
	`
)

var (
	dataSourcesConfig = dataSourceConfig + dataSourcePluralConfig
	//go:embed testdata/dummy-ca.pem
	DummyCACert    string
	resourceName   = "mongodbatlas_stream_connection.test"
	dataSourceName = "data.mongodbatlas_stream_connection.test"
)

func TestAccStreamRSStreamConnection_cluster(t *testing.T) {
	testCase := testCaseCluster(t, "")
	resource.ParallelTest(t, *testCase)
}

func testCaseCluster(t *testing.T, nameSuffix string) *resource.TestCase {
	t.Helper()
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, false)
		_, instanceName        = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName         = "conn-cluster" + nameSuffix
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfig + configureCluster(projectID, instanceName, connectionName, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkClusterAttributes(resourceName, clusterName),
					checkClusterAttributes(dataSourceName, clusterName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func configureCluster(projectID, instanceName, connectionName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			instance_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Cluster"
		 	cluster_name = %[4]q
			db_role_to_execute = {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
	`, projectID, instanceName, connectionName, clusterName)
}

func checkClusterAttributes(resourceName, clusterName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Cluster"),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName+"hello"),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.role", "atlasAdmin"),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.type", "BUILT_IN"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func checkStreamConnectionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["instance_name"], rs.Primary.Attributes["project_id"], rs.Primary.Attributes["connection_name"]), nil
	}
}

func checkStreamConnectionExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_connection" {
				continue
			}
			projectID := rs.Primary.Attributes["project_id"]
			instanceName := rs.Primary.Attributes["instance_name"]
			connectionName := rs.Primary.Attributes["connection_name"]
			_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, instanceName, connectionName).Execute()
			if err != nil {
				return fmt.Errorf("stream connection (%s:%s:%s) does not exist", projectID, instanceName, connectionName)
			}
		}
		return nil
	}
}

func CheckDestroyStreamConnection(state *terraform.State) error {
	if instanceDestroyedErr := acc.CheckDestroyStreamInstance(state); instanceDestroyedErr != nil {
		return instanceDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_connection" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		instanceName := rs.Primary.Attributes["instance_name"]
		connectionName := rs.Primary.Attributes["connection_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, instanceName, connectionName).Execute()
		if err == nil {
			return fmt.Errorf("stream connection (%s:%s:%s) still exists", projectID, instanceName, connectionName)
		}
	}
	return nil
}
