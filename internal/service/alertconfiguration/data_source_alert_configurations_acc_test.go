package alertconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc/todoacc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSAlertConfigurations_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_alert_configurations.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkCount(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckNoResourceAttr(dataSourceName, "total_count"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_withOutputTypes(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_alert_configurations.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		outputTypes    = []string{"resource_hcl", "resource_import"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOutputType(orgID, projectName, outputTypes),
				Check: resource.ComposeTestCheckFunc(
					checkCount(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.output.#", "2"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_invalidOutputTypeValue(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configOutputType(orgID, projectName, []string{"resource_hcl", "invalid_type"}),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_totalCount(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_alert_configurations.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configTotalCount(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkCount(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
				),
			},
		},
	})
}

func configBasic(orgID, projectName string) string {
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

func configOutputType(orgID, projectName string, outputTypes []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		data "mongodbatlas_alert_configurations" "test" {
			project_id = mongodbatlas_project.test.id
			output_type = %[3]s
		}
	`, orgID, projectName, strings.ReplaceAll(fmt.Sprintf("%+q", outputTypes), " ", ","))
}

func configTotalCount(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		data "mongodbatlas_alert_configurations" "test" {
			project_id = mongodbatlas_project.test.id

			list_options {
				include_count = true
			}
		}
	`, orgID, projectName)
}

func checkCount(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := todoacc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := config.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]

		alertResp, _, err := conn.AlertConfigurations.List(context.Background(), projectID, &matlas.ListOptions{
			PageNum:      0,
			ItemsPerPage: 100,
			IncludeCount: true,
		})

		if err != nil {
			return fmt.Errorf("the Alert Configurations List for project (%s) could not be read", projectID)
		}

		resultsCountAttr := rs.Primary.Attributes["results.#"]
		var resultsCount int
		if resultsCount, err = strconv.Atoi(resultsCountAttr); err != nil {
			return fmt.Errorf("%s results count is somehow not a number %s", resourceName, resultsCountAttr)
		}

		if resultsCount != len(alertResp) {
			return fmt.Errorf("%s results count (%d) did not match that of current Alert Configurations (%d)", resourceName, resultsCount, len(alertResp))
		}

		if totalCountAttr := rs.Primary.Attributes["total_count"]; totalCountAttr != "" {
			var totalCount int
			if totalCount, err = strconv.Atoi(totalCountAttr); err != nil {
				return fmt.Errorf("%s total count is somehow not a number %s", resourceName, totalCountAttr)
			}
			if totalCount != resultsCount {
				return fmt.Errorf("%s total count (%d) did not match that of results count (%d)", resourceName, totalCount, resultsCount)
			}
		}

		return nil
	}
}
