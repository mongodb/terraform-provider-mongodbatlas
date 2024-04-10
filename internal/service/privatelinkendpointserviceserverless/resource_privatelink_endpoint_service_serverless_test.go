package privatelinkendpointserviceserverless_test

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

func TestAccServerlessPrivateLinkEndpointService_basic(t *testing.T) {
	var (
		resourceName            = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceName          = "data.mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceEndpointsName = "data.mongodbatlas_privatelink_endpoints_service_serverless.test"
		orgID                   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName             = acc.RandomProjectName()
		instanceName            = acc.RandomClusterName()
		commentOrigin           = "this is a comment for serverless private link endpoint"
		commentUpdated          = "this is updated comment for serverless private link endpoint"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, instanceName, commentOrigin),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttr(datasourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "instance_name"),
				),
			},
			{
				Config: configBasic(orgID, projectName, instanceName, commentUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentUpdated),
					resource.TestCheckResourceAttr(datasourceName, "comment", commentUpdated),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "instance_name"),
				),
			},
			{
				Config:            configBasic(orgID, projectName, instanceName, commentOrigin),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServerlessPrivateLinkEndpointService_AWSEndpointCommentUpdate(t *testing.T) {
	var (
		resourceName            = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceEndpointsName = "data.mongodbatlas_privatelink_endpoints_service_serverless.test"
		orgID                   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName             = acc.RandomProjectName()
		instanceName            = acc.RandomClusterName()
		commentUpdated          = "this is updated comment for serverless private link endpoint"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAWSEndpoint(orgID, projectName, instanceName, false, ""),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "instance_name"),
				),
			},
			{
				Config: configAWSEndpoint(orgID, projectName, instanceName, true, commentUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentUpdated),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceEndpointsName, "instance_name"),
				),
			},
			{
				Config:            configAWSEndpoint(orgID, projectName, instanceName, true, commentUpdated),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_serverless" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		privateLink, _, err := acc.ConnV2().ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"]).Execute()
		if err == nil && privateLink != nil {
			return fmt.Errorf("endpoint_id (%s) still exists", ids["endpoint_id"])
		}
	}
	return nil
}

func configBasic(orgID, projectName, instanceName, comment string) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id   = mongodbatlas_project.test.id
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	}


	resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_serverless.test.project_id
		instance_name = mongodbatlas_privatelink_endpoint_serverless.test.instance_name
		endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
		provider_name = "AWS"
		comment = %[4]q
	}

	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q
		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true

		lifecycle {
			ignore_changes = [connection_strings_private_endpoint_srv]
		}
	}

	data "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		name         = mongodbatlas_serverless_instance.test.name
	}

	data "mongodbatlas_privatelink_endpoints_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
	  }

	data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		endpoint_id = mongodbatlas_privatelink_endpoint_service_serverless.test.endpoint_id
	}

	`, orgID, projectName, instanceName, comment)
}

func configAWSVPCEndpoint() string {
	return `

	# Create Primary VPC
resource "aws_vpc" "primary" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

# Create IGW
resource "aws_internet_gateway" "primary" {
  vpc_id = aws_vpc.primary.id
}

# Route Table
resource "aws_route" "primary-internet_access" {
  route_table_id         = aws_vpc.primary.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.primary.id
}

# Subnet-A
resource "aws_subnet" "primary-az1" {
  vpc_id                  = aws_vpc.primary.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1a"
}

# Subnet-B
resource "aws_subnet" "primary-az2" {
  vpc_id                  = aws_vpc.primary.id
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = false
  availability_zone       = "us-east-1b"
}

resource "aws_security_group" "primary_default" {
  name_prefix = "default-"
  description = "Default security group for all instances in ${aws_vpc.primary.id}"
  vpc_id      = aws_vpc.primary.id
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "tcp"
    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}`
}

func configAWSEndpoint(orgID, projectName, instanceName string, updateComment bool, comment string) string {
	peServiceServerless := `resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
	project_id                 = mongodbatlas_privatelink_endpoint_serverless.test.project_id
	instance_name              = mongodbatlas_serverless_instance.test.name
	endpoint_id                = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
	cloud_provider_endpoint_id = aws_vpc_endpoint.test.id
	provider_name              = "AWS"
  }`
	if updateComment {
		peServiceServerless = fmt.Sprintf(`resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
			project_id                 = mongodbatlas_privatelink_endpoint_serverless.test.project_id
			instance_name              = mongodbatlas_serverless_instance.test.name
			endpoint_id                = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
			cloud_provider_endpoint_id = aws_vpc_endpoint.test.id
			provider_name              = "AWS"
			comment                    = %[1]q
		  }`, comment)
	}

	return fmt.Sprintf(`

	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}

	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q
		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name         = "SERVERLESS"
		provider_settings_region_name           = "US_EAST_1"
		continuous_backup_enabled               = true
	  }

	  resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id   = mongodbatlas_project.test.id
		provider_name = "AWS"
		instance_name = mongodbatlas_serverless_instance.test.name
	  }

	  # aws_vpc config
	  %[4]s
	  
	  resource "aws_vpc_endpoint" "test" {
		vpc_id             = aws_vpc.primary.id
		service_name       = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_service_name
		vpc_endpoint_type  = "Interface"
		subnet_ids         = [aws_subnet.primary-az1.id, aws_subnet.primary-az2.id]
		security_group_ids = [aws_security_group.primary_default.id]
	  }
	  
	  %[5]s

	  data "mongodbatlas_privatelink_endpoints_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
	  }

	data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = mongodbatlas_privatelink_endpoint_service_serverless.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		endpoint_id = mongodbatlas_privatelink_endpoint_service_serverless.test.endpoint_id
	}

	`, orgID, projectName, instanceName, configAWSVPCEndpoint(), peServiceServerless)
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
		_, _, err := acc.ConnV2().ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("endpoint_id (%s) does not exist", ids["endpoint_id"])
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["instance_name"], ids["endpoint_id"]), nil
	}
}
