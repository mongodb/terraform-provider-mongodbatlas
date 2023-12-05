package privateendpointregionalmode_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkRSPrivateEndpointRegionalMode_conn(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		endpointResourceSuffix = "atlasple"
		resourceSuffix         = "atlasrm"
		resourceName           = fmt.Sprintf("mongodbatlas_private_endpoint_regional_mode.%s", resourceSuffix)

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		providerName = "AWS"
		region       = os.Getenv("AWS_REGION")

		clusterName         = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
		clusterResourceName = "global_cluster"
	)

	clusterResource := acc.ConfigClusterGlobal(clusterResourceName, orgID, projectName, clusterName, "false")
	clusterDataSource := testAccMongoDBAtlasPrivateEndpointRegionalModeClusterData(clusterResourceName, resourceSuffix, endpointResourceSuffix)
	endpointResources := testAccMongoDBAtlasPrivateLinkEndpointServiceConfigUnmanagedAWS(
		awsAccessKey, awsSecretKey, projectID, providerName, region, endpointResourceSuffix,
	)

	dependencies := []string{clusterResource, clusterDataSource, endpointResources}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfigWithDependencies(resourceSuffix, projectID, false, dependencies),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName, clusterResourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfigWithDependencies(resourceSuffix, projectID, true, dependencies),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName, clusterResourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccNetworkRSPrivateEndpointRegionalMode_basic(t *testing.T) {
	var (
		resourceSuffix = "atlasrm"
		resourceName   = fmt.Sprintf("mongodbatlas_private_endpoint_regional_mode.%s", resourceSuffix)
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix, projectID, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix, projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeClusterData(clusterResourceName, regionalModeResourceName, privateLinkResourceName string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cluster" %[1]q {
			project_id = mongodbatlas_cluster.%[1]s.project_id
			name       = mongodbatlas_cluster.%[1]s.name
			depends_on = [
				mongodbatlas_privatelink_endpoint_service.%[3]s,
				mongodbatlas_private_endpoint_regional_mode.%[2]s
			]
		}
	`, clusterResourceName, regionalModeResourceName, privateLinkResourceName)
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfigWithDependencies(resourceName, projectID string, enabled bool, dependencies []string) string {
	resources := make([]string, len(dependencies)+1)

	resources[0] = testAccMongoDBAtlasPrivateEndpointRegionalModeConfigNoProject(resourceName, projectID, enabled)
	copy(resources[1:], dependencies)

	return strings.Join(resources, "\n\n")
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfigNoProject(resourceName, projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" %[1]q {
			project_id   = %[2]q
			enabled      = %[3]t
		}
	`, resourceName, projectID, enabled)
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceName, projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" %[1]q {
			project_id   = %[2]q
			enabled      = %[3]t
		}
	`, resourceName, projectID, enabled)
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf("==========================================================================\n")
		conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		projectID := rs.Primary.ID

		_, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), projectID)

		if err == nil {
			return nil
		}

		return fmt.Errorf("regional mode for project_id (%s) does not exist", projectID)
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName, clusterResourceName string) resource.TestCheckFunc {
	resourceName := strings.Join([]string{"data", "mongodbatlas_cluster", clusterResourceName}, ".")
	return func(s *terraform.State) error {
		conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Could not find resource state for cluster (%s) on project (%s)", clusterName, projectID)
		}

		var rsPrivateEndpointCount int
		var err error

		if rsPrivateEndpointCount, err = strconv.Atoi(rs.Primary.Attributes["connection_strings.0.private_endpoint.#"]); err != nil {
			return fmt.Errorf("Connection strings private endpoint count is not a number")
		}

		c, _, _ := conn.Clusters.Get(context.Background(), projectID, clusterName)

		fmt.Printf("testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate %#v \n", rs.Primary.Attributes)
		fmt.Printf("cluster.ConnectionStrings %#v \n", advancedcluster.FlattenConnectionStrings(c.ConnectionStrings))

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

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy(s *terraform.State) error {
	conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_private_endpoint_regional_mode" {
			continue
		}

		setting, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Could not read regionalized private endpoint setting for project %q", rs.Primary.ID)
		}

		if setting.Enabled != false {
			return fmt.Errorf("Regionalized private endpoint setting for project %q was not properly disabled", rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasPrivateLinkEndpointServiceConfigUnmanagedAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, serviceResourceName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "%[5]s"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
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
