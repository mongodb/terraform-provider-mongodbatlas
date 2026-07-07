package metricintegration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

const (
	resourceName         = "mongodbatlas_metric_integration.test"
	dataSourceName       = "data.mongodbatlas_metric_integration.test"
	pluralDataSourceName = "data.mongodbatlas_metric_integrations.test"
	datasourcesConfig    = `
		data "mongodbatlas_metric_integration" "test" {
			project_id            = mongodbatlas_metric_integration.test.project_id
			metric_integration_id = mongodbatlas_metric_integration.test.metric_integration_id
		}

		data "mongodbatlas_metric_integrations" "test" {
			project_id = mongodbatlas_metric_integration.test.project_id
			depends_on = [mongodbatlas_metric_integration.test]
		}
	`
)

func TestAccMetricIntegration_basic(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		integrationType = "OTEL"
		providerType    = "CUSTOM"
		aggregation     = "DELTA"
		endpoint        = os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT")
		headerValue     = os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY")
		metricSelection = []string{"ATLAS_STREAM_PROCESSING"}
		extraHeader     = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); preCheckMetricIntegration(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, !extraHeader, true),
				Check:  checkBasic(integrationType, providerType, aggregation, endpoint, metricSelection, !extraHeader, true),
			},
			{
				Config: configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, extraHeader, false),
				Check:  checkBasic(integrationType, providerType, aggregation, endpoint, metricSelection, extraHeader, false),
			},
			{
				Config:                               configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue, metricSelection, extraHeader, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "metric_integration_id",
				ImportStateVerifyIgnore:              []string{"headers"}, // header values are redacted on GET
			},
		},
	})
}

func preCheckMetricIntegration(tb testing.TB) {
	tb.Helper()
	if os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT") == "" || os.Getenv("MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY") == "" {
		tb.Fatal("`MONGODB_ATLAS_METRIC_INTEGRATION_ENDPOINT` and `MONGODB_ATLAS_METRIC_INTEGRATION_API_KEY` must be set for acceptance testing")
	}
}

func configBasic(projectID, integrationType, providerType, aggregation, endpoint, headerValue string, metricSelection []string, extraHeader, withDS bool) string {
	selectionHCL := hcl.StringSliceToHCL(metricSelection)
	extraHeaderHCL := ""
	if extraHeader {
		extraHeaderHCL = `,
				{
					name  = "x-custom-header"
					value = "custom-value"
				}`
	}
	dsConfig := ""
	if withDS {
		dsConfig = datasourcesConfig
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_metric_integration" "test" {
			project_id              = %[1]q
			integration_type        = %[2]q
			provider_type           = %[3]q
			aggregation_temporality = %[4]q
			endpoint                = %[5]q
			metric_selection        = %[6]s

			headers = [
				{
					name  = "dd-api-key"
					value = %[7]q
				}%[8]s
			]
		}

		%[9]s
	`, projectID, integrationType, providerType, aggregation, endpoint, selectionHCL, headerValue, extraHeaderHCL, dsConfig)
}

func checkBasic(integrationType, providerType, aggregation, endpoint string, metricSelection []string, extraHeader, withDS bool) resource.TestCheckFunc {
	headerCount := "1"
	if extraHeader {
		headerCount = "2"
	}
	setChecks := []string{"project_id", "metric_integration_id"}
	mapChecks := map[string]string{
		"integration_type":        integrationType,
		"provider_type":           providerType,
		"aggregation_temporality": aggregation,
		"endpoint":                endpoint,
		"metric_selection.#":      strconv.Itoa(len(metricSelection)),
		"headers.#":               headerCount,
	}
	var checks []resource.TestCheckFunc
	var dsName *string
	if withDS {
		dsName = new(dataSourceName)
		checks = append(checks, resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)))
	}
	checks = append(checks, acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		// The metric integration API is preview-only and not yet part of the Atlas SDK, so existence
		// is verified with a raw HTTP call using the preview version header.
		// TODO: swap the raw HTTP call and preview header for the generated SDK method once the Atlas
		// SDK supports the preview media type (CLOUDP-417642).
		resp, err := acc.GetMetricIntegration(context.Background(), projectID, integrationID)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("metric integration for project_id %s with id %s does not exist, status %d", projectID, integrationID, resp.StatusCode)
	}
}

func checkDestroy(state *terraform.State) error {
	for name, rs := range state.RootModule().Resources {
		if name != resourceName {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		resp, err := acc.GetMetricIntegration(context.Background(), projectID, integrationID)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("metric integration for project_id %s with id %s still exists", projectID, integrationID)
		}
		return fmt.Errorf("checkDestroy, unexpected status %d for project_id %s with id %s", resp.StatusCode, projectID, integrationID)
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		integrationID := rs.Primary.Attributes["metric_integration_id"]
		if projectID == "" || integrationID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, integrationID), nil
	}
}
