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
		endpointResourceSuffix = "atlasple"
		resourceSuffix         = "atlasrm"
		resourceName           = fmt.Sprintf("mongodbatlas_private_endpoint_regional_mode.%s", resourceSuffix)

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		providerName    = "AWS"
		region          = os.Getenv("AWS_REGION")
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")

		clusterName = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
	)

	endpointResources := testAccMongoDBAtlasPrivateLinkEndpointServiceConfigCompleteAWS(
		awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, endpointResourceSuffix,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix, projectID, clusterName, endpointResources, endpointResourceSuffix, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix, projectID, clusterName, endpointResources, endpointResourceSuffix, true),
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

func testAccMongoDBAtlasPrivateEndpointRegionalModeClusterData(clusterResourceName, projectID, regionalModeResourceName, privateLinkResourceName string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cluster" %[1]q {
			project_id = %[2]q
			name       = mongodbatlas_cluster.%[1]s.name
			depends_on = [
				mongodbatlas_private_endpoint_regional_mode.%[3]s,
				mongodbatlas_privatelink_endpoint_service.%[4]s
			]
		}
	`, clusterResourceName, projectID, regionalModeResourceName, privateLinkResourceName)
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceName, projectID, clusterName, endpointResources, endpointResourceName string, enabled bool) string {
	clusterResourceName := "global_cluster"
	clusterResource := testAccMongoDBAtlasClusterConfigGlobal(clusterResourceName, projectID, clusterName, "false")
	clusterData := testAccMongoDBAtlasPrivateEndpointRegionalModeClusterData(clusterResourceName, projectID, resourceName, endpointResourceName)

	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" %[1]q {
			project_id   = %[2]q
			enabled      = %[3]t
		}

		%[4]s

		%[5]s

		%[6]s
	`, resourceName, projectID, enabled, clusterResource, clusterData, endpointResources)
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
