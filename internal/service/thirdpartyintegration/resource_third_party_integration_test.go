package thirdpartyintegration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSThirdPartyIntegration_basic(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		cfg               = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + cfg.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	cfg.Type = "OPS_GENIE"
	cfg.APIKey = apiKey

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *cfg,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", cfg.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", cfg.Region),
				),
			},
		},
	},
	)
}

func TestAccConfigRSThirdPartyIntegration_importBasic(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		cfg               = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + cfg.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	cfg.Type = "OPS_GENIE"
	cfg.APIKey = apiKey

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *cfg,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(resourceName, "region", cfg.Region),
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
	acc.SkipTestForCI(t)
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey            = os.Getenv("OPS_GENIE_API_KEY")
		cfg               = testAccCreateThirdPartyIntegrationConfig()
		updatedConfig     = testAccCreateThirdPartyIntegrationConfig()
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
		CheckDestroy:             testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasThirdPartyIntegrationResourceConfig(&seedInitialConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThirdPartyIntegrationExists(resourceName, &targetIntegration),
					resource.TestCheckResourceAttr(resourceName, "type", cfg.Type),
					resource.TestCheckResourceAttr(resourceName, "api_key", cfg.APIKey),
					resource.TestCheckResourceAttr(resourceName, "region", cfg.Region),
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

func testAccCheckMongoDBAtlasThirdPartyIntegrationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["type"]), nil
	}
}
