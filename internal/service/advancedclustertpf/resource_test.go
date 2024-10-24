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
		}

		data "mongodbatlas_advanced_cluster" "test" {
			group_id = "111111111111111111111111"  # Auto-generated schema hasn't been changed yet, that's why it's group_id not project_id
			cluster_name = "test"
		}

		data "mongodbatlas_advanced_clusters" "tests" {
				group_id = "111111111111111111111111"  # Auto-generated schema hasn't been changed yet, that's why it's group_id not project_id
		}
	`
}
