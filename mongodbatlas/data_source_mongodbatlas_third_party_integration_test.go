package mongodbatlas

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"crypto/rand"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

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

	NEWRELIC = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		license_key = "%[4]s"
		account_id  = "%[5]s"
		write_token = "%[6]s"
		read_token  = "%[7]s"
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

	FLOWDOCK = `
	resource "mongodbatlas_third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		flow_name = "%[4]s"
		api_token = "%[5]s"
		org_name =  "%[6]s"
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

func TestAccConfigDSThirdPartyIntegration_basic(t *testing.T) {
	SkipTestForCI(t) // TODO: Address failures in v1.4.6

	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		config            = testAccCreateThirdPartyIntegrationConfig()

		testExecutionName = "test_" + config.AccountID
		resourceName      = "data.mongodbatlas_third_party_integration." + testExecutionName

		seedConfig = thirdPartyConfig{
			Name:        testExecutionName,
			ProjectID:   projectID,
			Integration: *config,
		}
	)
	seedConfig.Integration.APIKey = apiKey

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationDataSourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", config.Type),
				),
			},
		},
	})
}

func testAccMongoDBAtlasThirdPartyIntegrationDataSourceConfig(config *thirdPartyConfig) string {
	// create the resource config first
	resourceConfig := testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(config)

	return fmt.Sprintf(`
	%[1]s

	data "mongodbatlas_third_party_integration" "%[2]s" {
		project_id = mongodbatlas_third_party_integration.%[2]s.project_id
		type = mongodbatlas_third_party_integration.%[2]s.type
	}
	`, resourceConfig, config.Name)
}

func testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(config *thirdPartyConfig) string {
	switch config.Integration.Type {
	case "PAGER_DUTY":
		return fmt.Sprintf(PAGERDUTY,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.LicenseKey,
		)
	case "DATADOG":
		return fmt.Sprintf(DATADOG,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.APIKey,
			config.Integration.Region,
		)
	case "NEW_RELIC":
		return fmt.Sprintf(NEWRELIC,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.LicenseKey,
			config.Integration.AccountID,
			config.Integration.WriteToken,
			config.Integration.ReadToken,
		)
	case "OPS_GENIE":
		return fmt.Sprintf(OPSGENIE,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.APIKey,
			config.Integration.Region,
		)
	case "VICTOR_OPS":
		return fmt.Sprintf(VICTOROPS,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.APIKey,
			config.Integration.RoutingKey,
		)
	case "WEBHOOK":
		return fmt.Sprintf(WEBHOOK,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.URL,
		)
	case "MICROSOFT_TEAMS":
		return fmt.Sprintf(MICROSOFTTEAMS,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.MicrosoftTeamsWebhookURL,
		)
	case "PROMETHEUS":
		return fmt.Sprintf(PROMETHEUS,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.UserName,
			config.Integration.Password,
			config.Integration.ServiceDiscovery,
			config.Integration.Scheme,
		)
	default:
		return fmt.Sprintf(Unknown3rdParty,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
		)
	}
}

func testAccCreateThirdPartyIntegrationConfig() *matlas.ThirdPartyIntegration {
	account := testGenString(6, numeric)
	return &matlas.ThirdPartyIntegration{
		Type: "OPS_GENIE",
		// Pager dutty 20-character strings
		LicenseKey:  testGenString(20, alphabet),
		APIToken:    fmt.Sprintf("xoxb-%s-%s-%s", testGenString(12, numeric), testGenString(12, numeric), testGenString(24, alphaNum)),
		TeamName:    "MongoSlackTestTeam " + account,
		ChannelName: "MongoSlackTestChannel " + account,
		// DataDog 40
		APIKey: testGenString(40, alphaNum),
		Region: "EU",

		AccountID:        account,
		WriteToken:       "write-test-" + testGenString(20, alphaNum),
		ReadToken:        "read-test-" + testGenString(20, alphaNum),
		RoutingKey:       testGenString(40, alphaNum),
		FlowName:         "MongoFlow test" + account,
		OrgName:          "MongoOrgTest " + account,
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

func testAccCheckThirdPartyIntegrationExists(resourceName string, integration *matlas.ThirdPartyIntegration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		if integrationResponse, _, err := conn.Integrations.Get(context.Background(), ids["project_id"], ids["type"]); err == nil {
			*integration = *integrationResponse
			return nil
		}

		return fmt.Errorf("third party integration (%s) does not exist", ids["project_id"])
	}
}
