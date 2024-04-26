package thirdpartyintegration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	resource.Test(t, *basicPagerDuty(t))
}

func TestAccThirdPartyIntegration_basicOpsGenie(t *testing.T) {
	resource.Test(t, *basicOpsGenie(t))
}

func TestAccThirdPartyIntegration_basicVictorOps(t *testing.T) {
	resource.Test(t, *basicVictorOps(t))
}

func TestAccThirdPartyIntegration_basicDatadog(t *testing.T) {
	resource.Test(t, *basicDatadog(t))
}

func TestAccThirdPartyIntegration_basicPrometheus(t *testing.T) {
	resource.Test(t, *basicPrometheus(t))
}

func TestAccThirdPartyIntegration_basicMicrosoftTeams(t *testing.T) {
	resource.Test(t, *basicMicrosoftTeams(t))
}

func TestAccThirdPartyIntegration_basicWebhook(t *testing.T) {
	resource.Test(t, *basicWebhook(t))
}

func basicOpsGenie(tb testing.TB) *resource.TestCase {
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
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttr(resourceName, "region", "US"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "region", "US"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configOpsGenie(projectID, updatedAPIKey),
				Check: resource.ComposeTestCheckFunc(
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

func basicPagerDuty(tb testing.TB) *resource.TestCase {
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
				Config: configPagerDuty(projectID, serviceKey),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "service_key", serviceKey),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configPagerDuty(projectID, updatedServiceKey),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "service_key", updatedServiceKey),
				),
			},
			importStep(resourceName),
		},
	}
}

func basicVictorOps(tb testing.TB) *resource.TestCase {
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
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttrSet(resourceName, "routing_key"),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttrSet(dataSourceName, "routing_key"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configVictorOps(projectID, updatedAPIKey),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
				),
			},
			importStep(resourceName),
		},
	}
}

func basicDatadog(tb testing.TB) *resource.TestCase {
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
				Config: configDatadog(projectID, apiKey, "US"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", apiKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "region", region),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configDatadog(projectID, updatedAPIKey, "US"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedAPIKey),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
			importStep(resourceName),
		},
	}
}

func basicPrometheus(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID        = acc.ProjectIDExecution(tb)
		username         = "someuser"
		updatedUsername  = "otheruser"
		password         = "somepassword"
		serviceDiscovery = "http"
		scheme           = "https"
		intType          = "PROMETHEUS"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configPrometheus(projectID, username, password, serviceDiscovery, scheme),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "user_name", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "service_discovery", serviceDiscovery),
					resource.TestCheckResourceAttr(resourceName, "scheme", scheme),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourceName, "user_name", username),
					resource.TestCheckResourceAttr(dataSourceName, "service_discovery", serviceDiscovery),
					resource.TestCheckResourceAttr(dataSourceName, "scheme", scheme),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configPrometheus(projectID, updatedUsername, password, serviceDiscovery, scheme),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "user_name", updatedUsername),
				),
			},
			importStep(resourceName),
		},
	}
}

func basicMicrosoftTeams(tb testing.TB) *resource.TestCase {
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
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "microsoft_teams_webhook_url", url),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configMicrosoftTeams(projectID, updatedURL),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "microsoft_teams_webhook_url", updatedURL),
				),
			},
			importStep(resourceName),
		},
	}
}

func basicWebhook(tb testing.TB) *resource.TestCase {
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
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", intType),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(dataSourceName, "type", intType),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				Config: configWebhook(projectID, updatedURL),
				Check: resource.ComposeTestCheckFunc(
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
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().Integrations.Get(context.Background(), ids["project_id"], ids["type"])
		if err == nil {
			return fmt.Errorf("third party integration service (%s) still exists", ids["type"])
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

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["type"]), nil
	}
}

var dataStr = `
data "mongodbatlas_third_party_integration" "test" {
	project_id = mongodbatlas_third_party_integration.test.project_id
	type = mongodbatlas_third_party_integration.test.type
}

data "mongodbatlas_third_party_integrations" "test" {
	project_id = mongodbatlas_third_party_integration.test.project_id
}
`

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
	) + dataStr
}

func configPagerDuty(projectID, serviceKey string) string {
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
	) + dataStr
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
	) + dataStr
}

func configDatadog(projectID, apiKey, region string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		api_key = "%[3]s"
		region  ="%[4]s"
	}
	`,
		projectID,
		"DATADOG",
		apiKey,
		region,
	) + dataStr
}

func configPrometheus(projectID, username, password, serviceDiscovery, scheme string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "%[2]s"
		user_name = "%[3]s"	
		password  = "%[4]s"
		service_discovery = "%[5]s" 
		scheme = "%[6]s"
		enabled = true
	}
	`,
		projectID,
		"PROMETHEUS",
		username,
		password,
		serviceDiscovery,
		scheme,
	) + dataStr
}

func configMicrosoftTeams(projectID, url string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "MICROSOFT_TEAMS"
		microsoft_teams_webhook_url = "%[2]s"	
	}
	`, projectID, url) + dataStr
}

func configWebhook(projectID, url string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_third_party_integration" "test" {
		project_id = "%[1]s"
		type = "WEBHOOK"
		url = "%[2]s"	
	}
	`, projectID, url) + dataStr
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if _, _, err := acc.Conn().Integrations.Get(context.Background(), ids["project_id"], ids["type"]); err == nil {
			return nil
		}
		return fmt.Errorf("third party integration (%s) does not exist", ids["project_id"])
	}
}
