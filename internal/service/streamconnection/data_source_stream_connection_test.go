package streamconnection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamDSStreamConnection_kafkaPlaintext(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_connection.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", kafkaNetworkingPublic, false)),
				Check:  kafkaStreamConnectionAttributeChecks(dataSourceName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, false),
			},
		},
	})
}

func TestAccStreamDSStreamConnection_kafkaSSL(t *testing.T) {
	var (
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
		dataSourceName = "data.mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingPublic, true)),
				Check:  kafkaStreamConnectionAttributeChecks(dataSourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true, false),
			},
		},
	})
}

func TestAccStreamDSStreamConnection_cluster(t *testing.T) {
	var (
		dataSourceName         = "data.mongodbatlas_stream_connection.test"
		projectID, clusterName = acc.ClusterNameExecution(t, false)
		instanceName           = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(clusterStreamConnectionConfig(projectID, instanceName, clusterName)),
				Check:  clusterStreamConnectionAttributeChecks(dataSourceName, clusterName),
			},
		},
	})
}

func TestAccStreamDSStreamConnection_sample(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_connection.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
		sampleName     = "sample_stream_solar"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(sampleStreamConnectionConfig(projectID, instanceName, sampleName)),
				Check:  sampleStreamConnectionAttributeChecks(dataSourceName, instanceName, sampleName),
			},
		},
	})
}

func streamConnectionDataSourceConfig(streamConnectionConfig string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_connection" "test" {
			project_id = mongodbatlas_stream_connection.test.project_id
			instance_name = mongodbatlas_stream_connection.test.instance_name
			connection_name = mongodbatlas_stream_connection.test.connection_name
		}
	`, streamConnectionConfig)
}
