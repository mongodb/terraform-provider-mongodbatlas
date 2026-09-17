package clusteroverloadsimulation_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_cluster_overload_simulation.test"
	dataSourceName = "data.mongodbatlas_cluster_overload_simulation.test"
)

func TestAccClusterOverloadSimulation_basic(t *testing.T) {
	clusterInfo := acc.GetClusterInfo(t, &acc.ClusterRequest{MongoDBMajorVersion: "9.0"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, 900),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkBasic(clusterInfo.ProjectID, clusterInfo.Name, 900),
					checkExists(resourceName),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"delete_on_create_timeout"},
				ImportStateVerifyIdentifierAttribute: "simulation_id",
			},
		},
	})
}

func configBasic(info *acc.ClusterInfo, durationSeconds int64) string {
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_cluster_overload_simulation" "test" {
			project_id       = %[2]q
			cluster_name     = %[3]s
			duration_seconds = %[4]d
		}

		data "mongodbatlas_cluster_overload_simulation" "test" {
			project_id    = mongodbatlas_cluster_overload_simulation.test.project_id
			cluster_name  = mongodbatlas_cluster_overload_simulation.test.cluster_name
			simulation_id = mongodbatlas_cluster_overload_simulation.test.simulation_id
		}
	`, info.TerraformStr, info.ProjectID, info.TerraformNameRef, durationSeconds)
}

func checkBasic(projectID, clusterName string, durationSeconds int64) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
		resource.TestCheckResourceAttr(resourceName, "duration_seconds", strconv.FormatInt(durationSeconds, 10)),
		resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
		resource.TestCheckResourceAttrSet(resourceName, "request_date"),
		resource.TestCheckResourceAttrSet(resourceName, "simulation_id"),
		resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
	}
	for _, attr := range []string{"project_id", "cluster_name", "duration_seconds", "expires_at", "request_date", "simulation_id", "state"} {
		checks = append(checks, resource.TestCheckResourceAttrPair(dataSourceName, attr, resourceName, attr))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if err := getClusterOverloadSimulation(rs); err != nil {
			return fmt.Errorf("cluster overload simulation does not exist: %w", err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_overload_simulation" {
			continue
		}
		if err := getClusterOverloadSimulation(rs); err == nil {
			return fmt.Errorf("cluster overload simulation %s still exists", rs.Primary.Attributes["simulation_id"])
		}
	}
	return nil
}

func getClusterOverloadSimulation(rs *terraform.ResourceState) error {
	_, _, err := acc.ConnV2().OverloadProtectionSimulationAPI.GetClusterOverloadSimulation(
		context.Background(),
		rs.Primary.Attributes["project_id"],
		rs.Primary.Attributes["cluster_name"],
		rs.Primary.Attributes["simulation_id"],
	).Execute()
	return err
}

func importStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["cluster_name"],
			rs.Primary.Attributes["simulation_id"],
		), nil
	}
}
