package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSearchNode_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_search_node.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-search-node")
		clusterName  = acctest.RandomWithPrefix("test-acc-search-node")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasSearchNodeDestroy,
		Steps: []resource.TestStep{
			newSearchNodeTestStep(resourceName, orgID, projectName, clusterName, "S20_HIGHCPU_NVME", 3),
			newSearchNodeTestStep(resourceName, orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
			{
				Config:            testAccMongoDBAtlasSearchNodeConfig(orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckSearchNodeImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func newSearchNodeTestStep(resourceName, orgID, projectName, clusterName, instanceSize string, searchNodeCount int) resource.TestStep {
	return resource.TestStep{
		Config: testAccMongoDBAtlasSearchNodeConfig(orgID, projectName, clusterName, instanceSize, searchNodeCount),
		Check: resource.ComposeTestCheckFunc(
			testAccCheckMongoDBAtlasSearchNodeExists(resourceName),
			resource.TestCheckResourceAttrSet(resourceName, "project_id"),
			resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
			resource.TestCheckResourceAttr(resourceName, "specs.0.instance_size", instanceSize),
			resource.TestCheckResourceAttr(resourceName, "specs.0.node_count", fmt.Sprintf("%d", searchNodeCount)),
			resource.TestCheckResourceAttrSet(resourceName, "state_name"),
		),
	}
}

func testAccMongoDBAtlasSearchNodeConfig(orgID, projectName, clusterName, instanceSize string, searchNodeCount int) string {
	clusterConfig := advancedClusterConfig(orgID, projectName, clusterName)
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_search_node" "test" {
			project_id = mongodbatlas_project.test.id
			cluster_name = %[2]q
			specs = [
				{
					instance_size = %[3]q
					node_count = %[4]d
				}
			]
		}
	`, clusterConfig, clusterName, instanceSize, searchNodeCount)
}

func advancedClusterConfig(orgID, projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q
		cluster_type = "REPLICASET"
		retain_backups_enabled = "true"

		replication_specs {
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
			}
		}
	}
	`, orgID, projectName, clusterName)
}

func testAccCheckSearchNodeImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := decodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["cluster_name"]), nil
	}
}

func testAccCheckMongoDBAtlasSearchNodeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := decodeStateID(rs.Primary.ID)

		connV2 := testAccProviderSdkV2.Meta().(*MongoDBClient).AtlasV2
		_, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(context.Background(), ids["project_id"], ids["cluster_name"]).Execute()
		if err != nil {
			return fmt.Errorf("search node deployment (%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasSearchNodeDestroy(state *terraform.State) error {
	if projectDestroyedErr := testAccCheckMongoDBAtlasProjectDestroy(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	if clusterDestroyedErr := testAccCheckMongoDBAtlasAdvancedClusterDestroy(state); clusterDestroyedErr != nil {
		return clusterDestroyedErr
	}

	connV2 := testAccProviderSdkV2.Meta().(*MongoDBClient).AtlasV2
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_search_node" {
			_, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]).Execute()
			// TODO probably need more logic to look into state.
			if err == nil {
				return fmt.Errorf("search node deployment (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
			}
		}
	}

	return nil
}
