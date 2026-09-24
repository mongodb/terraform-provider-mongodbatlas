package clusteradaptivesettings_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
)

const (
	resourceName   = "mongodbatlas_cluster_adaptive_settings.test"
	dataSourceName = "data.mongodbatlas_cluster_adaptive_settings.test"
)

func TestAccClusterAdaptiveSettings_basic(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				// The first apply verifies that omitting the optional overrides creates stable state.
				Config: configBasic(projectID, clusterName, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: checkNoOverridesResourceAndDataSource(),
			},
			{
				// Initialize an empty object when Atlas has no stored overrides.
				Config: configBasic(projectID, clusterName, `{}`),
				Check:  checkResourceAndDataSource(`{}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: configBasic(projectID, clusterName, `{
					LOAD_SHEDDING              = true
					SEARCH_LOAD_SHEDDING = true
				}`),
				Check: checkResourceAndDataSource(
					`{"LOAD_SHEDDING":true,"SEARCH_LOAD_SHEDDING":true}`,
				),
			},
			{
				// Change Load Shedding independently without changing Search Load Shedding.
				Config: configBasic(projectID, clusterName, `{
					LOAD_SHEDDING              = false
					SEARCH_LOAD_SHEDDING = true
				}`),
				Check: checkResourceAndDataSource(`{"LOAD_SHEDDING":false,"SEARCH_LOAD_SHEDDING":true}`),
			},
			{
				Config: configBasic(projectID, clusterName, `{
					LOAD_SHEDDING = false
				}`),
				Check: checkResourceAndDataSource(
					`{"LOAD_SHEDDING":false}`,
				),
			},
			{
				// Removing the attribute after a non-empty map resets the whole API map with null.
				Config: configBasic(projectID, clusterName, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: checkNoOverridesResourceAndDataSource(),
			},
			{
				Config: configBasic(projectID, clusterName, `{
					SEARCH_LOAD_SHEDDING = false
				}`),
				Check: checkResourceAndDataSource(
					`{"SEARCH_LOAD_SHEDDING":false}`,
				),
			},
			{
				// An explicitly configured empty map resets every effective key while remaining {} in state.
				Config: configBasic(projectID, clusterName, `{}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: checkResourceAndDataSource(
					`{}`,
				),
			},
			{
				// Leave an override set so import and destroy verify populated settings.
				Config: configBasic(projectID, clusterName, `{ SEARCH_LOAD_SHEDDING = false }`),
				Check:  checkResourceAndDataSource(`{"SEARCH_LOAD_SHEDDING":false}`),
			},
			{
				// Reset outside Terraform, then verify unchanged configuration restores the override.
				// This step would fail without the PostRead hook: the decoder skips the absent
				// field, state stays stale, and Terraform reports no drift.
				PreConfig: func() { resetOverrides(t, projectID, clusterName) },
				Config:    configBasic(projectID, clusterName, `{ SEARCH_LOAD_SHEDDING = false }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: checkResourceAndDataSource(`{"SEARCH_LOAD_SHEDDING":false}`),
			},
			{
				ResourceName:                         resourceName,
				ImportStateId:                        projectID + "/" + clusterName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
}

func resetOverrides(t *testing.T, projectID, clusterName string) {
	t.Helper()
	req := admin.NewAdaptiveSettingsUpdateRequest()
	req.NullFields = []string{"AdaptiveSettingsOverrides"}
	_, _, err := acc.ConnV2().ClustersAPI.UpdateClusterAdaptiveSettings(t.Context(), projectID, clusterName, req).Execute()
	require.NoError(t, err)
}

func configBasic(projectID, clusterName, overrides string) string {
	overridesLine := ""
	if overrides != "" {
		overridesLine = fmt.Sprintf("adaptive_settings_overrides = jsonencode(%s)", overrides)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id   = %[1]q
			cluster_name = %[2]q
			%[3]s
		}

		data "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id   = mongodbatlas_cluster_adaptive_settings.test.project_id
			cluster_name = mongodbatlas_cluster_adaptive_settings.test.cluster_name
		}
	`, projectID, clusterName, overridesLine)
}

func checkResourceAndDataSource(expectedOverrides string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "adaptive_settings_overrides", expectedOverrides),
		resource.TestCheckResourceAttrSet(resourceName, "effective_adaptive_settings"),
		resource.TestCheckResourceAttrPair(dataSourceName, "adaptive_settings_overrides", resourceName, "adaptive_settings_overrides"),
		resource.TestCheckResourceAttrPair(dataSourceName, "effective_adaptive_settings", resourceName, "effective_adaptive_settings"),
	)
}

func checkNoOverridesResourceAndDataSource() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkExists(resourceName),
		resource.TestCheckNoResourceAttr(resourceName, "adaptive_settings_overrides"),
		resource.TestCheckNoResourceAttr(dataSourceName, "adaptive_settings_overrides"),
		resource.TestCheckResourceAttrSet(resourceName, "effective_adaptive_settings"),
		resource.TestCheckResourceAttrPair(dataSourceName, "effective_adaptive_settings", resourceName, "effective_adaptive_settings"),
	)
}

func checkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		_, _, err := acc.ConnV2().ClustersAPI.GetClusterAdaptiveSettings(
			context.Background(),
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["cluster_name"],
		).Execute()
		if err != nil {
			return fmt.Errorf("cluster adaptive settings do not exist: %w", err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_adaptive_settings" {
			continue
		}
		settings, _, err := acc.ConnV2().ClustersAPI.GetClusterAdaptiveSettings(
			context.Background(),
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["cluster_name"],
		).Execute()
		if err != nil {
			return fmt.Errorf("checking cluster adaptive settings reset: %w", err)
		}
		if settings.AdaptiveSettingsOverrides != nil && len(*settings.AdaptiveSettingsOverrides) > 0 {
			return fmt.Errorf("cluster adaptive settings for %s were not reset", rs.Primary.Attributes["cluster_name"])
		}
	}
	return nil
}
