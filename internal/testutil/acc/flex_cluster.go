package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	FlexDataSourceName       = "data.mongodbatlas_flex_cluster.test"
	FlexDataSourcePluralName = "data.mongodbatlas_flex_clusters.test"
	FlexDataSource           = `
	data "mongodbatlas_flex_cluster" "test" {
		project_id = mongodbatlas_flex_cluster.test.project_id
		name       = mongodbatlas_flex_cluster.test.name
		
		depends_on = [mongodbatlas_flex_cluster.test]
	}
	data "mongodbatlas_flex_clusters" "test" {
		project_id = mongodbatlas_flex_cluster.test.project_id
		
		depends_on = [mongodbatlas_flex_cluster.test]
	}`
)

func CheckDestroyFlexCluster(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_advanced_cluster" && rs.Type != "mongodbatlas_flex_cluster" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		name := rs.Primary.Attributes["name"]
		_, _, err := ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
		if err == nil {
			return fmt.Errorf("flex cluster (%s:%s) still exists", projectID, name)
		}
	}
	return nil
}

func CheckExistsFlexCluster() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_advanced_cluster" && rs.Type != "mongodbatlas_flex_cluster" {
				continue
			}
			projectID := rs.Primary.Attributes["project_id"]
			name := rs.Primary.Attributes["name"]
			_, _, err := ConnV2().FlexClustersApi.GetFlexCluster(context.Background(), projectID, name).Execute()
			if err != nil {
				return fmt.Errorf("flex cluster (%s:%s) not found", projectID, name)
			}
		}
		return nil
	}
}
