package serverlessinstance_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccMigrationServerlessInstance_basic(t *testing.T) {
	var (
		serverlessInstance matlas.Cluster
		resourceName       = "mongodbatlas_serverless_instance.test"
		instanceName       = acctest.RandomWithPrefix("test-acc-serverless")
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acctest.RandomWithPrefix("test-acc-serverless")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            acc.ConfigServerlessInstanceBasic(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					checkConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName),
					checkExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "termination_protection_enabled", "false"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigServerlessInstanceBasic(orgID, projectName, instanceName, true),
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
