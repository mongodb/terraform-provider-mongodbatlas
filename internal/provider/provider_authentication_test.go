package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccSTSAssumeRole_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = "test-acc-tf-p-temp-delete"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSTSAssumeRole(t); acc.PreCheckRegularCredsAreEmpty(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_count"),
					resource.TestCheckResourceAttrSet(resourceName, "teams.#"),
				),
			},
		},
	})
}

func configBasic(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name  			 = %[2]q
		}
`, orgID, projectName)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if resp, _, err := acc.ConnV2().ProjectsApi.GetProjectByName(context.Background(), rs.Primary.Attributes["name"]).Execute(); err == nil {
			fmt.Printf("DEBUG HI: %s, org: %s\n", resp.GetId(), resp.GetOrgId())
			return nil
		}
		return fmt.Errorf("project (%s) does not exist", rs.Primary.ID)
	}
}
