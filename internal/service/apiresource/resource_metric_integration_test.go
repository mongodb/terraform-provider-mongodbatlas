package apiresource_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// Atlas masks header values in read responses, showing leading asterisks
// followed by the last few characters of the secret.
var (
	maskedHeaderRegexp        = regexp.MustCompile(`^\*+.{1,8}$`)
	maskedHeaderRotatedSuffix = regexp.MustCompile(`\*+n-v2$`)
)

// TestAccAPIResource_otelMetricIntegration_basic covers the OTEL Metric
// Integration preview endpoint (private preview, no typed provider yet —
// PR #4418 has the autogen in flight). Exercises the full demo surface:
//
//   - preview = true
//   - sensitive_body for a nested list (headers carry an auth token)
//   - id_attribute = ["metricIntegrationId"] (API response key, not the
//     misleading "id" in the codegen yaml)
//   - update_method = "PUT" override (endpoint rejects PATCH with 405)
//   - body change (aggregationTemporality)
//   - sensitive_body rotation (bearer token)
func TestAccAPIResource_otelMetricIntegration_basic(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		tokenInitial = "dummy-bearer-token-for-poc"
		tokenRotated = "rotated-bearer-token-v2"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOTELMetricIntegration(projectID),
		Steps: []resource.TestStep{
			{
				Config: configOTELMetricIntegration(projectID, "DELTA", tokenInitial),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsOTELMetricIntegration(resourceName, projectID),
					resource.TestCheckResourceAttr(resourceName, "preview", "true"),
					resource.TestCheckResourceAttr(resourceName, "update_method", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "body.aggregationTemporality", "DELTA"),
					resource.TestCheckResourceAttr(resourceName, "body.integrationType", "OTEL"),
					resource.TestCheckResourceAttr(resourceName, "body.providerType", "CUSTOM"),
					resource.TestCheckResourceAttrSet(resourceName, "output.metricIntegrationId"),
					// Atlas masks header values in read responses — confirms the
					// sensitive bearer token never round-trips in the clear.
					resource.TestMatchResourceAttr(resourceName, "output.headers.0.value", maskedHeaderRegexp),
					resource.TestCheckResourceAttrSet(dataSourceName, "output.endpoint"),
				),
			},
			{
				// Update #1: body field (aggregationTemporality).
				Config: configOTELMetricIntegration(projectID, "CUMULATIVE", tokenInitial),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsOTELMetricIntegration(resourceName, projectID),
					resource.TestCheckResourceAttr(resourceName, "body.aggregationTemporality", "CUMULATIVE"),
				),
			},
			{
				// Update #2: sensitive_body rotation. State doesn't expose the
				// rotated value (sensitive_body is write-only); we verify the
				// updated mask reflects the new last 4 chars ("on-v2").
				Config: configOTELMetricIntegration(projectID, "CUMULATIVE", tokenRotated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsOTELMetricIntegration(resourceName, projectID),
					resource.TestMatchResourceAttr(resourceName, "output.headers.0.value", maskedHeaderRotatedSuffix),
				),
			},
		},
	})
}

func configOTELMetricIntegration(projectID, temporality, token string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_resource" "test" {
			path          = "/api/atlas/v2/groups/%[1]s/metricIntegrations"
			id_attribute  = ["metricIntegrationId"]
			preview       = true
			update_method = "PUT"

			body = {
				integrationType        = "OTEL"
				providerType           = "CUSTOM"
				aggregationTemporality = %[2]q
				endpoint               = "https://httpbin.org/post"
				metricSelection        = ["ATLAS_STREAM_PROCESSING"]
			}

			sensitive_body = {
				headers = [
					{
						name  = "Authorization"
						value = "Bearer %[3]s"
					},
				]
			}

			# Headers come back masked from Atlas (no real secret in the response),
			# so the masked values can safely live in the visible output.
			response_export_values = ["metricIntegrationId", "endpoint", "aggregationTemporality", "headers"]
		}

		data "mongodbatlas_api_resource" "test" {
			path    = mongodbatlas_api_resource.test.id
			preview = true

			response_export_values = ["endpoint", "aggregationTemporality"]
		}
	`, projectID, temporality, token)
}

func checkExistsOTELMetricIntegration(rsName, projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("not found: %s", rsName)
		}
		integrationID := rs.Primary.Attributes["output.metricIntegrationId"]
		if integrationID == "" {
			return fmt.Errorf("checkExists: output.metricIntegrationId not set for %s", rsName)
		}
		if !otelMetricIntegrationExists(projectID, integrationID) {
			return fmt.Errorf("otel metric integration (%s/%s) does not exist", projectID, integrationID)
		}
		return nil
	}
}

func checkDestroyOTELMetricIntegration(projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_api_resource" {
				continue
			}
			integrationID := rs.Primary.Attributes["output.metricIntegrationId"]
			if integrationID == "" {
				continue
			}
			if otelMetricIntegrationExists(projectID, integrationID) {
				return fmt.Errorf("otel metric integration (%s/%s) still exists", projectID, integrationID)
			}
		}
		return nil
	}
}

// otelMetricIntegrationExists uses UntypedAPICall — this endpoint is preview
// private and not yet in the SDK.
func otelMetricIntegrationExists(projectID, integrationID string) bool {
	params := config.APICallParams{
		Method:        http.MethodGet,
		VersionHeader: "application/vnd.atlas.preview+json",
		RelativePath:  "/api/atlas/v2/groups/{projectId}/metricIntegrations/{integrationId}",
		PathParams: map[string]string{
			"projectId":     projectID,
			"integrationId": integrationID,
		},
	}
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), params, nil)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err == nil
}
