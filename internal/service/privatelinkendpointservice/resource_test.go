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

func TestAccNetworkRSPrivateLinkEndpointServiceAWS_Complete(t *testing.T) {
	testCase := basicAWSTestCase(t)
	resource.Test(t, *testCase)
}

func TestAccNetworkRSPrivateLinkEndpointServiceAWS_Failed(t *testing.T) {
	var (
		resourceSuffix = "test"

		providerName = "AWS"
		projectID    = acc.ProjectIDExecution(t)
		region       = os.Getenv("AWS_REGION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configFailAWS(
					projectID, providerName, region, resourceSuffix,
				),
				ExpectError: regexp.MustCompile("privatelink endpoint service is in a failed state: Interface endpoint vpce-11111111111111111 was not found."),
			},
		},
	})
}

func TestAccNetworkRSPrivateLinkEndpointService_deleteOnCreateTimeout(t *testing.T) {
	var (
		resourceSuffix = "test"
		providerName   = "AWS"
		region         = os.Getenv("AWS_REGION")
		// Create private link endpoint outside of test configuration to avoid cleanup issues
		projectID, privateLinkEndpointID = acc.PrivateLinkEndpointIDExecution(t, providerName, region)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDeleteOnCreateTimeoutWithExistingEndpoint(
					projectID, providerName, privateLinkEndpointID, resourceSuffix, "1s", true,
				),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func basicAWSTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs AWS configuration
	var (
		resourceSuffix = "test"
		resourceName   = fmt.Sprintf("mongodbatlas_privatelink_endpoint_service.%s", resourceSuffix)
		datasourceName = fmt.Sprintf("data.mongodbatlas_privatelink_endpoint_service.%s", resourceSuffix)

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		providerName    = "AWS"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		region          = os.Getenv("AWS_REGION")
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
		checkAttrs      = []string{"project_id", "private_link_id", "endpoint_service_id"}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrSetChecks(resourceName, checks, checkAttrs...)
	checks = acc.AddAttrSetChecks(datasourceName, checks, checkAttrs...)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(tb); acc.PreCheckAwsEnvPrivateLinkEndpointService(tb) },
		CheckDestroy:             checkDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		Steps: []resource.TestStep{
			{
				Config: configCompleteAWS(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, resourceSuffix,
				),
				Check: resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_link_id"},
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

func configCompleteAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, resourceSuffix string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region        = "%[5]s"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%[6]s"
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%[7]s"]
			security_group_ids = ["%[8]s"]
			
		}

		resource "mongodbatlas_privatelink_endpoint_service" %[9]q {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id   = aws_vpc_endpoint.ptfe_service.id
			private_link_id       = mongodbatlas_privatelink_endpoint.test.id
			provider_name         = "%[4]s"
		}

		data "mongodbatlas_privatelink_endpoint_service" %[9]q {
			project_id            = %[3]q
			private_link_id       =  mongodbatlas_privatelink_endpoint_service.%[9]s.private_link_id
			endpoint_service_id = mongodbatlas_privatelink_endpoint_service.%[9]s.endpoint_service_id
			provider_name = "%[4]s"
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, resourceSuffix)
}

func configFailAWS(projectID, providerName, region, resourceSuffix string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = %[1]q
			provider_name = %[2]q
			region        = %[3]q
		}

		resource "mongodbatlas_privatelink_endpoint_service" %[4]q {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id   = "vpce-11111111111111111"
			private_link_id       = mongodbatlas_privatelink_endpoint.test.id
			provider_name         = %[2]q
		}
	`, projectID, providerName, region, resourceSuffix)
}

func configDeleteOnCreateTimeoutWithExistingEndpoint(projectID, providerName, privateLinkEndpointID, resourceSuffix, timeout string, deleteOnTimeout bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint_service" %[4]q {
			project_id            = %[1]q
			private_link_id       = %[3]q
			endpoint_service_id   = "vpce-11111111111111111"
			provider_name         = %[2]q
			delete_on_create_timeout = %[6]t
			
			timeouts {
				create = %[5]q
			}
		}
	`, projectID, providerName, privateLinkEndpointID, resourceSuffix, timeout, deleteOnTimeout)
}
