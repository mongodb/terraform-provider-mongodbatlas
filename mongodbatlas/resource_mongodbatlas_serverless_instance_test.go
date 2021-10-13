package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasServerlessInstance_basic(t *testing.T) {
	var (
		serverlessInstance matlas.Cluster
		resourceName       = "mongodbatlas_serverless_instance.test"
		instanceName       = acctest.RandomWithPrefix("test-acc-serverless")
		projectID          = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfig(projectID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasServerlessInstance_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_serverless_instance.test"
		instanceName = acctest.RandomWithPrefix("test-acc-serverless")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfig(projectID, instanceName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasServerlessInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName string, serverlessInstance *matlas.Cluster) resource.TestCheckFunc {
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

		serverlessResponse, _, err := conn.ServerlessInstances.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			*serverlessInstance = *serverlessResponse
			return nil
		}

		return fmt.Errorf("serverless instance (%s) does not exist", ids["name"])
	}
}

func testAccCheckMongoDBAtlasServerlessInstanceDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_serverless_instance" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		serverlessInstance, _, err := conn.ServerlessInstances.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil && serverlessInstance != nil {
			return fmt.Errorf("serverless instance (%s) still exists", ids["name"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasServerlessInstanceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["name"]), nil
	}
}

func testAccMongoDBAtlasServerlessInstanceConfig(projectID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_serverless_instance" "test" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			
			provider_settings_backing_provider_name = "AWS"
			provider_settings_provider_name = "SERVERLESS"
			provider_settings_region_name = "US_EAST_1"
		}

	`, projectID, name)
}
