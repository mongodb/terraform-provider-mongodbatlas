package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkRSPrivateEndpointRegionalMode_conn(t *testing.T) {
	SkipTestExtCred(t)
	var (
		endpointResourceSuffix = "atlasple"
		resourceSuffix         = "atlasrm"
		resourceName           = fmt.Sprintf("mongodbatlas_private_endpoint_regional_mode.%s", resourceSuffix)

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		providerName = "AWS"
		region       = os.Getenv("AWS_REGION")

		clusterName         = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
		clusterResourceName = "global_cluster"
	)

	clusterResource := testAccMongoDBAtlasClusterConfigGlobal(clusterResourceName, orgID, projectName, clusterName, "false")
	clusterDataSource := testAccMongoDBAtlasPrivateEndpointRegionalModeClusterData(clusterResourceName, resourceSuffix, endpointResourceSuffix)
	endpointResources := testAccMongoDBAtlasPrivateLinkEndpointServiceConfigUnmanagedAWS(
		awsAccessKey, awsSecretKey, projectID, providerName, region, endpointResourceSuffix,
	)

	dependencies := []string{clusterResource, clusterDataSource, endpointResources}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy,
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
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix,orgID, projectName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceSuffix, orgID, projectName, true),
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

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(resourceName, orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[3]q
			org_id = %[2]q
		}
		resource "mongodbatlas_private_endpoint_regional_mode" %[1]q {
			project_id   = mongodbatlas_project.test.id
			enabled      = %[4]t
		}
	`, resourceName, orgID, projectName, enabled)
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf("==========================================================================\n")
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

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate(projectID, clusterName, clusterResourceName string) resource.TestCheckFunc {
	resourceName := strings.Join([]string{"data", "mongodbatlas_cluster", clusterResourceName}, ".")
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Could not find resource state for cluster (%s) on project (%s)", clusterName, projectID)
		}

		var rsPrivateEndpointCount int
		var err error

		if rsPrivateEndpointCount, err = strconv.Atoi(rs.Primary.Attributes["connection_strings.0.private_endpoint.#"]); err != nil {
			return fmt.Errorf("Connection strings private endpoint count is not a number")
		}

		cluster, _, _ := conn.Clusters.Get(context.Background(), projectID, clusterName)

		fmt.Printf("testAccCheckMongoDBAtlasPrivateEndpointRegionalModeClustersUpToDate %#v \n", rs.Primary.Attributes)
		fmt.Printf("cluster.ConnectionStrings %#v \n", flattenConnectionStrings(cluster.ConnectionStrings))

		if rsPrivateEndpointCount != len(cluster.ConnectionStrings.PrivateEndpoint) {
			return fmt.Errorf("Cluster PrivateEndpoint count does not match resource")
		}

		if rs.Primary.Attributes["connection_strings.0.standard"] != cluster.ConnectionStrings.Standard {
			return fmt.Errorf("Cluster standard connection_string does not match resource")
		}

		if rs.Primary.Attributes["connection_strings.0.standard_srv"] != cluster.ConnectionStrings.StandardSrv {
			return fmt.Errorf("Cluster standard connection_string does not match resource")
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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
