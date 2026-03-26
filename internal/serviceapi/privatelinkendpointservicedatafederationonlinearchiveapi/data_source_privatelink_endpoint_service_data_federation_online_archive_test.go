package privatelinkendpointservicedatafederationonlinearchiveapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	AWSregion      = "US_EAST_1"
	dataSourceName = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchiveDS_basicAWS(t *testing.T) {
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
				Config: dataSourceConfigBasicAWS(projectID, endpointID, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "provider_name"),
					checkDataSourceEncodedID(dataSourceName, projectID, endpointID),
				),
			},
		},
	})
}

func checkDataSourceEncodedID(resourceName, expectedProjectID, expectedEndpointID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("id is empty")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		if ids["project_id"] != expectedProjectID || ids["endpoint_id"] != expectedEndpointID {
			return fmt.Errorf("unexpected decoded ID map: %+v", ids)
		}
		return nil
	}
}

func dataSourceConfigBasicAWS(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= %[3]q
	}

	data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test.project_id
	  endpoint_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test.endpoint_id
	}
	`, projectID, endpointID, comment)
}
