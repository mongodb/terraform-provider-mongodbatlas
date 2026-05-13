package projectipaccesslist_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccQueryProjectIPAccessList_basic(t *testing.T) {
	projectID := acc.ProjectIDExecution(t)
	ipAddress := acc.RandomIP(179, 154, 226)
	listAddress := "mongodbatlas_project_ip_access_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: configIPEntry(projectID, ipAddress),
			},
			{
				Query:  true,
				Config: queryConfig(projectID),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listAddress, 1),
					querycheck.ExpectIdentity(listAddress, map[string]knownvalue.Check{
						"project_id": knownvalue.StringExact(projectID),
						"entry":      knownvalue.StringExact(ipAddress),
					}),
				},
			},
		},
	})
}

func configIPEntry(projectID, ipAddress string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project_ip_access_list" "test" {
  project_id = %q
  ip_address = %q
  comment    = "test entry for query acceptance test"
}
`, projectID, ipAddress)
}

func queryConfig(projectID string) string {
	return fmt.Sprintf(`
provider "mongodbatlas" {}

list "mongodbatlas_project_ip_access_list" "test" {
  provider = mongodbatlas

  config {
    project_id = %q
  }
}
`, projectID)
}
