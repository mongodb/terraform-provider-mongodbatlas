package thirdpartyintegration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// dummy keys used for credential values in third party notifications
const dummy32CharKey = "11111111111111111111111111111111"
const dummy32CharKeyUpdated = "11111111111111111111111111111112"
const dummy36CharKey = "11111111-1111-1111-1111-111111111111"
const dummy36CharKeyUpdated = "11111111-1111-1111-1111-111111111112"

const resourceName = "mongodbatlas_third_party_integration.test"
const dataSourceName = "data." + resourceName
const dataSourcePluralName = "data.mongodbatlas_third_party_integrations.test"

func TestAccThirdPartyIntegration_basicPagerDuty(t *testing.T) {
	// basic test also include testing of plural data source which is why it cannot run in parallel
	resource.Test(t, *basicPagerDutyTest(t))
}

func TestAccThirdPartyIntegration_opsGenie(t *testing.T) {
	resource.ParallelTest(t, *opsGenieTest(t))
}

func TestAccThirdPartyIntegration_victorOps(t *testing.T) {
	resource.ParallelTest(t, *victorOpsTest(t))
}

func TestAccThirdPartyIntegration_datadog(t *testing.T) {
	resource.ParallelTest(t, *datadogTest(t))
}

func TestAccThirdPartyIntegration_prometheus(t *testing.T) {
	resource.ParallelTest(t, *prometheusTest(t))
}

func TestAccThirdPartyIntegration_microsoftTeams(t *testing.T) {
	resource.ParallelTest(t, *microsoftTeamsTest(t))
}

func TestAccThirdPartyIntegration_webhook(t *testing.T) {
	resource.ParallelTest(t, *webhookTest(t))
}

func basicPagerDutyTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID         = acc.ProjectIDExecution(tb)
		serviceKey        = dummy32CharKey
		updatedServiceKey = dummy32CharKeyUpdated
		intType           = "PAGER_DUTY"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicPagerDuty(projectID, serviceKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "service_key", serviceKey),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.type", intType),
				),
			},
			{
				Config: configBasicPagerDuty(projectID, updatedServiceKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "service_key", updatedServiceKey),
				),
			},
			importStep(resourceName),
		},
	}
}

func opsGenieTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID     = acc.ProjectIDExecution(tb)
		apiKey        = dummy36CharKey
		updatedAPIKey = dummy36CharKeyUpdated
		intType       = "OPS_GENIE"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configOpsGenie(projectID, apiKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttr(resourceName, "region", "US"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "region", "US"),
				),
			},
			{
				Config: configOpsGenie(projectID, updatedAPIKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
					resource.TestCheckResourceAttr(resourceName, "region", "US"),
				),
			},
			importStep(resourceName),
		},
	}
}

func victorOpsTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID     = acc.ProjectIDExecution(tb)
		apiKey        = dummy36CharKey
		updatedAPIKey = dummy36CharKeyUpdated
		intType       = "VICTOR_OPS"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configVictorOps(projectID, apiKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttrSet(resourceName, "routing_key"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttrSet(dataSourceName, "routing_key"),
				),
			},
			{
				Config: configVictorOps(projectID, updatedAPIKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
				),
			},
			importStep(resourceName),
		},
	}
}

func datadogTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID     = acc.ProjectIDExecution(tb)
		apiKey        = dummy32CharKey
		updatedAPIKey = dummy32CharKeyUpdated
		region        = "US"
		intType       = "DATADOG"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDatadog(projectID, apiKey, "US", false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "send_collection_latency_metrics", "false"),
					resource.TestCheckResourceAttr(resourceName, "send_database_metrics", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "region", region),
					resource.TestCheckResourceAttr(dataSourceName, "send_collection_latency_metrics", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "send_database_metrics", "false"),
				),
			},
			{
				Config: configDatadog(projectID, apiKey, "US", true, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "send_collection_latency_metrics", "true"),
					resource.TestCheckResourceAttr(resourceName, "send_database_metrics", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "region", region),
					resource.TestCheckResourceAttr(dataSourceName, "send_collection_latency_metrics", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "send_database_metrics", "false"),
				),
			},
			{
				Config: configDatadog(projectID, updatedAPIKey, "US", true, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "send_collection_latency_metrics", "false"),
					resource.TestCheckResourceAttr(resourceName, "send_database_metrics", "true"),
				),
			},
			{
				Config: configDatadog(projectID, updatedAPIKey, "US", true, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "send_collection_latency_metrics", "true"),
					resource.TestCheckResourceAttr(resourceName, "send_database_metrics", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "send_collection_latency_metrics", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "send_database_metrics", "true"),
				),
			},
			importStep(resourceName),
		},
	}
}

func prometheusTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID        = acc.ProjectIDExecution(tb)
		username         = "someuser"
		updatedUsername  = "otheruser"
		password         = "somepassword"
		serviceDiscovery = "http"
		intType          = "PROMETHEUS"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configPrometheus(projectID, username, password, serviceDiscovery),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "user_name", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "service_discovery", serviceDiscovery),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "user_name", username),
					resource.TestCheckResourceAttr(dataSourceName, "service_discovery", serviceDiscovery),
				),
			},
			{
				Config: configPrometheus(projectID, updatedUsername, password, serviceDiscovery),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "user_name", updatedUsername),
				),
			},
			importStep(resourceName),
		},
	}
}

func microsoftTeamsTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID = acc.ProjectIDExecution(tb)
		intType   = "MICROSOFT_TEAMS"
		url       = "https://apps.webhook.office.com/webhookb2/" +
			"c9c5fafc-d9fe-4ffb-9773-77d804ea4372@c9656" +
			"3a8-841b-4ef9-af16-33548fffc958/IncomingWebhook" +
			"/484cccf0a678fffff86388b63203110a/42a0070b-5f35-ffff-be83-ac7e7f55d7d3"
		updatedURL = "https://apps.webhook.office.com/webhookb2/" +
			"c9c5fafc-d9fe-4ffb-9773-77d804ea4372@c9656" +
			"3a8-841b-4ef9-af16-33548fffc958/IncomingWebhook" +
			"/484cccf0a678fffff86388b63203110a/42a0070b-5f35-ffff-be83-ac7e7f55d7d4"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configMicrosoftTeams(projectID, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "microsoft_teams_webhook_url", url),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
				),
			},
			{
				Config: configMicrosoftTeams(projectID, updatedURL),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "microsoft_teams_webhook_url", updatedURL),
				),
			},
			importStep(resourceName),
		},
	}
}

func webhookTest(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID  = acc.ProjectIDExecution(tb)
		intType    = "WEBHOOK"
		url        = "https://www.mongodb.com/webhook"
		updatedURL = "https://www.mongodb.com/webhook/updated"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWebhook(projectID, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
				),
			},
			{
				Config: configWebhook(projectID, updatedURL),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "url", updatedURL),
				),
			},
			importStep(resourceName),
		},
	}
}

func importStep(resourceName string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateIdFunc: importStateIDFunc(resourceName),
		ImportState:       true,
		ImportStateVerify: false, // API Obfuscation will always make import mismatch
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_third_party_integration" {
			continue
		}
		attrs := rs.Primary.Attributes
		if attrs["project_id"] == "" {
			return fmt.Errorf("no project id is set")
		}
		if attrs["type"] == "" {
			return fmt.Errorf("no type is set")
		}
		_, _, err := acc.ConnV2().ThirdPartyIntegrationsApi.GetThirdPartyIntegration(context.Background(), attrs["project_id"], attrs["type"]).Execute()
		if err == nil {
			return fmt.Errorf("third party integration service (%s) still exists", attrs["type"])
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

		attrs := rs.Primary.Attributes
		if attrs["project_id"] == "" {
			return "", fmt.Errorf("no project id is set")
		}
		if attrs["type"] == "" {
			return "", fmt.Errorf("no type is set")
		}

		return fmt.Sprintf("%s-%s", attrs["project_id"], attrs["type"]), nil
	}
}

var singularDataStr = `
data "mongodbatlas_third_party_integration" "test" {
	project_id = mongodbatlas_third_party_integration.test.project_id
	type = mongodbatlas_third_party_integration.test.type
}

`

var pluralDataStr = `
data "mongodbatlas_third_party_integrations" "test" {
	project_id = mongodbatlas_third_party_integration.test.project_id
}`

func configOpsGenie(projectID, apiKey string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		api_key = "%[3]s"
		region  = "%[4]s"
	}`,
		projectID,
		"OPS_GENIE",
		apiKey,
		"US",
	) + singularDataStr
}

func configBasicPagerDuty(projectID, serviceKey string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		service_key = "%[3]s"
	}
	`,
		projectID,
		"PAGER_DUTY",
		serviceKey,
	) + singularDataStr + pluralDataStr
}

func configVictorOps(projectID, apiKey string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		api_key = "%[3]s"
		routing_key = "testing"
	}
	`,
		projectID,
		"VICTOR_OPS",
		apiKey,
	) + singularDataStr
}

func configDatadog(projectID, apiKey, region string, useOptionalAttr, sendCollectionLatencyMetrics, sendDatabaseMetrics bool) string {
	optionalConfigAttrs := ""
	if useOptionalAttr {
		optionalConfigAttrs = fmt.Sprintf(
			`send_collection_latency_metrics = %[1]t
		send_database_metrics = %[2]t`, sendCollectionLatencyMetrics, sendDatabaseMetrics)
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		api_key = "%[3]s"
		region  ="%[4]s"
		
		%[5]s
	}
	`,
		projectID,
		"DATADOG",
		apiKey,
		region,
		optionalConfigAttrs,
	) + singularDataStr
}

func configPrometheus(projectID, username, password, serviceDiscovery string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		user_name = "%[3]s"	
		password  = "%[4]s"
		service_discovery = "%[5]s" 
		enabled = true
	}
	`,
		projectID,
		"PROMETHEUS",
		username,
		password,
		serviceDiscovery,
	) + singularDataStr
}

func configMicrosoftTeams(projectID, url string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "MICROSOFT_TEAMS"
		microsoft_teams_webhook_url = "%[2]s"	
	}
	`, projectID, url) + singularDataStr
}

func configWebhook(projectID, url string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "WEBHOOK"
		url = "%[2]s"	
	}
	`, projectID, url) + singularDataStr
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		attrs := rs.Primary.Attributes
		if attrs["project_id"] == "" {
			return fmt.Errorf("no project id is set")
		}
		if attrs["type"] == "" {
			return fmt.Errorf("no type is set")
		}
		if _, _, err := acc.ConnV2().ThirdPartyIntegrationsApi.GetThirdPartyIntegration(context.Background(), attrs["project_id"], attrs["type"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("third party integration (%s) does not exist", attrs["project_id"])
	}
}
