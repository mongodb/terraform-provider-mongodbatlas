package privatelinkendpointservice_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName       = "mongodbatlas_privatelink_endpoint_service.this"
	datasourceName     = "data." + resourceName
	dummyVPCEndpointID = "vpce-11111111111111111"
)

func TestAccPrivateLinkEndpointService_completeAWS(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAwsEnvPrivateLinkEndpointService(t) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		Steps: []resource.TestStep{
			{
				Config: configCompleteAWS(projectID, vpcID, subnetID, securityGroupID),
				Check:  checkCompleteAWS(),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_link_id"},
			},
		},
	})
}

func TestAccPrivateLinkEndpointService_failedAWS(t *testing.T) {
	projectID := acc.ProjectIDExecution(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configFailedAWS(projectID, "EU_WEST_1"), // Different region to avoid project conflicts.
				ExpectError: regexp.MustCompile("privatelink endpoint service is in a failed state: Interface endpoint " + dummyVPCEndpointID + " was not found."),
			},
		},
	})
}

func TestAccPrivateLinkEndpointService_deleteOnCreateTimeout(t *testing.T) {
	projectID := acc.ProjectIDExecution(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				// Different region to avoid project conflicts.
				Config:      configDeleteOnCreateTimeout(projectID, "US_WEST_2"),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
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

		return fmt.Sprintf("%s--%s--%s--%s", ids["project_id"], ids["private_link_id"], ids["endpoint_service_id"], ids["provider_name"]), nil
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
		_, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["endpoint_service_id"], ids["private_link_id"]).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("the MongoDB Interface Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["endpoint_service_id"], rs.Primary.Attributes["project_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["endpoint_service_id"], ids["private_link_id"]).Execute()
		if err == nil {
			return fmt.Errorf("the MongoDB Private Endpoint(%s) still exists", ids["endpoint_service_id"])
		}
	}
	return nil
}

func configCompleteAWS(projectID, vpcID, subnetID, securityGroupID string) string {
	const region = "us-east-1" // Different region to avoid project conflicts.
	return fmt.Sprintf(`
		provider "aws" {
			region = %[5]q
		}

		resource "mongodbatlas_privatelink_endpoint" "this" {
			project_id    = %[1]q
			region        = %[5]q
			provider_name = "AWS"
		}

		data "mongodbatlas_privatelink_endpoint" "this" {
			project_id      = mongodbatlas_privatelink_endpoint.this.project_id
			provider_name   = mongodbatlas_privatelink_endpoint.this.provider_name
			private_link_id = mongodbatlas_privatelink_endpoint.this.private_link_id
		}

		resource "aws_vpc_endpoint" "this" {
			vpc_id             = %[2]q
			subnet_ids         = [%[3]q]
			security_group_ids = [%[4]q]
			service_name       = mongodbatlas_privatelink_endpoint.this.endpoint_service_name
			vpc_endpoint_type  = "Interface"
		}

		resource "mongodbatlas_privatelink_endpoint_service" "this" {
			project_id          = %[1]q
			endpoint_service_id = aws_vpc_endpoint.this.id
			private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
			provider_name       = "AWS"
		}

		data "mongodbatlas_privatelink_endpoint_service" "this" {
			project_id          = mongodbatlas_privatelink_endpoint_service.this.project_id
			private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
			endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
			provider_name       = mongodbatlas_privatelink_endpoint.this.provider_name
		}
	`, projectID, vpcID, subnetID, securityGroupID, region)
}

func checkCompleteAWS() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
		resource.TestCheckResourceAttrSet(resourceName, "endpoint_service_id"),
		resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
		resource.TestCheckResourceAttrSet(datasourceName, "private_link_id"),
		resource.TestCheckResourceAttrSet(datasourceName, "endpoint_service_id"),
	)
}

func configFailedAWS(projectID, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "this" {
			project_id    = %[1]q
			provider_name = "AWS"
			region        = %[2]q
		}

		resource "mongodbatlas_privatelink_endpoint_service" "this" {
			project_id          = mongodbatlas_privatelink_endpoint.this.project_id
			endpoint_service_id = %[3]q
			private_link_id     = mongodbatlas_privatelink_endpoint.this.id
			provider_name       = "AWS"
		}
	`, projectID, region, dummyVPCEndpointID)
}

func configDeleteOnCreateTimeout(projectID, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "this" {
			project_id    = %[1]q
			provider_name = "AWS"
			region        = %[2]q
		}

		resource "mongodbatlas_privatelink_endpoint_service" "this" {
			project_id               = mongodbatlas_privatelink_endpoint.this.project_id
			private_link_id          = mongodbatlas_privatelink_endpoint.this.private_link_id
			endpoint_service_id      = %[3]q
			provider_name            = "AWS"
			delete_on_create_timeout = true

			timeouts {
				create = "1s"
			}
		}
	`, projectID, region, dummyVPCEndpointID)
}
