package thirdpartyintegration_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccThirdPartyIntegration_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs Opsgenie config

	var (
		projectID            = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey               = os.Getenv("OPS_GENIE_API_KEY")
		cfg                  = createIntegrationConfig()
		testExecutionName    = "test_3rd_party_" + cfg.AccountID
		resourceName         = "mongodbatlas_third_party_integration." + testExecutionName
		dataSourceName       = "data." + resourceName
		dataSourcePluralName = "data.mongodbatlas_third_party_integrations.test"
	)

	cfg.Type = "OPS_GENIE"
	cfg.APIKey = apiKey

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *cfg,
	}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", cfg.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", cfg.Region),

					resource.TestCheckResourceAttr(dataSourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(dataSourceName, "api_key", cfg.APIKey),
					resource.TestCheckResourceAttr(dataSourceName, "region", cfg.Region),

					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false, // API Obfuscation will always make import mismatch
			},
		},
	}
}

func TestAccThirdPartyIntegration_updateBasic(t *testing.T) {
	acc.SkipTestForCI(t) // needs Opsgenie config

	var (
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		cfg               = createIntegrationConfig()
		updatedConfig     = createIntegrationConfig()
		testExecutionName = "test_3rd_party_" + cfg.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	// setting type
	cfg.Type = "OPS_GENIE"
	updatedConfig.Type = "OPS_GENIE"
	updatedConfig.Region = "US"
	cfg.APIKey = apiKey
	updatedConfig.APIKey = apiKey

	seedInitialConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *cfg,
	}

	seedUpdatedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *updatedConfig,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&seedInitialConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", cfg.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", cfg.Region),
				),
			},
			{
				Config: configBasic(&seedUpdatedConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", updatedConfig.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedConfig.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", updatedConfig.Region),
				),
			},
		},
	},
	)
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

const (
	Unknown3rdParty = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
	}
	`

	PAGERDUTY = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		service_key = "%[4]s"
	}
	`

	DATADOG = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
		region  ="%[5]s"
	}
	`

	OPSGENIE = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
		region  = "%[5]s"
	}
	`
	VICTOROPS = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
		routing_key = "%[5]s"
	}
	`

	MICROSOFTTEAMS = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		microsoft_teams_webhook_url = "%[4]s"	
	}
	`

	PROMETHEUS = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		user_name = "%[4]s"	
		password  = "%[5]s"
		service_discovery = "%[6]s" 
		scheme = "%[7]s"
		enabled = true
	}
	`

	WEBHOOK = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		url = "%[4]s"	
	}
`
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numeric  = "0123456789"
	alphaNum = alphabet + numeric
)

type thirdPartyConfig struct {
	Name        string
	ProjectID   string
	Integration matlas.ThirdPartyIntegration
}

func configBasic(cfg *thirdPartyConfig) string {
	dataStr := fmt.Sprintf(`
		data "mongodbatlas_third_party_integration" %[1]q {
			project_id = mongodbatlas_third_party_integration.%[1]s.project_id
			type = mongodbatlas_third_party_integration.%[1]s.type
		}

		data "mongodbatlas_third_party_integrations" "test" {
			project_id = mongodbatlas_third_party_integration.%[1]s.project_id
		}
	`, cfg.Name)

	switch cfg.Integration.Type {
	case "PAGER_DUTY":
		return fmt.Sprintf(PAGERDUTY,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.ServiceKey,
		) + dataStr
	case "DATADOG":
		return fmt.Sprintf(DATADOG,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.APIKey,
			cfg.Integration.Region,
		) + dataStr
	case "OPS_GENIE":
		return fmt.Sprintf(OPSGENIE,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.APIKey,
			cfg.Integration.Region,
		) + dataStr
	case "VICTOR_OPS":
		return fmt.Sprintf(VICTOROPS,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.APIKey,
			cfg.Integration.RoutingKey,
		) + dataStr
	case "WEBHOOK":
		return fmt.Sprintf(WEBHOOK,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.URL,
		) + dataStr
	case "MICROSOFT_TEAMS":
		return fmt.Sprintf(MICROSOFTTEAMS,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.MicrosoftTeamsWebhookURL,
		) + dataStr
	case "PROMETHEUS":
		return fmt.Sprintf(PROMETHEUS,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
			cfg.Integration.UserName,
			cfg.Integration.Password,
			cfg.Integration.ServiceDiscovery,
			cfg.Integration.Scheme,
		) + dataStr
	default:
		return fmt.Sprintf(Unknown3rdParty,
			cfg.Name,
			cfg.ProjectID,
			cfg.Integration.Type,
		) + dataStr
	}
}

func createIntegrationConfig() *matlas.ThirdPartyIntegration {
	account := testGenString(6, numeric)
	return &matlas.ThirdPartyIntegration{
		Type:        "OPS_GENIE",
		TeamName:    "MongoSlackTestTeam " + account,
		ChannelName: "MongoSlackTestChannel " + account,
		// DataDog 40
		APIKey:           testGenString(40, alphaNum),
		Region:           "EU",
		ReadToken:        "read-test-" + testGenString(20, alphaNum),
		RoutingKey:       testGenString(40, alphaNum),
		URL:              "https://www.mongodb.com/webhook",
		Secret:           account,
		UserName:         "PROM_USER",
		Password:         "PROM_PASSWORD",
		ServiceDiscovery: "http",
		Scheme:           "https",
		Enabled:          false,
		MicrosoftTeamsWebhookURL: "https://apps.webhook.office.com/webhookb2/" +
			"c9c5fafc-d9fe-4ffb-9773-77d804ea4372@c9656" +
			"3a8-841b-4ef9-af16-33548fffc958/IncomingWebhook" +
			"/484cccf0a678fffff86388b63203110a/42a0070b-5f35-ffff-be83-ac7e7f55d7d3",
	}
}

func testGenString(length int, charSet string) string {
	sequence := make([]byte, length)
	upperBound := big.NewInt(int64(len(charSet)))
	for i := range sequence {
		n, _ := rand.Int(rand.Reader, upperBound)
		sequence[i] = charSet[int(n.Int64())]
	}
	return string(sequence)
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
