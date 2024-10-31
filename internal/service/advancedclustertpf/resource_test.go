package advancedclustertpf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(),
			},
		},
	})
}

func configBasic() string {
	return `
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = "111111111111111111111111"
			name = "test"
			cluster_type = "TYPE"
		}
		data "mongodbatlas_advanced_cluster" "test" {
			project_id = "111111111111111111111111"
			name = "test"

			depends_on = [mongodbatlas_advanced_cluster.test]	
		}

		data "mongodbatlas_advanced_clusters" "tests" {
				project_id = "111111111111111111111111"

				depends_on = [mongodbatlas_advanced_cluster.test]	
		}
	`
}
