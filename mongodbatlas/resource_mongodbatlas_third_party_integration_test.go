package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSThirdPartyIntegration_basic(t *testing.T) {
	SkipTest(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		config            = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + config.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	config.Type = "OPS_GENIE"
	config.APIKey = apiKey

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *config,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", config.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", config.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", config.Region),
				),
			},
		},
	},
	)
}

func TestAccConfigRSThirdPartyIntegration_importBasic(t *testing.T) {
	SkipTest(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		config            = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + config.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	config.Type = "OPS_GENIE"
	config.APIKey = apiKey

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *config,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", config.Type),
					resource.TestCheckResourceAttr(resourceName, "region", config.Region),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasThirdPartyIntegrationImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false, // API Obfuscation will always make import mismatch
			},
		},
	},
	)
}

func TestAccConfigRSThirdPartyIntegration_updateBasic(t *testing.T) {
	SkipTest(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		config            = testAccCreateThirdPartyIntegrationConfig()
		updatedConfig     = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + config.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	// setting type
	config.Type = "OPS_GENIE"
	updatedConfig.Type = "OPS_GENIE"
	updatedConfig.Region = "US"
	config.APIKey = apiKey
	updatedConfig.APIKey = apiKey

	seedInitialConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *config,
	}

	seedUpdatedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *updatedConfig,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedInitialConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", config.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", config.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", config.Region),
				),
			},
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedUpdatedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", updatedConfig.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", updatedConfig.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", updatedConfig.Region),
				),
			},
		},
	},
	)
}

func testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_third_party_integration" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.Integrations.Get(context.Background(), ids["project_id"], ids["type"])

		if err == nil {
			return fmt.Errorf("third party integration service (%s) still exists", ids["type"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasThirdPartyIntegrationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["type"]), nil
	}
}
