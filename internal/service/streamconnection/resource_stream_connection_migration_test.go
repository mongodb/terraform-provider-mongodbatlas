package streamconnection_test

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA
	mig.CreateAndRunTest(t, testCaseKafkaPlaintext(t))
}

func TestMigStreamRSStreamConnection_cluster(t *testing.T) {
	var (
		resourceName           = "mongodbatlas_stream_connection.test"
		projectID, clusterName = acc.ClusterNameExecution(t)
		instanceName           = acc.RandomName()
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            clusterStreamConnectionConfig(projectID, instanceName, clusterName),
				Check:             clusterStreamConnectionAttributeChecks(resourceName, clusterName),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   clusterStreamConnectionConfig(projectID, instanceName, clusterName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
