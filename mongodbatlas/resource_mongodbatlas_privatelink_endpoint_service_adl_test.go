package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceMongoDBAtlasPrivateLinkEndpointServiceADL_basic(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_privatelink_endpoint_service_adl.adl-test"
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		endpointID    = "1"
		commentOrigin = "this is a comment for adl private link endpoint"
		commentUpdate = "this is a comment for adl private link endpoint [UPDATED]"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceADLConfig(projectID, endpointID, commentOrigin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "type", "DATA_LAKE"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceADLConfig(projectID, endpointID, commentUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "type", "DATA_LAKE"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentUpdate),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasPrivateLinkEndpointServiceADL_importBasic(t *testing.T) {
	var (
		//privateLink   = matlas.PrivateLinkEndpointADL{}
		resourceName  = "mongodbatlas_privatelink_endpoint_service_adl.test"
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		endpointID    = ""
		commentOrigin = "this is a comment for adl private link endpoint"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			/*	{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceADLConfig(projectID, endpointID, commentOrigin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(resourceName, &privateLink),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "type", "DATA_LAKE"),
					resource.TestCheckResourceAttr(resourceName, "provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
				),
			},*/
			{
				Config:            testAccMongoDBAtlasPrivateLinkEndpointServiceADLConfig(projectID, endpointID, commentOrigin),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_adl" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		privateLink, _, err := conn.PrivateLinkEndpointsADL.Get(context.Background(), ids["project_id"], ids["endpoint_id"])
		if err == nil && privateLink != nil {
			return fmt.Errorf("endpoint_id (%s) still exists", ids["endpoint_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasPrivateLinkEndpointServiceADLConfig(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint_service_adl" "adl_test" {
			project_id   = "%[1]s"
			endpoint_id  = "%[2]s"
			comment      = "%[3]s"
			type		 = "DATA_LAKE"
			provider_name	 = "AWS"
		}
	`, projectID, endpointID, comment)
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.PrivateLinkEndpointsADL.Get(context.Background(), ids["project_id"], ids["endpoint_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("endpoint_id (%s) does not exist", ids["endpoint_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["endpoint_id"]), nil
	}
}
