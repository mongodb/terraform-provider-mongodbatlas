package privatelinkendpoint_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkDSPrivateLinkEndpoint_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	resourceName := "data.mongodbatlas_privatelink_endpoint.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(projectID, providerName, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
				),
			},
		},
	})
}

func configDS(projectID, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		data "mongodbatlas_privatelink_endpoint" "test" {
			project_id      = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id = mongodbatlas_privatelink_endpoint.test.id
			provider_name = "%[2]s"
		}
	`, projectID, providerName, region)
}
