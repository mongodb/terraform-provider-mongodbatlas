package privateendpointregionalmode_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccPrivateEndpointRegionalMode_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccPrivateEndpointRegionalMode_conn(t *testing.T) {
	var (
		endpointResourceSuffix                 = "atlasple"
		resourceSuffix                         = "atlasrm"
		resourceName                           = fmt.Sprintf("mongodbatlas_private_endpoint_regional_mode.%s", resourceSuffix)
		awsAccessKey                           = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey                           = os.Getenv("AWS_SECRET_ACCESS_KEY")
		providerName                           = "AWS"
		region                                 = os.Getenv("AWS_REGION_LOWERCASE")
		privatelinkEndpointServiceResourceName = fmt.Sprintf("mongodbatlas_privatelink_endpoint_service.%s", endpointResourceSuffix)
		clusterDependsOn                       = fmt.Sprintf("%s, %s", resourceName, privatelinkEndpointServiceResourceName)
		spec1                                  = acc.ReplicationSpec(&acc.ReplicationSpecRequest{Region: os.Getenv("AWS_REGION_UPPERCASE"), ProviderName: providerName, ZoneName: "Zone 1"})
		spec2                                  = acc.ReplicationSpec(&acc.ReplicationSpecRequest{Region: "US_WEST_2", ProviderName: providerName, ZoneName: "Zone 2"})
		clusterInfo                            = acc.GetClusterInfo(t, &acc.ClusterRequest{Geosharded: true, DiskSizeGb: 80, ResourceDependencyName: clusterDependsOn}, spec1, spec2)
		clusterName                            = clusterInfo.ClusterName
		projectID                              = clusterInfo.ProjectID
		clusterResourceName                    = clusterInfo.ClusterNameStr
		endpointResources                      = testConfigUnmanagedAWS(
			awsAccessKey, awsSecretKey, projectID, providerName, region, endpointResourceSuffix,
		)
		dependencies = []string{clusterInfo.ClusterTerraformStr, endpointResources}
	)

	configBefore := configWithDependencies(resourceSuffix, projectID, false, dependencies)
	configAfter := configWithDependencies(resourceSuffix, projectID, true, dependencies)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAwsEnv(t); acc.PreCheckAwsRegionCases(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBefore,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkModeClustersUpToDate(projectID, clusterName, clusterResourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: configAfter,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkModeClustersUpToDate(projectID, clusterName, clusterResourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		resourceName   = "mongodbatlas_private_endpoint_regional_mode.test"
		dataSourceName = "data.mongodbatlas_private_endpoint_regional_mode.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),

					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "false"),
				),
			},
			{
				Config: configBasic(orgID, projectName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),

					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
				),
			},
		},
	}
}

func configWithDependencies(resourceName, projectID string, enabled bool, dependencies []string) string {
	resources := make([]string, len(dependencies)+1)

	resources[0] = configNoProject(resourceName, projectID, enabled)
	copy(resources[1:], dependencies)

	return strings.Join(resources, "\n\n")
}

func configNoProject(resourceName, projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" %[1]q {
			project_id   = %[2]q
			enabled      = %[3]t
		}
	`, resourceName, projectID, enabled)
}

func configBasic(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id   = mongodbatlas_project.project.id
			enabled      = %[3]t
		}

		data "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id = mongodbatlas_project.project.id
			depends_on = [ mongodbatlas_private_endpoint_regional_mode.test ]
		}
	`, orgID, projectName, enabled)
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
		projectID := rs.Primary.ID
		_, _, err := acc.Conn().PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), projectID)
		if err == nil {
			return nil
		}
		return fmt.Errorf("regional mode for project_id (%s) does not exist", projectID)
	}
}

func checkModeClustersUpToDate(projectID, clusterName, clusterResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterResourceName]
		if !ok {
			return fmt.Errorf("Could not find resource state for cluster (%s) on project (%s)", clusterName, projectID)
		}
		var rsPrivateEndpointCount int
		var err error
		if rsPrivateEndpointCount, err = strconv.Atoi(rs.Primary.Attributes["connection_strings.0.private_endpoint.#"]); err != nil {
			return fmt.Errorf("Connection strings private endpoint count is not a number")
		}
		c, _, _ := acc.Conn().Clusters.Get(context.Background(), projectID, clusterName)
		if rsPrivateEndpointCount != len(c.ConnectionStrings.PrivateEndpoint) {
			return fmt.Errorf("Cluster PrivateEndpoint count does not match resource")
		}
		if rs.Primary.Attributes["connection_strings.0.standard"] != c.ConnectionStrings.Standard {
			return fmt.Errorf("Cluster standard connection_string does not match resource")
		}
		if rs.Primary.Attributes["connection_strings.0.standard_srv"] != c.ConnectionStrings.StandardSrv {
			return fmt.Errorf("Cluster standard connection_string does not match resource")
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_private_endpoint_regional_mode" {
			continue
		}
		setting, _, _ := acc.Conn().PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), rs.Primary.ID)
		if setting != nil && setting.Enabled != false {
			return fmt.Errorf("Regionalized private endpoint setting for project %q was not properly disabled", rs.Primary.ID)
		}
	}
	return nil
}

func testConfigUnmanagedAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, serviceResourceName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = %[5]q
			access_key = %[1]q
			secret_key = %[2]q
		}
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = %[3]q
			provider_name = %[4]q
			region        = %[5]q
		}
		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = aws_vpc.primary.id
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = [aws_subnet.primary-az1.id]
			security_group_ids = [aws_security_group.primary_default.id]
			
		}
		resource "mongodbatlas_privatelink_endpoint_service" %[6]q {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id   =  aws_vpc_endpoint.ptfe_service.id
			private_link_id       = mongodbatlas_privatelink_endpoint.test.id
			provider_name         = %[4]q
		}
		resource "aws_vpc" "primary" {
			cidr_block           = "10.0.0.0/16"
			enable_dns_hostnames = true
			enable_dns_support   = true
		}

		resource "aws_internet_gateway" "primary" {
			vpc_id = aws_vpc.primary.id
		}

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
			availability_zone       = "%[5]sa"
		  }
		  
		  # Subnet-B
		  resource "aws_subnet" "primary-az2" {
			vpc_id                  = aws_vpc.primary.id
			cidr_block              = "10.0.2.0/24"
			map_public_ip_on_launch = false
			availability_zone       = "%[5]sb"
		  }
		  
		  resource "aws_security_group" "primary_default" {
			name_prefix = "default-"
			description = "Default security group for all instances in vpc"
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
		  }
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, serviceResourceName)
}
