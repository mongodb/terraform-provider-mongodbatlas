package federateddatabaseinstance_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigFederatedDatabaseInstance_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_federated_database_instance.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		name         = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyFederatedDatabaseInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configFirstSteps(name, projectName, orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage_stores.0.read_preference.0.tag_sets.#"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_stores.0.read_preference.0.tag_sets.0.tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_databases.0.collections.0.data_sources.0.database", "sample_airbnb"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configFirstSteps(name, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
