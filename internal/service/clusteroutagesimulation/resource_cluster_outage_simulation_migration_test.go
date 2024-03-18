package clusteroutagesimulation_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigOutageSimulationCluster_SingleRegion_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster_outage_simulation.test_outage"
		projectID    = mig.ProjectIDGlobal(t)
		clusterName  = acc.RandomClusterName()
		config       = configSingleRegion(projectID, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(resourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigOutageSimulationCluster_MultiRegion_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster_outage_simulation.test_outage"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		clusterName  = acc.RandomClusterName()
		config       = configMultiRegion(projectName, orgID, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(resourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
