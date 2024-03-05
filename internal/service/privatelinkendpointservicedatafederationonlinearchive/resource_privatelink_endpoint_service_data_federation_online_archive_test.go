package privatelinkendpointservicedatafederationonlinearchive_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName       = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	projectID          = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID         = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
	defaultComment     = "Terraform Acceptance Test"
	defaultAtlasRegion = "US_EAST_1"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	// Skip because private endpoints are deleted daily from dev environment
	acc.SkipTestForCI(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(projectID, endpointID, defaultComment),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", defaultComment),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
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
func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_updateComment(t *testing.T) {
	// Skip because private endpoints are deleted daily from dev environment
	acc.SkipTestForCI(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(projectID, endpointID, defaultComment),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", defaultComment),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			{
				Config: resourceConfigBasic(projectID, endpointID, "Terraform Acceptance Test Updated"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform Acceptance Test Updated"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			{
				Config: resourceConfigBasic(projectID, endpointID, ""),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
				),
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basicWithRegionDnsName(t *testing.T) {
	// Skip because private endpoints are deleted daily from dev environment
	acc.SkipTestForCI(t)
	customerEndpointDNSName := asCustomerEndpointDNSName(endpointID)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicWithRegionDNSName(projectID, endpointID, defaultComment),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform Acceptance Test"),
					resource.TestCheckResourceAttr(resourceName, "region", defaultAtlasRegion),
					resource.TestCheckResourceAttr(resourceName, "customer_endpoint_dns_name", customerEndpointDNSName),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
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

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["endpoint_id"]), nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataFederationApi.GetDataFederationPrivateEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"]).Execute()
		if err == nil {
			return fmt.Errorf("Private endpoint service data federation online archive still exists")
		}
	}
	return nil
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private endpoint service data federation online archive not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Private endpoint service data federation online archive ID not set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataFederationApi.GetDataFederationPrivateEndpoint(context.Background(), ids["project_id"], ids["endpoint_id"]).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

func resourceConfigBasic(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= %[3]q
	}
	`, projectID, endpointID, comment)
}

func asCustomerEndpointDNSName(endpointID string) string {
	// Found in `.AWS.dev.us-east-1`:  https://github.com/10gen/mms/blob/85ec3df92711014b17643c05a61f5c580786556c/server/conf/data-lake-endpoint-services.json
	serviceName := "vpce-svc-0a7247db33497082e"
	return fmt.Sprintf("%s-8charsra.%s.us-east-1.vpce.amazonaws.com", endpointID, serviceName)
}

func resourceConfigBasicWithRegionDNSName(projectID, endpointID, comment string) string {
	customerEndpointDNSName := asCustomerEndpointDNSName(endpointID)
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id					= %[1]q
	  endpoint_id					= %[2]q
	  provider_name					= "AWS"
	  comment						= %[3]q
	  region						= "US_EAST_1"
	  customer_endpoint_dns_name 	= %[4]q
	}
	`, projectID, endpointID, comment, customerEndpointDNSName)
}
