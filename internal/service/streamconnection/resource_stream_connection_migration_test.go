package streamconnection_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		instanceName = acc.RandomName()
	)
	mig.SkipIfVersionBelow(t, "1.14.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckPreviewFlag(t); acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false),
				Check:             kafkaStreamConnectionAttributeChecks(resourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false),
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

func TestMigStreamRSStreamConnection_kafkaSSL(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		instanceName = acc.RandomName()
	)
	mig.SkipIfVersionBelow(t, "1.14.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckPreviewFlag(t); acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true),
				Check:             kafkaStreamConnectionAttributeChecks(resourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true),
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

func TestMigStreamRSStreamConnection_cluster(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		clusterInfo  = acc.GetClusterInfo(nil)
		instanceName = acc.RandomName()
	)
	mig.SkipIfVersionBelow(t, "1.15.2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckPreviewFlag(t); acc.PreCheckBasic(t) },
		CheckDestroy: CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            clusterStreamConnectionConfig(clusterInfo.ProjectIDStr, instanceName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
				Check:             clusterStreamConnectionAttributeChecks(resourceName, clusterInfo.ClusterName),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   clusterStreamConnectionConfig(clusterInfo.ProjectIDStr, instanceName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
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
