package advancedclustertpf_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_move_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName),
			},
			{
				Config: configMoveSecond(projectID, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

// TODO: We temporarily use mongodbatlas_database_user instead of mongodbatlas_cluster to set up the initial environment
func configMoveFirst(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "oldtpf" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
			auth_database_name = "admin"
			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`, projectID, clusterName)
}

func configMoveSecond(projectID, clusterName string) string {
	return `
		moved {
			from = mongodbatlas_database_user.oldtpf
			to   = mongodbatlas_advanced_cluster.test
		}
	` + configBasic(projectID, clusterName, "")
}
