package clusteradaptivesettings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/require"
)

const (
	resourceName     = "mongodbatlas_cluster_adaptive_settings.test"
	dataSourceName   = "data.mongodbatlas_cluster_adaptive_settings.test"
	apiVersionHeader = "application/vnd.atlas.2025-03-12+json"
	readPath         = "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/adaptiveSettings"
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
				Config: configWithoutOverrides(projectID, clusterName),
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
					SEARCH_OVERLOAD_PROTECTION = true
				}`),
				Check: checkResourceAndDataSource(
					`{"LOAD_SHEDDING":true,"SEARCH_OVERLOAD_PROTECTION":true}`,
				),
			},
			{
				// Change Load Shedding independently without changing Search Load Shedding.
				Config: configBasic(projectID, clusterName, `{
					LOAD_SHEDDING              = false
					SEARCH_OVERLOAD_PROTECTION = true
				}`),
				Check: checkResourceAndDataSource(`{"LOAD_SHEDDING":false,"SEARCH_OVERLOAD_PROTECTION":true}`),
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
				Config: configWithoutOverrides(projectID, clusterName),
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
					SEARCH_OVERLOAD_PROTECTION = false
				}`),
				Check: checkResourceAndDataSource(
					`{"SEARCH_OVERLOAD_PROTECTION":false}`,
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
				Config: configBasic(projectID, clusterName, `{ SEARCH_OVERLOAD_PROTECTION = false }`),
				Check:  checkResourceAndDataSource(`{"SEARCH_OVERLOAD_PROTECTION":false}`),
			},
			{
				// Reset outside Terraform, then verify unchanged configuration restores the override.
				PreConfig: func() { resetOverrides(t, projectID, clusterName) },
				Config:    configBasic(projectID, clusterName, `{ SEARCH_OVERLOAD_PROTECTION = false }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: checkResourceAndDataSource(`{"SEARCH_OVERLOAD_PROTECTION":false}`),
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

func TestAccClusterAdaptiveSettings_invalidOverrides(t *testing.T) {
	tests := map[string]struct {
		overrides string
		errorText string
	}{
		"null entry": {overrides: `{ SEARCH_OVERLOAD_PROTECTION = null }`, errorText: `Override "SEARCH_OVERLOAD_PROTECTION" is null`},
		"JSON null":  {overrides: `null`, errorText: "Use jsonencode with an object"},
		"array":      {overrides: `[]`, errorText: "Use jsonencode with an object"},
		"scalar":     {overrides: `false`, errorText: "Use jsonencode with an object"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Steps: []resource.TestStep{{
					Config:      configBasic("111111111111111111111111", "test", test.overrides),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(test.errorText),
				}},
			})
		})
	}
}

func resetOverrides(t *testing.T, projectID, clusterName string) {
	t.Helper()
	resp, err := acc.MongoDBClient.UntypedAPICall(t.Context(), config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  readPath,
		PathParams:    map[string]string{"projectId": projectID, "clusterName": clusterName},
		Method:        "PATCH",
	}, []byte(`{"adaptiveSettingsOverrides":null}`))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	require.NoError(t, err)
}

func configBasic(projectID, clusterName, overrides string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id                  = %[1]q
			cluster_name                = %[2]q
			adaptive_settings_overrides = jsonencode(%[3]s)
		}

		data "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id   = mongodbatlas_cluster_adaptive_settings.test.project_id
			cluster_name = mongodbatlas_cluster_adaptive_settings.test.cluster_name
		}
	`, projectID, clusterName, overrides)
}

func configWithoutOverrides(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id   = %[1]q
			cluster_name = %[2]q
		}

		data "mongodbatlas_cluster_adaptive_settings" "test" {
			project_id   = mongodbatlas_cluster_adaptive_settings.test.project_id
			cluster_name = mongodbatlas_cluster_adaptive_settings.test.cluster_name
		}
	`, projectID, clusterName)
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
		_, err := readAdaptiveSettings(rs)
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
		body, err := readAdaptiveSettings(rs)
		if err != nil {
			return fmt.Errorf("checking cluster adaptive settings reset: %w", err)
		}
		var settings struct {
			Overrides map[string]json.RawMessage `json:"adaptiveSettingsOverrides"`
		}
		if err := json.Unmarshal(body, &settings); err != nil {
			return fmt.Errorf("decoding cluster adaptive settings: %w", err)
		}
		if len(settings.Overrides) != 0 {
			return fmt.Errorf("cluster adaptive settings for %s were not reset", rs.Primary.Attributes["cluster_name"])
		}
	}
	return nil
}

func readAdaptiveSettings(rs *terraform.ResourceState) ([]byte, error) {
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  readPath,
		PathParams: map[string]string{
			"projectId":   rs.Primary.Attributes["project_id"],
			"clusterName": rs.Primary.Attributes["cluster_name"],
		},
		Method: "GET",
	}, nil)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("empty cluster adaptive settings response")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading cluster adaptive settings response: %w", err)
	}
	return body, nil
}
