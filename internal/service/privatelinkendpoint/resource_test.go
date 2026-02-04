package privatelinkendpoint_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccPrivateLinkEndpoint_basicAWS(t *testing.T) {
	resource.ParallelTest(t, *basicAWSTestCase(t, "us-east-1"))
}

func basicAWSTestCase(tb testing.TB, region string) *resource.TestCase {
	tb.Helper()
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		projectID    = acc.ProjectIDExecution(tb)
		providerName = constant.AWS
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, providerName, region, nil),
				Check:  checkBasic(resourceName, providerName, region, nil),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccPrivateLinkEndpoint_basicAzure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		projectID    = acc.ProjectIDExecution(t)
		region       = "US_EAST_2"
		providerName = constant.AZURE
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, providerName, region, nil),
				Check:  checkBasic(resourceName, providerName, region, nil),
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

func TestAccPrivateLinkEndpoint_basicGCP(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		projectID    = acc.ProjectIDExecution(t)
		region       = "us-central1"
		providerName = constant.GCP
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, providerName, region, nil),
				Check:  checkBasic(resourceName, providerName, region, conversion.Pointer(false)),
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

func TestAccPrivateLinkEndpoint_deleteOnCreateTimeout(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		region       = "eu-west-1"
		providerName = constant.AWS
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configDeleteOnCreateTimeout(projectID, providerName, region, "1s", true),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func TestAccPrivateLinkEndpoint_gcpPortMappingEnabled(t *testing.T) {
	resource.ParallelTest(t, *basicGCPTestCaseWithPortMapping(t, true))
}

func TestAccPrivateLinkEndpoint_gcpPortMappingDisabled(t *testing.T) {
	resource.ParallelTest(t, *basicGCPTestCaseWithPortMapping(t, false))
}

func basicGCPTestCaseWithPortMapping(tb testing.TB, portMappingEnabled bool) *resource.TestCase {
	tb.Helper()
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		projectID    = acc.ProjectIDExecution(tb)
		providerName = constant.GCP
		region       = "us-west3"
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, providerName, region, conversion.Pointer(portMappingEnabled)),
				Check:  checkBasic(resourceName, providerName, region, conversion.Pointer(portMappingEnabled)),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s-%s", ids["project_id"], ids["private_link_id"], ids["provider_name"], ids["region"]), nil
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if _, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpointService(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("the MongoDB Private Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["private_link_id"], rs.Primary.Attributes["project_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpointService(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"]).Execute()
		if err == nil {
			return fmt.Errorf("the MongoDB Private Endpoint(%s) still exists", ids["private_link_id"])
		}
	}
	return nil
}

func configDeleteOnCreateTimeout(projectID, providerName, region, timeout string, deleteOnTimeout bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = %[1]q
			provider_name = %[2]q
			region        = %[3]q
			delete_on_create_timeout = %[5]t
			
			timeouts {
				create = %[4]q
			}
		}
	`, projectID, providerName, region, timeout, deleteOnTimeout)
}

// configBasic is a helper function to create a basic configuration for a private link endpoint.
// IMPORTANT: Use a different region in each test to avoid project conflicts. Legacy and port-mapped GCP can use the same region.
func configBasic(projectID, providerName, region string, portMappingEnabled *bool) string {
	portMappingEnabledStr := ""
	if portMappingEnabled != nil {
		portMappingEnabledStr = fmt.Sprintf("port_mapping_enabled = %t", *portMappingEnabled)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = %[1]q
			provider_name = %[2]q
			region        = %[3]q
			%[4]s
		}
	`, projectID, providerName, region, portMappingEnabledStr)
}

func checkBasic(resourceName, providerName, region string, portMappingEnabled *bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
		resource.TestCheckResourceAttrSet(resourceName, "region"),
		resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
		resource.TestCheckResourceAttr(resourceName, "region", region),
	}
	if portMappingEnabled != nil {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "port_mapping_enabled", strconv.FormatBool(*portMappingEnabled)))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
