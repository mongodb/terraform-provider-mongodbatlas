package streamworkspaceprivateendpoint_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_stream_workspace_private_endpoint.test"
	dataSourceName       = "data.mongodbatlas_stream_workspace_private_endpoint.test"
	pluralDataSourceName = "data.mongodbatlas_stream_workspace_private_endpoints.test"
	datasourcesConfig    = `
		data "mongodbatlas_stream_workspace_private_endpoint" "test" {
			project_id  = mongodbatlas_stream_workspace_private_endpoint.test.project_id
			endpoint_id = mongodbatlas_stream_workspace_private_endpoint.test.endpoint_id
		}

		data "mongodbatlas_stream_workspace_private_endpoints" "test" {
			project_id = mongodbatlas_stream_workspace_private_endpoint.test.project_id
		}
	`
)

// TestAccStreamWorkspacePrivateEndpoint_basic registers a real customer VPC endpoint, checks the
// singular and plural data sources, then imports. Requires a pre-created AWS interface endpoint
// whose DNS name belongs to that endpoint and region.
func TestAccStreamWorkspacePrivateEndpoint_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a real AWS VPC endpoint

	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("AWS_VPC_ENDPOINT_ID")
		dnsName    = os.Getenv("AWS_VPC_ENDPOINT_DNS_NAME")
		region     = os.Getenv("AWS_VPC_ENDPOINT_REGION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckStreamWorkspacePrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configEndpoint(projectID, endpointID, region, dnsName) + datasourcesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "customer_dns_name", dnsName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_provider", resourceName, "cloud_provider"),
					resource.TestCheckResourceAttrSet(pluralDataSourceName, "results.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// The resource has no id attribute; endpoint_id is its identifier.
				ImportStateVerifyIdentifierAttribute: "endpoint_id",
				// delete_on_create_timeout is a create-time behavior flag with no API representation.
				ImportStateVerifyIgnore: []string{"delete_on_create_timeout"},
			},
		},
	})
}

func configEndpoint(projectID, endpointID, region, dnsName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_workspace_private_endpoint" "test" {
			project_id        = %[1]q
			endpoint_id       = %[2]q
			region            = %[3]q
			customer_dns_name = %[4]q
		}
	`, projectID, endpointID, region, dnsName)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf("%s/%s", a["project_id"], a["endpoint_id"]), nil
	}
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_workspace_private_endpoint" {
				continue
			}
			a := rs.Primary.Attributes
			_, _, err := acc.ConnV2().StreamsAPI.GetWorkspacePrivateEndpoint(context.Background(), a["project_id"], a["endpoint_id"]).Execute()
			if err != nil {
				return fmt.Errorf("workspace private endpoint (%s) does not exist", a["endpoint_id"])
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_workspace_private_endpoint" {
			continue
		}
		a := rs.Primary.Attributes
		_, _, err := acc.ConnV2().StreamsAPI.GetWorkspacePrivateEndpoint(context.Background(), a["project_id"], a["endpoint_id"]).Execute()
		if err == nil {
			return fmt.Errorf("workspace private endpoint (%s) still exists", a["endpoint_id"])
		}
	}
	return nil
}
