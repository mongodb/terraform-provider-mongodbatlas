package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSAlertConfigurations_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_alert_configurations.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAlertConfigurations(projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationsCount(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAlertConfigurations(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_alert_configurations" "test" {
			project_id = "%s"

			list_options {
				page_num = 0
			}
		}
	`, projectID)
}

func testAccCheckMongoDBAtlasAlertConfigurationsCount(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]

		alertResp, _, err := conn.AlertConfigurations.List(context.Background(), projectID, &matlas.ListOptions{
			PageNum:      0,
			ItemsPerPage: 100,
			IncludeCount: true,
		})

		if err != nil {
			return fmt.Errorf("the Alert Configurations List for project (%s) could not be read", projectID)
		}

		resultsNumber := rs.Primary.Attributes["results.#"]
		var dataSourceResultsCount int

		if dataSourceResultsCount, err = strconv.Atoi(resultsNumber); err != nil {
			return fmt.Errorf("%s results count is somehow not a number %s", resourceName, resultsNumber)
		}

		apiResultsCount := len(alertResp)
		if dataSourceResultsCount != len(alertResp) {
			return fmt.Errorf("%s results count (%v) did not match that of current Alert Configurations (%d)", resourceName, dataSourceResultsCount, apiResultsCount)
		}

		return nil
	}
}
