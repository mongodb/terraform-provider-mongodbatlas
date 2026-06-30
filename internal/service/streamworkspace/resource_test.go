package streamworkspace_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	region         = "VIRGINIA_USA"
	cloudProvider  = "AWS"
	failoverRegion = "DUBLIN_IRL"
)

func TestAccStreamWorkspaceRS_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_stream_workspace.test"
		dataSourceName = "data.mongodbatlas_stream_workspace.test"
		pluralDSName   = "data.mongodbatlas_stream_workspaces.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
	)
	attrsMap := map[string]string{
		"workspace_name":                     workspaceName,
		"data_process_region.region":         region,
		"data_process_region.cloud_provider": cloudProvider,
		"stream_config.max_tier_size":        "SP30",
		"stream_config.tier":                 "SP10",
	}
	attrsSet := []string{"project_id", "hostnames.#"}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspaceResourceWithDataSourcesConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamsWorkspaceExists(resourceName),
					acc.CheckRSAndDS(resourceName, &dataSourceName, &pluralDSName, attrsSet, attrsMap),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamsWorkspaceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStreamWorkspaceRS_withFailoverRegions(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	attrsMap := map[string]string{
		"workspace_name":                     workspaceName,
		"data_process_region.region":         region,
		"data_process_region.cloud_provider": cloudProvider,
		"failover_regions.#":                 "1",
		"failover_regions.0.cloud_provider":  cloudProvider,
		"failover_regions.0.region":          failoverRegion,
		"stream_config.tier":                 "SP10",
	}
	attrsSet := []string{"project_id", "hostnames.#"}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, failoverRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamsWorkspaceExists(resourceName),
					acc.CheckRSAndDS(resourceName, nil, nil, attrsSet, attrsMap),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamsWorkspaceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStreamWorkspaceRS_updateWithFailoverRegions(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			// Step 1: create without failover_regions
			{
				Config: streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamsWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "failover_regions.#", "0"),
				),
			},
			// Step 2: add failover_regions via update (null → value, allowed)
			{
				Config: streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, failoverRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamsWorkspaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "failover_regions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "failover_regions.0.region", failoverRegion),
					resource.TestCheckResourceAttr(resourceName, "failover_regions.0.cloud_provider", cloudProvider),
				),
			},
			// Step 3: no-op plan — failover_regions unchanged, expect no diff
			{
				Config:   streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, failoverRegion),
				PlanOnly: true,
			},
		},
	})
}

func TestAccStreamWorkspaceRS_failoverWriteOnceRequiresReplace(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			// Step 1: create with DUBLIN_IRL failover
			{
				Config: streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, "DUBLIN_IRL"),
				Check:  checkStreamsWorkspaceExists(resourceName),
			},
			// Step 2: verify that changing failover_regions produces a non-empty plan (replacement required).
			// PlanOnly avoids applying with OREGON_USA which may be unsupported in some environments.
			// The specific DestroyBeforeCreate action is covered by the unit test for failoverRegionsWriteOnce.
			{
				Config:             streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, "OREGON_USA"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccStreamWorkspaceRS_simultaneousDataProcessRegionAndFailoverChange(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider),
			},
			{
				Config:      streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, "OREGON_USA", cloudProvider, failoverRegion),
				ExpectError: regexp.MustCompile(`data_process_region and failover_regions cannot both be changed`),
			},
		},
	})
}

func checkStreamsWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		projectID := rs.Primary.Attributes["project_id"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamWorkspace(context.Background(), projectID, workspaceName).Execute()
		if err != nil {
			return fmt.Errorf("stream workspace (%s:%s) does not exist: %s", projectID, workspaceName, err)
		}
		return nil
	}
}

func checkStreamsWorkspaceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["workspace_name"]), nil
	}
}

func streamsWorkspaceWithFailoverRegionsConfig(projectID, workspaceName, region, cloudProvider, failoverRegion string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_workspace" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			data_process_region = {
				region = %[3]q
				cloud_provider = %[4]q
			}
			failover_regions = [
				{
					cloud_provider = %[4]q
					region = %[5]q
				}
			]
			stream_config = {
				tier = "SP10"
			}
		}
	`, projectID, workspaceName, region, cloudProvider, failoverRegion)
}

func streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return streamsWorkspaceWithStreamConfigConfig(projectID, workspaceName, region, cloudProvider, "SP10", "SP30")
}

func streamsWorkspaceWithStreamConfigConfig(projectID, workspaceName, region, cloudProvider, tier, maxTierSize string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_workspace" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			data_process_region = {
				region = %[3]q
				cloud_provider = %[4]q
			}
			stream_config = {
				tier = %[5]q
				max_tier_size = %[6]q
			}
		}
	`, projectID, workspaceName, region, cloudProvider, tier, maxTierSize)
}

func streamsWorkspaceResourceWithDataSourcesConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspace" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
			workspace_name = mongodbatlas_stream_workspace.test.workspace_name
		}

		data "mongodbatlas_stream_workspaces" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
		}
	`, streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider))
}
