package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSAlertConfigurations_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_alert_configurations.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAlertConfigurations(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationsCount(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAlertConfigurations(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		data "mongodbatlas_alert_configurations" "test" {
			project_id = mongodbatlas_project.test.id

			list_options {
				page_num = 0
			}
		}
	`, orgID, projectName)
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
