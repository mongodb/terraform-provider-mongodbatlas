package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourcePrivatelinkEndpointServiceDataFederetionDataArchives = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives.test"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchivesDSPlural_basic(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasic(projectID, endpointID),
				Check: resource.ComposeAggregateTestCheckFunc(
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
