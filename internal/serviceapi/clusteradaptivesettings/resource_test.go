package clusteradaptivesettings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
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
				Config: configBasic(projectID, clusterName, `{
					OVERLOAD_PROTECTION        = true
					SEARCH_OVERLOAD_PROTECTION = true
				}`),
				Check: checkResourceAndDataSource(
					`{"OVERLOAD_PROTECTION":true,"SEARCH_OVERLOAD_PROTECTION":true}`,
				),
			},
			{
				Config: configBasic(projectID, clusterName, `{
					OVERLOAD_PROTECTION = false
				}`),
				Check: checkResourceAndDataSource(
					`{"OVERLOAD_PROTECTION":false}`,
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
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
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

func importStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return fmt.Sprintf("%s/%s",
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["cluster_name"],
		), nil
	}
}
