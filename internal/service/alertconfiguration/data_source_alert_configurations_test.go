package alertconfiguration_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAlertConfigurations_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicPluralDS(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkCount(dataSourcePluralName),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckNoResourceAttr(dataSourcePluralName, "total_count"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_withOutputTypes(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		outputTypes = []string{"resource_hcl", "resource_import"}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configOutputType(projectID, outputTypes),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkCount(dataSourcePluralName),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.output.#", "2"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_invalidOutputTypeValue(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configOutputType(projectID, []string{"resource_hcl", "invalid_type"}),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
		},
	})
}

func TestAccConfigDSAlertConfigurations_totalCount(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configTotalCount(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkCount(dataSourcePluralName),
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "total_count"),
				),
			},
		},
	})
}

func configBasicPluralDS(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_alert_configurations" "test" {
			project_id = %[1]q

			list_options {
				page_num = 0
			}
		}
	`, projectID)
}

func configOutputType(projectID string, outputTypes []string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_alert_configurations" "test" {
			project_id = %[1]q
			output_type = %[2]s
		}
	`, projectID, strings.ReplaceAll(fmt.Sprintf("%+q", outputTypes), " ", ","))
}

func configTotalCount(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_alert_configurations" "test" {
			project_id = %[1]q

			list_options {
				include_count = true
			}
		}
	`, projectID)
}

func checkCount(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]

		alertResp, _, err := acc.ConnV2().AlertConfigurationsApi.ListAlertConfigurations(context.Background(), projectID).Execute()

		if err != nil {
			return fmt.Errorf("the Alert Configurations List for project (%s) could not be read", projectID)
		}

		resultsCountAttr := rs.Primary.Attributes["results.#"]
		var resultsCount int
		if resultsCount, err = strconv.Atoi(resultsCountAttr); err != nil {
			return fmt.Errorf("%s results count is somehow not a number %s", resourceName, resultsCountAttr)
		}

		if resultsCount != len(alertResp.GetResults()) {
			return fmt.Errorf("%s results count (%d) did not match that of current Alert Configurations (%d)", resourceName, resultsCount, len(alertResp.GetResults()))
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
