package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasPrivateEndpointLinkAWS_basic(t *testing.T) {
	//SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_private_endpoint_service_link.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		vpcID        = os.Getenv("AWS_VPC_ENDPOINT_ID")
		region       = "us-east-1"
		providerName = "AWS"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//checkPeeringEnvAWS(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkConfigBasicAWS(projectID, providerName, vpcID, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_service_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasPrivateEndpointLinkAWS_Complete(t *testing.T) {
	//SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_private_endpoint_service_link.test"

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		providerName    = "AWS"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		region          = os.Getenv("AWS_REGION")
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) /*; checkAwsEnv(t); checkPeeringEnvAWS(t) */ },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkConfigCompleteAWS(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_endpoint_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasPrivateEndpointLinkAWS_import(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_private_endpoint_service_link.test"

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		providerName    = "AWS"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		region          = os.Getenv("AWS_REGION")
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkConfigCompleteAWS(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_endpoint_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasPrivateEndpointLinkImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasPrivateEndpointLinkAzure_basic(t *testing.T) {
	//SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_private_endpoint_service_link.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		azureID      = os.Getenv("AZURE_SUBCRIPTION_ID")
		ipAddress    = ""
		region       = "US_EAST_2"
		providerName = "AZURE"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//checkPeeringEnvAWS(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkConfigBasicAzure(projectID, providerName, region, azureID, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_service_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasPrivateEndpointLinkAzure_import(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_private_endpoint_service_link.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		azureID      = os.Getenv("AZURE_SUBCRIPTION_ID")
		ipAddress    = ""
		region       = "US_EAST_2"
		providerName = "AZURE"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkConfigBasicAzure(projectID, providerName, region, azureID, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_endpoint_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasPrivateEndpointLinkImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateEndpointLinkImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["private_link_id"], ids["interface_endpoint_id"]), nil
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["interface_endpoint_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("the MongoDB Interface Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["interface_endpoint_id"], rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointLinkDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_private_endpoint" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["interface_endpoint_id"])
		if err == nil {
			return fmt.Errorf("the MongoDB Private Endpoint(%s) still exists", ids["interface_endpoint_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasPrivateEndpointLinkConfigCompleteAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "us-east-1"
			access_key = "%s"
			secret_key = "%s"
		}

		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%[4]s"
			region        = "%s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%s"
			service_name       = "${mongodbatlas_private_endpoint.test.endpoint_service_name}"
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%s"]
			security_group_ids = ["%s"]
		}

		resource "mongodbatlas_private_endpoint_service_link" "test" {
			project_id            = "${mongodbatlas_private_endpoint.test.project_id}"
			private_link_id       =  "${aws_vpc_endpoint.ptfe_service.id}"
			endpoint_service_id = "${mongodbatlas_private_endpoint.test.private_link_id}"
			provider_name = "%[4]s"
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID)
}

func testAccMongoDBAtlasPrivateEndpointLinkConfigBasicAWS(projectID, providerName, vpcId, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%[1]s"
			provider_name = "%[2]s"
			region        = "%[3]s"
		}
		resource "mongodbatlas_private_endpoint_service_link" "test" {
			project_id            = "%[1]s"
			private_link_id       = "%[4]s"
			provider_name = "%[2]s"
			endpoint_service_id = "${mongodbatlas_private_endpoint.test.private_link_id}"
		}
	`, projectID, providerName, region, vpcId)
}

func testAccMongoDBAtlasPrivateEndpointLinkConfigBasicAzure(projectID, providerName, region, azureID, ipAddress string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%[1]s"
			provider_name = "%[2]s"
			region        = "%[3]s"
		}
		resource "mongodbatlas_private_endpoint_service_link" "test" {
			project_id            = "%[1]s"
			private_link_id       = "%[4]s"
			provider_name = "%[2]s"
			endpoint_service_id = "${mongodbatlas_private_endpoint.test.private_link_id}"
			private_endpoint_ip_address       = "%[5]s"
		}
	`, projectID, providerName, region, azureID, ipAddress)
}
