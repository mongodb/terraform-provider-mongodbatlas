package clusteroverloadsimulation_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName     = "mongodbatlas_cluster_overload_simulation.test"
	dataSourceName   = "data.mongodbatlas_cluster_overload_simulation.test"
	apiVersionHeader = "application/vnd.atlas.preview+json"
	readPath         = "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/overloadSimulations/{simulationId}"
)

func TestAccClusterOverloadSimulation_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{MongoDBMajorVersion: "9.0"})
		firstID     string
	)

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
					captureSimulationID(resourceName, &firstID),
				),
			},
			{
				Config: configBasic(&clusterInfo, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkBasic(clusterInfo.ProjectID, clusterInfo.Name, 3600),
					checkExists(resourceName),
					checkSimulationReplaced(resourceName, &firstID),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "simulation_id",
			},
		},
	})
}

func TestAccClusterOverloadSimulation_invalidDuration(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "mongodbatlas_cluster_overload_simulation" "test" {
						project_id       = "111111111111111111111111"
						cluster_name     = "Cluster0"
						duration_seconds = 901
					}
				`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("must be one of"),
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
		resource.TestCheckResourceAttrSet(resourceName, "state"),
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
		resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), readRequest(rs), nil)
		if resp != nil && resp.Body != nil {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				return fmt.Errorf("closing cluster overload simulation response: %w", closeErr)
			}
		}
		if err != nil {
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
		resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), readRequest(rs), nil)
		if resp != nil && resp.Body != nil {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				return fmt.Errorf("closing cluster overload simulation response: %w", closeErr)
			}
		}
		if validate.StatusNotFound(resp) {
			continue
		}
		if err != nil {
			return fmt.Errorf("checking cluster overload simulation destruction: %w", err)
		}
		return fmt.Errorf("cluster overload simulation %s still exists", rs.Primary.Attributes["simulation_id"])
	}
	return nil
}

func readRequest(rs *terraform.ResourceState) config.APICallParams {
	return config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  readPath,
		PathParams: map[string]string{
			"projectId":    rs.Primary.Attributes["project_id"],
			"clusterName":  rs.Primary.Attributes["cluster_name"],
			"simulationId": rs.Primary.Attributes["simulation_id"],
		},
		Method: "GET",
	}
}

func captureSimulationID(name string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		*target = rs.Primary.Attributes["simulation_id"]
		if *target == "" {
			return fmt.Errorf("simulation_id is empty for %s", name)
		}
		return nil
	}
}

func checkSimulationReplaced(name string, previousID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *previousID == "" {
			return errors.New("previous simulation_id is empty")
		}
		var currentID string
		if err := captureSimulationID(name, &currentID)(s); err != nil {
			return err
		}
		if currentID == *previousID {
			return fmt.Errorf("expected simulation %s to be replaced", *previousID)
		}
		return nil
	}
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
