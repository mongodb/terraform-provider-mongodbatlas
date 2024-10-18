package flexcluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceType = "mongodbatlas_flex_cluster"
	resourceName = "mongodbatlas_flex_cluster.flex_cluster"
)

func TestAccFlexClusterRS_basic(t *testing.T) {
	tc := basicTestCase(t)
	resource.ParallelTest(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID = acc.ProjectIDExecution(t)
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, false),
				Check:  checksFlexCluster(),
			},
			{
				Config: configBasic(projectID, true),
				Check:  checksFlexCluster(),
			},
			{
				Config:            configBasic(projectID, true),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func configBasic(projectID string, terminationProtectionEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_flex_cluster" "flex_cluster" {
			project_id = %[1]q
			name       = "flexClusterName"
			provider_settings = {
				backing_provider_name = "AWS"
				region_name           = "US_EAST_1"
			}
			termination_protection_enabled = %[2]t
		}`, projectID, terminationProtectionEnabled)
}

func checksFlexCluster() resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists()}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == resourceType {
				projectID := rs.Primary.Attributes["project_id"]
				name := rs.Primary.Attributes["name"]
				_, _, err := acc.ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
				if err != nil {
					return fmt.Errorf("flex cluster (%s:%s) not found", projectID, name)
				}
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceType {
			projectID := rs.Primary.Attributes["project_id"]
			name := rs.Primary.Attributes["name"]
			_, _, err := acc.ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
			if err == nil {
				return fmt.Errorf("flex cluster (%s:%s) still exists", projectID, name)
			}
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

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]), nil
	}
}
