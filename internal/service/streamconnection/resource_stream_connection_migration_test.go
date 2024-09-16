package streamconnection_test

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
		config       = kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false)
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false, true),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigStreamRSStreamConnection_kafkaSSL(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
		config       = kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true)
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigStreamRSStreamConnection_cluster(t *testing.T) {
	var (
		resourceName           = "mongodbatlas_stream_connection.test"
		projectID, clusterName = acc.ClusterNameExecution(t)
		instanceName           = acc.RandomName()
		config                 = clusterStreamConnectionConfig(projectID, instanceName, clusterName)
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             clusterStreamConnectionAttributeChecks(resourceName, clusterName),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
