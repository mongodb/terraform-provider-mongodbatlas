package customdnsconfigurationclusteraws_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const dataSourceName = "data.mongodbatlas_custom_dns_configuration_cluster_aws.test"

func TestAccConfigDSCustomDNSConfigurationAWS_basic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acc.RandomProjectName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(orgID, projectName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
				),
			},
		},
	})
}

func configDS(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id     = mongodbatlas_project.test.id
			enabled       = %[3]t
		}

		data "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
			project_id      = mongodbatlas_custom_dns_configuration_cluster_aws.test.id
		}
	`, orgID, projectName, enabled)
}
