package projectapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_project_api.test"

func TestAccProjectAPI_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, false),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectName, true),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectName, false),
				Check:  checkBasic(),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func configBasic(orgID, projectName string, withTags bool) string {
	tags := ""
	if withTags {
		tags = `
		tags = [
			{
				key   = "firstKey"
				value = "firstValue"
			},
			{
				key   = "secondKey"
				value = "secondValue"
			},
		]`
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project_api" "test" {
			org_id           = %q
			name             = %q
			%s
		}
	`, orgID, projectName, tags)
}

func checkBasic() resource.TestCheckFunc {
	// adds checks for computed attributes not defined in config
	setAttrsChecks := []string{"id", "created", "cluster_count"}
	checks := acc.AddAttrSetChecks(resourceName, nil, setAttrsChecks...)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["id"]
		if groupID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().ProjectsApi.GetGroup(context.Background(), groupID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("project(%s) does not exist", groupID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_api" {
			continue
		}
		groupID := rs.Primary.Attributes["id"]
		if groupID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().ProjectsApi.GetGroup(context.Background(), groupID).Execute()
		if err == nil {
			return fmt.Errorf("project (%s) still exists", groupID)
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
		groupID := rs.Primary.Attributes["id"]
		if groupID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return groupID, nil
	}
}
