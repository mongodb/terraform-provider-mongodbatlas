package privatelinkendpointservicedatafederationonlinearchiveapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test"

func TestAccPrivatelinkEndpointServiceDataFederationOnlineArchiveAPI_basic(t *testing.T) {
	var (
		projectID      = acc.ProjectIDExecution(t)
		endpointID     = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		comment        = "Terraform Acceptance Test"
		commentUpdated = "Terraform Acceptance Test Updated"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, endpointID, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			{
				Config: configBasic(projectID, endpointID, commentUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", commentUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			{
				Config: configBasic(projectID, endpointID, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
			},
		},
	})
}

func configBasic(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
			project_id    = %[1]q
			endpoint_id   = %[2]q
			provider_name = "AWS"
			comment       = %[3]q
		}
	`, projectID, endpointID, comment)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		endpointID := rs.Primary.Attributes["endpoint_id"]
		if projectID == "" || endpointID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().DataFederationApi.GetPrivateEndpointId(context.Background(), projectID, endpointID).Execute()
		if err != nil {
			return fmt.Errorf("private endpoint (%s/%s) does not exist: %s", projectID, endpointID, err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		endpointID := rs.Primary.Attributes["endpoint_id"]
		if projectID == "" || endpointID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().DataFederationApi.GetPrivateEndpointId(context.Background(), projectID, endpointID).Execute()
		if err == nil {
			return fmt.Errorf("private endpoint (%s/%s) still exists", projectID, endpointID)
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
		projectID := rs.Primary.Attributes["project_id"]
		endpointID := rs.Primary.Attributes["endpoint_id"]
		if projectID == "" || endpointID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, endpointID), nil
	}
}
