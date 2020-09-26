package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	Unknown3rdParty = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
	}
	`

	PAGERDUTY = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		service_key = "%[4]s"
	}
	`
	SLACK = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"

		api_token = "%[4]s"
		team_name = "%[5]s"
		channel_name = "%[6]s"
	}
	`

	DATADOG = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
		region  ="%[5]s"
	}
	`

	NEWRELIC = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		license_key = "%[4]s"
		account_id  = "%[5]s"
		write_token = "%[6]s"
		read_token  = "%[7]s"
	}
	`

	OPSGENIE = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
		region  = "%[5]s"
	}
	`
	VICTOROPS = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		api_key = "%[4]s"
	}
	`

	FLOWDOCK = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		flow_name = "%[4]s"
		api_token = "%[5]s"
		org_name =  "%[6]s"
	}
	`

	WEBHOOK = `
	resource "third_party_integration" "%[1]s" {
		project_id = "%[2]s"
		type = "%[3]s"
		url = "%[4]s"	
	}
	`
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numeric  = "0123456789"
	alphaNum = alphabet + numeric
)

var (
	seeder = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type thirdPartyConfig struct {
	Name        string
	ProjectID   string
	Integration matlas.ThirdPartyIntegration
}

func TestAccdataSourceMongoDBAtlasThirdPartyIntegration_basic(t *testing.T) {
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		config            = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_" + config.AccountID
		resourceName      = "data.mongodbatlas_third_party_integration." + testExecutionName

		seedConfig = thirdPartyConfig{
			Name:        testExecutionName,
			ProjectID:   orgID,
			Integration: *config,
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationDataSourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", config.Type),
					resource.TestCheckResourceAttr(resourceName, "team_name", config.TeamName),
					resource.TestCheckResourceAttr(resourceName, "channel_name", config.ChannelName),
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
	case "SLACK":
		return fmt.Sprintf(SLACK,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.APIToken,
			config.Integration.TeamName,
			config.Integration.ChannelName,
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
		)

	case "FLOWDOCK":
		return fmt.Sprintf(FLOWDOCK,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.FlowName,
			config.Integration.APIToken,
			config.Integration.OrgName,
		)
	case "WEBHOOK":
		return fmt.Sprintf(WEBHOOK,
			config.Name,
			config.ProjectID,
			config.Integration.Type,
			config.Integration.URL,
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
		Type: "SLACK",
		// Pager dutty 20-character strings
		LicenseKey: testGenString(20, alphabet),
		// Slack xoxb-333649436676-799261852869-clFJVVIaoJahpORboa3Ba2al
		APIToken:    fmt.Sprintf("xoxb-%s-%s-%s", testGenString(12, numeric), testGenString(12, numeric), testGenString(24, alphaNum)),
		TeamName:    "MongoSlackTestTeam " + account,
		ChannelName: "MongoSlackTestChannel " + account,
		// DataDog 40
		APIKey: testGenString(40, alphaNum),
		Region: "EU",

		AccountID:  account,
		WriteToken: "write-test-" + testGenString(20, alphaNum),
		ReadToken:  "read-test-" + testGenString(20, alphaNum),
		RoutingKey: testGenString(40, alphaNum),
		FlowName:   "MongoFlow test" + account,
		OrgName:    "MongoOrgTest " + account,
		URL:        "https://www.mongodb.com/webhook",
		Secret:     account,
	}
}

func testGenString(length int, charSet string) string {
	sequence := make([]byte, length)
	for i := range sequence {
		sequence[i] = charSet[seeder.Intn(len(charSet))]
	}
	return string(sequence)
}

func testAccCheckThirdPartyIntegrationExists(resourceName string, integration *matlas.ThirdPartyIntegration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

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
