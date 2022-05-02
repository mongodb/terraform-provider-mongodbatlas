package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasPrivateEndpointRegionalMode_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_private_endpoint_regional_mode.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID, clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID, clusterName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID, clusterName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id   = %[1]q
			enabled      = %[2]t
		}

		%[3]s
	`, projectID, enabled, testAccMongoDBAtlasClusterConfigGlobal(projectID, clusterName, "false"))
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		status, _, _ := conn.Clusters.Status(context.Background(), projectID, clusterName)

		if status.ChangeStatus == matlas.ChangeStatusPending {
			return fmt.Errorf("cluster (%s) for project (%s) still has changes PENDING", clusterName, projectID)
		}

		rs, ok := s.RootModule().Resources["mongodbatlas_cluster.global_cluster"]

		if !ok {
			return fmt.Errorf("Could not find resource state for cluster (%s) on project (%s)", clusterName, projectID)
		}

		cluster, _, _ := conn.Clusters.Get(context.Background(), projectID, clusterName)

		if reflect.DeepEqual(rs.Primary.Attributes["connection_strings"], flattenConnectionStrings(cluster.ConnectionStrings)) {
			return nil
		}

		fmt.Printf("resource.Primary.Attributes['connection_strings'] %#v \n", rs.Primary.Attributes["connection_strings"])
		fmt.Printf("cluster.ConnectionStrings %#v \n", flattenConnectionStrings(cluster.ConnectionStrings))

		return fmt.Errorf("Connection strings not equal in resource state and response from apif or cluster (%s) on project (%s)", clusterName, projectID)
	}
}
