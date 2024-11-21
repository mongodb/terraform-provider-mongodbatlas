package advancedclustertpf_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_move_preferred(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	skipIfMoveTestsNotSet(t)
	t.Setenv(advancedclustertpf.MoveModeEnvVarName, advancedclustertpf.MoveModeValPreferred)
	// TODO: temporary no parallel tests so t.Setenv can be used
	resource.Test(t, resource.TestCase{
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

// Skip for now as these needs to either run with mock data or be refactored to use a live cluster
func skipIfMoveTestsNotSet(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_RUN_MOVED_TESTS") == "" {
		t.Skip("Move tests are skipped unless TF_RUN_MOVED_TESTS is set")
	}
}

func TestAccAdvancedCluster_move_rawstate(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	skipIfMoveTestsNotSet(t)
	t.Setenv(advancedclustertpf.MoveModeEnvVarName, advancedclustertpf.MoveModeValRawState)
	resource.Test(t, resource.TestCase{
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

func TestAccAdvancedCluster_move_json(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	skipIfMoveTestsNotSet(t)
	t.Setenv(advancedclustertpf.MoveModeEnvVarName, advancedclustertpf.MoveModeValJSON)
	resource.Test(t, resource.TestCase{
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

func TestAccAdvancedCluster_move_invalid(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName),
			},
			{
				Config:      configMoveSecond(projectID, clusterName),
				ExpectError: regexp.MustCompile("Unable to Move Resource State"),
			},
			{
				Config: configMoveFirst(projectID, clusterName),
			},
		},
	})
}

// TODO: We temporarily use mongodbatlas_database_user instead of mongodbatlas_cluster to set up the initial environment
func configMoveFirst(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "oldtpf" {
			project_id         = %[1]q
			username           = %[2]q # TODO: temporarily we use the username in database_user source as the cluster name in destination
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
