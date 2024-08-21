package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchiveDS_basic(t *testing.T) {
	var (
		projectID               = acc.ProjectIDExecution(t)
		endpointID              = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		customerEndpointDNSName = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_DNS_NAME")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfigBasic(projectID, endpointID, customerEndpointDNSName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "comment", comment),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "region", atlasRegion),
					resource.TestCheckResourceAttr(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "customer_endpoint_dns_name", customerEndpointDNSName),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "type"),
					resource.TestCheckResourceAttrSet(dataSourcePrivatelinkEndpointServiceDataFederetionDataArchive, "provider_name"),
				),
			},
		},
	})
}

func dataSourcesConfigBasic(projectID, endpointID, customerEndpointDNSName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id					= %[1]q
	  endpoint_id					= %[2]q
	  provider_name					= "AWS"
	  comment						= "Terraform Acceptance Test"
	  region						= %[3]q
	  customer_endpoint_dns_name	= %[4]q
	}

	data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
		project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
		endpoint_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.endpoint_id
	}
	`, projectID, endpointID, atlasRegion, customerEndpointDNSName)
}
