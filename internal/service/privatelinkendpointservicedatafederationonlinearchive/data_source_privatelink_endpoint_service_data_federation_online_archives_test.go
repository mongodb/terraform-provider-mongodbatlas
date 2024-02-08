package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives.test"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchivesDSPlural_basic(t *testing.T) {
	// Skip because private endpoints are deleted daily from dev environment
	acc.SkipTestForCI(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasic(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives, "project_id", projectID),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives, "results.#"),
				),
			},
		},
	})
}

func dataSourceConfigBasic(projectID, endpointID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= "Terraform Acceptance Test"
	}

	data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "test" {
		project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
	}
	`, projectID, endpointID)
}
