package metricintegration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName     = "mongodbatlas_metric_integration.test"
	apiVersionHeader = "application/vnd.atlas.preview+json"
	otelAPIKeyEnvVar = "OTEL_API_KEY"
)

func TestAccMetricIntegrationAPI_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		apiKey    = os.Getenv(otelAPIKeyEnvVar)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); preCheckOTEL(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, apiKey, "DELTA"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "integration_type", "OTEL"),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "aggregation_temporality", "DELTA"),
					resource.TestCheckResourceAttr(resourceName, "endpoint", "https://otlp.datadoghq.com/v1/metrics"),
					resource.TestCheckResourceAttr(resourceName, "metric_selection.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "metric_selection.*", "ATLAS_STREAM_PROCESSING"),
					resource.TestCheckResourceAttr(resourceName, "headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "headers.0.name", "dd-api-key"),
				),
			},
			{
				Config: configBasic(projectID, apiKey, "CUMULATIVE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "aggregation_temporality", "CUMULATIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"headers.0.value"},
			},
		},
	})
}

func configBasic(projectID, apiKey, aggregationTemporality string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_metric_integration" "test" {
			group_id                 = %[1]q
			integration_type         = "OTEL"
			provider_type            = "CUSTOM"
			aggregation_temporality  = %[4]q
			endpoint                 = "https://otlp.datadoghq.com/v1/metrics"
			metric_selection         = ["ATLAS_STREAM_PROCESSING"]
			headers = [
				{
					name  = "dd-api-key"
					value = %[3]q
				}
			]
		}
	`, projectID, projectID, apiKey, aggregationTemporality)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		id := rs.Primary.Attributes["id"]
		if groupID == "" || id == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		resp, err := metricIntegrationAPICall(context.Background(), groupID, id)
		if err != nil {
			return fmt.Errorf("metric integration (%s/%s) does not exist: %w", groupID, id, err)
		}
		resp.Body.Close()
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_metric_integration" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		id := rs.Primary.Attributes["id"]
		if groupID == "" || id == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		resp, err := metricIntegrationAPICall(context.Background(), groupID, id)
		if err == nil {
			resp.Body.Close()
			return fmt.Errorf("metric integration (%s/%s) still exists", groupID, id)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		id := rs.Primary.Attributes["id"]
		if groupID == "" || id == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", groupID, id), nil
	}
}

func metricIntegrationAPICall(ctx context.Context, groupID, id string) (*http.Response, error) {
	return acc.MongoDBClient.UntypedAPICall(ctx, config.APICallParams{
		VersionHeader: apiVersionHeader,
		RelativePath:  "/api/atlas/v2/groups/{groupId}/metricIntegrations/{id}",
		PathParams:    map[string]string{"groupId": groupID, "id": id},
		Method:        "GET",
	}, nil)
}

func preCheckOTEL(t *testing.T) {
	t.Helper()
	if os.Getenv(otelAPIKeyEnvVar) == "" {
		t.Skipf("%s must be set for this acceptance test", otelAPIKeyEnvVar)
	}
}

