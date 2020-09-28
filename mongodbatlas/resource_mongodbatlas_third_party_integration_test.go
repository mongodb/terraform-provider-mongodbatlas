package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasThirdPartyIntegration_basic(t *testing.T) {
	var (
		targetIntegration = matlas.ThirdPartyIntegration{}
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		config            = testAccCreateThirdPartyIntegrationConfig()
		testExecutionName = "test_3rd_party_" + config.AccountID
		resourceName      = "mongodbatlas_third_party_integration." + testExecutionName
	)

	config.Type = "OPS_GENIE"

	seedConfig := thirdPartyConfig{
		Name:        testExecutionName,
		ProjectID:   projectID,
		Integration: *config,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy,
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

func testAccCheckMongoDBAtlasThirdPartyIntegrationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)
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
