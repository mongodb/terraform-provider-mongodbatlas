package orglogintegration_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_org_log_integration.test"
	dataSourceName       = "data.mongodbatlas_org_log_integration.test"
	pluralDataSourceName = "data.mongodbatlas_org_log_integrations.test"
	datasourcesConfig    = `
		data "mongodbatlas_org_log_integration" "test" {
			org_id         = mongodbatlas_org_log_integration.test.org_id
			integration_id = mongodbatlas_org_log_integration.test.integration_id
		}

		data "mongodbatlas_org_log_integrations" "test" {
			org_id     = mongodbatlas_org_log_integration.test.org_id
			depends_on = [mongodbatlas_org_log_integration.test]
		}
	`
)

var logTypesEvents = []string{"EVENTS"}

func TestAccOrgLogIntegration_basicOTel(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		// Dummy endpoints that pass the API validity check.
		endpoint0   = "https://192.0.2.1/v1/logs"
		endpoint1   = "https://192.0.2.2/v1/logs"
		headersHCL0 = `otel_supplied_headers = [{name = "header-0", value = "val-0"}]`
		headersHCL1 = `otel_supplied_headers = [{name = "header-0", value = "val-0-updated"},{name = "header-1", value = "val-1"}]`
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicOTel(orgID, endpoint0, headersHCL0, true),
				Check:  checkBasicOTel(endpoint0, true),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(
						pluralDataSourceName,
						"otel_endpoint",
						knownvalue.StringExact(endpoint0),
						map[string]knownvalue.Check{
							"type": knownvalue.StringExact("OTEL_LOG_EXPORT"),
						},
					),
				},
			},
			{
				Config: configBasicOTel(orgID, endpoint1, headersHCL1, false),
				Check:  checkBasicOTel(endpoint1, false),
			},
			{
				Config:                               configBasicOTel(orgID, endpoint1, headersHCL1, false),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "integration_id",
				ImportStateVerifyIgnore:              []string{"otel_supplied_headers"}, // otel_supplied_headers values are redacted on GET
			},
		},
	})
}

func configBasicOTel(orgID, endpoint, headersHCL string, withDS bool) string {
	dsConfig := ""
	if withDS {
		dsConfig = datasourcesConfig
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_org_log_integration" "test" {
			org_id        = %[1]q
			type          = "OTEL_LOG_EXPORT"
			log_types     = ["EVENTS"]
			otel_endpoint = %[2]q
			%[3]s
		}

		%[4]s
	`, orgID, endpoint, headersHCL, dsConfig)
}

func checkBasicOTel(endpoint string, withDS bool) resource.TestCheckFunc {
	setChecks := []string{"integration_id"}
	mapChecks := map[string]string{
		"type":          "OTEL_LOG_EXPORT",
		"otel_endpoint": endpoint,
		"log_types.#":   strconv.Itoa(len(logTypesEvents)),
		"log_types.0":   logTypesEvents[0],
	}
	var dsName *string
	if withDS {
		dsName = new(dataSourceName)
	}
	headerChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrWith(resourceName, "otel_supplied_headers.#", acc.IntGreatThan(0)),
	}
	if withDS {
		headerChecks = append(headerChecks, resource.TestCheckResourceAttrWith(dataSourceName, "otel_supplied_headers.#", acc.IntGreatThan(0)))
	}
	checks := []resource.TestCheckFunc{
		acc.CheckRSAndDS(resourceName, dsName, nil, setChecks, mapChecks, checkExists(resourceName)),
	}
	checks = append(checks, headerChecks...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// TODO(CLOUDP-433802): switch from preview SDK to prod SDK once Org Log Integrations API is GA
func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		orgID := rs.Primary.Attributes["org_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if orgID == "" || integrationID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().PushBasedLogExportAPI.GetOrgLogIntegration(context.Background(), orgID, integrationID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("org log integration for org_id %s with id %s does not exist", orgID, integrationID)
	}
}

// TODO(CLOUDP-433802): switch from preview SDK to prod SDK once Org Log Integrations API is GA
func checkDestroy(state *terraform.State) error {
	for name, rs := range state.RootModule().Resources {
		if name != resourceName {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if orgID == "" || integrationID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().PushBasedLogExportAPI.GetOrgLogIntegration(context.Background(), orgID, integrationID).Execute()
		if err == nil {
			return fmt.Errorf("org log integration for org_id %s with id %s still exists", orgID, integrationID)
		}
		return nil
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		integrationID := rs.Primary.Attributes["integration_id"]
		if orgID == "" || integrationID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", orgID, integrationID), nil
	}
}
