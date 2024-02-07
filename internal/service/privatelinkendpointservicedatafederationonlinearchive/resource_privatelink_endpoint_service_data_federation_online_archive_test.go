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
	resourceName = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	endpointID   = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(projectID, endpointID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
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

func resourceConfigBasic(projectID, endpointID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= "Terraform Acceptance Test"
	}
	`, projectID, endpointID)
}
